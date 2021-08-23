package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/open-cmi/cmmns"
	"github.com/open-cmi/cmmns/middleware"

	"github.com/gin-gonic/gin"
)

var configfile string = ""

func main() {
	flag.StringVar(&configfile, "config", configfile, "config file")
	flag.Parse()

	if configfile == "" {
		flag.Usage()
		return
	}

	err := cmmns.Init(configfile)
	if err != nil {
		fmt.Printf("service init failed\n")
		return
	}
	// init router
	r := gin.Default()

	cmmns.NauthInit(r)
	middleware.Init(r)
	cmmns.AuthInit(r)

	useSocket := false
	if useSocket {
		const sockAddr = "/tmp/cmmns.sock"
		os.Remove(sockAddr)
		unixAddr, err := net.ResolveUnixAddr("unix", sockAddr)
		if err != nil {
			fmt.Println(err)
			return
		}

		listener, err := net.ListenUnix("unix", unixAddr)
		if err != nil {
			fmt.Println("listening error:", err)
		}
		fmt.Printf("listening unix socket: %s\n", sockAddr)
		http.Serve(listener, r)
	} else {
		r.Run(":30002")
	}
}
