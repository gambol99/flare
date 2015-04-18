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

// the interface to a store
type Store interface {
	// delete the keys in the namespace
	Delete(string, bool) error
	// list all the keys under a namespace
	List(string, bool) ([]string, error)
	// retrieve a key from the store
	Get(string) (string, error)
	// set a key in the store
	Set(string, string) error
	// check if a key exists in the store
	Exists(string) (bool, error)
	// flush and destroy all the entries
	Flush() error
	// add a listener for changes to the namespace
	AddListener(StoreEventChannel)
	// add watch on the namespace
	Watch(string) error
	// shutdown the resources
	Close() error
}

// the actions which can occur in the store event]
const (
	CHANGED = 0
	DELETED = 1
	ADDED   = 2
)

// the definition for a change in the store
type StoreEvent struct {
	// the key / id the event occurred
	ID string
	// the value of the key
	Value string
	// the type of event, added, deleted, changed
	Action int
}

// a channel which the events are received upon
type StoreEventChannel chan *StoreEvent
