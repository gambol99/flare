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

    "github.com/gambol99/flare/pkg/utils"

    log "github.com/Sirupsen/logrus"
)

//
//  The processor is the handles for events coming in from the rules store and container service
//
func (r *Flare) watchEvents() {
    // atomic switch
    var lock utils.AtomicSwitch

    // step: watch for the shutdown channel
    go func() {
        log.Debugf("Flare waiting for the shutdown signal")
        <- r.shutdown_channel
        lock.SwitchOn()
    }()

    // step: create the goroutine to handle incoming events from the container and
    // rules service
    go func() {
        for lock.IsSwitched() {
            select {
                // an event from the container service
                case event := <- r.container_events:
                    log.Infof("Container event: %s, processing the container: %s", )
                    if event.HasStarted() {
                        go r.processContainerStarted(event)
                    }
                    if event.HasDestroyed() {
                        go r.processContainerDestruction(event)
                    }
                // a event from the rules service
                case event := <- r.rules_events:
                    log.Infof("Rules event: %s processing the event now", event)

                // a signal to shutdown the service
                case <- r.shutdown_channel:


            }
        }
    }()


}

// Process the container destruction and remove any chains which has been pushed
//  event:      the container event itself
func (r *Flare) processContainerStarted(event ContainerEvent) {
    // step: we inspect the container and look for any flare specs
    log.Debugf("Retrieving the container environment, id: %s", event.ID)
    environment, err := r.containers.Environment(event.ID)
    if err != nil {
        log.Errorf("Failed to retrieve the environment from the container: %s, error: %s",
            event.ID, err)

        // step: @TODO we should be push this into a background and retry
        return
    }
    _ = environment
    // step: iterate the environment and look for


}

// Process the container destruction and remove any chains which has been pushed
//  event:      the container event itself
func (r *Flare) processContainerDestruction(event ContainerEvent) {



}