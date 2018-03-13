package constant

import (
	"errors"
)

// errors
var (
	// 4xx
	ErrHttpEmptyRequest  = errors.New("nats proxy: Request cannot be nil")
	ErrRouterInvalidPath = errors.New("invalid path cannot process")
	ErrResourceNotFound  = errors.New("resource not found")
	ErrMethodNotAllowed  = errors.New("method not allowed")

	// 5xx
	ErrRegistryEmptyNode   = errors.New("registry: require at least one node")
	ErrSelectNoRegistry    = errors.New("selector: registry can not be empty")
	ErrSelectNotFound      = errors.New("selector: service not found")
	ErrSelectNoneAvailable = errors.New("selector: none available")
	ErrEmptyMsg            = errors.New("message cannot be nil")
)

// configurations
const (
	GOGO_CONFIG_PATH = "/etc/gogo/"
	CONFIG_FILE_NAME = ".config.json"

	ORGANIZATION = "gogo"

	// Message types
	REQUEST     = "request"
	RESPONSE    = "response"
	HEALTHCHECK = "healthCheck"

	// Service configs
	CONFIG_NATS_ADDRESS                         = "nats_addr"
	CONFIG_CONSUL_ADDRRESS                      = "consul_addr"
	CONFIG_HC_SCRIPT                            = "hc_script"
	CONFIG_HC_INTERVAL                          = "hc_interval"
	CONFIG_HC_DEREGISTER_CRITICAL_SERVICE_AFTER = "hc_deregister_critical_service_after"
	CONFIG_HC_LOAD_CRITICAL_THRESHOLD           = "hc_load_critical_threshold"
	CONFIG_HC_LOAD_WARNING_THRESHOLD            = "hc_load_warning_threshold"
	CONFIG_HC_MEMORY_CRITICAL_THRESHOLD         = "hc_memory_critical_threshold"
	CONFIG_HC_MEMORY_WARNING_THRESHOLD          = "hc_memory_warning_threshold"
	CONFIG_HC_CPU_CRITICAL_THRESHOLD            = "hc_cpu_critical_threshold"
	CONFIG_HC_CPU_WARNING_THRESHOLD             = "hc_cpu_warning_threshold"

	// health check script usage
	HC_SCRIPT_ARGS = "-subj="

	// Default value for health checks configs
	DEFAULT_HC_SCRITP                           = "gghc"
	DEFAULT_HC_INTERVAL                         = "1m"
	DEFALT_HC_DEREGISTER_CRITICAL_SERVICE_AFTER = "5m"
	DEFALT_HC_LOAD_CRITICAL_THRESHOLD           = "0.9"
	DEFALT_HC_LOAD_WARNING_THRESHOLD            = "0.8"
	DEFALT_HC_MEMORY_CRITICAL_THRESHOLD         = "5"
	DEFALT_HC_MEMORY_WARNING_THRESHOLD          = "15"
	DEFALT_HC_CPU_CRITICAL_THRESHOLD            = "5"
	DEFALT_HC_CPU_WARNING_THRESHOLD             = "15"
)

// health check
const (
	MEMORY_WARNING  = "MW"
	MEMORY_CRITICAL = "MC"
	LOAD_WARNING    = "LW"
	LOAD_CRITICAL   = "LC"

	OK       = 0
	Warning  = 1
	Critical = 2
)
