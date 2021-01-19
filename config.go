package main

import (
	"github.com/ansd/lastpass-go"
	swarm "github.com/docker/docker/api/types/swarm"
)

// LastpassSecrets ...
type LastpassSecrets map[string]*lastpass.Account

// DockerSecrets ...
type DockerSecrets map[string]swarm.Secret

// ConfigFile ...
type ConfigFile struct {
	Log struct {
		Level string
	}
	LastPass struct {
		TwoFactor string
	}
	Secrets struct {
		Groups []string
		Lists  []string
	}
}
