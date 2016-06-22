package main

import (
	"flag"
	"log"
	"os"

	"github.com/cpg1111/maestrod/config"
	"github.com/cpg1111/maestrod/datastore"
	"github.com/cpg1111/maestrod/server"
)

var (
	configPath    = flag.String("config-path", "/etc/maestrod/conf.toml", "path to the config file to load, defaults to /etc/maestrod/conf.toml")
	runtime       = flag.String("runtime", "", "type of runtime, defaults to native or the value specifcied in the configuration file, other options are: docker, kubernetes, rkt, EC2, GCE, DO, libvirt")
	hostIP        = flag.String("host-ip", "", "host ip for the server to bind to, defaults to 127.0.0.1 or the value specifcied in the configuration file")
	port          = flag.Uint("port", 0, "port number for the server to listen on, defaults to 8484 or the value specifcied in the configuration file")
	workspaceDir  = flag.String("workspace-dir", "", "working directory for maestro cloning and building, defaults to /tmp/maestro")
	datastoreType = flag.String("datastore-type", "redis", "type of data store to persist configuration in, defaults to redis")
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
	if *port != 0 {
		conf.Server.Port = *port
	}
	if *workspaceDir != "" {
		conf.Server.WorkspaceDir = *workspaceDir
	}
	if *datastoreType != "" {
		conf.Server.DataStoreType = *datastoreType
	}
	var datastoreHost string
	var datastorePort string
	if conf.Server.DataStoreEnvIP != "" {
		datastoreHost = os.Getenv(conf.Server.DataStoreEnvIP)
	} else {
		datastoreHost = conf.Server.DataStoreStaticIP
	}
	if conf.Server.DataStoreEnvPort != "" {
		datastorePort = os.Getenv(conf.Server.DataStoreEnvPort)
	} else {
		datastorePort = conf.Server.DataStoreStaticPort
	}
	var store datastore.Datastore
	switch conf.Server.DataStoreType {
	case "redis":
		store = datastore.NewRedis(datastoreHost, datastorePort, conf.Server.DataStorePWD)
		log.Println("redis datastore created")
	default:
		log.Fatal("specifcied datastore currently not supported, please create an issue @ https://github.com/cpg1111/maestrod")
	}
	server.Run(&conf.Server, &store)
}
