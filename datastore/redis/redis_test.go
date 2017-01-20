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

package redis

import (
	"encoding/json"
	"os"
	"testing"
)

var store = New(os.Getenv("REDIS_SERVICE_HOST"), os.Getenv("REDIS_SERVICE_PORT"), "")

type redisResult struct {
	Data interface{}
}

func checkContent(t *testing.T, res string, doneChan chan bool) {
	t.Log(res)
	decodedRes := &redisResult{}
	decodeErr := json.Unmarshal(([]byte)(res), decodedRes)
	if decodeErr != nil {
		t.Error(decodeErr)
	}
	if decodedRes.Data != "test" {
		t.Errorf("Expected test fournd %s", decodedRes.Data)
	}
	doneChan <- true
}

func saveContent(t *testing.T, key, val string) {
	cmd := store.store.Set(key, val, 0)
	_, cmdErr := cmd.Result()
	if cmdErr != nil {
		t.Error(cmdErr)
	}
}

func TestRedisSave(t *testing.T) {
	doneChan := make(chan bool)
	store.Save("testSaveData", "test", func(err error) {
		if err != nil {
			t.Error(err)
		}
		cmd := store.store.Get("testSaveData")
		res, resErr := cmd.Result()
		if resErr != nil {
			t.Error(resErr)
		}
		checkContent(t, res, doneChan)
	})
	_ = <-doneChan
}

func TestRedisFind(t *testing.T) {
	doneChan := make(chan bool)
	saveContent(t, "testFindData", "test")
	store.Find("testFindData", func(res []byte, err error) {
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
		strRes := (string)(res)
		if strRes != "test" {
			t.Errorf("expected test, found %s on find", strRes)
		}
		doneChan <- true
	})
	_ = <-doneChan
}

func TestRedisRemove(t *testing.T) {
	doneChan := make(chan bool)
	saveContent(t, "testRemoveData", "test")
	store.Remove("testRemoveData", func(err error) {
		if err != nil {
			t.Error(err)
		}
		cmd := store.store.Get("testRemoveData")
		res, resErr := cmd.Result()
		if resErr != nil {
			t.Error(resErr)
		}
		if res != "" {
			t.Errorf("expected empty string found %s for remove", res)
		}
		doneChan <- true
	})
	_ = <-doneChan
}

func TestRedisUpdate(t *testing.T) {
	doneChan := make(chan bool)
	saveContent(t, "testUpdateData", "test")
	store.Update("testUpdateData", "updated_test", func(err error) {
		if err != nil {
			t.Error(err)
		}
		cmd := store.store.Get("testUpdateData")
		res, resErr := cmd.Result()
		if resErr != nil {
			t.Error(resErr)
		}
		decodedRes := &redisResult{}
		decodeErr := json.Unmarshal(([]byte)(res), decodedRes)
		if decodeErr != nil {
			t.Error(decodeErr)
		}
		if decodedRes.Data != "updated_test" {
			t.Errorf("Expected updated_test found %s for update", decodedRes.Data)
		}
		doneChan <- true
	})
	_ = <-doneChan
}

func TestRedisFindAndUpdate(t *testing.T) {
	doneChan := make(chan bool)
	saveContent(t, "testFindUpdateData", "test")
	store.FindAndUpdate("testFindUpdateData", "updated_test", func(res []byte, err error) {
		if err != nil {
			t.Error(err)
		}
		strRes := (string)(res)
		if strRes != "{\"Data\":\"updated_test\"}" {
			t.Errorf("expected updated_test found %s for findAndUpdate", strRes)
		}
		doneChan <- true
	})
	_ = <-doneChan
}
