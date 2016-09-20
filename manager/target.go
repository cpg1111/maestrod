package manager

import (
	"fmt"

	"github.com/cpg1111/maestrod/config"
)

// GetTarget creates an address to speak with the runtime target through the config
// if TargetEnvHost and TargetEnvPort are populated then it uses those for host and port
// otherwise it uses TargetHost and TargetPort, it requires TargetProtocol
func GetTarget(conf *config.Server) string {
	var (
		host string
		port string
	)
	if len(conf.TargetEnvHost) > 0 {
		host = conf.TargetEnvHost
	} else {
		host = conf.TargetHost
	}
	if len(conf.TargetEnvPort) > 0 {
		port = conf.TargetEnvPort
	} else {
		port = conf.TargetPort
	}
	return fmt.Sprintf("%s://%s:%s", conf.TargetProtocol, host, port)
}
