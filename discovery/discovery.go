package discovery

import (
	"fmt"
	"log"
	"os"
	consul "github.com/hashicorp/consul/api"
	health "github.com/nzgogo/micro/health"
)


// Discovery service
type Discovery struct {
	Client *consul.Client
	Config *consul.Config
	ServList []ServiceInfo
}

type ServiceInfo struct {
	ServiceID string
	IP        string
	Port      int
	Load      int
}
type ServiceList []ServiceInfo

func NewDiscovery() *Discovery{
	config := consul.DefaultConfig()

	// create the client
	client, _ := consul.NewClient(config)

	cr := &Discovery{
		Client: client,
		Config: config,
	}

	return cr
}

func CheckErr(err error) {
	if err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}

func (r *Discovery)GetKeyValue(service_name string, id string) []byte {
	key := service_name + "." + id  + ".health"

	kv, _, err := r.Client.KV().Get(key, nil)
	if kv == nil {
		return nil
	}
	CheckErr(err)

	return kv.Value
}

func (r *Discovery)GetLoads(service_name string){

	for _,service := range r.ServList {
		s := r.GetKeyValue(service_name, service.ServiceID)
		if s != nil {
			responseMsg := health.Decode(s)
			//if err == nil {
				service.Load = responseMsg.ServiceLoad
			//}
		}
		fmt.Println("service node updated ip:", service.IP, " port:", service.Port, " serviceid:", service.ServiceID, " load:", service.Load)

	}


}

func (r *Discovery)DiscoverServices(healthyOnly bool, service_name string) {

	services, _, err := r.Client.Catalog().Services(&consul.QueryOptions{})
	CheckErr(err)

	fmt.Println("--do discover ---:", r.Config.Address)

	for name := range services {
		servicesData, _, err := r.Client.Health().Service(name, "", healthyOnly,
			&consul.QueryOptions{})
		CheckErr(err)
		for _, entry := range servicesData {
			if service_name != entry.Service.Service {
				continue
			}
			for _, health := range entry.Checks {
				if health.ServiceName != service_name {
					continue
				}
				fmt.Println("  health nodeid:", health.Node, " service_name:", health.ServiceName, " service_id:", health.ServiceID, " status:", health.Status, " ip:", entry.Service.Address, " port:", entry.Service.Port)

				var node ServiceInfo
				node.IP = entry.Service.Address
				node.Port = entry.Service.Port
				node.ServiceID = health.ServiceID

				r.ServList = append(r.ServList, node)
			}
		}
	}
	return
}



