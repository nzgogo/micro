package registry

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"testing"
	//"os/exec"
	consul "github.com/hashicorp/consul/api"
	//"strings"
	//"log"
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

//func runConsulAgent(){
//	cmd := exec.Command("./run_consul.sh")
//	cmd.Stdin = strings.NewReader("some input")
//	var out bytes.Buffer
//	cmd.Stdout = &out
//	err := cmd.Run()
//	if err != nil {
//		log.Fatal(err)
//	}
//}

func newNode(id, addr string, port int) *Node {
	return &Node{id,addr,port}
}

func TestRegister(t *testing.T) {
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

	go newMockServer(r, l1)
	go newMockServer(r, l2)
	go newMockServer(r, l3)


	client := NewRegistry()
	err = client.Init()
	if err != nil {
		t.Fatalf("NewRegistry faild. error: %v", err)
	}

	err = client.Register(&Service{
		Name: "order",
		Version : "v-0.1",
		Nodes: []*Node{
			newNode("123","localhost", 50000),
		},
	})
	if err != nil {
		t.Fatalf("Register faild. error: %v", err)
	}

	err = client.Register(&Service{
		Name: "order",
		Version : "v-0.1",
		Nodes: []*Node{
			newNode("456","localhost", 50001),
		},
	})
	if err != nil {
		t.Fatalf("Register faild. error: %v", err)
	}

	err = client.Register(&Service{
		Name: "porn",
		Version : "v-0.1",
		Nodes: []*Node{
			newNode("789","localhost", 50002),
		},
	})
	if err != nil {
		t.Fatalf("Register faild. error: %v", err)
	}



}

//func TestGetService(t *testing.T) {
//	cr, cl := newTestRegistry(&mockRegistry{
//		err: errors.New("client-error"),
//		url: "/v1/health/service/service-name",
//	})
//	defer cl()
//
//	if _, err := cr.GetService("test-service"); err == nil {
//		t.Fatalf("Expected error not to be `nil`")
//	} else {
//		t.Log("%v",err)
//	}
//}