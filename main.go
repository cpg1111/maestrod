package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/cpg1111/maestrod/config"
	"github.com/cpg1111/maestrod/datastore"
	etcd2 "github.com/cpg1111/maestrod/datastore/etcd/v2"
	etcd3 "github.com/cpg1111/maestrod/datastore/etcd/v3"
	"github.com/cpg1111/maestrod/datastore/mongodb"
	"github.com/cpg1111/maestrod/datastore/redis"
	"github.com/cpg1111/maestrod/gitactivity"
	"github.com/cpg1111/maestrod/lifecycle"
	"github.com/cpg1111/maestrod/manager"
	"github.com/cpg1111/maestrod/manager/docker"
	"github.com/cpg1111/maestrod/manager/k8s"
	"github.com/cpg1111/maestrod/statecom"
)

var (
	configPath    = flag.String("config-path", "/etc/maestrod/conf.toml", "path to the config file to load, defaults to /etc/maestrod/conf.toml")
	runtime       = flag.String("runtime", "k8s", "type of runtime, defaults to k8s or the value specifcied in the configuration file, other options are: docker, spawnd, rkt, EC2, libvirt")
	hostIP        = flag.String("host-ip", "", "host ip for the gitactivity to bind to, defaults to 127.0.0.1 or the value specifcied in the configuration file")
	port          = flag.Uint("port", 0, "port number for the gitactivity to listen on, defaults to 8484 or the value specifcied in the configuration file")
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
	var (
		store    datastore.Datastore
		storeErr error
	)
	switch conf.Server.DataStoreType {
	case "redis":
		store = redis.New(datastoreHost, datastorePort, conf.Server.DataStorePWD)
		storeErr = nil
		log.Println("redis datastore created")
	case "etcd2":
		store, storeErr = etcd2.NewV2(datastoreHost, datastorePort)
		log.Println("etcd2 datastore created")
	case "etcd3":
		store, storeErr = etcd3.NewV3(datastoreHost, datastorePort)
		log.Println("etcd3 datastore created")
	case "mongodb":
		store, storeErr = mongodb.New(datastoreHost, datastorePort, conf.Server.DataStoreUser, conf.Server.DataStorePWD)
		log.Println("mongodb datastore create")
	default:
		log.Fatal("specifcied datastore currently not supported, please create an issue @ https://github.com/cpg1111/maestrod")
	}
	if storeErr != nil {
		log.Fatal(storeErr)
	}
	return &store
}

func getManager(conf *config.Config) manager.Driver {
	switch conf.Server.Runtime {
	case "docker":
		driver, dErr := docker.New("v1.23", conf.Server.MaestroVersion)
		if dErr != nil {
			log.Fatal(dErr)
		}
		return *driver
	case "kubernetes":
	case "k8s":
		driver := k8s.New(conf.Server.MaestroVersion, &conf.Server)
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
	gitactivity.Run(&conf.Server, store, queue)
	statecom.Run(conf.Server.Host, conf.Server.ServerCertPath, conf.Server.ServerKeyPath, (int)(conf.Server.StateComPort), store)
	running := &lifecycle.Running{}
	managerDriver := getManager(conf)
	var err error
	for err == nil {
		err = lifecycle.Check(conf, queue, running, managerDriver)
		time.Sleep(3 * time.Second)
	}
	log.Fatal(err)
}
