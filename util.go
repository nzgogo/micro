package gogo

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/constant"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

// URLToServiceName builds the service name for a internal transport use from an URL
func URLToServiceName(host string, path string) string {
	str := strings.Split(path, "/")
	if len(str) < 4 {
		return ""
	}
	return constant.ORGANIZATION + "-" + str[2] + "-" + str[3]
}

// URLToServiceVersion builds the service version for a internal transport use from an URL
func URLToServiceVersion(path string) string {
	str := strings.Split(path, "/")
	if len(str) < 4 {
		return ""
	}
	return str[1]
}

func readConfigFile(srvName string) map[string]string {
	filename := srvName + constant.CONFIG_FILE_NAME
	currentFolder := "./"
	etcFolder := constant.GOGO_CONFIG_PATH

	if _, err := os.Stat(currentFolder + filename); os.IsNotExist(err) {
		if _, err := os.Stat(etcFolder + filename); os.IsNotExist(err) {
			return make(map[string]string)
		} else {
			filename = etcFolder + filename
		}
	} else {
		filename = currentFolder + filename
	}

	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return make(map[string]string)
	}
	configMap := make(map[string]interface{})
	err = codec.Unmarshal(fileBytes, &configMap)
	if err != nil {
		return make(map[string]string)
	}
	return parseConfigFile(configMap)
}

func parseConfigFile(raw map[string]interface{}) map[string]string {
	configs := make(map[string]string)

	for k, v := range raw {
		switch vv := v.(type) {
		case string:
			configs[k] = vv
		case []interface{}:
			configs[k] = selectOneOption(vv)
		}
	}

	return configs
}

func selectOneOption(options []interface{}) string {
	rand.Seed(time.Now().UnixNano())

	if v, ok := options[rand.Intn(len(options))].(string); ok {
		return v
	}
	return ""
}

func packHealthCheck(config map[string]string, srvSubject string) *consul.AgentServiceCheck {
	if config[constant.CONFIG_HC_SCRIPT] == "" {
		config[constant.CONFIG_HC_SCRIPT] = constant.DEFAULT_HC_SCRITP
	}
	if config[constant.CONFIG_HC_INTERVAL] == "" {
		config[constant.CONFIG_HC_INTERVAL] = constant.DEFAULT_HC_INTERVAL
	}
	if config[constant.CONFIG_HC_DEREGISTER_CRITICAL_SERVICE_AFTER] == "" {
		config[constant.CONFIG_HC_DEREGISTER_CRITICAL_SERVICE_AFTER] = constant.DEFALT_HC_DEREGISTER_CRITICAL_SERVICE_AFTER
	}
	if config[constant.CONFIG_HC_LOAD_CRITICAL_THRESHOLD] == "" {
		config[constant.CONFIG_HC_LOAD_CRITICAL_THRESHOLD] = constant.DEFALT_HC_LOAD_CRITICAL_THRESHOLD
	}
	if config[constant.CONFIG_HC_LOAD_WARNING_THRESHOLD] == "" {
		config[constant.CONFIG_HC_LOAD_WARNING_THRESHOLD] = constant.DEFALT_HC_LOAD_WARNING_THRESHOLD
	}
	if config[constant.CONFIG_HC_MEMORY_CRITICAL_THRESHOLD] == "" {
		config[constant.CONFIG_HC_MEMORY_CRITICAL_THRESHOLD] = constant.DEFALT_HC_MEMORY_CRITICAL_THRESHOLD
	}
	if config[constant.CONFIG_HC_MEMORY_WARNING_THRESHOLD] == "" {
		config[constant.CONFIG_HC_MEMORY_WARNING_THRESHOLD] = constant.DEFALT_HC_MEMORY_WARNING_THRESHOLD
	}
	if config[constant.CONFIG_HC_CPU_CRITICAL_THRESHOLD] == "" {
		config[constant.CONFIG_HC_CPU_CRITICAL_THRESHOLD] = constant.DEFALT_HC_CPU_CRITICAL_THRESHOLD
	}
	if config[constant.CONFIG_HC_CPU_WARNING_THRESHOLD] == "" {
		config[constant.CONFIG_HC_CPU_WARNING_THRESHOLD] = constant.DEFALT_HC_CPU_WARNING_THRESHOLD
	}

	arg := constant.HC_SCRIPT_ARGS + srvSubject

	return &consul.AgentServiceCheck{
		Args:                           []string{config[constant.CONFIG_HC_SCRIPT], arg},
		Interval:                       config[constant.CONFIG_HC_INTERVAL],
		DeregisterCriticalServiceAfter: config[constant.CONFIG_HC_DEREGISTER_CRITICAL_SERVICE_AFTER],
	}
}

func healthCheck(configs map[string]string) (int, map[string]string) {
	if configs == nil {
		log.Println("Missing service configurations")
		return constant.Warning, nil
	}
	var status = constant.OK
	feedback := make(map[string]string)
	//check cpu
	//cpuCriticalThreshold, err1 := strconv.ParseFloat(configs[CONFIG_HC_CPU_CRITICAL_THRESHOLD], 64)
	//cpuWarningThreshold, err2 := strconv.ParseFloat(configs[CONFIG_HC_CPU_WARNING_THRESHOLD], 64)
	//if err1 == nil || err2 == nil {
	//	cpuPercent, err := cpu.Percent(0, false)
	//	if err != nil {
	//		log.Println("Failed to get CPU information.")
	//		retMsg += " Failed get CPU information. "
	//		status |= Warning
	//	} else {
	//		isCritical := false
	//		cpstr := strconv.FormatFloat(cpuPercent[0], 'f', 2, 64)
	//		if err1 == nil {
	//			if 100-cpuPercent[0] < cpuCriticalThreshold {
	//				msg := " CPU is critical. Percentage of CPU used: " + cpstr + "%"
	//				log.Println(msg)
	//				retMsg += msg
	//				status |= Critical
	//				isCritical = true
	//			}
	//		}
	//		if err2 == nil {
	//			if 100-cpuPercent[0] < cpuWarningThreshold {
	//				msg := " CPU is warning. Percentage of CPU used: " + cpstr + "%"
	//				if !isCritical {
	//					log.Println(msg)
	//					retMsg += msg
	//				}
	//				status |= Warning
	//			}
	//		}
	//	}
	//}

	//check memory usage
	memoryCriticalThreshold, err1 := strconv.ParseFloat(configs[constant.CONFIG_HC_MEMORY_CRITICAL_THRESHOLD], 64)
	memoryWarningThreshold, err2 := strconv.ParseFloat(configs[constant.CONFIG_HC_MEMORY_WARNING_THRESHOLD], 64)
	if err1 == nil || err2 == nil {
		v, err := mem.VirtualMemory()
		if err != nil {
			log.Println("Failed to get memory information.")
			feedback[constant.MEMORY_WARNING] = "Failed to get CPU information. "
			status |= constant.Warning
		} else {
			memoryPercent := v.UsedPercent
			mpstr := strconv.FormatFloat(memoryPercent, 'f', 2, 64)
			isCritical := false
			if err1 == nil {
				if 100-memoryPercent < memoryCriticalThreshold {
					msg := "Memory is critical. Percentage of Memory used: " + mpstr + "%"
					log.Println(msg)
					feedback[constant.MEMORY_CRITICAL] = msg
					status |= constant.Critical
					isCritical = true
				}
			}
			if err2 == nil {
				if 100-memoryPercent < memoryWarningThreshold {
					msg := "Memory is warning. Percentage of Memory used: " + mpstr + "%"
					if !isCritical {
						log.Println(msg)
						feedback[constant.MEMORY_WARNING] = msg
					}
					status |= constant.Warning

				}
			}
		}
	}

	//check load
	loadCriticalThreshold, err1 := strconv.ParseFloat(configs[constant.CONFIG_HC_LOAD_CRITICAL_THRESHOLD], 64)
	loadWarningThreshold, err2 := strconv.ParseFloat(configs[constant.CONFIG_HC_LOAD_WARNING_THRESHOLD], 64)
	if err1 == nil || err2 == nil {
		l, err := load.Avg()
		if err != nil {
			log.Println("Failed to get load information.")
			feedback[constant.LOAD_WARNING] = "Failed to get load information. "
			status |= constant.Warning
		} else {
			cpuInfo, err := cpu.Info()
			if err != nil {
				log.Println("Failed to get CPU information.")
				feedback[constant.LOAD_WARNING] = "Failed get CPU information. "
				status |= constant.Warning
			} else {
				coreCnt := int32(0)
				for _, p := range cpuInfo {
					coreCnt += p.Cores
				}

				load := l.Load5 / float64(coreCnt)
				lstr := strconv.FormatFloat(load, 'f', 2, 64)
				isCritical := false
				if err1 == nil {
					if load > loadCriticalThreshold {
						msg := "Overload critical. System loads: " + lstr
						log.Println(msg)
						feedback[constant.LOAD_CRITICAL] = msg
						status |= constant.Critical
						isCritical = true
					}
				}
				if err2 == nil {
					if load > loadWarningThreshold {
						msg := "Overload warning. System loads: " + lstr
						if !isCritical {
							log.Println(msg)
							feedback[constant.LOAD_WARNING] = msg
						}
						status |= constant.Warning
					}
				}
			}
		}
	}
	if status > constant.Critical {
		status = constant.Critical
	}
	return status, feedback
}
