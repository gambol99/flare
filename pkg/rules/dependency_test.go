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
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
    graph := NewDependencies()
    assert.NoError(t, graph.Connect("A", "B"))
    assert.NoError(t, graph.Connect("C", "D"))
    assert.True(t, graph.Exists("A"))
    assert.True(t, graph.Exists("B"))
    assert.True(t, graph.Exists("C"))
    assert.True(t, graph.Exists("D"))
    assert.False(t, graph.Exists("Z"))
}

func TestDisconnect(t *testing.T) {
    graph := NewDependencies()
    assert.NoError(t, graph.Connect("A", "B"))
    assert.NoError(t, graph.Connect("A", "C"))
    assert.NoError(t, graph.Connect("A", "F"))
    edges, found, err := graph.Edges("A")
    assert.NoError(t, err)
    assert.True(t, found)
    if !assert.Equal(t, 3, len(edges)) {
        t.Logf("Edges: %s", edges)
    }
    // lets drop a connection
    assert.NoError(t, graph.Disconnect("A", "F"))
    edges, found, err = graph.Edges("A")
    assert.NoError(t, err)
    assert.True(t, found)
    if !assert.Equal(t, 2, len(edges)) {
        t.Logf("Edges: %s", edges)
    }
}

func TestEdges(t *testing.T) {
    graph := NewDependencies()
    assert.NoError(t, graph.Connect("A", "B"))
    assert.NoError(t, graph.Connect("B", "C"))
    assert.NoError(t, graph.Connect("A", "D"))
    edges, found, err := graph.Edges("A")
    assert.NoError(t, err)
    assert.True(t, found)
    if !assert.Equal(t, 2, len(edges)) {
        t.Logf("Edges: %s", edges)
    }
    assert.NoError(t, graph.Connect("A", "F"))
    assert.NoError(t, graph.Connect("A", "G"))
    edges, found, err = graph.Edges("A")
    assert.NoError(t, err)
    assert.True(t, found)
    if !assert.Equal(t, 4, len(edges)) {
        t.Logf("Edges: %s", edges)
    }
}