package etcd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cpg1111/maestrod/datastore"

	etcdv2 "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

// Etcd2 is a struct for the Etcd v2 driver
type Etcd2 struct {
	datastore.Datastore
	Client *etcdv2.Client
	Key    etcdv2.KeysAPI
}

func getEndpoint(host, port string) string {
	return fmt.Sprintf("http://%s:%s", host, port)
}

// NewV2 returns a pointer to an Etcd2 driver or an error
func NewV2(host, port string) (*Etcd2, error) {
	cfg := etcdv2.Config{
		Endpoints:               []string{getEndpoint(host, port)},
		Transport:               etcdv2.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	client, cliErr := etcdv2.New(cfg)
	if cliErr != nil {
		return nil, cliErr
	}
	keyAPI := etcdv2.NewKeysAPI(client)
	return &Etcd2{
		Client: &client,
		Key:    keyAPI,
	}, nil
}

// Save saves data in etcd
func (e Etcd2) Save(key string, data interface{}, callback datastore.NoResultCallback) {
	go func(e Etcd2, k string, d interface{}, c datastore.NoResultCallback) {
		value, marshErr := json.Marshal(d)
		if marshErr != nil {
			c(marshErr)
			return
		}
		_, setErr := e.Key.Set(context.Background(), k, (string)(value), nil)
		c(setErr)
	}(e, key, data, callback)
}

// Find finds data in etcd
func (e Etcd2) Find(queryStr string, callback datastore.ResultCallback) {
	go func(e Etcd2, q string, c datastore.ResultCallback) {
		resp, getErr := e.Key.Get(context.Background(), q, nil)
		c(([]byte)(resp.Node.Value), getErr)
	}(e, queryStr, callback)
}

// Remove removes data in etcd
func (e Etcd2) Remove(queryStr string, callback datastore.NoResultCallback) {
	go func(e Etcd2, q string, c datastore.NoResultCallback) {
		_, delErr := e.Key.Delete(context.Background(), q, nil)
		c(delErr)
	}(e, queryStr, callback)
}

// Update updates data in etcd
func (e Etcd2) Update(queryStr string, update interface{}, callback datastore.NoResultCallback) {
	go func(e Etcd2, q string, u interface{}, c datastore.NoResultCallback) {
		value, marshErr := json.Marshal(u)
		if marshErr != nil {
			callback(marshErr)
			return
		}
		_, updateErr := e.Key.Update(context.Background(), q, (string)(value))
		c(updateErr)
	}(e, queryStr, update, callback)
}

// FindAndUpdate updates data, then returns it
func (e Etcd2) FindAndUpdate(queryStr string, update interface{}, callback datastore.ResultCallback) {
	go func(e Etcd2, q string, u interface{}, c datastore.ResultCallback) {
		value, marshErr := json.Marshal(u)
		if marshErr != nil {
			c(nil, marshErr)
			return
		}
		resp, updateErr := e.Key.Update(context.Background(), q, (string)(value))
		c(([]byte)(resp.Node.Value), updateErr)
	}(e, queryStr, update, callback)
}
