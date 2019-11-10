package dashboard

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Add a method to get a copy of the message from the processor
// Finish rest of dashboard build out

func HandleMessage(message map[string]interface{}) {
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Got %s", r.URL.Path)

	if err != nil {
		log.Fatal(err)
	}
}

func StartServer(wg *sync.WaitGroup) {
	println("Dashboard initializing")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/index.html")
	})

	fs := http.FileServer(http.Dir("web/static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Use a goroutine so we don't block while listening for requests
	go func() {
		// Tell the WaitGroup we're done after this function finishes
		defer wg.Done()

		log.Fatal(http.ListenAndServe(":80", nil))
	}()
}
