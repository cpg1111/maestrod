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
	"errors"
	"fmt"
	"log"

	"github.com/cpg1111/maestrod/datastore"

	redis "gopkg.in/redis.v3"
)

// RedisStore is the datastore backed by redis
type RedisStore struct {
	datastore.Datastore
	store *redis.Client
}

type redisData struct {
	Data interface{}
}

// NewRedis returns a pointer to a new RedisStore instance
func New(host, port, password string) *RedisStore {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
	}
	return &RedisStore{
		store: redis.NewClient(options),
	}
}

// Save saves data in redis and takes a NoResultCallback
func (r RedisStore) Save(key string, data interface{}, callback datastore.NoResultCallback) {
	go func() {
		rData := &redisData{
			Data: data,
		}
		dataStr, marshalErr := json.Marshal(rData)
		if marshalErr != nil {
			callback(marshalErr)
			return
		}
		cmd := r.store.Set(key, dataStr, 0)
		res, resErr := cmd.Result()
		if resErr != nil {
			callback(resErr)
			return
		}
		log.Println("redis message: ", res)
		callback(nil)
	}()
}

// Find finds data in redis and takes a ResultCallback
func (r RedisStore) Find(queryStr string, callback datastore.ResultCallback) {
	go func() {
		cmd := r.store.Get(queryStr)
		res, resErr := cmd.Result()
		callback(([]byte)(res), resErr)
	}()
}

// Remove removes data from redis and takes a NoResultCallback
func (r RedisStore) Remove(queryStr string, callback datastore.NoResultCallback) {
	go func() {
		cmd := r.store.Set(queryStr, nil, 0)
		res, resErr := cmd.Result()
		if resErr != nil {
			callback(resErr)
			return
		}
		log.Println(res)
		callback(nil)
	}()
}

// Update updates a key in redis with a new value
func (r RedisStore) Update(queryStr string, update interface{}, callback datastore.NoResultCallback) {
	resChan := make(chan string)
	errChan := make(chan error)
	go func() {
		gCMD := r.store.Get(queryStr)
		gRes, gResErr := gCMD.Result()
		if gResErr != nil {
			errChan <- gResErr
			return
		}
		resChan <- gRes
	}()
	go func() {
		for {
			select {
			case errMsg := <-errChan:
				if errMsg != nil {
					log.Fatal(errMsg)
					return
				}
			case resMsg := <-resChan:
				if resMsg == "" {
					callback(errors.New("no object found with that key"))
					return
				}
				rData := &redisData{
					Data: update,
				}
				dataStr, marshalErr := json.Marshal(rData)
				if marshalErr != nil {
					callback(marshalErr)
					return
				}
				sCMD := r.store.Set(queryStr, dataStr, 0)
				res, resErr := sCMD.Result()
				if resErr != nil {
					callback(resErr)
					return
				}
				log.Println(res)
				callback(nil)
			}
		}
	}()
}

// Find and update does the same as Update but passes the data to the callback as well
func (r RedisStore) FindAndUpdate(queryStr string, update interface{}, callback datastore.ResultCallback) {
	r.Update(queryStr, update, func(err error) {
		if err != nil {
			callback(nil, err)
			return
		}
		r.Find(queryStr, callback)
	})
}
