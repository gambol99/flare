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
	"fmt"

	"github.com/gambol99/flare/pkg/utils"
)

const (
	// an action to accept the matching rule
	ACTION_ACCEPT = 0
	// an action to drop the matching rule
	ACTION_DROP = 1
	// the rule destination is a service tag, i.e. needs to be pull from service discovery
	TYPE_SERVICE = "service"
	// the rule destination is a ip address/subnet
	TYPE_ADDRESS = "address"
	// the group destination
	TYPE_GROUP = "type"
	// tcp protocol
	PROTOCOL_TCP = "tcp"
	// tcp protocol
	PROTOCOL_UDP = "udp"
	// tcp protocol
	PROTOCOL_ALL = "all"
)

// Validate the rule for any errors
func (rule FlareRule) IsValid() error {
	// check the rule type
	if !utils.IsValidArgument(rule.RuleType, TYPE_SERVICE, TYPE_ADDRESS, TYPE_GROUP) {
		return fmt.Errorf("the rule type: %s is invalid", rule.RuleType)
	}
	// step: switch on the rule type
	if rule.RuleType == TYPE_GROUP {
		if utils.IsEmpty(rule.ID) {
			return fmt.Errorf("The rule is a group reference, but no group ID specified")
		}
		return nil

	} else {
		// make sure we have a description
		if utils.IsEmpty(rule.Description) {
			return fmt.Errorf("the rule does not have a description")
		}
		// check we have a destination
		if utils.IsEmpty(rule.Destination) || !utils.IsValidIpAddress(rule.Destination) {
			return fmt.Errorf("the destination: %s is invalid for an address type", rule.Destination)
		}
		// check the protocol
		if !utils.IsValidArgument(rule.Protocol, PROTOCOL_ALL, PROTOCOL_TCP, PROTOCOL_UDP) {
			return fmt.Errorf("the protocol: %s is invalid", rule.Protocol)
		}
		// check the ports
		if rule.Ports == nil || len(rule.Ports) <= 0 {
			return fmt.Errorf("the rule does not have any port defined")
		}
		// check each of the ports
		for _, port := range rule.Ports {
			if !utils.IsPort(port) {
				return fmt.Errorf("the port: %d is invalid", port)
			}
		}
		// check the action
		if !utils.IsValidArgument(rule.Action, ACTION_ACCEPT, ACTION_DROP) {
			return fmt.Errorf("the action: %d is invalid", rule.Action)
		}
	}
	return nil
}

// Compare the rule against myself and return true if equal
func (rule FlareRule) Compare(src *FlareRule) bool {
	if rule.RuleType != src.RuleType {
		return false
	}
	// check if it's a group rule, we only need the id the same
	if rule.RuleType == TYPE_GROUP {
		if rule.ID == src.ID {
			return false
		}
		return true
	}
	// check the description
	if rule.Description != src.Description {
		return false
	}
	// check the action
	if rule.Action != src.Action {
		return false
	}
	// check the inverse flag
	if rule.Inverse != src.Inverse {
		return false
	}
	// check the ports count
	if len(rule.Ports) != len(src.Ports) {
		return false
	}
	// check each of the ports
	for i := 0; i < len(rule.Ports); i++ {
		if rule.Ports[i] != src.Ports[i] {
			return false
		}
	}

	// check the destination
	if rule.Destination != src.Destination {
		return false
	}
	// check the protocol
	if rule.Protocol != src.Protocol {
		return false
	}
	return true
}

// Add in port the port range
//  port:       the port you wish to add
func (rule *FlareRule) Port(port int) *FlareRule {
	rule.Ports = append(rule.Ports, port)
	return rule
}

func (rule *FlareRule) Addr(address string) *FlareRule {
	rule.Destination = address
	return rule
}

func (rule *FlareRule) SetID(id string) *FlareRule {
	rule.ID = id
	return rule
}

// Set the description on the rule
func (rule *FlareRule) Comment(description string) *FlareRule {
	rule.Description = description
	return rule
}

// Invert the logic of the rule
func (rule *FlareRule) Invert() *FlareRule {
	rule.Inverse = true
	return rule
}

func (rule *FlareRule) Group() *FlareRule {
	rule.RuleType = TYPE_GROUP
	return rule
}

func (rule *FlareRule) Service() *FlareRule {
	rule.RuleType = TYPE_SERVICE
	return rule
}

func (rule *FlareRule) Address() *FlareRule {
	rule.RuleType = TYPE_ADDRESS
	return rule
}

// Set the action of the rule drop
func (rule *FlareRule) Drop() *FlareRule {
	rule.Action = ACTION_DROP
	return rule
}

func (rule *FlareRule) Accept() *FlareRule {
	rule.Action = ACTION_ACCEPT
	return rule
}

func (rule *FlareRule) TCP() *FlareRule {
	rule.Protocol = "tcp"
	return rule
}

func (rule *FlareRule) UDP() *FlareRule {
	rule.Protocol = "udp"
	return rule
}

func (rule *FlareRule) All() *FlareRule {
	rule.Protocol = "all"
	return rule
}

func (rule *FlareRule) AnyAny() *FlareRule {
	return rule.All().DestinationAny()
}

func (rule *FlareRule) DestinationAny() *FlareRule {
	rule.Destination = "0.0.0.0/0"
	return rule
}

func (rule FlareRule) String() string {
	return fmt.Sprintf("action: %s, destination: %s:[%s], protocol: %s",
		rule.ActionType(), rule.Destination, rule.Ports, rule.Protocol)
}

// Create a empty rule
func NewRule() *FlareRule {
	return &FlareRule{
		Ports:  make([]int, 0),
		Action: ACTION_DROP,
	}
}

func (rule FlareRule) ActionType() string {
	if rule.Action == ACTION_DROP {
		return "DROP"
	}
	return "ACCEPT"
}
