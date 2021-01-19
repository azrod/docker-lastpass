package main

import (
	"context"
	"strings"

	b64 "encoding/base64"

	"github.com/ansd/lastpass-go"
	swarm "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	dockerClient "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

type docker struct {
	*dockerClient.Client
}

func (docker *docker) Connect() (err error) {

	docker.Client, err = dockerClient.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error("Init docker sock")
	}

	return err
}

// CreateUsername ...
func (docker *docker) CreateUsername(la *lastpass.Account) (err error) {

	opts := swarm.SecretSpec{}
	opts.Annotations.Name = "lastpass_" + strings.ReplaceAll(la.Name, " ", "-") + "_Username"
	opts.Data = []byte(b64.StdEncoding.EncodeToString([]byte(la.Username)))

	_, err = docker.SecretCreate(context.Background(), opts)
	if err != nil {
		return err
	}

	return nil
}

// CreatePassword ...
func (docker *docker) CreatePassword(la *lastpass.Account) (err error) {

	opts := swarm.SecretSpec{}
	opts.Annotations.Name = "lastpass_" + strings.ReplaceAll(la.Name, " ", "-") + "_Password"
	opts.Data = []byte(b64.StdEncoding.EncodeToString([]byte(la.Password)))

	_, err = docker.SecretCreate(context.Background(), opts)
	if err != nil {
		return err
	}

	return nil
}
