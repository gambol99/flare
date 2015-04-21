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
    "github.com/gyuho/goraph/graph"

    log "github.com/Sirupsen/logrus"
    "fmt"
)

// The implementation of the dependency graph
type DependencyGraph struct {
    // the graphing data
    data *graph.Data
}

// Create a new dependency graph
func NewDependencies() Dependencies {
    return &DependencyGraph{
        data: graph.New(),
    }
}

// Create a node in the graph
//  id:         the id or reference to the node
func (r *DependencyGraph) Add(id string) {
    r.data.AddNode(graph.NewNode(id))
}

// Create a connection between two nodes in the tree
//  source:     the source node, where were connecting from
//  dest:       the destination of the destination
func (r *DependencyGraph) Connect(src, dest string) error {
    log.Debugf("Creating a dependency between source: %s and destination: %s", src, dest)
    r.data.Connect(
        graph.NewNode(src),
        graph.NewNode(dest),
        1,
    )
    return nil
}

// Break the connection / relationship between two nodes in the tree
//  src:        the source node, where were connecting from
//  dest:       the destination of the destination
func (r *DependencyGraph) Disconnect(src, dest string) error {
    log.Debugf("Disconnect the relationship between, src: %s, dest: %s", src, dest)
    // step: check the src and dest exists
    if !r.Exists(src) {
        return fmt.Errorf("the source node: %s does not exist in the graph data", src)
    }
    if !r.Exists(dest) {
        return fmt.Errorf("the dest node: %s does not exist in the graph data", dest)
    }
    r.data.DeleteEdge(r.data.GetNodeByID(src), r.data.GetNodeByID(dest))
    return nil
}

// Delete the
//  src:        delete the src node and any relationships which are connected
func (r *DependencyGraph) Delete(src string) error {
    log.Debugf("Deleting the node: %s and relationship it may have", src)

    return nil
}

// Checks to see if a node exists
//  name:       the name of the node we are looking for
func (r *DependencyGraph) Exists(name string) bool {
    log.Debugf("Looking for node: %s in graph data", name )
    if node := r.data.GetNodeByID(name); node != nil {
        return true
    }
    return false
}

// Discovery any dependencies on the source node and produce a list of these
// dependents if there are any
//  src:        the source node you are looking for connections to
func (r *DependencyGraph) Edges(src string) ([]string, bool, error) {
    log.Debugf("Looking for relationships to: %s", src)
    // step: check the node exists
    if !r.Exists(src) {
        return nil, false, fmt.Errorf("the node: %s does not exist", src)
    }
    list := make([]string, 0)
    for node, _ := range r.data.GetNodeByID(src).WeightTo {
        list = append(list, node.ID)
    }
    found := false
    if len(list) > 0 {
        found = true
    }

    return list, found, nil
}

// Returns the size of the graph database
func (r *DependencyGraph) Size() int {
    return r.data.GetNodeSize()
}
