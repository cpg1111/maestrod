package main

import (
	"flag"
	"log"

	"github.com/cpg1111/maestrod/config"
	"github.com/cpg1111/maestrod/server"
)

var (
	configPath   = flag.String("config-path", "/etc/maestrod/conf.toml", "path to the config file to load, defaults to /etc/maestrod/conf.toml")
	runtime      = flag.String("runtime", "", "type of runtime, defaults to native or the value specifcied in the configuration file, other options are: docker, kubernetes, rkt, EC2, GCE, DO, libvirt")
	hostIP       = flag.String("host-ip", "", "host ip for the server to bind to, defaults to 127.0.0.1 or the value specifcied in the configuration file")
	port         = flag.Uint("port", 0, "port number for the server to listen on, defaults to 8484 or the value specifcied in the configuration file")
	workspaceDir = flag.String("workspace-dir", "", "working directory for maestro cloning and building, defaults to /tmp/maestro")
)

func main() {
	flag.Parse()
	conf, loadErr := config.Load(*configPath)
	if loadErr != nil {
		log.Fatal(loadErr)
	}
	if *hostIP != "" {
		conf.Server.Host = *hostIP
	}
	if conf.Server.Port != 0 {
		conf.Server.Port = *port
	}
	if conf.Server.WorkspaceDir != "" {
		conf.Server.WorkspaceDir = *workspaceDir
	}
	server.Run(&conf.Server)
}
