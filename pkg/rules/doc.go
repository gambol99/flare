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

	"github.com/gambol99/flare/pkg/api"
)

const (
	NAMESPACE_GROUPS = "groups"
	NAMESPACE_FLARE  = "nodes"
)

//
//	An internal wrapper for the groups
//
type RulesGroup struct {
	// the name of the group
	ID string
	// the reference to the underlying flare group
	Group *api.FlareRuleGroup
}

// The type of notifications we can send
const (
	RULES_CHANGED = 0
	RULES_ADDED   = 1
	RULES_DELETED = 2
)

//
//	Rule events sent to upstream listeners
//
type RulesEvent struct {
	// the id of the group in question
	ID string
	// the type of event
	NotifyType int
}

// A channel which events are sent down
type RulesEventChannel chan *RulesEvent

//
// Methods related to the group dependencies - essentially manages the graphing
// data structure
//
type Dependencies interface {
	// add a node to the grap
	Add(string)
	// add a connection between two groups
	Connect(string, string) error
	// remove connection between two groups
	Disconnect(string, string) error
	// delete a group from the dependency
	Delete(string) error
	// find out if anyone is referencing this groups and produce a list of them
	Edges(string) ([]string, bool, error)
	// get the number of nodes
	Size() int
	// check if a node exists
	Exists(string) bool
}

//
// The store interface - just to retrieve, place and watch for changes in rules
//
type RulesService interface {
	// retrieve a group from backend
	Get(name string) (*api.FlareRuleGroup, error)
	// Add or update a group in the store
	Add(group *api.FlareRuleGroup) error
	// Delete a group from the backend
	Delete(id string) (error)
	// flush all the rules from the backend can cache
	Flush() error
	// retrieve a list of the groups from the store
	ListGroups() ([]string, error)
	// check if a group exists
	IsGroup(ID string) (bool, error)
	// add our self as a listener to changes in groups
	AddListener(RulesEventChannel) error
	// remove the channel from listening
	DeleteListener(RulesEventChannel)
	// Sync the rules cache with the backend
	Sync() error
	// provider a list of listener
	Listeners() []RulesEventChannel
	// release and shutdown resources
	Close() error
}
