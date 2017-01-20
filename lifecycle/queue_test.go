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

package lifecycle

import (
	"os"
	"testing"

	"github.com/cpg1111/maestrod/datastore"
	"github.com/cpg1111/maestrod/datastore/redis"
)

var store = redis.New(os.Getenv("REDIS_SERVICE_HOST"), os.Getenv("REDIS_SERVICE_PORT"), "")

var castedStore = (datastore.Datastore)(*store)

var queue = NewQueue(&castedStore)

func TestAdd(t *testing.T) {
	queue.Add("test", "test", "asdfasdfasdf", "asdfew")
	if !(queue.Queue[0].Project == "test" && queue.Queue[0].Branch == "test" && queue.Queue[0].PrevCommit == "asdfasdfasdf" && queue.Queue[0].CurrCommit == "asdfew") {
		t.Error("Add did not add all the fields")
	}
}

func TestPop(t *testing.T) {
	queue.Add("test", "test", "asdfasdfasdf", "asdfew")
	running := &Running{}
	next := queue.Pop(running, 1)
	if next == nil {
		t.Error("queue should've popped an entry")
	}
	another := queue.Pop(running, 1)
	if another != nil {
		t.Error("queue's second pop should've returned nil")
	}
}
