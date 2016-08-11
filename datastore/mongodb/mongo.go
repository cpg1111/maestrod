package mongodb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cpg1111/maestrod/datastore"

	mgo "gopkg.in/mgo.v2"
)

type MongoStore struct {
	datastore.Datastore
	store *mgo.Session
	db    *mgo.Database
}

type mongoQuery struct {
	Collection string      `json:"collection"`
	Query      interface{} `json:"query"`
}

type mongoRes struct{}

func New(host, port, username, password string) (*MongoStore, error) {
	info := &mgo.DialInfo{
		Addrs:    []string{fmt.Sprintf("%s:%s", host, port)},
		Direct:   false,
		Timeout:  3 * time.Second,
		FailFast: true,
		Database: "maestrod",
		Source:   "admin",
		Username: username,
		Password: password,
	}
	session, sessErr := mgo.DialWithInfo(info)
	if sessErr != nil {
		return nil, sessErr
	}
	db := session.DB("maestrod")
	collectionInfo := &mgo.CollectionInfo{
		Capped: false,
	}
	db.C("configs").Create(collectionInfo)
	db.C("projects").Create(collectionInfo)
	db.C("queueSnapShots").Create(collectionInfo)
	return &MongoStore{
		store: session,
		db:    db,
	}, nil
}

func (m MongoStore) Save(key string, data interface{}, callback datastore.NoResultCallback) {
	go func() {
		collection := m.db.C(key)
		queryErr := collection.Insert(data)
		callback(queryErr)
	}()
}

func (m MongoStore) Find(queryStr string, callback datastore.ResultCallback) {
	go func() {
		query := &mongoQuery{}
		unmarshErr := json.Unmarshal(([]byte)(queryStr), query)
		if unmarshErr != nil {
			callback(nil, unmarshErr)
			return
		}
		collection := m.db.C(query.Collection)
		resQuery := collection.Find(query.Query)
		res := &mongoRes{}
		queryErr := resQuery.One(res)
		callback(res, queryErr)
	}()
}

func (m MongoStore) Remove(queryStr string, callback datastore.NoResultCallback) {
	go func() {
		query := &mongoQuery{}
		unmarshErr := json.Unmarshal(([]byte)(queryStr), query)
		if unmarshErr != nil {
			callback(unmarshErr)
			return
		}
		collection := m.db.C(query.Collection)
		queryErr := collection.Remove(query.Query)
		callback(queryErr)
	}()
}

func (m MongoStore) Update(queryStr string, update interface{}, callback datastore.NoResultCallback) {
	go func() {
		query := &mongoQuery{}
		unmarshErr := json.Unmarshal(([]byte)(queryStr), query)
		if unmarshErr != nil {
			callback(unmarshErr)
			return
		}
		collection := m.db.C(query.Collection)
		queryErr := collection.Update(query.Query, update)
		callback(queryErr)
	}()
}

func (m MongoStore) FindAndUpdate(queryStr string, update interface{}, callback datastore.ResultCallback) {
	doneChan := make(chan bool)
	go m.Update(queryStr, update, func(err error) {
		if err != nil {
			callback(nil, err)
			doneChan <- true
			return
		}
		updated := &mongoQuery{}
		unmarshErr := json.Unmarshal(([]byte)(queryStr), updated)
		if unmarshErr != nil {
			callback(nil, unmarshErr)
			doneChan <- true
			return
		}
		updated.Query = update
		newQueryBytes, marshErr := json.Marshal(updated)
		if marshErr != nil {
			callback(nil, marshErr)
			doneChan <- true
			return
		}
		m.Find((string)(newQueryBytes), callback)
		doneChan <- true
	})
	_ = <-doneChan
}
