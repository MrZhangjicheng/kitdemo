package host

import (
	"fmt"
	"testing"
)

func TestExtract(t *testing.T) {
	hostport := ":0"
	res, _ := Extract(hostport, nil)
	fmt.Println(res)
}
