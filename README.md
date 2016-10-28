# maestrod
Manager Daemon for maestro https://github.com/cpg1111/maestro

## Building

If you have glide (https://github.com/Masterminds/glide)

```
make
make install
```

Otherwise if you have docker

```
make docker
```

## Running

Maestrod requires a "runtime" to run, the currently supported runtimes are:

- Kubernetes (https://github.com/kubernetes/kubernetes)

- Docker Engine (https://github.com/docker/docker)

Maestrod also requires a Key Value datastore, the currently supported datastores are:

- Etcd v2 or v3 (https://github.com/coreos/etcd)

- MongoDB (https://github.com/mongodb/mongo)

- Redis (https://github.com/antirez/redis)

(Once Go 1.8.X is released both runtimes and datastores will be standard Go plugins that anyone can create)

To configure Maestrod you have a toml file as follows:

```
[Server]
Runtime=<the runtime to be used> # Currently can be: k8s, kubernetes or docker REQUIRED
MaxBuilds=<the number of max active builds you want> # REQUIRED
DataStoreType=<the datastore to be used> # Currently can be: etcd2, etcd3, mongodb, or redis REQUIRED
DataStoreUser=<any user to auth with datastore>
DataStorePWD=<any password to auth with datastore> # plain text for now, but it will take a hash in the future
DataStoreEnvIP=<env var name for the datastore IP addr> 
DataStoreStaticIP=<hardcoded IP addr of the datastore>
DataStoreEnvPort=<env var name for the datastore port number>
DataStoreStaticPort=<hardcoded port number of the datastore>
TargetProtocol=<protocol to speak to the runtime> # REQUIRE
TargetHost=<hardcoded IP addr for the runtime>
TargetEnvHost=<env var for the runtime IP addr>
TargetPort=<hardcoded port number for the runtime>
TargetEnvPort=<env var for the runtime port number>
ClientCertPath=<path to a certificate for clients to the runtime>
ClientKeyPath=<path to cert key for clients to the runtime>
ServerCertPath=<path to cert for serving the webhooks over https>
ServerKeyPath=<path to cert key for serving the webhooks over https>
MaestroVersion=<version of maestro to run>
Host=<host to bind to>
InsecurePort=<Insecure port to listen on>
SecurePort=<Secure port to listen on>
StateComPort=<Port for state communication between maestro and maestrod>
WorkspaceDir=<The directory in which each maestro worker clones their project into>

[[Projects]]
Name=<project name>
MaestroConfPath=<path to a maestro conf>
DeployBranches=<which branches to run deployment on when pushed>

[[Mounts]]
Kind=<type> # currently only hostPath is supported
Path=<path both on the host and container>
Name=<name of volume mount>

```

Then run

```
maestrod \
--config-path=<path to config>
--runtime=<conf override for runtime>
--host-ip=<host to bind to config override>
--port=<port to listen to config override>
--workspace-dir=<path to clone projects into config override>
--datastore-type=<datastore config override>
```

Then add `http(s)://<public maestro ip>:<maestro port>/push` as a webhook for repo pushes and 

`http(s)://<public maestro ip>:<maestro port>/pullrequest` for pull request hooks.

