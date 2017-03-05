package consul

import (
	"fmt"
	"github.com/docker/libkv"
	libkvStore "github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"
	"github.com/tboerger/redirects/config"
	"github.com/tboerger/redirects/store"
	"time"
)

func init() {
	consul.Register()
}

// data is a basic struct that iplements the Store interface.
type data struct {
	store  libkvStore.Store
	prefix string
}

// New initializes a new Consul store.
func New(s libkvStore.Store, prefix string) store.Store {
	return &data{
		store:  s,
		prefix: prefix,
	}
}

// Load initializes the Consul storage.
func Load() store.Store {
	prefix := config.Consul.Prefix

	s, err := libkv.NewStore(
		libkvStore.CONSUL,
		config.Consul.Endpoints,
		&libkvStore.Config{
			ConnectionTimeout: config.Consul.Timeout * time.Second,
		},
	)

	// TODO: Handle this error properly
	if err != nil {
		panic(fmt.Sprintf("TODO: Failed to init store. %s", err))
	}

	if ok, _ := s.Exists(prefix); !ok {
		err := s.Put(
			prefix,
			nil,
			&libkvStore.WriteOptions{
				IsDir: true,
			},
		)

		// TODO: Handle this error properly
		if err != nil {
			panic(fmt.Sprintf("TODO: Failed to create prefix. %s", err))
		}
	}

	return New(
		s,
		prefix,
	)
}