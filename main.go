package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hello world")
		d, err := io.ReadAll(r.Body)

		if err != nil {
			// w.WriteHeader(http.StatusBadRequest)
			// w.Write([]byte("oops"))
			//or

			http.Error(w, "Oooops", http.StatusBadRequest)
			return
		}
		log.Printf("Data %s\n", d)

		fmt.Fprintf(w, "Hello %s", d)
	})

	http.HandleFunc("/goodbye", func(w http.ResponseWriter, r *http.Request) {
		log.Println("goodbye")
	})

	http.ListenAndServe(":9090", nil)
}
