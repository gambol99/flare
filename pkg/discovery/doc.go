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

// a channel to receive updates on
type ServiceEvent chan string

//  The interface to a discovery provider
type Discovery interface {
	// retrieve a list of services
	Services() ([]*Service, error)
	// get the service from the provider
	Get(string) (*Service, error)
	// get the endpoints for a service
	Endpoints(string, bool) ([]*Endpoint, error)
	// add a listener for changes in the services
	AddServiceListener(ServiceEvent)
	// watch a service for changes
	Watch(string)
	// remove the service from being watched
	Unwatch(string)
	// a list of those service being watched
	Watching() []string
	// close and release the resource
	Close() error
}

// the normalized definition for a service
type Service struct {
	// the identifier for the service
	ID string
	// the name of a service
	Name string
	// the tags associated to the service
	Tags []string
	// a list of endpoints for the service
	Endpoints []*Endpoint
}

// the definition for a endpoint
type Endpoint struct {
	// the ip address of the endpoint
	Address string
	// the port of the endpoint
	Port int
}
