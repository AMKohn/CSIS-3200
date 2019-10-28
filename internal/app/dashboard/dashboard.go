package dashboard

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Got %s", r.URL.Path)

	if err != nil {
		log.Fatal(err)
	}
}

func Main() {
	println("Dashboard initializing")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}