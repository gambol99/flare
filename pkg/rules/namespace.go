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
    "fmt"
    "bytes"
)

// Generate a namespace (i.e. key) from the section arguments
//  keys:       an array of section, i.e. elements making up the namespace
func (r *RulesStore) getNamespace(keys ...string) string {
    var namespace bytes.Buffer
    namespace.WriteString(fmt.Sprintf("/%s", r.cfg.StoreNamespace))
    for _, key := range keys {
        if key != "" {
            namespace.WriteString(fmt.Sprintf("/%s", key))
        }
    }
    return namespace.String()
}
