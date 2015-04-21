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
)

func (group FlareRuleGroup) String() string {
	return fmt.Sprintf("id: %s, active: %t,rules: %d", group.ID, group.Active, len(group.Rules))
}

func (group FlareRuleGroup) IsValid() error {
	// step: validate the group
	if group.ID == "" {
		return fmt.Errorf("the group does not have an id")
	}

	if group.Rules == nil {
		return fmt.Errorf("the group: %s does not have any rules", group.ID)
	}

	// step: iterate and validate each of the rules
	for _, rule := range group.Rules {
		if err := rule.IsValid(); err != nil {
			return fmt.Errorf("invalid group rule in: %s, error: %s", group.ID, err)
		}
	}
	return nil
}

// Create an empty group
func NewFlareGroup() *FlareRuleGroup {
	return &FlareRuleGroup{
		Active: true,
		Rules:  make([]*FlareRule, 0),
	}
}

// Provide us with a list of references to other groups if any
func (group *FlareRuleGroup) References() []string {
	list := make([]string, 0)
	for _, rule := range group.Rules {
		if rule.RuleType == TYPE_GROUP {
			list = append(list, rule.ID)
		}
	}
	return list
}

// Compare one group to another and return true if the same
func (group *FlareRuleGroup) Compare(src *FlareRuleGroup) bool {
	// check the id
	if group.ID != src.ID {
		return false
	}
	// check the active status
	if group.Active != src.Active {
		return false
	}
	// check the description
	if group.Description != src.Description {
		return false
	}
	// check the number of rules
	if len(group.Rules) != len(src.Rules) {
		return false
	}
	// check each of the rules

	return true
}

// Set the Id of the group at return our self
func (group *FlareRuleGroup) SetID(id string) *FlareRuleGroup {
	group.ID = id
	return group
}

// Add a rule to the group and return it
func (group *FlareRuleGroup) Rule() *FlareRule {
	rule := NewRule()
	group.Rules = append(group.Rules, rule)
	return rule
}

// Remove a reference to the group
//	name:		the name of the group you wish to remove
func (group *FlareRuleGroup) RemoveGroup(name string) error {
	for index, rule := range group.Rules {
		if rule.RuleType == TYPE_GROUP && rule.ID == name {
			group.Rules = group.Rules[:index+copy(group.Rules[index:], group.Rules[index+1:])]
			return nil
		}
	}
	return fmt.Errorf("failed to find the reference to group: %s in group", name)
}
