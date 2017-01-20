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

package etcd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cpg1111/maestrod/datastore"

	etcdv3 "github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

// Etcd3 is a struct for the etcd v3 driver
type Etcd3 struct {
	datastore.Datastore
	Client *etcdv3.Client
	Key    etcdv3.KV
}

func getEndpoint(host, port string) string {
	return fmt.Sprintf("http://%s:%s", host, port)
}

// NewV3 returns a pointer to an Etcd3 struct or an error
func NewV3(host, port string) (*Etcd3, error) {
	cfg := etcdv3.Config{
		Endpoints:   []string{getEndpoint(host, port)},
		DialTimeout: time.Second,
	}
	client, clientErr := etcdv3.New(cfg)
	if clientErr != nil {
		return nil, clientErr
	}
	return &Etcd3{
		Client: client,
		Key:    etcdv3.NewKV(client),
	}, nil
}

// Save saves data in etcd
func (e Etcd3) Save(key string, data interface{}, callback datastore.NoResultCallback) {
	go func() {
		value, marshErr := json.Marshal(data)
		if marshErr != nil {
			callback(marshErr)
			return
		}
		_, putErr := e.Key.Put(context.Background(), key, (string)(value))
		callback(putErr)
	}()
}

// Find finds data in etcd
func (e Etcd3) Find(queryStr string, callback datastore.ResultCallback) {
	go func() {
		resp, respErr := e.Key.Get(context.Background(), queryStr)
		if respErr != nil {
			callback(nil, respErr)
			return
		}
		callback(resp.Kvs[0].Value, nil)
	}()
}

// Remove removes data in etcd
func (e Etcd3) Remove(queryStr string, callback datastore.NoResultCallback) {
	go func() {
		_, delErr := e.Key.Delete(context.Background(), queryStr)
		fmt.Println(delErr)
		callback(delErr)
	}()
}

// Update updates data in etcd
func (e Etcd3) Update(queryStr string, update interface{}, callback datastore.NoResultCallback) {
	go func() {
		value, marshErr := json.Marshal(update)
		if marshErr != nil {
			callback(marshErr)
			return
		}
		_, updateErr := e.Key.Put(context.Background(), queryStr, (string)(value))
		callback(updateErr)
	}()
}

// FindAndUpdate updates data and returns it from etcd
func (e Etcd3) FindAndUpdate(queryStr string, update interface{}, callback datastore.ResultCallback) {
	go func() {
		value, marshErr := json.Marshal(update)
		if marshErr != nil {
			callback(nil, marshErr)
			return
		}
		_, updateErr := e.Key.Put(context.Background(), queryStr, (string)(value))
		if updateErr != nil {
			callback(nil, updateErr)
			return
		}
		resp, getErr := e.Key.Get(context.Background(), queryStr)
		if getErr != nil {
			callback(nil, getErr)
		}
		callback(resp.Kvs[0].Value, nil)
	}()
}
