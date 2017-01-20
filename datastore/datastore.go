/*
Copyright 2016 Christian Grabowski All rights reserved.
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

package datastore

// ResultCallback is a callback that is expecting both a resulting data interface and an error
type ResultCallback func([]byte, error)

// NoResultCallback is a callback just expecting an error
type NoResultCallback func(error)

// Datastore is an interface for possible datastores
type Datastore interface {
	Save(key string, data interface{}, callback NoResultCallback)
	Find(queryStr string, callback ResultCallback)
	Remove(queryStr string, callback NoResultCallback)
	Update(queryStr string, update interface{}, callback NoResultCallback)
	FindAndUpdate(queryStr string, update interface{}, callback ResultCallback)
}
