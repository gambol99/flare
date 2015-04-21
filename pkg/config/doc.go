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

package config

import "sync"

//
// The configuration
//
type Configuration interface {
}

type ConfigutationOptions interface {
	// add an option to list
	Add(string, interface{})
	// set a option
	Set(string) error
	// get a option back from the map
	Get(string, interface{}) interface{}
	// the string
	String() string
}

// Configuration options for various components
type configOptions struct {
	sync.RWMutex
	// the options map
	arguments map[string]interface{}
}
