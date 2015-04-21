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

import (
	"fmt"
)

func (event StoreEvent) String() string {
	return fmt.Sprintf("id: %s, value: %s, action: %s",
		event.ID, event.Value, event.ActionType())
}

// Get the string representation of action
func (event StoreEvent) ActionType() string {
	switch event.Action {
	case CHANGED:
		return "changed"
	case DELETED:
		return "deleted"
	case ADDED:
		return "added"
	default:
		return "unknown"
	}
}

// Check if it's a directory
func (event StoreEvent) IsDirectory() bool {
	return event.Directory
}

func (event StoreEvent) IsFile() bool {
	if event.Directory {
		return false
	}
	return true
}

// Check if the event is a changed event
func (event StoreEvent) IsChanged() bool {
	if event.Action == CHANGED {
		return true
	}
	return false
}

// Check if the event is a deleted event
func (event StoreEvent) IsDeleted() bool {
	if event.Action == DELETED {
		return true
	}
	return false
}

// Check if the event is a added event
func (event StoreEvent) IsAdded() bool {
	if event.Action == ADDED {
		return true
	}
	return false
}