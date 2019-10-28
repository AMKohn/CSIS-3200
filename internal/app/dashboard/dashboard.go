package dashboard

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Got %s", r.URL.Path)

	if err != nil {
		log.Fatal(err)
	}
}

func StartServer(wg *sync.WaitGroup) {
	println("Dashboard initializing")
	http.HandleFunc("/", handler)

	// Use a goroutine so we don't block while listening for requests
	go func() {
		// Tell the WaitGroup we're done after this function finishes
		defer wg.Done()

		log.Fatal(http.ListenAndServe(":80", nil))
	}()
}