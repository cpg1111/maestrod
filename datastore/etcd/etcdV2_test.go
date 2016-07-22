package etcd

import (
	"os"
	"testing"
)

var etcd, etcdErr = NewV2(os.Getenv("ETC_SERVICE_HOST"), os.Getenv("ETCD_SERVICE_PORT"))

func TestEtcdSave(t *testing.T) {
	if etcdErr != nil {
		t.Error(etcdErr)
	}
}
