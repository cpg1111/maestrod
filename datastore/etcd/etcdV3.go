package etcd

import (
	"encoding/json"
	"time"

	"github.com/cpg1111/maestrod/datastore"

	etcdv3 "github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

type Etcd3 struct {
	datastore.Datastore
	Client *etcdv3.Client
	Key    etcdv3.KV
}

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

func (e Etcd3) Save(key string, data interface{}, callback datastore.NoResultCallback) {
	value, marshErr := json.Marshal(data)
	if marshErr != nil {
		callback(marshErr)
		return
	}
	_, putErr := e.Key.Put(context.Background(), key, (string)(value), nil)
	callback(putErr)
}

func (e Etcd3) Find(queryStr string, callback datastore.ResultCallback) {
	resp, respErr := e.Key.Get(context.Background(), queryStr, nil)
	if respErr != nil {
		callback(nil, respErr)
		return
	}
	var values []interface{}
	for i := range resp.Kvs {
		if (string)(resp.Kvs[i].Key) == queryStr {
			values := append(values, resp.Kvs[i])
		}
	}
	callback(values, nil)
}

func (e Etcd3) Remove(queryStr string, callback datastore.NoResultCallback) {
	_, delErr := e.Key.Delete(context.Background(), queryStr, nil)
	callback(delErr)
}

func (e Etcd3) Update(queryStr string, update interface{}, callback datastore.NoResultCallback) {
	value, marshErr := json.Marshal(update)
	if marshErr != nil {
		callback(marshErr)
		return
	}
	_, updateErr := e.Key.Put(context.Background(), queryStr, value, nil)
	callback(updateErr)
}

func (e Etcd3) FindAndUpdate(queryStr string, update interface{}, callback datastore.ResultCallback) {
	value, marshErr := json.Marshal(update)
	if marshErr != nil {
		callback(nil, marshErr)
		return
	}
	resp, updateErr := e.Key.Put(context.Background(), queryString, value, nil)
	var values []interface{}
	for i := range resp.Kvs {
		if i == queryString {
			values := append(values, resp.Kvs[i])
		}
	}
}
