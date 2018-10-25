package resource

func init() {
	ResourcesMap = map[string]Resource{}
}

var ResourcesMap map[string]Resource

type Resource interface {
	Create(interface{}) (interface{}, error)
	List(*ListOptions) ([]interface{}, error)
	Update(interface{}, interface{}) (interface{}, error)
	Delete(interface{}) (interface{}, error)
	Get(string, *GetOptions) (interface{}, error)
}

type GetOptions map[string]string

type ListOptions struct {
	Filter map[string]string
}
