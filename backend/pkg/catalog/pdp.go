package catalog

import (
	"log"
	"net/http"
	_ "time" // show "loading..." in the frontend

	"github.com/golang/protobuf/proto"
)

func (c *context) NewPDPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// time.Sleep(2 * time.Second) // show "loading..." in the frontend

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")

		uuid := r.URL.Query().Get("uuid")
		if uuid == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		m, close := c.read()
		defer close()

		p, ok := m.products[uuid]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		respond(w, p)
	}
}

func respond(w http.ResponseWriter, response proto.Message) {
	w.Header().Set("Content-Type", "application/protobuf")

	bytes, err := proto.Marshal(response)
	if err != nil {
		log.Printf("failed to serialize: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to send result: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
