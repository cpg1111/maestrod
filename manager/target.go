/*
Copyright 2016 Christian Grabowski All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
