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

package discovery

import (
	"fmt"
	"net/url"

	"github.com/gambol99/flare/pkg/config"

	log "github.com/Sirupsen/logrus"
)

func New(cfg *config.FlareConfiguration) (Discovery, error) {
	log.Debugf("Creating a new discovery agent from url: %s", cfg.DiscoveryURL)
	// parse the url location
	location, err := url.Parse(cfg.DiscoveryURL)
	if err != nil {
		log.Errorf("unable to parse the discovery url, error: %s", err)
		return nil, err
	}
	switch location.Scheme {
	case "consul":
		return NewConsulAgent(location)
	case "file":
		return NewFileAgent(location)
	default:
		return nil, fmt.Errorf("unsupport discovery provider: %s", location.Scheme)
	}
	return nil, fmt.Errorf("unsupport discovery provider")
}
