package k8s

import (
	"plugin"

	"github.com/cpg1111/maestrod/config"
)

func PluginDriver(maestroVersion string, conf *config.Server) *Driver {
	return New(maestroVersion, conf)
}
