package selector

import (
	"testing"

	"github.com/nzgogo/micro/registry"
)

func TestStrategies(t *testing.T) {
	testData := []*registry.Service{
		&registry.Service{
			Name:    "test1",
			Version: "latest",
			Nodes: []*registry.Node{
				&registry.Node{
					ID: "test1-1",
				},
				&registry.Node{
					ID: "test1-2",
				},
			},
		},
		&registry.Service{
			Name:    "test1",
			Version: "default",
			Nodes: []*registry.Node{
				&registry.Node{
					ID: "test1-3",
				},
				&registry.Node{
					ID: "test1-4",
				},
			},
		},
	}

	for name, strategy := range map[string]Strategy{"random": Random, "roundrobin": RoundRobin} {
		next := strategy(testData)
		counts := make(map[string]int)

		for i := 0; i < 100; i++ {
			node, err := next()
			if err != nil {
				t.Fatal(err)
			}
			counts[node.ID]++
		}

		t.Logf("%s: %+v\n", name, counts)
	}
}
