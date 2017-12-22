package registry
//This is a test file that includes all function tests covered in registry.go.
//To get this Test working, you need to run a consul agent in dev mode first
//Just run command line "consul agent -dev" before you go test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"testing"
	consul "github.com/hashicorp/consul/api"
)

func newHealthCheck(node, name, status string) *consul.HealthCheck {
	return &consul.HealthCheck{
		Node:        node,
		Name:        name,
		Status:      status,
		ServiceName: name,
	}
}

func newServiceEntry(node, address, name, version string, checks []*consul.HealthCheck) *consul.ServiceEntry {
	return &consul.ServiceEntry{
		Node: &consul.Node{Node: node, Address: name},
		Service: &consul.AgentService{
			Service: name,
			Address: address,
			Tags:    []string{version},
		},
		Checks: checks,
	}
}

type mockRegistry struct {
	body   []byte
	status int
	err    error
	url    string
}

func encodeData(obj interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func newMockServer(rg *mockRegistry, l net.Listener) error {
	mux := http.NewServeMux()
	mux.HandleFunc(rg.url, func(w http.ResponseWriter, r *http.Request) {
		if rg.err != nil {
			http.Error(w, rg.err.Error(), 500)
			return
		}
		w.WriteHeader(rg.status)
		w.Write(rg.body)
	})
	return http.Serve(l, mux)
}

func newNode(id string) *Node {
	return &Node{id}
}

//This test covers Init, Register, Deregister, GetService, ListServices. This test needs a running consul agent to work
func TestRegistry(t *testing.T) {
	{
		//go runConsulAgent()
		l1, err := net.Listen("tcp", "localhost:50000")
		if err != nil {
			// blurgh?!!
			panic(err.Error())
		}

		l2, err := net.Listen("tcp", "localhost:50001")
		if err != nil {
			// blurgh?!!
			panic(err.Error())
		}

		l3, err := net.Listen("tcp", "localhost:50002")
		if err != nil {
			// blurgh?!!
			panic(err.Error())
		}

		r := &mockRegistry{
			status: 200,
			body: []byte("Fuck off. I am retired, don't ask me to do a damn thing!"),
			url: "/v1/health/service/service-name",
		}

		//run three servers
		go newMockServer(r, l1)
		go newMockServer(r, l2)
		go newMockServer(r, l3)
	}


	client := NewRegistry()

	//register three servers
	{
		err := client.Init()
		if err != nil {
			t.Fatalf("NewRegistry faild. error: %v", err)
		}

		err = client.Register(&Service{
			Name: "order",
			Version : "v-0.1",
			Nodes: []*Node{
				newNode("123"),
			},
		})
		if err != nil {
			t.Fatalf("Register order-123 faild. error: %v", err)
		}

		err = client.Register(&Service{
			Name: "order",
			Version : "v-0.1",
			Nodes: []*Node{
				newNode("456"),
			},
		})
		if err != nil {
			t.Fatalf("Register order-456 faild. error: %v", err)
		}

		err = client.Register(&Service{
			Name: "porn",
			Version : "v-0.1",
			Nodes: []*Node{
				newNode("789"),
			},
		})
		if err != nil {
			t.Fatalf("Register porn-789 faild. error: %v", err)
		}
	}

	//list all registered services
	{
		services, err := client.ListServices()
		if err != nil {
			t.Fatalf("ListServices faild. error: %v", err)
		}

		pcnt := false
		ocnt := false
		for _,service := range services {
			if service.Name == "porn"{
				pcnt = true
			}
			if service.Name == "order"{
				ocnt = true
			}

		}
		if !pcnt {
			t.Fatalf("could not find porn service")
		}
		if !ocnt {
			t.Fatalf("could not find order service")
		}
	}

	//deregister service porn test
	{
		err := client.Deregister(&Service{
			Nodes: []*Node{
				newNode("789"),
			},
		})
		if err != nil {
			t.Fatalf("Deregister porn-789 faild. error: %v", err)
		}

		services, err := client.ListServices()
		if err != nil {
			t.Fatalf("ListServices in Deregister faild. error: %v", err)
		}

		pcnt := false
		for _,service := range services {
			if service.Name == "porn"{
				pcnt = true
			}
		}
		if pcnt {
			t.Fatalf("Expected service porn-789 deregistered")
		}
	}

	//get service test
	{
		services, err := client.GetService("order")
		if err != nil {
			t.Fatalf("GetService faild. error: %v", err)
		}
		cnt := 0
		for _,service := range services {
			if service.Name == "order"{
				cnt = len(service.Nodes)
			}
		}
		if cnt != 2{
			t.Fatalf("GetService faild. Cound not find order service")
		}

	}
}

//These test functions below are unit test that do not dependant on a running consul agent
func newConsulTestRegistry(r *mockRegistry) (*registry, func()) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		// blurgh?!!
		panic(err.Error())
	}
	cfg := consul.DefaultConfig()
	cfg.Address = l.Addr().String()
	cl, _ := consul.NewClient(cfg)

	go newMockServer(r, l)

	return &registry{
		Client:   cl,
		register: make(map[string]uint64),
	}, func() {
		l.Close()
	}
}

func newServiceList(svc []*consul.ServiceEntry) []byte {
	bts, _ := encodeData(svc)
	return bts
}

func TestConsul_GetService_WithError(t *testing.T) {
	cr, cl := newConsulTestRegistry(&mockRegistry{
		err: errors.New("client-error"),
		url: "/v1/health/service/service-name",
	})
	defer cl()

	if _, err := cr.GetService("test-service"); err == nil {
		t.Fatalf("Expected error not to be `nil`")
	} else {
		t.Log("%v",err)
	}
}

func TestConsul_GetService_WithHealthyServiceNodes(t *testing.T) {
	// warning is still seen as healthy, critical is not
	svcs := []*consul.ServiceEntry{
		newServiceEntry(
			"node-name-1", "node-address-1", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-1", "service-name", "passing"),
				newHealthCheck("node-name-1", "service-name", "warning"),
			},
		),
		newServiceEntry(
			"node-name-2", "node-address-2", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-2", "service-name", "passing"),
				newHealthCheck("node-name-2", "service-name", "warning"),
			},
		),
	}

	cr, cl := newConsulTestRegistry(&mockRegistry{
		status: 200,
		body:   newServiceList(svcs),
		url:    "/v1/health/service/service-name",
	})
	defer cl()

	svc, err := cr.GetService("service-name")
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	if exp, act := 1, len(svc); exp != act {
		t.Fatalf("Expected len of svc to be `%d`, got `%d`.", exp, act)
	}

	if exp, act := 2, len(svc[0].Nodes); exp != act {
		t.Fatalf("Expected len of nodes to be `%d`, got `%d`.", exp, act)
	}
}

func TestConsul_GetService_WithUnhealthyServiceNode(t *testing.T) {
	// warning is still seen as healthy, critical is not
	svcs := []*consul.ServiceEntry{
		newServiceEntry(
			"node-name-1", "node-address-1", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-1", "service-name", "passing"),
				newHealthCheck("node-name-1", "service-name", "warning"),
			},
		),
		newServiceEntry(
			"node-name-2", "node-address-2", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-2", "service-name", "passing"),
				newHealthCheck("node-name-2", "service-name", "critical"),
			},
		),
	}

	cr, cl := newConsulTestRegistry(&mockRegistry{
		status: 200,
		body:   newServiceList(svcs),
		url:    "/v1/health/service/service-name",
	})
	defer cl()

	svc, err := cr.GetService("service-name")
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	if exp, act := 1, len(svc); exp != act {
		t.Fatalf("Expected len of svc to be `%d`, got `%d`.", exp, act)
	}

	if exp, act := 1, len(svc[0].Nodes); exp != act {
		t.Fatalf("Expected len of nodes to be `%d`, got `%d`.", exp, act)
	}
}

func TestConsul_GetService_WithUnhealthyServiceNodes(t *testing.T) {
	// warning is still seen as healthy, critical is not
	svcs := []*consul.ServiceEntry{
		newServiceEntry(
			"node-name-1", "node-address-1", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-1", "service-name", "passing"),
				newHealthCheck("node-name-1", "service-name", "critical"),
			},
		),
		newServiceEntry(
			"node-name-2", "node-address-2", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-2", "service-name", "passing"),
				newHealthCheck("node-name-2", "service-name", "critical"),
			},
		),
	}

	cr, cl := newConsulTestRegistry(&mockRegistry{
		status: 200,
		body:   newServiceList(svcs),
		url:    "/v1/health/service/service-name",
	})
	defer cl()

	svc, err := cr.GetService("service-name")
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	if exp, act := 1, len(svc); exp != act {
		t.Fatalf("Expected len of svc to be `%d`, got `%d`.", exp, act)
	}

	if exp, act := 0, len(svc[0].Nodes); exp != act {
		t.Fatalf("Expected len of nodes to be `%d`, got `%d`.", exp, act)
	}
}