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

package api

import (
	"github.com/gambol99/flare/pkg/utils"
)

//
// The representation of the rule in the store
//
type FlareRule struct {
	utils.Validation `json:"validation,omitempty"`
	// the identifier of the group
	ID string `json:"id,omitempty", yaml:"id,omitempty"`
	// the type of rule i.e. a service definition or an ip rule
	RuleType string `json:"type", yaml:"type"`
	// a des-cription of the rule
	Description string `json:"description,omitempty", yaml:"description,omitempty"`
	// the destination - this is either a service definition or an ip address
	Destination string `json:"destination,omitempty", yaml:"destination,omitempty"` // consul://mysql@dc1 || 10.100.0.0/24 || 10.110.0.10
	// the protocol
	Protocol string `json:"protocol,omitempty", yaml:"protocol,omitempty"`
	// inverse the rule
	Inverse bool `json:"inverse,omitempty", yaml:"inverse,omitempty"`
	// a list of ports related to the rule above, note; only relevant when using ip addresses not services
	Ports []int `json:"ports,omitempty", yaml:"ports,omitempty"`
	// the action of the rule - i.e. drop or accept
	Action int `json:"action,omitempty", yaml:"action,omitempty"`
}

//
// The definition / structure for a collection of rules
//
type FlareRuleGroup struct {
	utils.Validation `json:"validation,omitempty"`
	// the id or identifier of the group - essentially a namespace
	ID string `json:"id", yaml:"id"`
	// a description for the group
	Description string `json:"desc", yaml:"description"`
	// is the group active
	Active bool `json:"active", yaml:"active"`
	// a collection of rules composing the group
	Rules []*FlareRule `json:"rules", yaml:"rules"`
}
