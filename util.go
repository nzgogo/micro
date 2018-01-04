package gogo

import (
	"fmt"
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
	str := strings.Split(path,"/")

	return str[1]+".core."+str[3]
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

// Mainly used by router.HttpMatch. Http request path contains
// information(key and subpath) that used to look up service in kv store
func PathToKeySubpath(path string) (key, subpath string){
	i := 0
	for m:=0;m<3;m++{
		x := strings.Index(path[i+1:],"/")
		if x < 0 {
			break
		}
		i += x
		i++
	}

	return path[:i], path[i:]
}