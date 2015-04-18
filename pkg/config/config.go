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

import (
	"fmt"
	"strings"

	"github.com/gambol99/flare/pkg/utils"
)

type FlareConfiguration struct {
	utils.Validation

	// the metrics endpoint url
	MetricsURL string
	// the port metrics should be running
	MetricsPort int

	// the interface the containers are attached
	DockerInterface string
	// the path to the docker socket
	DockerSocket string

	// the default policy for traffic i.e. ACCEPT, DROP
	DefaultPolicy string
	// a default security group - applyed to all containers
	DefaultSecurityGroup string

	// the logging level - INFO, DEBUG etc
	LogLevel string
	// the location / path for a bootstrap file
	Bootstrap string
	// the prefix for flare specs in the environment variables
	FlarePrefix string

	// the location for the store i.e. file://some/path or etcd://localhost:4001,host1:4001 etc
	StoreURL string
	// options for the store agent
	StoreOptions ConfigutationOptions
	// the arguments for the above
	StoreArguments []string
	// the store namespace i.e. prefix
	StoreNamespace string
	// the encoding of the backend
	StoreEncoding string

	// the discovery provider url
	DiscoveryURL string
	// options for the discovery agent
	DiscoveryOptions ConfigutationOptions
}

// Checks to see if the configuration is valid
func (config FlareConfiguration) IsValid() error {
	if !utils.IsValidURL(config.StoreURL) {
		return fmt.Errorf("the store url: %s is invalid, please check", config.StoreURL)
	}
	// check the discovery url
	if config.DiscoveryURL != "" {
		if !utils.IsValidURL(config.DiscoveryURL) {
			return fmt.Errorf("invalid discovery url: %s", config.DiscoveryURL)
		}
	}

	// check we have a flare prefix
	if utils.IsEmpty(config.FlarePrefix) {
		return fmt.Errorf("a flare prefix is required before we can continue")
	}

	// check the encoding
	config.StoreEncoding = strings.ToLower(config.StoreEncoding)
	if config.StoreEncoding == "" || !utils.IsValidArgument(config.StoreEncoding, "json", "yaml") {
		return fmt.Errorf("the encode: %s is invalid", config.StoreEncoding)
	}

	// check the interface exists
	if valid, err := utils.IsNetworkInterface(config.DockerInterface); err != nil {
		return fmt.Errorf("unable to get interface details, error: %s", err)
	} else if !valid {
		return fmt.Errorf("the interface: %s does not exist", config.DockerInterface)
	}

	// check the bootstrap exists
	if config.Bootstrap != "" {
		if exists, err := utils.FileExists(config.Bootstrap); err != nil {
			return fmt.Errorf("unable to check for boostrap file: %s, error: %s", config.Bootstrap, err)
		} else if !exists {
			return fmt.Errorf("the bootstrap file: %s does not exist", config.Bootstrap)
		}
	}

	// ensure the case
	config.LogLevel = strings.ToUpper(config.LogLevel)
	config.DefaultPolicy = strings.ToUpper(config.DefaultPolicy)
	if !utils.IsValidArgument(config.LogLevel, "INFO", "DEBUG", "NONE") {
		return fmt.Errorf("invalid log level: %s, supported are info, debug, none", config.LogLevel)
	}

	if !utils.IsValidArgument(config.DefaultPolicy, "DROP", "ACCEPT") {
		return fmt.Errorf("invalid defailt policy: %s, supported are drop, accept", config.DefaultPolicy)
	}
	return nil
}

func (c FlareConfiguration) String() string {
	return fmt.Sprintf(`
 Interface: '%s'
 Bootstrap: '%s'
 Docker Socket: '%s'
 Default Policy: '%s'
 Default Security Group: '%s'
 Default Flare Prefix: '%s'
 Log Level: '%s'
 StoreNamespace: %s
 Store URL: '%s'
 Store Options: '%s'
 Store Encoding: %s
 Discovery URL: '%s'
 Discovery Options: '%s'
`, c.DockerInterface, c.Bootstrap, c.DockerSocket, c.DefaultPolicy, c.DefaultSecurityGroup, c.LogLevel,
		c.StoreURL, c.StoreNamespace, c.StoreOptions, c.StoreEncoding, c.DiscoveryURL,
		c.DiscoveryOptions)
}

func NewConfigurationOptions() ConfigutationOptions {
	return &configOptions{
		arguments: make(map[string]interface{}),
	}
}

// Get the option value or return the default value provided
//  id:     the id or key of the option
//  df:     the default value to return
func (r *configOptions) Get(id string, df interface{}) interface{} {
	r.RLock()
	defer r.RUnlock()
	if value, found := r.arguments[id]; found {
		return value
	}
	return df
}

// Add a key value pair into the options - essentially used by the command line parser
//  value:      the value to be added to the set
func (r *configOptions) Set(value string) error {
	items := strings.Split(value, "=")
	if len(items) != 2 {
		return fmt.Errorf("invalid option: %s, must be in KEY=VALUE format", value)
	}
	r.arguments[items[0]] = items[1]
	return nil
}

// Add a key/pair to the options
//  id:     the id or key of the option
//  value:  the value of the key option
func (r *configOptions) Add(id string, value interface{}) {
	r.Lock()
	defer r.Unlock()
	r.arguments[id] = value
}

func (r *configOptions) String() string {
	return fmt.Sprintf("%s", r.arguments)
}

// Generate a default configuration
func DefaultConfig() *FlareConfiguration {
	return &FlareConfiguration{
		DockerInterface:      "docker0",
		DockerSocket:         "/var/run/docker.sock",
		DefaultPolicy:        "DROP",
		DefaultSecurityGroup: "",
		LogLevel:             "INFO",
		StoreNamespace:       "flare",
		StoreURL:             "etcd://127.0.0.1:4001",
		StoreOptions:         NewConfigurationOptions(),
		StoreEncoding:        "json",
		DiscoveryURL:         "",
		DiscoveryOptions:     NewConfigurationOptions(),
		Bootstrap:            "",
		FlarePrefix:          "FLARE_",
	}
}
