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
	"github.com/gambol99/flare/pkg/config"
	"github.com/gambol99/flare/pkg/rules"

	log "github.com/Sirupsen/logrus"
	"github.com/gambol99/flare/pkg/utils"
)

// the implementation for the flare service
type Flare struct {
	// the rules store
	store rules.RulesService
	// the configuration
	cfg *config.FlareConfiguration
	//
	// the shutdown signal
	shutdown_channel utils.ShutdownChannel
	// the reference to the container service
	containers ContainerService
	// the channel which container events are sent
	container_events ContainerEventsChannel
	// the channel we receive events from the rule service
	rules_events rules.RulesEventChannel
	// the map of the containers which has been processed
}

// Creates a new Flare Service
//	cfg:		the configuration for the service
func New(cfg *config.FlareConfiguration) (FlareService, error) {
	log.Infof("Initializing the Flare Service, configuation: %s", cfg)

	var err error
	flare := new(Flare)
	flare.cfg = cfg
	flare.container_events =make(ContainerEventsChannel, 10)
	flare.shutdown_channel = make(utils.ShutdownChannel)

	// step: create the container service
	flare.containers, err = NewContainerService(cfg.DockerSocket)
	if err != nil {
		log.Errorf("Failed to create the container service, error: %s", err)
		return nil, err
	}

	// step: create the rules service
	flare.store, err = rules.New(cfg)
	if err != nil {
		log.Errorf("failed to create the rules store, error: %s", err)
		return nil, err
	}

	// step: we need to perform a preprocessing of the container and look for anyone
	// which might require us


	// step: add the events listeners
	flare.store.AddListener(flare.rules_events)
	flare.containers.AddListener(flare.container_events)

	// step: perform


	// step: start processing the events

	return flare, nil
}
