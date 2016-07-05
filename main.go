package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/cpg1111/maestrod/config"
	"github.com/cpg1111/maestrod/datastore"
	"github.com/cpg1111/maestrod/lifecycle"
	"github.com/cpg1111/maestrod/manager"
	"github.com/cpg1111/maestrod/manager/docker"
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

func getConf() *config.Config {
	flag.Parse()
	conf, loadErr := config.Load(*configPath)
	if loadErr != nil {
		log.Fatal(loadErr)
	}
	if *hostIP != "" {
		conf.Server.Host = *hostIP
	}
	if *port != 0 {
		conf.Server.InsecurePort = *port
	}
	if *workspaceDir != "" {
		conf.Server.WorkspaceDir = *workspaceDir
	}
	if *datastoreType != "" {
		conf.Server.DataStoreType = *datastoreType
	}
	return conf
}

func getDataStore(conf *config.Config) *datastore.Datastore {
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
	return &store
}

func getManager(conf *config.Config) manager.Driver {
	switch conf.Server.Runtime {
	case "docker":
		driver, dErr := docker.New(conf.Server.TargetHost, "v1.23", conf.Server.MaestroVersion)
		if dErr != nil {
			log.Fatal(dErr)
		}
		return *driver
	default:
		log.Fatal("specifcied runtime is not supported yet, please create an issue @ https://github.com/cpg1111/maestrod")
	}
	return nil
}

func main() {
	conf := getConf()
	store := getDataStore(conf)
	queue := lifecycle.NewQueue(store)
	server.Run(&conf.Server, store, queue)
	running := &lifecycle.Running{}
	managerDriver := getManager(conf)
	var err error
	for err == nil {
		err = lifecycle.Check(conf, queue, running, managerDriver)
		time.Sleep(3 * time.Second)
	}
	log.Fatal(err)
}
