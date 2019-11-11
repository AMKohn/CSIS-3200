package dashboard

import (
	"log"
	"net/http"
	"sync"
)

// Add a method to get a copy of the message from the processor
// Finish rest of dashboard build out

func HandleMessage(message map[string]interface{}) {
}

func StartServer(wg *sync.WaitGroup) {
	println("Dashboard initializing")

	// Handle API requests under /api
	http.HandleFunc("/api/", apiHandler)

	// Handle static files under /static
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Return the dashboard for all other requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	// Use a goroutine so we don't block while listening for requests
	go func() {
		// Tell the WaitGroup we're done after this function finishes
		defer wg.Done()

		log.Fatal(http.ListenAndServe(":80", nil))
	}()
}
