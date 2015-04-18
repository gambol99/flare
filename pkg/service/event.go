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

package service

func (r *ContainerEvent) ContainerID(id string) *ContainerEvent {
    r.ID = id
    return r
}

func (r *ContainerEvent) SetStatus(status string) *ContainerEvent {
    r.Status = status
    return r
}

func (r ContainerEvent) HasStarted() bool {
    if r.Status == DOCKER_START {
        return true
    }
    return false
}

func (r ContainerEvent) HasDestroyed() bool {
    if r.Status == DOCKER_DESTROY {
        return true
    }
    return false
}
