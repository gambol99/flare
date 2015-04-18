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
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Encode the content into either json or yaml, depending on the encode option
//	data:			the structure we should encode
func (rule *RulesStore) encode(data interface{}) (string, error) {
	switch rule.cfg.StoreEncoding {
	case "json":
		var buffer bytes.Buffer
		if err := json.NewEncoder(&buffer).Encode(data); err != nil {
			return "", fmt.Errorf("failed to encode the payload into json, error: %s", err)
		}
		return buffer.String(), nil

	default:
		out, err := yaml.Marshal(data)
		if err != nil {
			return "", fmt.Errorf("failed to encode the payload into yaml, error: %s", err)
		}
		return string(out), nil
	}
}

// Decodes the content (yaml or json) into the data structure
// 	content:		the content to be decoded
//	data:			the data structure we should decode into
func (rule *RulesStore) decode(content string, data interface{}) error {
	switch rule.cfg.StoreEncoding {
	case "json":
		if err := json.NewDecoder(strings.NewReader(content)).Decode(data); err != nil {
			return fmt.Errorf("failed to decode the payload into json, error: %s", err)
		}
		return nil
	default:
		if err := yaml.Unmarshal([]byte(content), data); err != nil {
			return fmt.Errorf("failed to encode the payload into yaml, error: %s", err)
		}
		return nil
	}
}
