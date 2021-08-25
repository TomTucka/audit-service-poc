package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: audit-service -post-timestream")
	}

	postTimestream := flag.Bool("post-timetstream", false, "Post to timestream")
	flag.Parse()

	if postTimestream {
		cmd.PostTimestream()
	}
}
