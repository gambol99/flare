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

package service

import (
	"strings"
	"sync"

	"github.com/gambol99/flare/pkg/utils"

	log "github.com/Sirupsen/logrus"
	docker "github.com/fsouza/go-dockerclient"
	"fmt"
	"regexp"
)

const (
	CONTAINER_STARTED = 0 << iota
	CONTAINER_DIED
)

const (
	DOCKER_START   = "start"
	DOCKER_DIE     = "die"
	DOCKER_DESTROY = "destroy"
)

type Docker struct {
	sync.RWMutex
	// the docker client */
	client *docker.Client
	// the shutdown signal
	shutdown utils.ShutdownChannel
	// a map those listening
	listeners map[ContainerEventsChannel]bool

}

// Create a new container service
//	socket:			the path of the docker socket
func NewContainerService(socket string) (ContainerService, error) {
	log.Infof("Creating a docker store service, socket: %s", socket)
	store := new(Docker)
	store.listeners = make(map[ContainerEventsChannel]bool, 0)
	// step: lets create the docker client
	if client, err := docker.NewClient("unix://" + socket); err != nil {
		log.Errorf("Failed to create a docker client, socket: %s, error: %s", socket, err)
		return nil, err
	} else {
		store.client = client
		store.shutdown = make(utils.ShutdownChannel)
		if err := store.client.Ping(); err != nil {
			log.Errorf("Failed to ping via the docker client, errorr: %s", err)
			return nil, err
		}
		// step: lets create the docker events
		if err := store.EventProcessor(); err != nil {
			log.Errorf("Failed to start the events processor, error: %s", err)
			return nil, err
		}
	}
	return store, nil
}

// Retrieve a list of container current running
func (r *Docker) List() ([]docker.APIContainers, error) {
	if containers, err := r.client.ListContainers(docker.ListContainersOptions{}); err != nil {
		log.Errorf("Failed to retrieve a list of container from docker, error: %s", err)
		return nil, err
	} else {
		return containers, nil
	}
}

// Retrieve via inspection the container information
//  id:     the container id you are looking for
func (r *Docker) Get(id string) (*docker.Container, error) {
	if container, err := r.client.InspectContainer(id); err != nil {
		log.Errorf("Failed to retrieve a container: %s from docker, error: %s", id, err)
		return nil, err
	} else {
		return container, nil
	}
}

// We docker inspect and pull the environment variables of the container
func (r *Docker) Environment(id string) (map[string]string, error) {
	// step: check the container exists
	if found, err := r.Exists(id); err != nil {
		return nil, fmt.Errorf("unable to retrieve the status of the container: %s, error: %s", id, err)
	} else if !found {
		return nil, fmt.Errorf("the container: %s does not exsit", id)
	}

	// step: docker inspect the container
	info, err := r.client.InspectContainer(id)
	if err != nil {
		log.Errorf("Failed to inspect the container: %s, error: %s", id, err)
		return nil, err
	}

	// step: iterate the
	environment := make(map[string]string, 0)
	for _, kv := range info.Config.Env {
		if found, _ := regexp.MatchString(`^(.*)=(.*)$`, kv); found {
			elements := strings.SplitN(kv, "=", 2)
			environment[elements[0]] = elements[1]
		} else {
			log.Debugf("Invalid environment variable: %s, skipping", kv)
		}
	}

	return environment, nil
}

// Check to see if a container exists
//  id:     the container id you are looking for
func (r *Docker) Exists(id string) (bool, error) {
	if _, err := r.client.InspectContainer(id); err != nil {
		if strings.HasPrefix("No such container", err.Error()) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// The EventProcessor listens out for events from docker and passes them
// upstream to the appropriate channel
func (r *Docker) EventProcessor() error {

	log.Infof("Starting the Docker Events Processor")
	update_channel := make(chan *docker.APIEvents, 5)
	// step: add ourself a
	if err := r.client.AddEventListener(update_channel); err != nil {
		log.Errorf("Failed to add ourselve as a docker events listener, error: %s", err)
		return err
	}
	// step: start the events processor
	go func() {
		log.Infof("Starting the events processor for docker events")
		for {
			select {
			case event := <-update_channel:
				log.Infof("Receivied a docker event, id: %s, status: %s", event.ID[:12], event.Status)
				// step: we ONLY care about start and destroy events
				if event.Status != DOCKER_DESTROY || event.Status != DOCKER_START {
					log.Debugf("Skipping the event: %s, not relevant to us", event.Status)
					break
				}
				// step: we create a container event
				container_event := new(ContainerEvent).
					ContainerID(event.ID).SetStatus(event.Status)

				// step: send the event upstream
				r.pushNotification(container_event)
			case <-r.shutdown:
				log.Infof("Recieved a shutdown signal from above, closing up resources")
				r.client.RemoveEventListener(update_channel)
				log.Infof("Exitting the events processor loop")
				return
			}
		}
	}()
	return nil
}

// Notify the channel when something has happened to the containers
// Params:
//		channel:	the channel to send the event upon
func (r *Docker) AddListener(channel ContainerEventsChannel) error {
	r.Lock()
	defer r.Unlock()
	log.Debugf("Adding the listen: %v to the container listeners", channel)
	r.listeners[channel] = true
	return nil
}

// Push the events upstream to the listeners
func (r *Docker) pushNotification(event *ContainerEvent) {
	r.RLock()
	defer r.RUnlock()
	for listener, _ := range r.listeners {
		go func() {
			listener <- *event
		}()
	}
}
