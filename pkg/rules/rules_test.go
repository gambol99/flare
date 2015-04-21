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
	"sync"
	"testing"

	"github.com/gambol99/flare/pkg/api"
	"github.com/gambol99/flare/pkg/config"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	rulesLock sync.Once
	rules     RulesService
)

func createTestRuleService(t *testing.T) RulesService {
	rulesLock.Do(func() {
		var err error
		log.SetLevel(log.PanicLevel)
		cfg := config.DefaultConfig()
		cfg.StoreURL = "etcd://127.0.0.1:4001"
		cfg.Bootstrap = "file://./tests/rules.yaml"
		cfg.DiscoveryURL = "consul://127.0.0.1:8500"
		rules, err = New(cfg)
		if err != nil {
			t.Logf("failed the create a rules service, skipping the tests")
		}
	})
	if rules == nil {
		t.SkipNow()
	}
	return rules
}

func createCleanRulesService(t *testing.T) RulesService {
	agent := createTestRuleService(t)
	flushTestBackend(t)
	return agent
}

func flushTestBackend(t *testing.T) {
	agent := createTestRuleService(t).(*RulesStore)
	agent.backend.DeleteAll("/flare/groups")
}

func checkAddTestGroup(t *testing.T, agent RulesService, group *api.FlareRuleGroup) {
	if !assert.NoError(t, agent.Add(group)) {
		t.FailNow()
	}
	if !assert.NotNil(t, group) {
		t.FailNow()
	}
}

func checkIsGroup(t *testing.T, agent RulesService, id string) bool {
	found, err := agent.IsGroup(id)
	if !assert.NoError(t, err) {
		t.Logf("Failed to isGroup(%s), error: %s", id, err)
		t.FailNow()
	}
	return found
}

func TestNewRuleStore(t *testing.T) {
	rules := createTestRuleService(t)
	assert.NotNil(t, rules)
}

func TestAdd(t *testing.T) {
	group := api.NewFlareGroup()
	group.ID = "security_group_authentication"
	group.Rule().Address().Addr("10.0.100.32").TCP().Port(389).Accept().Comment("allow ldap authentication")
	group.Rule().Address().Addr("10.0.100.32").TCP().Port(669).Accept().Comment("allow ldaps authentication")
	group.ID = "security_group_authentication_two"
	group.Rule().Address().Addr("10.0.100.32").TCP().Port(389).Accept().Comment("allow ldap authentication")
	group.Rule().Address().Addr("10.0.100.32").TCP().Port(669).Accept().Comment("allow ldaps authentication")

	err := rules.Add(group)
	assert.Nil(t, err)

	err = rules.Add(group)
	assert.Nil(t, err)
}

func TestIsGroup(t *testing.T) {
	agent := createTestRuleService(t)
	group := api.NewFlareGroup()
	group.ID = "security_group_authentication"
	group.Rule().Address().Addr("10.0.100.32").TCP().Port(389).Accept().Comment("allow ldap authentication")
	assert.Nil(t, agent.Add(group))
	assert.True(t, checkIsGroup(t, agent, "security_group_authentication"))
	assert.False(t, checkIsGroup(t, agent, "security_group_authentication_fake"))
}

func TestGet(t *testing.T) {
	agent := createCleanRulesService(t)
	groupA := api.NewFlareGroup()
	groupA.SetID("test_group1").Rule().Address().Addr("10.0.100.32").TCP().Port(80).Accept().Comment("service")
	groupB := api.NewFlareGroup()
	groupB.SetID("test_group2").Rule().Address().Addr("10.0.100.32").TCP().Port(443).Accept().Comment("service")
	checkAddTestGroup(t, agent, groupA)
	checkAddTestGroup(t, agent, groupB)

	group, err := agent.Get(groupA.ID)
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test_group1", group.ID)

	group, err = agent.Get(groupB.ID)
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test_group2", group.ID)
}

//func TestDelete(t *testing.T) {
//	agent := createCleanRulesService(t)
//	group := api.NewFlareGroup()
//	group.SetID("test_group1").Rule().Address().Addr("10.0.100.32").TCP().Port(80).Accept().Comment("service")
//	checkAddTestGroup(t, agent, group)
//	if !checkIsGroup(t, agent, group.ID) {
//		t.FailNow()
//	}
//	assert.NoError(t, agent.Delete(group.ID))
//	if checkIsGroup(t, agent, group.ID) {
//		t.FailNow()
//	}
//}

func TestDeleteGroupMembership(t *testing.T) {
	agent := createCleanRulesService(t)

	groupA := api.NewFlareGroup()
	groupA.SetID("test_group1").Rule().Address().Addr("10.0.100.32").TCP().Port(80).Accept().Comment("service")
	groupB := api.NewFlareGroup()
	groupB.SetID("test_group2").Rule().Address().Addr("10.0.100.32").TCP().Port(443).Accept().Comment("service")
	groupC := api.NewFlareGroup()
	groupC.SetID("test_group3").Rule().Group().SetID("test_group_fake")

	checkAddTestGroup(t, agent, groupA)
	checkAddTestGroup(t, agent, groupB)
	checkAddTestGroup(t, agent, groupB)

}

func TestGroups(t *testing.T) {
	agent := createCleanRulesService(t)
	groupA := api.NewFlareGroup()
	groupA.SetID("test_group1").Rule().Address().Addr("10.0.100.32").TCP().Port(80).Accept().Comment("service")
	groupB := api.NewFlareGroup()
	groupB.SetID("test_group2").Rule().Address().Addr("10.0.100.32").TCP().Port(443).Accept().Comment("service")
	checkAddTestGroup(t, agent, groupA)
	checkAddTestGroup(t, agent, groupB)

	groups, err := agent.ListGroups()
	assert.Nil(t, err)
	if !assert.NotNil(t, groups) {
		t.FailNow()
	}
	if !assert.NotEmpty(t, groups) {
		t.FailNow()
	}
	assert.Equal(t, 2, len(groups))
}
