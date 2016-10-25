package docker

import (
	"plugin"

	"github.com/cpg1111/maestrod/config"
)

func PluginDriver(maestroVersion string, conf *config.Server) *Driver {
	driver, dErr := New("v1.23", maestroVersion)
	if dErr != nil {
		panic(dErr)
	}
}
