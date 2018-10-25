package client

import (
	"context"
	"encoding/json"
	"path/filepath"
	"time"

	"go.etcd.io/etcd/client"
)

const ResourcePrefix = "/api"

func Create(resourceType string, resourceName string, data interface{}) (interface{}, error) {
	cfg := client.Config{
		Endpoints: []string{
			"http://127.0.0.1:2379",
		},
		Transport: client.DefaultTransport,
	}

	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	kAPI := client.NewKeysAPI(c)

	dataBytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	return kAPI.Set(ctx, filepath.Join(ResourcePrefix, resourceType, resourceName), string(dataBytes), &client.SetOptions{
		PrevExist: client.PrevNoExist,
	})
}
