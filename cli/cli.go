package cli

import (
	"flag"
	"fmt"
	"pkg/explorer"
	"pkg/rest"
	"runtime"
)

func usage() {
	fmt.Println("Welcome To Blockmonkey")
	fmt.Println("Please use the Flowing Commands :")
	fmt.Println("-port=4000 : Set Port!")
	fmt.Println("-mode=rest : Choose 'html', 'rest'")
	runtime.Goexit()
}

func Start() {
	port := flag.Int("port", 4000, "Set Port Of Server!")
	mode := flag.String("mode", "rest", "Choose 'html', 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		//start rest
		rest.Start(*port)
	case "html":
		//start html Explorer
		explorer.Start(*port)
	default:
		usage()
	}
}
