package gogoapi

import (
	"fmt"
	"regexp"
	"strings"

)

var (
	pathrgxp = regexp.MustCompile(":[A-z,0-9,$,-,_,.,+,!,*,',(,),\\,]{1,}")
)

// URLToNats builds the channel name
// from an URL and Method of http.Request
func URLToNats(host string, path string) string {
	subURI := strings.Replace(path, "/", ".", -1)
	subHost:=strings.Replace(host, "/", ".", -1)

	return subURI+subHost
}

// SubscribeURLToNats buils the subscription
// channel name with placeholders (started with ":").
// The placeholders are than used to obtain path variables
func SubscribeURLToNats(method string, urlPath string) string {
	subURL := pathrgxp.ReplaceAllString(urlPath, "*")
	// subURL = lastpathrgxp.ReplaceAllString(subURL, ".*")
	subURL = strings.Replace(subURL, "/", ".", -1)
	subURL = fmt.Sprintf("%s:%s", method, subURL)
	return subURL
}