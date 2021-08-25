package main

import (
	"flag"
	"fmt"
	"github.com/tomtucka/audit-service-poc/cmd"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: audit-service -post-timestream")
	}

	var postTimeStream bool

	flag.BoolVar(&postTimeStream, "post-timestream", false, "Post to timestream")
	flag.Parse()

	if postTimeStream {
		fmt.Println("Posting timestream")
		cmd.PostTimestream()
	} else {
		fmt.Println("Unable to post to timestream")
	}
}
