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
	"testing"

	"github.com/stretchr/testify/assert"
)

func createFakeGroup() *FlareRuleGroup {
	group := NewFlareGroup()
	group.ID = "security_group_one"
	group.Active = true
	group.Description = "my description"
	group.Rule().Address().Addr("10.0.100.32").TCP().Port(389).Accept().Comment("allow ldap authentication")
	group.Rule().Address().Addr("10.0.100.32").TCP().Port(669).Accept().Comment("allow ldaps authentication")
	return group
}

func TestNewFlareGroup(t *testing.T) {
	rule := NewFlareGroup()
	assert.NotNil(t, rule)
}

func TestIsValid(t *testing.T) {
	group := createFakeGroup()
	assert.NoError(t, group.IsValid())
	group.ID = ""
	assert.Error(t, group.IsValid())
}

func TestGroupCompare(t *testing.T) {
	groupA := createFakeGroup()
	groupB := createFakeGroup()
	assert.True(t, groupA.Compare(groupB))
	groupA.ID = "nothing"
	assert.False(t, groupA.Compare(groupB))
	groupA.ID = "security_group_one"
	groupB.Rule().Address().Addr("10.0.100.32").TCP().Port(389).Accept().Comment("allow ldap authentication")
	assert.False(t, groupA.Compare(groupB))
}

func TestReferences(t *testing.T) {
	group := NewFlareGroup()
	group.ID = "security_group_references"
	group.Active = true
	group.Description = "my description"
	group.Rule().Group().SetID("security_group_one")
	list := group.References()
	assert.NotEmpty(t, list)
	assert.Equal(t, 1, len(list))
}
