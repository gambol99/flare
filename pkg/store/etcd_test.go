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
	"sync"
	"testing"
	"time"

	"github.com/gambol99/flare/pkg/config"
	"github.com/gambol99/flare/pkg/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	lock      sync.Once
	etcdAgent Store
)

func createTestEctdAgent(t *testing.T) Store {
	lock.Do(func() {
		log.SetLevel(log.PanicLevel)
		cfg := config.DefaultConfig()
		cfg.StoreURL = "etcd://127.0.0.1:4001"
		cfg.StoreNamespace = "/"
		agent, err := New(cfg.StoreURL, cfg.StoreNamespace, cfg.StoreOptions)
		if err != nil {
			t.Logf("Unable to create a agent for etcd, error: %s", err)
		}
		// perform a simply test - not sure how else to do this??
		if _, err := agent.Get("/"); err != nil {
			t.Logf("skipping the tests, the etcd host is probably offline, error: %s", err)
		}
		etcdAgent = agent
	})
	if etcdAgent == nil {
		t.SkipNow()
	}
	return etcdAgent
}

func TestSet(t *testing.T) {
	agent := createTestEctdAgent(t)
	err := agent.Set("/my/test/key", "hello")
	assert.Nil(t, err)
	value, err := agent.Get("/my/test/key")
	assert.Nil(t, err)
	assert.Equal(t, "hello", value)
}

func TestGet(t *testing.T) {
	agent := createTestEctdAgent(t)
	value, err := agent.Get("/my/test/key")
	assert.Nil(t, err)
	assert.Equal(t, "hello", value)
}

func TestDelete(t *testing.T) {
	agent := createTestEctdAgent(t)
	assert.Nil(t, agent.Set("/my/test/key", "hello"))
	assert.Nil(t, agent.Delete("/my/test/key", false))
}

func TestDeleteRecursive(t *testing.T) {
	agent := createTestEctdAgent(t)
	assert.Nil(t, agent.Set("/my/test/key", "hello"))
	assert.Nil(t, agent.Set("/my/test/key1", "hello"))
	assert.Nil(t, agent.Set("/my/test/key2", "hello"))
	assert.Nil(t, agent.Set("/my/test/key4", "hello"))
	assert.Nil(t, agent.Delete("/my/test", true))
}

func TestList(t *testing.T) {
	agent := createTestEctdAgent(t)
	assert.Nil(t, agent.Set("/flare/test/groups/1", "1"))
	assert.Nil(t, agent.Set("/flare/test/groups/2", "1"))
	assert.Nil(t, agent.Set("/flare/test/groups/3", "1"))
	list, err := agent.List("/flare/test/groups", false)
	assert.Nil(t, err)
	if !assert.NotNil(t, list) {
		t.SkipNow()
	}
	assert.Equal(t, 3, len(list))
	assert.Equal(t, "/flare/test/groups/1", list[0])
}

func TestListRecursive(t *testing.T) {
	agent := createTestEctdAgent(t)
	assert.Nil(t, agent.Delete("/flare/test", true))
	assert.Nil(t, agent.Set("/flare/test/groups/1", "1"))
	assert.Nil(t, agent.Set("/flare/test/groups/2", "1"))
	assert.Nil(t, agent.Set("/flare/test/groups/3", "1"))
	assert.Nil(t, agent.Set("/flare/test/groups/folder/1", "1"))
	assert.Nil(t, agent.Set("/flare/test/groups/folder/2", "1"))
	assert.Nil(t, agent.Set("/flare/test/groups/folder/3", "1"))
	list, err := agent.List("/flare/test/groups", true)
	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 6, len(list))
	assert.Nil(t, agent.Set("/flare/test/groups/folder/4", "1"))
	list, err = agent.List("/flare/test/groups", true)
	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.NotEqual(t, 6, len(list))
	assert.Equal(t, 7, len(list))
}

func TestAddListener(t *testing.T) {
	agent := createTestEctdAgent(t)
	ch := make(StoreEventChannel, 4)
	agent.AddListener(ch)
	agent.Watch("/")
	value := utils.RandomUUID(32)
	agent.Set("/test", value)
	select {
	case <-time.After(time.Duration(500) * time.Millisecond):
		t.Logf("Timed out waiting to the event to occur")
		t.FailNow()
	case event := <-ch:
		assert.NotNil(t, event)
		assert.Equal(t, event.Value, value)
	}
}
