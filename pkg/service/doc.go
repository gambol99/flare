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
	docker "github.com/fsouza/go-dockerclient"
)

const (
	VERSION = "0.0.1"
	AUTHOR  = "Rohith Jayawardene"
	EMAIL   = "gambol99@gmail.com"
)

// the flare service
type FlareService interface {
}

// A container events
type ContainerEvent struct {
	// the if of the container in question
	ID string
	// the type of event
	Status string
}

// a channel used to send events about containers
type ContainerEventsChannel chan ContainerEvent

//
//  The interface for interacting with the containers via docker API
//
type ContainerService interface {
	// pull a list of containers
	List() ([]docker.APIContainers, error)
	// retrieve a information on a specific docker
	Get(string) (*docker.Container, error)
	// check if a container exists
	Exists(string) (bool, error)
	// get the environment of the container
	Environment(string) (map[string]string, error)
	// listen for container creations
	AddListener(ContainerEventsChannel) error
}
