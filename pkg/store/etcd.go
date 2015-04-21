/*
Copyright 2014 Rohith All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package store

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gambol99/flare/pkg/config"
	"github.com/gambol99/flare/pkg/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/go-etcd/etcd"
)

// the etcd implementation of a store
type EtcdStoreAgent struct {
	// a lock for the watcher map
	sync.RWMutex
	// the base root key
	prefix string
	// the etcd client
	client *etcd.Client
	// the etcd hosts extracted from the url
	hosts []string
	// stop channel for the client
	stopChannel utils.ShutdownChannel
	// the update channel we send our changes to
	eventsChannel StoreEventChannel
	/* a map of keys presently being watched */
	watchedKeys map[string]bool
}

func EtcdOptions() string {
	return `
    Options:
        cacert:     the ca certificate used to commiunicate with ectd
        key:        the certificate key for etcd
    `
}

const (
	ETCD_CACERT = "cacert"
	ETCC_KEY    = "key"
	ETCD_CAFILE = "cafile"
)

// Create a new Etcd Agent for the store
//  location:       the url for the etcd store
//  options:        a map of options for the etcd client
func NewEtcdStore(location *url.URL, namespace string, options config.ConfigutationOptions) (Store, error) {
	log.Infof("Creating a Etcd Store agent, url: %s, namespace: %s", location, namespace)

	var err error
	store := new(EtcdStoreAgent)
	store.prefix = namespace
	store.watchedKeys = make(map[string]bool, 0)
	store.hosts = make([]string, 0)
	store.stopChannel = make(utils.ShutdownChannel)

	// step: are we using tls or not
	store_protocol := "http"
	if options.Get(ETCD_CACERT, "") != "" {
		store_protocol = "https"
	}

	// step: lets get the hosts from the url
	for _, host := range strings.Split(location.Host, ",") {
		store.hosts = append(store.hosts, fmt.Sprintf("%s://%s", store_protocol, host))
	}

	log.Infof("Creating a Etcd Agent, hosts: %s", store.hosts)

	// step: create the etcd client
	if options.Get(ETCD_CACERT, "") != "" {
		store.client, err = etcd.NewTLSClient(store.hosts,
			options.Get(ETCD_CACERT, "").(string),
			options.Get(ETCC_KEY, "").(string),
			options.Get(ETCD_CAFILE, "").(string))
		if err != nil {
			log.Errorf("unable to create a TLS connection to etcd: %s, error: %s", location, err)
			return nil, err
		}
	} else {
		store.client = etcd.NewClient(store.hosts)
	}

	// step: start watching out for events
	go store.watchEvents()

	log.Debugf("Starting off the Etcd event processor on prefix: %s", store.prefix)

	return store, nil
}

func (r *EtcdStoreAgent) watchEvents() {
	log.Infof("Starting the event watcher for the etcd client, channel: %v", r.eventsChannel)
	// the kill switch for the goroutine
	kill_off := false

	// routine: waits on the shutdown signal for the client and flicks the kill switch
	go func() {
		log.Infof("Waiting on a shutdown signal from consumer, channel: %v", r.eventsChannel)
		// step: wait for the shutdown signal
		<-r.stopChannel
		log.Infof("Flicking the kill switch for watcher, channel: %v", r.prefix, r.eventsChannel)
		kill_off = true
	}()

	// routine: loops around watching until flick the switch
	go func() {
		wait_index := uint64(0)
		// step: look until we hit the kill switch
		for {
			if kill_off {
				break
			}
			// step: apply a watch on the key and wait
			response, err := r.client.Watch(r.prefix, wait_index, true, nil, nil)
			log.Debugf("Received a event from etcd store, key: %v", response)
			if err != nil {
				log.Errorf("Failed to attempting to watch the key: %s, error: %s", r.prefix, err)
				time.Sleep(time.Duration(3) * time.Second)
				wait_index = uint64(0)
				continue
			}
			// step: have we been requested to quit
			if kill_off {
				continue
			}
			// step: update the wait index
			wait_index = response.Node.ModifiedIndex + 1
			// step: cool - we have a notification - lets check if this key is being watched
			go r.processNodeChange(response)
		}
		log.Infof("Exitted the k/v watcher routine, channel: %v", r.eventsChannel)
	}()
}

// Make a directory structure
//	path		the directory you want to create
func (r *EtcdStoreAgent) Mkdir(path string) error {
	log.Debugf("Creating the directory: %s", path)
	response, err := r.client.CreateDir(path, 0)
	if err != nil {
		if strings.Contains(err.Error(), "Key already exists") {
			if response, err = r.client.Get(path, false, false); err != nil {
				return fmt.Errorf("unable to retrieve the path: %s, error: %s", path, err)
			} else if !response.Node.Dir {
				return fmt.Errorf("a non directory entry already exists for: %s", path)
			}
			return nil
		}
		log.Errorf("Failed to create the directory: %s, error: %s", path, err)
		return err
	}
	return nil
}

// List all the keys under a directory
// 	path:		the key for the directory
//	recursive:	do you want the listing to include subdirectories
func (r *EtcdStoreAgent) List(path string, recursive bool) ([]string, error) {
	list := make([]string, 0)
	if recursive {
		list, err := r.recursiveList(path, &list)
		if err != nil {
			return nil, err
		}
		return list, nil
	} else {
		// get the root node
		response, err := r.client.Get(path, true, false)
		if err != nil {
			return list, err
		}

		// is the node a directory?
		if !response.Node.Dir {
			return list, fmt.Errorf("the key: %s is not a directory")
		}

		// iterate and get the children
		for _, node := range response.Node.Nodes {
			if node.Dir {
				continue
			}
			list = append(list, node.Key)
		}
		return list, nil
	}
}

// Retrieve a value from the store
//  id:     the key you wish to retrieve
func (r *EtcdStoreAgent) Get(id string) (string, error) {
	log.Debugf("Retrieving key: %s from the etcd store", id)
	response, err := r.client.Get(id, false, true)
	if err != nil {
		log.Errorf("failed to get the key: %s, error: %s", id, err)
		return "", err
	}
	return response.Node.Value, nil
}

// Add the key to the etcd store
//  id:     the key you wish to add to the store
//  value:  the value of the key
func (r *EtcdStoreAgent) Set(id string, value string) error {
	log.Debugf("Setting the key: %s, value: %s in etcd", id, value)
	return utils.Attempt(3, func() error {
		if _, err := r.client.Set(id, value, 0); err != nil {
			log.Errorf("unable to set the key: %s, error: %s", id, err)
			return err
		}
		return nil
	})
}

// Delete a or all keys under the namespace
// 	path:		the path to the namespace
//	recursive:	whether or not to be recursive
func (r *EtcdStoreAgent) Delete(path string, recursive bool) error {
	log.Debugf("Deleting the key/s on/under: %s, recursive: %t", path, recursive)
	_, err := r.client.Delete(path, recursive)
	return err
}

// Delete all the child of a directory
// 	path:		the path of directory
func (r *EtcdStoreAgent) DeleteAll(path string) error {
	log.Debugf("Deleting all the entries under path: %s", path)
	if found, err := r.Exists(path); err != nil {
		return err
	} else if !found {
		fmt.Errorf("the path: %s does not exist", path)
	}

	is_directory, err := r.isDirectory(path)
	if err != nil {
		return err
	} else if !is_directory {
		return fmt.Errorf("path: %s is not a directory", path)
	}

	listings, err := r.List(path, false)
	if err != nil {
		return err
	}
	// delete the files
	for _, file := range listings {
		if _, err := r.client.Delete(file, true); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown the resources
func (r *EtcdStoreAgent) Close() error {
	log.Infof("Shutting down the Etcd Agent")
	r.stopChannel <- true
	return nil
}

// Checks to see if the key exists in the store
//  id:     the id / key which you are checking for
func (r *EtcdStoreAgent) Exists(id string) (bool, error) {
	log.Debugf("Checking if key: %s exists in etcd store", id)
	if _, err := r.client.Get(id, false, false); err != nil {
		if strings.Contains(err.Error(), "Key not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Destroy and delete all the keys under the store prefix
func (r *EtcdStoreAgent) Flush() error {
	log.Infof("Flushing and deleting all the keys under prefix: %s", r.prefix)
	_, err := r.client.Delete(r.prefix, true)
	if err != nil {
		return err
	}
	return nil
}

// Add a watch on the namespace
//  key:    the namespace to start watching
func (r *EtcdStoreAgent) Watch(key string) error {
	log.Debugf("Adding a watch on the key space: %s", key)
	if !r.isValidKey(key) {
		return fmt.Errorf("invalid key: %s", key)
	}
	r.Lock()
	defer r.Unlock()
	r.watchedKeys[key] = true
	return nil
}

// Add a event listener to the changes on keys
//  ch:     the channel you wish event's to appear on
func (r *EtcdStoreAgent) AddListener(ch StoreEventChannel) {
	r.Lock()
	defer r.Unlock()
	log.Debugf("Adding a listener to etcd event, channel: %v", ch)
	r.eventsChannel = ch
}

func (r *EtcdStoreAgent) recursiveList(path string, paths *[]string) ([]string, error) {
	response, err := r.client.Get(path, true, true)
	if err != nil {
		return nil, err
	}
	for _, node := range response.Node.Nodes {
		if node.Dir {
			r.recursiveList(node.Key, paths)
		} else {
			*paths = append(*paths, node.Key)
		}
	}
	return *paths, nil
}

func (r *EtcdStoreAgent) isDirectory(path string) (bool, error) {
	response, err := r.client.Get(path, false, false)
	if err != nil {
		return false, err
	}
	return response.Node.Dir, nil
}

func (r *EtcdStoreAgent) isValidKey(key string) bool {
	if strings.HasPrefix(key, "/") {
		return true
	}
	return false
}

func (r *EtcdStoreAgent) processNodeChange(response *etcd.Response) {
	// step: are there any keys being watched
	if len(r.watchedKeys) <= 0 || r.eventsChannel == nil {
		log.Debugf("Skipping the event on: %s as no one is listening yet", response.Node.Key)
		return
	}

	// step: iterate the list and find out if our key is being watched
	path := response.Node.Key
	log.Debugf("Checking if key: %s is being watched", path)
	for watch_key := range r.watchedKeys {
		if strings.HasPrefix(path, watch_key) {
			log.Debugf("Sending notification of change on key: %s, channel: %v, event: %v", path, r.eventsChannel, response)
			// step: we create an event and send upstream
			event := new(StoreEvent)
			event.ID = response.Node.Key
			event.Directory = response.Node.Dir
			event.Value = response.Node.Value
			switch response.Action {
			case "set":
				event.Action = CHANGED
			case "delete":
				event.Action = DELETED
			}
			// step: send the event upstream via the channel
			r.eventsChannel <- event
			return
		}
	}
}
