package main

import (
	"flag"
	"fmt"
	"github.com/tomtucka/audit-service-poc/cmd"
	"html/template"
	"net/http"
	"os"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("<h1>Hello World!</h1>"))
	tpl.Execute(w, nil)
}

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
		for sum < 10000000 {
			fmt.Println(sum)
			cmd.PostTimestream()
			sum += 1
		}
	}

	if *getTimestream {
		cmd.GetTimestream()
	}

	if *serve {
		port := os.Getenv("PORT")
		if port == "" {
			port = "4000"
		}

		mux := http.NewServeMux()

		mux.HandleFunc("/", indexHandler)
		http.ListenAndServe(":"+port, mux)
	}
}
