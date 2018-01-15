package gogo

import (
	"regexp"
	"strings"
)

var (
	pathrgxp = regexp.MustCompile(":[A-z,0-9,$,-,_,.,+,!,*,',(,),\\,]{1,}")
)

// URLToIntnlTrans builds the channel name for a internal transport use from an URL
// TODO regexp
func URLToIntnlTrans(host string, path string) string {
	//subURI := strings.Replace(path, "/", ".", -1)
	//subHost:=strings.Replace(host, "/", ".", -1)
	str := strings.Split(path, "/")

	return "gogo-" + str[2] + "-" + str[3]
}
