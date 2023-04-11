package k8s_headless_svc

import (
	"fmt"
	"go-micro.dev/v4/registry"
)

// about services within the registry.
type k8sSvcWatcher struct {
}

func (k *k8sSvcWatcher) Next() (*registry.Result, error) {

	return &registry.Result{}, nil
}
func (k *k8sSvcWatcher) Stop() {

}

type Service struct {
	Namespace string //k8s namespace
	SvcName   string //service 的name
	PodPort   int32  //service 对应的endpoint 的port 也就是pod的container port
}
type k8sSvcRegister struct {
	k8sService []*Service //当前服务可能需要依赖多个其他服务

	opts *registry.Options
}

func (k *k8sSvcRegister) Init(opts ...registry.Option) error {
	for _, o := range opts {
		o(k.opts)
	}
	return nil
}
func (k *k8sSvcRegister) Options() registry.Options {
	return registry.Options{}
}
func (k *k8sSvcRegister) Register(*registry.Service, ...registry.RegisterOption) error {
	//解析dns 返回pod id
	// because we use k8s svc replace Register func ,so this func do nothing
	return nil
}
func (k *k8sSvcRegister) Deregister(*registry.Service, ...registry.DeregisterOption) error {
	//解析dns 返回pod id
	// because we use k8s svc replace Register func ,so this func do nothing

	return nil
}
func (k *k8sSvcRegister) GetService(string, ...registry.GetOption) ([]*registry.Service, error) {

	var service []*registry.Service
	var nodes []*registry.Node
	//
	ipMaps, err := getDnsForPodIP(k.k8sService)
	if err != nil {
		return []*registry.Service{}, err
	}
	for svcName, ips := range ipMaps {
		for _, ip := range ips {
			nodes = append(nodes, &registry.Node{Address: ip})

		}
		service = append(service, &registry.Service{Name: svcName, Version: "latest", Nodes: nodes})

	}

	//nodes = append(nodes, &registry.Node{Address: "127.0.0.1:8080"})
	//service = append(service, &registry.Service{Name: "user", Version: "latest", Nodes: nodes})
	fmt.Println("1111111")
	fmt.Println(service)
	for _, s := range service {
		fmt.Println(s.Name)
	}
	fmt.Println("1111111")

	return service, nil
}
func (k *k8sSvcRegister) ListServices(...registry.ListOption) ([]*registry.Service, error) {

	var service []*registry.Service
	var nodes []*registry.Node
	//
	ipMaps, err := getDnsForPodIP(k.k8sService)
	if err != nil {
		return []*registry.Service{}, err
	}
	for svcName, ips := range ipMaps {
		for _, ip := range ips {
			nodes = append(nodes, &registry.Node{Address: ip})

		}
		service = append(service, &registry.Service{Name: svcName, Version: "latest", Nodes: nodes})

	}
	fmt.Println("222222")

	fmt.Println(service)
	fmt.Println("2222")

	return service, nil
}

func (k *k8sSvcRegister) Watch(option ...registry.WatchOption) (registry.Watcher, error) {
	return &k8sSvcWatcher{}, nil
}
func (k *k8sSvcRegister) String() string {
	return "k8s-headless-svc"
}

// NewRegistry creates a kubernetes registry.
func NewRegistry(k8sService []*Service, opts ...registry.Option) registry.Registry {
	k := k8sSvcRegister{
		k8sService: k8sService,
		opts:       &registry.Options{},
	}
	return &k
}
