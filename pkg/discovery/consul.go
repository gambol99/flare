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

package discovery

import (
	"fmt"
	"net/url"

	"sync"

	log "github.com/Sirupsen/logrus"
	consulapi "github.com/hashicorp/consul/api"
)

// the consul discovery provider implementation
type ConsulAgent struct {
	sync.RWMutex
	// the consul api client
	client *consulapi.Client
	// the current wait index
	wait_index uint64
	// the kill off
	kill_off bool
	// a map of service being watched
	watching map[string]int
}

func NewConsulAgent(location *url.URL) (Discovery, error) {
	log.Infof("Initializing the Consul Agent, url: %s", location)
	var err error
	service := new(ConsulAgent)
	service.watching = make(map[string]int, 0)
	cfg := consulapi.DefaultConfig()
	cfg.Address = location.Host
	service.client, err = consulapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create a consul agent, error: %s", err)
	}
	return service, nil
}

func (r *ConsulAgent) Services() ([]*Service, error) {
	log.Debugf("retrieving a list of services from consul")

	return nil, nil
}

func (r *ConsulAgent) Get(string) (*Service, error) {

	return nil, nil
}

func (r *ConsulAgent) Endpoints(string, bool) ([]*Endpoint, error) {

	return nil, nil
}

func (r *ConsulAgent) AddServiceListener(ServiceEvent) {

}

func (r *ConsulAgent) Watch(name string) {
	log.Debugf("Adding a consul watch on the service: %s", name)
	r.Lock()
	defer r.Unlock()

}

func (r *ConsulAgent) Unwatch(string) {

}

func (r *ConsulAgent) Watching() []string {

	return nil
}

func (r *ConsulAgent) Close() error {

	return nil
}

func (r *ConsulAgent) getService(name string) (*Service, error) {
	log.Debugf("retreiving the service: %s from consul", name)
	//catalog := r.client.Catalog()

	return nil, nil
}
