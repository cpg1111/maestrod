package mongodb

import (
	"encoding/json"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"
)

var mongo, mErr = New(os.Getenv("MONGO_SERVICE_HOST"), os.Getenv("MONGO_SERVICE_PORT"), "", "")

type testData struct {
	Msg string `json:"msg", bson:"msg"`
}

func fmtQuery(msg string) (string, error) {
	query := &mongoQuery{
		Collection: "test",
		Query: testData{
			Msg: msg,
		},
	}
	queryBytes, queryErr := json.Marshal(query)
	return (string)(queryBytes), queryErr
}

func TestMongoSave(t *testing.T) {
	if mErr != nil {
		t.Error(mErr)
		return
	}
	data := testData{
		Msg: "hello save",
	}
	mongo.db.C("test").Create(&mgo.CollectionInfo{Capped: false})
	doneChan := make(chan bool)
	mongo.Save("test", data, func(err error) {
		if err != nil {
			t.Error(err)
			doneChan <- true
		}
		resQuery := mongo.db.C("test").Find(data)
		res := &testData{}
		resErr := resQuery.One(res)
		if resErr != nil {
			t.Error(resErr)
			doneChan <- true
		}
		if res.Msg != "hello save" {
			t.Errorf("expected hello save found %s", res.Msg)
			doneChan <- true
		}
		doneChan <- true
	})
	_ = <-doneChan
}

func TestMongoFind(t *testing.T) {
	if mErr != nil {
		t.Error(mErr)
		return
	}
	mongo.db.C("test").Create(&mgo.CollectionInfo{Capped: false})
	mongo.db.C("test").Insert(testData{Msg: "hello find"})
	doneChan := make(chan bool)
	queryString, queryErr := fmtQuery("hello find")
	if queryErr != nil {
		t.Error(queryErr)
	}
	mongo.Find(queryString, func(data interface{}, err error) {
		if err != nil {
			t.Error(err)
			doneChan <- true
		}
		if data == nil {
			t.Error("unable to find specified data")
			doneChan <- true
		}
		doneChan <- true
	})
	_ = <-doneChan
}

func TestMongoRemove(t *testing.T) {
	if mErr != nil {
		t.Error(mErr)
	}
	mongo.db.C("test").Create(&mgo.CollectionInfo{Capped: false})
	mongo.db.C("test").Insert(testData{Msg: "hello remove"})
	doneChan := make(chan bool)
	queryString, queryErr := fmtQuery("hello remove")
	if queryErr != nil {
		t.Error(queryErr)
	}
	mongo.Remove(queryString, func(err error) {
		if err != nil {
			t.Error(err)
			doneChan <- true
		}
		query := mongo.db.C("test").Find(testData{
			Msg: "hello remove",
		})
		res := &testData{}
		resErr := query.One(res)
		if res.Msg != "" || resErr == nil {
			t.Error("Found a removed record")
			doneChan <- true
		}
		doneChan <- true
	})
	_ = <-doneChan
}

func TestMongoUpdate(t *testing.T) {
	if mErr != nil {
		t.Error(mErr)
	}
	mongo.db.C("test").Create(&mgo.CollectionInfo{Capped: false})
	mongo.db.C("test").Insert(testData{Msg: "hello update"})
	doneChan := make(chan bool)
	queryString, queryErr := fmtQuery("hello update")
	if queryErr != nil {
		t.Error(queryErr)
	}
	mongo.Update(queryString, testData{Msg: "hello update 2"}, func(err error) {
		if err != nil {
			t.Error(err)
			doneChan <- true
		}
		query := mongo.db.C("test").Find(testData{Msg: "hello update 2"})
		res := &testData{}
		resErr := query.One(res)
		if resErr != nil {
			t.Error(resErr)
			doneChan <- true
		}
		doneChan <- true
	})
	_ = <-doneChan
}

func TestMongoFindAndUpdate(t *testing.T) {
	if mErr != nil {
		t.Error(mErr)
	}
	mongo.db.C("test").Create(&mgo.CollectionInfo{Capped: false})
	mongo.db.C("test").Insert(testData{Msg: "hello findAndUpdate"})
	doneChan := make(chan bool)
	queryString, queryErr := fmtQuery("hello findAndUpdate")
	if queryErr != nil {
		t.Error(queryErr)
	}
	mongo.FindAndUpdate(queryString, testData{Msg: "hello findAndUpdate 2"}, func(res interface{}, err error) {
		if res == nil {
			t.Error("did not find updated result")
			doneChan <- true
		}
		if err != nil {
			t.Error(err)
			doneChan <- true
		}
		doneChan <- true
	})
	_ = <-doneChan
}
