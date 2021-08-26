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

	postTimestream := flag.Bool("post-timestream", false, "Post to timestream")
	flag.Parse()
	fmt.Println("postTimeStream has value ", *postTimestream)


	if *postTimestream {
		fmt.Println("Posting timestream")
		cmd.PostTimestream()
	} else {
		fmt.Println("Unable to post timestream")
	}
}
