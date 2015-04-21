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

package rules

import (
    "fmt"

    log "github.com/Sirupsen/logrus"
)

// Add a listener to events regarding changes in the rules
//	ch:		the channel we should send the event on
func (r *RulesStore) AddListener(ch RulesEventChannel) error {
    if ch == nil {
        return fmt.Errorf("you have not passed a channel to send upon")
    }
    r.Lock()
    defer r.Unlock()
    log.Infof("Adding the listener: %v to the rules service", ch)
    r.listeners[ch] = true
    return nil
}

// Remove a listen from the rule store
//	ch:		the channel which we are removing
func (r *RulesStore) DeleteListener(ch RulesEventChannel) {
    log.Infof("Removing the listener: %s from the rules sevice", ch)
    r.Lock()
    defer r.Unlock()
    delete(r.listeners, ch)
}

// Retrieve a list of those listening to events
func (r *RulesStore) Listeners() []RulesEventChannel {
    list := make([]RulesEventChannel, 0)
    r.RLock()
    defer r.RUnlock()
    for listener, _ := range r.listeners {
        list = append(list, listener)
    }
    return list;
}

// Handle the events from the backend store
func (r *RulesStore) eventProcessor() {

    for {
        select {
        case <-r.shutdown_channel:
            log.Infof("Killing off the Rules event processor")
            return
        case event := <-r.events:
            log.Infof("Received event from the backend store: %s", event)


        }
    }
}

// Sending a notification upstream and protected us from panics
//	event:		a pointer to a rules events
func (r *RulesStore) pushNotification(event *RulesEvent) error {
    r.RLock()
    defer r.Unlock()
    log.Debugf("Sending the event: %s upstream to listeners", event)

    // step: iterate the list of listeners and send the notifications
    for listener, _  := range r.listeners {
        log.Debugf("Sending upstream via channel: %s", listener)
        go func() {
            listener <- event
        }()
    }
    return nil
}

func (r RulesEvent) String() string {
    return fmt.Sprintf("id: %s, type: %s, ", r.ID, r.NotifyTypeInString())
}

// Set the rule notification to changed
func (r *RulesEvent) Changed() *RulesEvent {
    r.NotifyType = RULES_CHANGED
    return r
}

// Set the rule notification to added
func (r *RulesEvent) Added() *RulesEvent {
    r.NotifyType = RULES_ADDED
    return r
}

// Set the rule notification to deleted
func (r *RulesEvent) Deleted() *RulesEvent {
    r.NotifyType = RULES_DELETED
    return r
}

// Check to see if it's an changed event
func (r *RulesEvent) IsChanged() bool {
    if r.NotifyType == RULES_CHANGED {
        return true
    }
    return false
}

// Check to see if it's an added event
func (r *RulesEvent) IsAdded() bool {
    if r.NotifyType == RULES_ADDED {
        return true
    }
    return false
}

// Check to see if it's an deleted event
func (r *RulesEvent) IsDeleted() bool {
    if r.NotifyType == RULES_DELETED {
        return true
    }
    return false
}

func (r RulesEvent) NotifyTypeInString() string {
    switch r.NotifyType {
        case RULES_CHANGED:
            return "changed"
        case RULES_ADDED:
            return "added"
        case RULES_DELETED:
            return "deleted"
        default:
            return "unknown"
    }
}
