package client

import (
	"conductor/pkg/resource"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/client"
	"path/filepath"
	"time"
	"github.com/google/uuid"
)

const ResourcePrefix = "/api"

func Create(resourceType string, resourceName string, data interface{}, dataResource *resource.BaseResource) (interface{}, error) {
	cfg := client.Config{
		Endpoints: []string{
			"http://127.0.0.1:2379",
		},
		Transport: client.DefaultTransport,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Errorf("client new error")
		return nil, err
	}

	kAPI := client.NewKeysAPI(c)
	path := filepath.Join(ResourcePrefix, resourceType, resourceName)
	log.Errorf("uuidgen  error")
	var UUID string
	UUID = uuid.New().String()
	dataResource.UUID = UUID
	dataResource.Version = 0
	dataResource.SelfLink = path
	dataResource.CreationTimestamp = time.Now()


	dataBytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Errorf("client unmarshalal error")

		return nil, err
	}

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))

	resp, err := kAPI.Set(ctx, path, string(dataBytes), &client.SetOptions{
		PrevExist: client.PrevNoExist,
	})
	if err != nil {
		log.Infof("error in creating data of type %s: %+v \n", resourceType, err)
		return nil, err
	} else {
		// print common key info
		log.Infof("Get is done. Metadata is %q\n", resp)
	}
	return resp.Node.Value, nil

}
func Update(resourceType string, resourceName string, data interface{}) (interface{}, error) {
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
	return kAPI.Update(ctx, filepath.Join(ResourcePrefix, resourceType, resourceName), string(dataBytes))
}

func Get(resourceType string, filter string) (interface{}, error) {
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


	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	resp, err := kAPI.Get(ctx, filepath.Join(ResourcePrefix, resourceType, filter), nil)
	if err != nil {
		log.Infof("error in finding data for key %+v: %+v \n", filepath.Join(ResourcePrefix, resourceType, filter), err)
		return nil, err
	} else {
		// print common key info
		log.Infof("Get is done. Metadata is %q\n", resp)
	}
	return resp.Node.Value, nil
}

func Delete(resourceType string, filter string) (interface{}, error) {
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


	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	resp, err := kAPI.Delete(ctx, filepath.Join(ResourcePrefix, resourceType, filter), nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		// print common key info
		log.Infof("Delete is done. Metadata is %q\n", resp)
	}
	return resp.Node.Value, nil
}

