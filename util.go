package gogo

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
	"github.com/nzgogo/micro/codec"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

const (
	OK       = 0
	Warning  = 1
	Critical = 2
)

const (
	GOGO_CONFIG_PATH = "/etc/gogo/"
	CONFIG_FILE_NAME = ".config.json"

	MEMORY_WARNING = "MW"
	MEMORY_CRITICAL = "MC"
	LOAD_WARNING = "LW"
	LOAD_CRITICAL = "LC"
)

// URLToIntnlTrans builds the channel name for a internal transport use from an URL
func URLToIntnlTrans(host string, path string) string {
	str := strings.Split(path, "/")
	return ORGANIZATION + "-" + str[2] + "-" + str[3]
}

func readConfigFile(srvName string) map[string]string {
	filename := srvName + CONFIG_FILE_NAME
	currentFolder := "./"
	etcFolder := GOGO_CONFIG_PATH

	if _, err := os.Stat(currentFolder + filename); os.IsNotExist(err) {
		if _, err := os.Stat(etcFolder + filename); os.IsNotExist(err) {
			return make(map[string]string)
		} else {
			filename = etcFolder + filename
		}
	} else {
		filename = currentFolder + filename
	}

	fileBytes, _ := ioutil.ReadFile(filename)
	configMap := make(map[string]string)
	err := codec.Unmarshal(fileBytes, &configMap)

	if err != nil {
		return make(map[string]string)
	}

	return configMap
}

func packHealthCheck(config map[string]string, srvSubject string) *consul.AgentServiceCheck {
	if config[CONFIG_HC_SCRIPT] == "" {
		config[CONFIG_HC_SCRIPT] = DEFAULT_HC_SCRITP
	}
	if config[CONFIG_HC_INTERVAL] == "" {
		config[CONFIG_HC_INTERVAL] = DEFAULT_HC_INTERVAL
	}
	if config[CONFIG_HC_DEREGISTER_CRITICAL_SERVICE_AFTER] == "" {
		config[CONFIG_HC_DEREGISTER_CRITICAL_SERVICE_AFTER] = DEFALT_HC_DEREGISTER_CRITICAL_SERVICE_AFTER
	}
	if config[CONFIG_HC_LOAD_CRITICAL_THRESHOLD] == "" {
		config[CONFIG_HC_LOAD_CRITICAL_THRESHOLD] = DEFALT_HC_LOAD_CRITICAL_THRESHOLD
	}
	if config[CONFIG_HC_LOAD_WARNING_THRESHOLD] == "" {
		config[CONFIG_HC_LOAD_WARNING_THRESHOLD] = DEFALT_HC_LOAD_WARNING_THRESHOLD
	}
	if config[CONFIG_HC_MEMORY_CRITICAL_THRESHOLD] == "" {
		config[CONFIG_HC_MEMORY_CRITICAL_THRESHOLD] = DEFALT_HC_MEMORY_CRITICAL_THRESHOLD
	}
	if config[CONFIG_HC_MEMORY_WARNING_THRESHOLD] == "" {
		config[CONFIG_HC_MEMORY_WARNING_THRESHOLD] = DEFALT_HC_MEMORY_WARNING_THRESHOLD
	}
	if config[CONFIG_HC_CPU_CRITICAL_THRESHOLD] == "" {
		config[CONFIG_HC_CPU_CRITICAL_THRESHOLD] = DEFALT_HC_CPU_CRITICAL_THRESHOLD
	}
	if config[CONFIG_HC_CPU_WARNING_THRESHOLD] == "" {
		config[CONFIG_HC_CPU_WARNING_THRESHOLD] = DEFALT_HC_CPU_WARNING_THRESHOLD
	}

	arg := "-subj=" + srvSubject

	return &consul.AgentServiceCheck{
		//Notes: "health check",
		Args:                           []string{config[CONFIG_HC_SCRIPT], arg},
		Interval:                       config[CONFIG_HC_INTERVAL],
		DeregisterCriticalServiceAfter: config[CONFIG_HC_DEREGISTER_CRITICAL_SERVICE_AFTER],
	}
}

func healthCheck(configs map[string]string) (int, map[string]string) {
	var status = OK
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
	memoryCriticalThreshold, err1 := strconv.ParseFloat(configs[CONFIG_HC_MEMORY_CRITICAL_THRESHOLD], 64)
	memoryWarningThreshold, err2 := strconv.ParseFloat(configs[CONFIG_HC_MEMORY_WARNING_THRESHOLD], 64)
	if err1 == nil || err2 == nil {
		v, err := mem.VirtualMemory()
		if err != nil {
			log.Println("Failed to get memory information.")
			feedback[MEMORY_WARNING] = "Failed to get CPU information. "
			status |= Warning
		} else {
			memoryPercent := v.UsedPercent
			mpstr := strconv.FormatFloat(memoryPercent, 'f', 2, 64)
			isCritical := false
			if err1 == nil {
				if 100-memoryPercent < memoryCriticalThreshold {
					msg := "Memory is critical. Percentage of Memory used: " + mpstr + "%"
					log.Println(msg)
					feedback[MEMORY_CRITICAL] = msg
					status |= Critical
					isCritical = true
				}
			}
			if err2 == nil {
				if 100-memoryPercent < memoryWarningThreshold {
					msg := "Memory is warning. Percentage of Memory used: " + mpstr + "%"
					if !isCritical {
						log.Println(msg)
						feedback[MEMORY_WARNING] = msg
					}
					status |= Warning

				}
			}
		}
	}

	//check load
	loadCriticalThreshold, err1 := strconv.ParseFloat(configs[CONFIG_HC_LOAD_CRITICAL_THRESHOLD], 64)
	loadWarningThreshold, err2 := strconv.ParseFloat(configs[CONFIG_HC_LOAD_WARNING_THRESHOLD], 64)
	if err1 == nil || err2 == nil {
		l, err := load.Avg()
		if err != nil {
			log.Println("Failed to get load information.")
			feedback[LOAD_WARNING] =  "Failed to get load information. "
			status |= Warning
		} else {
			cpuInfo, err := cpu.Info()
			if err != nil {
				log.Println("Failed to get CPU information.")
				feedback[LOAD_WARNING] = "Failed get CPU information. "
				status |= Warning
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
						feedback[LOAD_CRITICAL] =  msg
						status |= Critical
						isCritical = true
					}
				}
				if err2 == nil {
					if load > loadWarningThreshold {
						msg := "Overload warning. System loads: " + lstr
						if !isCritical {
							log.Println(msg)
							feedback[LOAD_WARNING] =  msg
						}
						status |= Warning
					}
				}
			}
		}
	}
	if status > Critical {status = Critical}
	return status, feedback
}
