package gogo

import (
	"io/ioutil"
	"github.com/nzgogo/micro/codec"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"os"
	"strings"
	"strconv"
	"log"
)

var (
	OK = 0
	Warning = 1
	Critical = 2
)

// URLToIntnlTrans builds the channel name for a internal transport use from an URL TODO regexp
func URLToIntnlTrans(host string, path string) string {
	str := strings.Split(path, "/")
	return "gogo-" + str[2] + "-" + str[3]
}

func readConfigFile() map[string]string {
	filename := "./config.json"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return make(map[string]string)
	}

	fileBytes, _ := ioutil.ReadFile(filename)
	configMap := make(map[string]string)
	err := codec.Unmarshal(fileBytes, &configMap)

	if err != nil {
		return make(map[string]string)
	}

	return configMap
}

func healthCheck(configs map[string]string) int {
	//check cpu
	cpuPercent,err:=cpu.Percent(0,false)
	if err!=nil{
		log.Println("Failed get CPU information.")
		return Critical
	}
	cpuCriticalThreshold,_:=strconv.ParseFloat(configs["cpu_critical_threshold"], 64)
	cpuWarningThreshold,_:=strconv.ParseFloat(configs["cpu_warning_threshold"], 64)
	if 100-cpuPercent[0] < cpuCriticalThreshold {
		log.Println("CPU is Critical: ")
		if c,err:=cpu.Percent(0,true); err==nil{
			log.Println(c)
		}
		return Critical
	}
	if 100-cpuPercent[0] < cpuWarningThreshold {
		log.Println("CPU is Warning: ")
		if c,err:=cpu.Percent(0,true); err==nil{
			log.Println(c)
		}
		return Warning
	}

	//check memory usage
	v , err := mem.VirtualMemory()
	if err!=nil{
		log.Println("Failed get memory information.")
		return Critical
	}

	memoryPercent := v.UsedPercent
	memoryCriticalThreshold,_:=strconv.ParseFloat(configs["memory_critical_threshold"], 64)
	memoryWarningThreshold,_:=strconv.ParseFloat(configs["memory_warning_threshold"], 64)
	if 100-memoryPercent < memoryCriticalThreshold {
		log.Println("Memory is Critical: ")
		log.Println(v)
		return Critical
	}
	if 100-memoryPercent < memoryWarningThreshold {
		log.Println("Memory is Warning: ")
		log.Println(v)
		return Warning
	}

	//check load
	l, err := load.Avg()
	if err!=nil{
		log.Println("Failed get load information.")
		return Critical
	}
	loads := l.Load5
	loadCriticalThreshold,_:=strconv.ParseFloat(configs["load_critical_threshold"], 64)
	loadWarningThreshold,_:=strconv.ParseFloat(configs["load_warning_threshold"], 64)
	if loads > loadCriticalThreshold {
		log.Println("Load Critical")
		log.Println(l)
		if m,err:=load.Misc(); err==nil {
			log.Println(m)
		}
		return Critical
	}
	if loads > loadWarningThreshold {
		log.Printf("Load Warning")
		log.Println(l)
		if m,err:=load.Misc(); err==nil {
			log.Println(m)
		}
		return Warning
	}
	return OK
}
