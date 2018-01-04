package gogo
import (
	"golang.org/x/net/context"
	"micro/codec"
	"net/http"
)

// HttpHandlerFunc represents a single method of a http handler. It's used primarily
// for the wrappers.
type HttpHandlerFunc func(w http.ResponseWriter, r *http.Request) error

// HttpHandlerWrapper wraps the HttpHandlerFunc and returns the equivalent
type HttpHandlerWrapper func(HttpHandlerFunc) HttpHandlerFunc

// HandlerFunc represents a single method of a service handler. It's used primarily
// for the wrappers (after api interpreter and before service handler).
type HandlerFunc func(ctx context.Context, req codec.Request) error

// HandlerWrapper wraps the HandlerFunc and returns the equivalent
type HandlerWrapper func(HandlerFunc) HandlerFunc
