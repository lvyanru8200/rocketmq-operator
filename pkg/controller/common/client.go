package common

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func GetOperatorNamespace() string {
	nsBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		if os.IsNotExist(err) {
			return "defaultNamespace"
		}
		log.Panicf("err: %+v", err)
	}
	ns := strings.TrimSpace(string(nsBytes))
	return ns
}
