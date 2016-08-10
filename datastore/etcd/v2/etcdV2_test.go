package etcd

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/context"
)

var etcd, etcdErr = NewV2(os.Getenv("ETCD2_SERVICE_HOST"), os.Getenv("ETCD2_SERVICE_PORT"))

type Etcd2TestData struct {
	Message string
}

func TestEtcd2Save(t *testing.T) {
	if etcdErr != nil {
		t.Error(etcdErr)
	}
	_, _ = etcd.Key.Delete(context.Background(), "testSave", nil)
	done := make(chan bool)
	etcd.Save("testSave", Etcd2TestData{Message: "test"}, func(err error) {
		if err != nil {
			t.Error(err)
			done <- true
		}
		resp, respErr := etcd.Key.Get(context.Background(), "testSave", nil)
		if respErr != nil {
			t.Error(respErr)
			done <- true
		}
		result := &Etcd2TestData{}
		unmarshErr := json.Unmarshal(([]byte)(resp.Node.Value), result)
		if unmarshErr != nil {
			t.Error(unmarshErr)
			done <- true
		}
		if result.Message != "test" {
			t.Errorf("Expected test, found %s", result.Message)
			done <- true
		}
		done <- true
	})
	_ = <-done
}

func TestEtcd2Find(t *testing.T) {
	testData := &Etcd2TestData{
		Message: "test",
	}
	_, _ = etcd.Key.Delete(context.Background(), "testFind", nil)
	testValue, marshErr := json.Marshal(testData)
	if marshErr != nil {
		t.Error(marshErr)
	}
	_, setErr := etcd.Key.Create(context.Background(), "testFind", (string)(testValue))
	if setErr != nil {
		t.Error(setErr)
	}
	done := make(chan bool)
	etcd.Find("testFind", func(val []byte, err error) {
		if err != nil {
			t.Error(err)
			done <- true
		}
		result := &Etcd2TestData{}
		unmarshErr := json.Unmarshal(val, result)
		if unmarshErr != nil {
			t.Error(unmarshErr)
			done <- true
		}
		if result.Message != "test" {
			t.Errorf("expected test foudn %s", result.Message)
			done <- true
		}
		done <- true
	})
	_ = <-done
}

func TestEtcd2Remove(t *testing.T) {
	testData := &Etcd2TestData{
		Message: "test",
	}
	_, _ = etcd.Key.Delete(context.Background(), "testRemove", nil)
	testValue, marshErr := json.Marshal(testData)
	if marshErr != nil {
		t.Error(marshErr)
	}
	_, setErr := etcd.Key.Set(context.Background(), "testRemove", (string)(testValue), nil)
	if setErr != nil {
		t.Error(setErr)
	}
	done := make(chan bool)
	etcd.Remove("testRemove", func(err error) {
		if err != nil {
			t.Error(err)
			done <- true
		}
		res, resErr := etcd.Key.Get(context.Background(), "testRemove", nil)
		correctErr := strings.Index(resErr.Error(), "Key not found") > -1
		if resErr != nil && !correctErr {
			t.Error(resErr)
			done <- true
		}
		if res != nil && res.Node != nil && res.Node.Value != "" {
			t.Error("did not remove value testRemove from etcd")
			done <- true
		}
		done <- true
	})
	_ = <-done
}

func TestEtcd2Update(t *testing.T) {
	testData := &Etcd2TestData{
		Message: "test",
	}
	_, _ = etcd.Key.Delete(context.Background(), "testUpdate", nil)
	testValue, marshErr := json.Marshal(testData)
	if marshErr != nil {
		t.Error(marshErr)
	}
	_, setErr := etcd.Key.Create(context.Background(), "testUpdate", (string)(testValue))
	if setErr != nil {
		t.Error(setErr)
	}
	newData := &Etcd2TestData{
		Message: "update",
	}
	done := make(chan bool)
	t.Log(newData)
	etcd.Update("testUpdate", newData, func(err error) {
		if err != nil {
			t.Error(err)
			done <- true
		}
		resp, resErr := etcd.Key.Get(context.Background(), "testUpdate", nil)
		if resErr != nil {
			t.Error(resErr)
			done <- true
		}
		result := &Etcd2TestData{}
		unmarshErr := json.Unmarshal(([]byte)(resp.Node.Value), result)
		if unmarshErr != nil {
			t.Error(unmarshErr)
			done <- true
		}
		if result.Message != "update" {
			t.Errorf("expected update, found %s", result.Message)
			done <- true
		}
		done <- true
	})
	_ = <-done
}

func TestEtcd2FindAndUpdate(t *testing.T) {
	testData := &Etcd2TestData{
		Message: "test",
	}
	_, _ = etcd.Key.Delete(context.Background(), "testFindUpdate", nil)
	testValue, marshErr := json.Marshal(testData)
	if marshErr != nil {
		t.Error(marshErr)
	}
	_, setErr := etcd.Key.Create(context.Background(), "testFindUpdate", (string)(testValue))
	if setErr != nil {
		t.Error(setErr)
	}
	newData := &Etcd2TestData{
		Message: "update",
	}
	done := make(chan bool)
	etcd.FindAndUpdate("testUpdate", newData, func(val []byte, err error) {
		if err != nil {
			t.Error(err)
			done <- true
		}
		result := &Etcd2TestData{}
		unmarshErr := json.Unmarshal(val, result)
		if unmarshErr != nil {
			t.Error(unmarshErr)
		}
		if result.Message != "update" {
			t.Errorf("expected 'update', but found %s", result.Message)
			done <- true
		}
		done <- true
	})
	_ = <-done
}
