package config

import (
	"errors"

	"github.com/coredhcp/coredhcp/logger"
)

var _log = logger.GetLogger("etcd")

type EtcdConfig struct {
	Endpoints []string
}

func (c *Config) ParseEtcdConfig() error {

	if !c.v.IsSet("Etcd.Endpoints") {
		return errors.New("Missing etcd endpoints")
	}

	endpoints := c.v.GetStringSlice("Etcd.Endpoints")
	if len(endpoints) == 0 {
		return errors.New("Endpoint list is empty")
	}

	log.Println("Etcd: found endpoints:", endpoints)

	return nil
}
