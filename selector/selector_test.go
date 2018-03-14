package selector

import (
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/nzgogo/micro/registry"
)

type mockRegistry struct {
	body   []byte
	status int
	err    error
	url    string
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

func newNode(id string) *registry.Node {
	return &registry.Node{id}
}

func runConsulAgent() {
	cmd := exec.Command("consul", "agent", "-dev")
	err := cmd.Start()
	if err != nil {
		panic(err.Error())
	}
}

func stopConsulAgent() {
	cmd := "echo $(ps cax | grep consul | grep -o '^[ ]*[0-9]*')"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		panic(err.Error())
	}
	pid, err := strconv.Atoi(string(out[:len(out)-1]))
	if err != nil {
		panic(err.Error())
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		panic(err.Error())
	}
	p.Kill()
	if err != nil {
		panic(err.Error())
	}
}

func TestSelector(t *testing.T) {
	//setup a consul agent
	{
		runConsulAgent()
	}
	time.Sleep(300 * time.Millisecond)

	//Setup three services
	{
		l1, err := net.Listen("tcp", "localhost:50000")
		if err != nil {
			panic(err.Error())
		}

		l2, err := net.Listen("tcp", "localhost:50001")
		if err != nil {
			panic(err.Error())
		}

		l3, err := net.Listen("tcp", "localhost:50002")
		if err != nil {
			panic(err.Error())
		}

		r := &mockRegistry{
			status: 200,
			body:   []byte("Fuck off. I am retired, don't ask me to do a damn thing!"),
			url:    "/v1/health/service/service-name",
		}

		//run three servers
		go newMockServer(r, l1)
		go newMockServer(r, l2)
		go newMockServer(r, l3)
	}

	client := registry.NewRegistry()

	//register three services
	{
		err := client.Init()
		if err != nil {
			t.Fatalf("NewRegistry faild. error: %v", err)
		}

		err = client.Register(&registry.Service{
			Name:    "gogo.core.api",
			Version: "v-0.1",
			Nodes: []*registry.Node{
				newNode("123"),
			},
		})
		if err != nil {
			t.Fatalf("Register gogo.core.api-123 faild. error: %v", err)
		}

		err = client.Register(&registry.Service{
			Name:    "gogo.core.api",
			Version: "v-0.2",
			Nodes: []*registry.Node{
				newNode("456"),
			},
		})
		if err != nil {
			t.Fatalf("Register gogo.core.api-456 faild. error: %v", err)
		}

		err = client.Register(&registry.Service{
			Name:    "porn",
			Version: "v-0.1",
			Nodes: []*registry.Node{
				newNode("789"),
			},
		})
		if err != nil {
			t.Fatalf("Register porn service faild. error: %v", err)
		}
	}

	//selector test
	{
		slt := NewSelector(client, SetStrategy(RoundRobin))
		err := slt.Init()
		if err != nil {
			t.Fatalf("NewSelector init failed. error: %v", err)
		}
		subject, err := slt.Select("gogo.core.api", "v-0.1")
		if err != nil {
			t.Fatalf("Selector failed. error: %v", err)
		}
		t.Logf("Success! Service subject: %s:\n", subject)
	}
	//kill the consul agent process
	{
		stopConsulAgent()
	}
}
