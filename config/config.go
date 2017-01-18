package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	gs "cloud.google.com/go/storage"
	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	awscreds "github.com/aws/aws-sdk-go/aws/credentials"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"google.golang.org/api/option"
)

// Mount is a struct in the config for Runtime mounts
type Mount struct {
	Kind      string
	ID        string
	FSType    string
	Server    string
	Endpoints string
	Path      string
	Name      string
	ReadOnly  bool
}

// Server is a struct for the runtime server
type Server struct {
	Runtime             string
	RuntimeTLSClient    bool
	RuntimeTLSServer    bool
	MaxBuilds           int
	DataStoreType       string
	DataStoreUser       string
	DataStorePWD        string
	DataStoreEnvIP      string
	DataStoreStaticIP   string
	DataStoreEnvPort    string
	DataStoreStaticPort string
	TargetProtocol      string
	TargetHost          string
	TargetPort          string
	TargetEnvHost       string
	TargetEnvPort       string
	ClientCertPath      string
	ClientKeyPath       string
	ServerCertPath      string
	ServerKeyPath       string
	MaestroVersion      string
	Host                string
	SecurePort          uint
	InsecurePort        uint
	StateComPort        uint
	WorkspaceDir        string
}

// Project is a struct in the config for each project for maestrod to spin up
type Project struct {
	Name            string   `json:"name"`
	MaestroConfPath string   `json:"confPath"`
	DeployBranches  []string `json:"deployBranches"`
}

// Config is the struct of the config file
type Config struct {
	Server   Server
	Projects []Project
	Mounts   []Mount
}

type remoteConfig struct {
	Storage string
	Bucket  string
	Object  string
}

func decode(r io.Reader) (*Config, error) {
	var conf Config
	if _, pErr := toml.DecodeReader(r, &conf); pErr != nil {
		return nil, pErr
	}
	return &conf, nil
}

func parseRemote(path string) *remoteConfig {
	storageIdx := strings.Index(path, "://")
	pathSlice := strings.Split(path[storageIdx+1:], "/")
	obj := pathSlice[1]
	if len(pathSlice) > 2 {
		for i := 2; i < len(pathSlice); i++ {
			obj = fmt.Sprintf("%s/%s", obj, pathSlice[i])
		}
	}
	return &remoteConfig{
		Storage: path[0:storageIdx],
		Bucket:  pathSlice[0],
		Object:  obj,
	}
}

func loadS3(path string) (*Config, error) {
	remote := parseRemote(path)
	creds := awscreds.NewEnvCredentials()
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}
	config := &aws.Config{
		Region:           aws.String(os.Getenv("AWS_S3_REGION")),
		Endpoint:         aws.String("s3.amazonaws.com"),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
		LogLevel:         aws.LogLevel(aws.LogLevelType(0)),
	}
	session := awssession.New(config)
	s3Client := s3.New(session)
	query := &s3.GetObjectInput{
		Bucket: aws.String(remote.Bucket),
		Key:    aws.String(remote.Object),
	}
	resp, err := s3Client.GetObject(query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return decode(resp.Body)
}

func loadGStorage(path string) (*Config, error) {
	remote := parseRemote(path)
	opts := options.WithServiceAccountFile(os.Getenv("GCLOUD_SVC_ACCNT_FILE"))
	ctx := context.Background()
	gsClient := gs.NewClient(ctx, opts)
	bucket := gsClient.Bucket(remote.Bucket)
	obj := bucket.Object(remote.Object)
	rdr, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rdr.Close()
	return decode(rdr)
}

func loadLocal(path string) (*Config, error) {
	conf, readErr := os.OpenFile(path, os.O_RDONLY, 0644)
	if readErr != nil {
		return nil, readErr
	}
	return decode(conf)
}

// Load loads a config file and returns a pointer to a config struct
func Load(path string) (*Config, error) {
	if strings.Contains(path, "s3://") {
		return loadS3(path)
	}
	if strings.Contains(path, "gs://") {
		return loadGStorage(path)
	}
	return loadLocal(path)
}
