package gogo

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

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

// URLToIntnlTrans builds the channel name for a internal transport use from an URL TODO regexp
func URLToIntnlTrans(host string, path string) string {
	str := strings.Split(path, "/")
	return "gogo-" + str[2] + "-" + str[3]
}

func readConfigFile(srvName string) map[string]string {
	filename := srvName + ".config.json"
	currentFolder := "./"
	etcFolder := "/etc/gogo/"

	if _, err := os.Stat(currentFolder + filename); os.IsNotExist(err) {
		if _, err := os.Stat(etcFolder + filename); os.IsNotExist(err) {
			return make(map[string]string)
		} else {
			filename = etcFolder + filename
		}
	} else {
		filename = etcFolder + filename
	}

	fileBytes, _ := ioutil.ReadFile(filename)
	configMap := make(map[string]string)
	err := codec.Unmarshal(fileBytes, &configMap)

	if err != nil {
		return make(map[string]string)
	}

	return configMap
}

func healthCheck(configs map[string]string) (int, []byte) {
	var retMsg = ""
	var status = OK

	//check cpu
	cpuCriticalThreshold, err1 := strconv.ParseFloat(configs["hc_cpu_critical_threshold"], 64)
	cpuWarningThreshold, err2 := strconv.ParseFloat(configs["hc_cpu_warning_threshold"], 64)
	if err1 == nil || err2 == nil {
		cpuPercent, err := cpu.Percent(0, false)
		if err != nil {
			log.Println("Failed get CPU information.")
			retMsg += " Failed get CPU information. "
			status |= Warning
		} else {
			isCritical := false
			cpstr := strconv.FormatFloat(cpuPercent[0], 'f', 2, 64)
			if err1 == nil {
				if 100-cpuPercent[0] < cpuCriticalThreshold {
					msg := " CPU is critical. Percentage of CPU used: " + cpstr + "%"
					log.Println(msg)
					retMsg += msg
					status |= Critical
					isCritical = true
				}
			}
			if err2 == nil {
				if 100-cpuPercent[0] < cpuWarningThreshold {
					msg := " CPU is warning. Percentage of CPU used: " + cpstr + "%"
					if !isCritical {
						log.Println(msg)
						retMsg += msg
					}
					status |= Warning
				}
			}
		}
	}

	//check memory usage
	memoryCriticalThreshold, err1 := strconv.ParseFloat(configs["hc_memory_critical_threshold"], 64)
	memoryWarningThreshold, err2 := strconv.ParseFloat(configs["hc_memory_warning_threshold"], 64)
	if err1 == nil || err2 == nil {
		v, err := mem.VirtualMemory()
		if err != nil {
			log.Println("Failed get memory information.")
			retMsg += " Failed get CPU information. "
			status |= Warning
		} else {
			memoryPercent := v.UsedPercent
			mpstr := strconv.FormatFloat(memoryPercent, 'f', 2, 64)
			isCritical := false
			if err1 == nil {
				if 100-memoryPercent < memoryCriticalThreshold {
					msg := " Memory is critical. Percentage of Memory used: " + mpstr + "%"
					log.Println(msg)
					retMsg += msg
					status |= Critical
					isCritical = true
				}
			}
			if err2 == nil {
				if 100-memoryPercent < memoryWarningThreshold {
					msg := " Memory is warning. Percentage of Memory used: " + mpstr + "%"
					if !isCritical {
						log.Println(msg)
						retMsg += msg
					}
					status |= Warning

				}
			}
		}
	}

	//check load
	loadCriticalThreshold, err1 := strconv.ParseFloat(configs["hc_load_critical_threshold"], 64)
	loadWarningThreshold, err2 := strconv.ParseFloat(configs["hc_load_warning_threshold"], 64)
	if err1 == nil || err2 == nil {
		l, err := load.Avg()
		if err != nil {
			log.Println("Failed get load information.")
			retMsg += " Failed get load information. "
			status |= Warning
		} else {
			load := l.Load5
			lstr := strconv.FormatFloat(load, 'f', 2, 64)
			isCritical := false
			if err1 == nil {
				if load > loadCriticalThreshold {
					msg := " Overload critical. System loads: " + lstr
					log.Println(msg)
					retMsg += msg
					status |= Critical
					isCritical = true
				}
			}
			if err2 == nil {
				if load > loadWarningThreshold {
					msg := " Overload warning. System loads: " + lstr
					if !isCritical {
						log.Println(msg)
						retMsg += msg
					}
					status |= Warning
				}
			}
		}

	}

	return status, []byte(retMsg)
}
