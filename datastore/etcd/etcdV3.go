package etcd

import (
	"github.com/cpg1111/maestrod/datastore"

	etcdv3 "github.com/coreos/etcd/clientv3"
)

type Etcd3 struct {
	datastore.Datastore
	Client etcdv3.Client
}

func (e Etcd3) Save(key string, data interface{}, callback datastore.NoResultCallback) {

}

func (e Etcd3) Find(queryStr string, callback datastore.ResultCallback) {

}

func (e Etcd3) Remove(queryStr string, callback datastore.NoResultCallback) {

}

func (e Etcd3) Update(queryStr string, update interface{}, callback datastore.NoResultCallback) {

}

func (e Etcd3) FindAndUpdate(queryStr string, update interface{}, callback datastore.ResultCallback) {

}
