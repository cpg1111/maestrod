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
	done := make(chan bool)
	etcd3.Save("testFind", testData, func(err error) {
		if err != nil {
			t.Error(err)
			done <- true
		}
		t.Log(etcd3)
		resp, respErr := etcd3.Key.Get(context.Background(), "testFind")
		if respErr != nil {
			t.Log("actual error")
			t.Error(respErr)
			done <- true
		}
		t.Log(resp, len(resp.Kvs))
		for i := range resp.Kvs {
			if (string)(resp.Kvs[i].Value) != "{\"Message\":\"test\"}" {
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
	if marshErr != nil {
		t.Error(marshErr)
	}
	_, putErr := etcd3.Key.Put(context.Background(), "testFind", (string)(testValue))
	if putErr != nil {
		t.Error(putErr)
	}
	done := make(chan bool)
	etcd3.Find("testFind", func(val []byte, err error) {
		if err != nil {
			t.Error(err)
			done <- true
		}
		result := &EtcdV3TestData{}
		unmarshErr := json.Unmarshal(val, result)
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

func TestEtcd3Remove(t *testing.T) {
	if etcd3Err != nil {
		t.Error(etcd3Err)
		return
	}
	testData := &EtcdV3TestData{
		Message: "test",
	}
	testValue, marshErr := json.Marshal(testData)
	if marshErr != nil {
		t.Error(marshErr)
	}
	_, putErr := etcd3.Key.Put(context.Background(), "testRemove", (string)(testValue))
	if putErr != nil {
		t.Error(putErr)
	}
	done := make(chan bool)
	etcd3.Remove("testRemove", func(err error) {
		if err != nil {
			t.Error(err)
			done <- true
			return
		}
		resp, getErr := etcd3.Key.Get(context.Background(), "testRemove")
		if getErr != nil {
			t.Error(getErr)
			done <- true
			return
		}
		if resp != nil && len(resp.Kvs) > 0 {
			t.Error("did not remove testRemove")
		}
		done <- true
	})
	_ = <-done
}

func TestEtcd3Update(t *testing.T) {
	if etcd3Err != nil {
		t.Error(etcd3Err)
		return
	}
	testData := &EtcdV3TestData{
		Message: "test",
	}
	testValue, marshErr := json.Marshal(testData)
	if marshErr != nil {
		t.Error(marshErr)
		return
	}
	_, putErr := etcd3.Key.Put(context.Background(), "testUpdate", (string)(testValue))
	if putErr != nil {
		t.Error(putErr)
		return
	}
	update := &EtcdV3TestData{
		Message: "update",
	}
	done := make(chan bool)
	etcd3.Update("testUpdate", update, func(err error) {
		if err != nil {
			t.Error(err)
			done <- true
			return
		}
		resp, respErr := etcd3.Key.Get(context.Background(), "testUpdate")
		if respErr != nil {
			t.Error(respErr)
			done <- true
			return
		}
		if resp == nil || len(resp.Kvs) == 0 || resp.Kvs[0] == nil || len(resp.Kvs[0].Value) == 0 {
			t.Error("no value found")
			done <- true
			return
		}
		result := &EtcdV3TestData{}
		unmarshErr := json.Unmarshal(resp.Kvs[0].Value, result)
		if unmarshErr != nil {
			t.Error(unmarshErr)
			done <- true
			return
		}
		if result.Message != "update" {
			t.Errorf("expected 'update' but found %s", result.Message)
			done <- true
			return
		}
		done <- true
	})
	_ = <-done
}

func TestEtcdV3FindAndUpdate(t *testing.T) {
	if etcd3Err != nil {
		t.Error(etcd3Err)
		return
	}
	testData := &EtcdV3TestData{
		Message: "test",
	}
	testValue, marshErr := json.Marshal(testData)
	if marshErr != nil {
		t.Error(marshErr)
		return
	}
	_, putErr := etcd3.Key.Put(context.Background(), "testFindUpdate", (string)(testValue))
	if putErr != nil {
		t.Error(putErr)
		return
	}
	update := &EtcdV3TestData{
		Message: "update",
	}
	done := make(chan bool)
	etcd3.FindAndUpdate("testFindUpdate", update, func(val []byte, err error) {
		if err != nil {
			t.Error(err)
			done <- true
			return
		}
		result := &EtcdV3TestData{}
		unmarshErr := json.Unmarshal(val, result)
		if unmarshErr != nil {
			t.Error(unmarshErr)
			done <- true
			return
		}
		if result.Message != "update" {
			t.Errorf("expected 'update', found %s", result.Message)
			done <- true
			return
		}
		done <- true
	})
	_ = <-done
}
