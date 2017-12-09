// NatsProxy serves as a proxy between gnats and http.
// It automatically translates the HTTP requests to nats messages.
// The url and method of the HTTP request serves as the name of the nats channel, where the message is sent.

package natsproxy

import (
	"bytes"
	"net/http"
	"regexp"
	"micro/transport"
	"strconv"
)

// HookFunc is the function that is
// used to modify response just before its
// transformed to HTTP response
type HookFunc func(Response)

type NatsProxy struct {
	//conn  *nats.Conn
	client *transport.Client
	hooks map[string]hookGroup
}

type hookGroup struct {
	regexp *regexp.Regexp
	hooks  []HookFunc
}

// NewNatsProxy creates an
// initialized NatsProxy
func NewNatsProxy(client *transport.Client) (*NatsProxy, error) {
	if err := client.TestConnection(); err != nil {
		return nil, err
	}
	return &NatsProxy{
		client,
		make(map[string]hookGroup, 0),
	}, nil
}

func (np *NatsProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	request, err := NewRequestFromHTTP(req) 	// Transform the HTTP request to NATS proxy request.
	if err != nil {
		http.Error(rw, "Cannot process request", http.StatusInternalServerError)
		return
	}

	var response Response
	// Post request to message queue
	respErr := np.client.Request(request, response)
	if respErr != nil {
		http.Error(rw, "No response", http.StatusInternalServerError)
		return
	}

	// Apply hook if regex match
	for _, hG := range np.hooks {
		if hG.regexp.MatchString(req.URL.Path) {
			for _, hook := range hG.hooks {
				hook(response)
			}
		}
	}
	writeResponse(rw, response)
}

// AddHook add the hook to modify,
// process response just before
// its transformed to HTTP form.
func (np *NatsProxy) AddHook(urlRegex string, hook HookFunc) error {
	hG, ok := np.hooks[urlRegex]
	if !ok {
		regexp, err := regexp.Compile(urlRegex)
		if err != nil {
			return err
		}
		hooks := make([]HookFunc, 1)
		hooks[0] = hook
		np.hooks[urlRegex] = hookGroup{
			regexp,
			hooks,
		}
	} else {
		hG.hooks = append(hG.hooks, hook)
	}
	return nil
}

func writeResponse(rw http.ResponseWriter, response Response) {
	// TODO Header
	statusCode, _ := strconv.Atoi(response.Header["StatusCode"])
	// Write the response code
	rw.WriteHeader(statusCode)

	// Write the bytes of response to a response writer.
	// TODO benchmark
	bytes.NewBuffer(response.Body).WriteTo(rw)
}

