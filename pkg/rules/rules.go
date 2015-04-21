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
	"sync"

	"github.com/gambol99/flare/pkg/api"
	"github.com/gambol99/flare/pkg/config"
	"github.com/gambol99/flare/pkg/store"
	"github.com/gambol99/flare/pkg/utils"

	log "github.com/Sirupsen/logrus"
	"path/filepath"
)

// the implementation of the rule store
type RulesStore struct {
	// locking for the struct
	sync.RWMutex
	// the store we're pulling the rules from
	backend store.Store
	// the channel we receive events from the store from
	events store.StoreEventChannel
	// the configuration
	cfg *config.FlareConfiguration
	// a map of the groups currently managed
	groups map[string]*api.FlareRuleGroup
	// the graph data used to manage the group dependencies
	dependencies Dependencies
	// the shutdown channel
	shutdown_channel utils.ShutdownChannel
	// a collection of those listening to events
	listeners map[RulesEventChannel]bool
}

// Create a new rules store
//  cfg:    the configuration for the rules store
func New(cfg *config.FlareConfiguration) (RulesService, error) {
	log.Infof("Initializing the Rules Store, store url: %s", cfg.StoreURL)

	var err error
	service := new(RulesStore)
	service.cfg = cfg
	service.dependencies = NewDependencies()
	service.listeners = make(map[RulesEventChannel]bool, 0)
	service.groups = make(map[string]*api.FlareRuleGroup, 0)
	service.events = make(store.StoreEventChannel, 10)
	service.shutdown_channel = make(utils.ShutdownChannel)

	// step: we create the backend
	if service.backend, err = store.New(cfg.StoreURL, cfg.StoreNamespace, cfg.StoreOptions); err != nil {
		log.Errorf("failed to create the backend store: %s, error: %s", cfg.StoreOptions, err)
		return nil, err
	}

	// step: ensure the namespace exists
	if err := service.backend.Mkdir(service.getNamespace(NAMESPACE_GROUPS)); err != nil {
		return nil, fmt.Errorf("failed to create the namespace for groups, error: %s", err)
	}
	// ensure the namespace for the nodes
	if err := service.backend.Mkdir(service.getNamespace(NAMESPACE_FLARE)); err != nil {
		return nil, fmt.Errorf("failed to create the namespace for nodes, error: %s", err)
	}

	// step: add our self as listener for events from the store
	service.backend.AddListener(service.events)

	// step: handle events from the backend store
	go service.eventProcessor()

	return service, nil
}

// Add or update a group
//	group:		the group you wish to update / add
func (r *RulesStore) Add(group *api.FlareRuleGroup) error {
	log.Infof("Adding group: %s", group.ID)
	// step: we first validate the group and all the rules
	if err := group.IsValid(); err != nil {
		log.Errorf("the group: %s is invalid, error: %s", group.ID, err)
		return err
	}

	// step: we need to check if the group references any others
	group_list, err := r.ListGroups()
	if err != nil {
		log.Errorf("Failed to get a list of groups, error: %s", err)
		return err
	}
	references := group.References()
	for _, reference := range references {
		if !utils.IsValidArgument(reference, group_list) {
			return fmt.Errorf("a rule in the group references group: %s which does not exists", reference)
		}
	}

	// step: we marshall the content
	content, err := r.encode(group)
	if err != nil {
		return fmt.Errorf("Failed to encode the data, error: %s", err)
	}

	// step: add group to the backend
	namespace := r.getNamespace(NAMESPACE_GROUPS, group.ID)

	// step: insert the group
	if err := r.backend.Set(namespace, content); err != nil {
		log.Errorf("unable to insert the group: %s into the backend, error: %s", group.ID, err)
		return err
	}

	// step: add to the cache
	r.Lock()
	defer r.Unlock()
	r.groups[group.ID] = group
	r.dependencies.Add(group.ID)

	return nil
}

// Delete a group from the backend
//	name:		the name of the security group you are deleting
func (r *RulesStore) Delete(name string) error {
	log.Infof("Attemping to remove the group: %s from the backend", name)
	// step: we need to make this group is not referenced by other groups
	edges, found, err := r.dependencies.Edges(name)
	if err != nil {
		log.Errorf("Failed to find the edges on node: %s", name)
		return err
	}
	// step: do we have any relationships?
	if found {
		log.Debugf("We have %d relationships to the node: %s", len(edges), edges)
		log.Warnf("The group: %s is referenced by %d other security groups", name, len(edges))

		// step: iterate the edges and remove the group from there
		for _, group_id := range edges {
			log.Warnf("Removing the reference to the group: %s from group: %s", name, group_id)
			group, err := r.Get(group_id)
			if err != nil {
				log.Errorf("Failed to retrieve the group: %s in order to remove reference to: %s", group_id, name)
				return err
			}
			if err := group.RemoveGroup(name); err != nil {
				return err
			}
		}
	} else {
		log.Debugf("No relationships to the node: %s", name)
	}
	return nil
}

// Flu shes all the rules - effectively deletes everything
func (r *RulesStore) Flush() error {
	log.Warnf("Deleting all the rules from backend and cache")
	r.Lock()
	defer r.Unlock()


	return nil
}

// Retrieve a group from the backend, decode and return the structure
// 	name:		the name of the group we are interested in
func (r *RulesStore) Get(name string) (*api.FlareRuleGroup, error) {
	log.Debugf("Retreiving the group: %s from the backend: %s", name)

	// step: pull it in from the backend
	content, err := r.backend.Get(r.getNamespace(NAMESPACE_GROUPS, name))
	if err != nil {
		log.Errorf("Failed to get the content from the backend, error: %s", err)
		return nil, err
	}

	// step: decode the content
	group := new(api.FlareRuleGroup)
	if err := r.decode(content, group); err != nil {
		log.Errorf("Failed to decode the content: %s, error: %s", content, err)
		return nil, err
	}

	return group, nil
}

// Produce a list of the security group presently in the backend
func (r *RulesStore) ListGroups() ([]string, error) {
	log.Debugf("Retrieving all the groups from the backend")

	// step: generate the namspace
	namespace := r.getNamespace(NAMESPACE_GROUPS)

	// step: retrieve a list of groups under the namespace
	groups, err := r.backend.List(namespace, true)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// Check to see if the group specified exists
//	ID:		the name of the group you are looking for
func (r *RulesStore) IsGroup(ID string) (bool, error) {
	log.Debugf("Looking for group: %s in the list of groups", ID)
	// step: make sure we have a group id
	if ID == "" {
		return false, fmt.Errorf("you have not supplied a group id to check")
	}
	// step: we get a list of groups
	groups, err := r.ListGroups()
	if err != nil {
		return false, err
	}
	// step: iterate the group names a look for a match
	for _, name := range groups {
		if filepath.Base(name) == ID {
			return true, nil
		}
	}
	return false, nil
}

// Syncs the rules with the rules cache
func (r *RulesStore) Sync() error {
	log.Debugf("Loading the rule group from backend")

	// step: we get a list of groups
	groups, err := r.backend.List(r.getNamespace(NAMESPACE_GROUPS), true)
	if err != nil {
		log.Errorf("Failed to retrieve a list of the groups from the backend, error: %s", err)
		return err
	}
	log.Infof("Found %d groups in the backend", len(groups))

	// step: create a temporary cache - we can perform an atomic swap later
	cache := make(map[string]*api.FlareRuleGroup, 0)
	depends := NewDependencies()

	// step: iterate the groups and build dependency tree
	for _, name := range groups {
		log.Debugf("Loading the group: %s from backend", name)

		// step: grab the group from the backend
		group, err := r.Get(name)
		if err != nil {
			return fmt.Errorf("failed to load the group: %s, %s", name, err)
		}

		// step: we need to validate the group before we accept it
		if err := group.IsValid(); err != nil {
			return fmt.Errorf("The group: %s is invalid, violation: %s", err)
		}

		// step: we add the node to the graph database
		depends.Add(name)

		// step: does it reference another group which does not exists
		for _, rule := range group.Rules {
			// check: if the rule a group type??
			if rule.RuleType == api.TYPE_GROUP {
				log.Debugf("Checking if the group reference: %s in group: %s exists", name, rule.ID)
				if !utils.IsValidArgument(rule.ID, groups) {
					log.Errorf("The reference to group: %s in group: %s does not exist", rule.ID, name)
					return err
				}
				// add the group dependency
				depends.Connect(name, rule.ID)
			}
		}

		// step: add the group to the cache
		cache[name] = group
	}

	// step: we need to lock and swap out the cache and dependencies
	r.Lock()
	defer r.Unlock()
	r.groups = cache
	r.dependencies = depends
	return nil
}

// Close up and release any resource from the service - also pass the shutdown signal down
// the chain to any dependencies
func (r *RulesStore) Close() error {
	log.Infof("Shutting down the Rules Store")
	// step: shutdown the event processor
	r.shutdown_channel <- true
	// step: shutdown the store
	r.backend.Close()
	return nil
}

func (r *RulesStore) bootstrapRulesStore(filename string) error {
	log.Infof("Bootstrapping the rules store, from bootstrap file: %s", filename)

	return nil
}
