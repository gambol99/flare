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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidIpAddress(t *testing.T) {
	valid := IsValidIpAddress("10.0.0.1")
	assert.True(t, valid, "the result should have been true for the ip address")
	valid = IsValidIpAddress("10.0.111.256")
	assert.False(t, valid, "the result should have been false for this ip adress")
	valid = IsValidIpAddress("10.0.111.0/24")
	assert.False(t, valid, "the result should have been false for this ip adress")
}

func TestGetInterfaces(t *testing.T) {
	list, err := GetInterfaces()
	assert.Nil(t, err, "we weren't able to get a list of the interfaces")
	assert.NotNil(t, list, "the list of interfaces was nil")
	assert.NotEmpty(t, list, "the list of interfaces was empty")
}

func TestIsNetworkInterface(t *testing.T) {
	list, err := GetInterfaces()
	assert.Nil(t, err, "we weren't able to get a list of the interfaces")
	for _, name := range list {
		valid, err := IsNetworkInterface(name)
		assert.Nil(t, err, "we weren't able to get a list of the interfaces")
		assert.True(t, valid, "the result should have been true")
	}
	valid, err := IsNetworkInterface("doesnotexist")
	assert.Nil(t, err, "we should have recieved an error here, the inteface does not exist")
	assert.False(t, valid, "the result should have been false, interface does not exist")
}

func TestIsEmpty(t *testing.T) {
	assert.True(t, IsEmpty(""))
	assert.False(t, IsEmpty("kdjskdjkjdkjs"))
}

func TestIsPort(t *testing.T) {
	assert.True(t, IsPort(0))
	assert.True(t, IsPort(676))
	assert.True(t, IsPort(32323))
	assert.False(t, IsPort(-1))
	assert.False(t, IsPort(65536))
}

func TestIsValidURL(t *testing.T) {
	assert.True(t, IsValidURL("http://127.0.0.1"), "the url should have been valid")
	assert.True(t, IsValidURL("etcd://127.0.0.1:4001"), "the url should have been valid")
	assert.True(t, IsValidURL("file://var/run/bootstrap.json"), "the url should have been valid")
	//assert.False(t, IsValidURL("111kdlskdlskdsldk"), "the url should have been invalid")
}

func TestIsValidArgument(t *testing.T) {
	assert.True(t, IsValidArgument("info", "info", "debug", "warn"))
	assert.False(t, IsValidArgument("info", "debug", "warn"))
	assert.True(t, IsValidArgument(1, "debug", 1, "warn", "info"))
}

func TestFileExists(t *testing.T) {
	exists, err := FileExists("/etc/passwd")
	assert.Nil(t, err)
	assert.True(t, exists)
	exists, err = FileExists("/etc/pasklasklakslaksd")
	assert.Nil(t, err)
	assert.False(t, exists)
}
