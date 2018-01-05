package router

import (
	"testing"
	"fmt"
)

func TestPathToKeySubpath (t *testing.T) {
	path := "/gogox/v1/greeter/hello"

	key,subpath, err := PathToKeySubpath(path)
	fmt.Printf("key, subpath: %v,%v, %v",key,subpath,err)
	if key!="/gogox/v1/greeter" || subpath!="/hello" || err != nil{
		t.Fatalf("%v", err)
	}

}