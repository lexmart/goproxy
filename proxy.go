package main

import (
	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
	"log"
	"net/http"
	"io/ioutil"
	"os"
	"strings"
	"fmt"
)



func main() {
	ip := os.Args[1]
	port := os.Args[2]
	authFile := os.Args[3]

	authData, err := ioutil.ReadFile(authFile)	
	if err != nil {
		log.Fatal("Failed to read auth file")
	}

	authMap := map[string]string{}
	for _, line := range strings.Split(string(authData), "\n") {
		parts := strings.Split(line, ",")
		if len(parts) == 2 {
			user := strings.TrimSpace(parts[0])
			pass := strings.TrimSpace(parts[1])
			authMap[user] = pass
		}
	}

	proxy := goproxy.NewProxyHttpServer()
	auth.ProxyBasic(proxy, "proxy", func(user, passwd string) bool {
		if actualPasswd, ok := authMap[user]; ok {
			return actualPasswd == passwd
		}
		return false
	})

	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "")
		return 
	})

	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(ip + ":" + port, proxy))
}
