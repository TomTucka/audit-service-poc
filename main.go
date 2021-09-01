package main

import (
	"flag"
	"fmt"
	"github.com/tomtucka/audit-service-poc/cmd"
	"log"
	"net/http"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: audit-service -post-timestream")
	}

	postTimestream := flag.Bool("post-timestream", false, "Post to timestream")
	getTimestream := flag.Bool("get-timestream", false, "Get timestream")
	serve := flag.Bool("serve", false, "serving")

	flag.Parse()

	if *postTimestream {
		sum := 1
		for sum < 100 {
			fmt.Println(sum)
			cmd.PostTimestream()
			sum += 1
		}
	}

	if *getTimestream {
		cmd.GetTimestream("20")
	}

	if *serve {
		handleRequests()
	}
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/sample", returnGetTimeStreamResults)
	http.HandleFunc("/deputy/", returnResultById)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnGetTimeStreamResults(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: timestream response")
	fmt.Fprint(w, cmd.GetTimestream("20"))
}

func returnResultById(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: timestream by deputy id response")
	id := strings.TrimPrefix(r.URL.Path, "/deputy/")
	fmt.Println(id)
	fmt.Fprint(w, cmd.GetTimestream(id))
}
