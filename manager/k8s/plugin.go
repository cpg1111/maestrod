package k8s

import (
	"github.com/cpg1111/maestrod/config"
	"github.com/cpg1111/maestrod/manager"
	d "github.com/cpg1111/maestrod/manager/k8s/driver"
)

func PluginDriver(maestroVersion string, conf *config.Config) manager.Driver {
	return d.New(maestroVersion, conf)
}
