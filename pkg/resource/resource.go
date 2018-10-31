package resource

import "time"

func init() {
	ResourcesMap = map[string]Resource{}
}

var ResourcesMap map[string]Resource

type BaseResource struct {
	SelfLink string `json:"self"`
	Version uint64
	UUID string
	CreationTimestamp time.Time
	DeletionTimestamp time.Time
	LastUpdateTimestamp time.Time

	Data *struct{}
}

type Resource interface {
	Create(interface{}) (interface{}, error)
	List(*ListOptions) ([]interface{}, error)
	Update(interface{}) (interface{}, error)
	Delete(string, *DeleteOptions) (interface{}, error)
	Get(string, *GetOptions) (interface{}, error)

	GetBaseResource(interface{}) *BaseResource
}

type GetOptions map[string]string
type DeleteOptions map[string]string

type ListOptions struct {
	Filter map[string]string
}

