package selector

import (
	"testing"
	"micro/registry"
	//"github.com/nzgogo/micro/registry"
)

func TestFilterVersion(t *testing.T) {
	testData := []struct {
		services []*registry.Service
		version  string
		count    int
	}{
		{
			services: []*registry.Service{
				&registry.Service{
					Name:    "test",
					Version: "1.0.0",
				},
				&registry.Service{
					Name:    "test",
					Version: "1.1.0",
				},
			},
			version: "1.0.0",
			count:   1,
		},
		{
			services: []*registry.Service{
				&registry.Service{
					Name:    "test",
					Version: "1.0.0",
				},
				&registry.Service{
					Name:    "test",
					Version: "1.1.0",
				},
			},
			version: "2.0.0",
			count:   0,
		},
	}

	for _, data := range testData {
		filter := FilterVersion(data.version)
		services := filter(data.services)

		if len(services) != data.count {
			t.Fatalf("Expected %d services, got %d", data.count, len(services))
		}

		var seen bool

		for _, service := range services {
			if service.Version != data.version {
				t.Fatalf("Expected version %s, got %s", data.version, service.Version)
			}
			seen = true
		}

		if seen == false && data.count > 0 {
			t.Fatalf("Expected %d services but seen is %t; result %+v", data.count, seen, services)
		}
	}
}