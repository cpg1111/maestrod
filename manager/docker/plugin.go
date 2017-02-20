package docker

import (
	"github.com/cpg1111/maestrod/config"
	"github.com/cpg1111/maestrod/manager"
	d "github.com/cpg1111/maestrod/manager/docker/driver"
)

func PluginDriver(maestroVersion string, conf *config.Config) manager.Driver {
	driver, dErr := d.New("v1.23", maestroVersion)
	if dErr != nil {
		panic(dErr)
	}
	return driver
}
