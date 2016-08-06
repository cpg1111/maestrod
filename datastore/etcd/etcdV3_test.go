package etcd

import (
	"encoding/json"
	"os"
	"testing"

	"golang.org/x/net/context"
)

var etcd3, etcd3Err = NewV3(os.Getenv("ETCD3_SERVICE_HOST"), os.Getenv("ETCD3_SERVICE_PORT"))

type EtcdV3TestData struct {
	Message string
}

func TestEtcd3Save(t *testing.T) {
	if etcd3Err != nil {
		t.Error(etcd3Err)
		return
	}
	testData := &EtcdV3TestData{
		Message: "test",
	}
	testValue, marshErr := json.Marshal(testData)
	done := make(chan bool)
	etcd3.Save("testFind", testValue, func(err error) {
		if err != nil {
			t.Error(err)
			done <- true
		}
		resp, respErr := etcd.Key.Get(context.Background(), "testFind", nil)
		if respErr != nil {
			t.Error(respErr)
			done <- true
		}
		for i := range resp.Kvs {
			if (string)(resp.Kvs[i].Value) != "test" {
				t.Errorf("expected test, found %s", (string)(resp.Kvs[i].Value))
				done <- true
				return
			}
		}
		done <- true
	})
	_ = <-done
}

func TestEtcd3Find(t *testing.T) {
	if etcd3Err != nil {
		t.Error(etcd3Err)
		return
	}
	testData := &EtcdV3TestData{
		Message: "test",
	}
	testValue, marshErr := json.Marshal(testData)
	_, putErr := etcd3.Key.Put(context.Background(), "testFind", testValue, nil)
	if putErr != nil {
		t.Error(putErr)
	}
	done := make(chan bool)
	etcd3.Find("testFind", func(val interface{}, err error) {
		if err != nil {
			t.Error(err)
			done <- true
		}
		result := &EtcdV3TestData{}
		unmarshErr := json.Unmarshal(([]byte)(val), result)
		if unmarshErr != nil {
			t.Error(unmarshErr)
		}
		if result.Message != "test" {
			t.Errorf("expected test, found %s", result.Message)
			done <- true
		}
		done <- true
	})
	_ = <-done
}
