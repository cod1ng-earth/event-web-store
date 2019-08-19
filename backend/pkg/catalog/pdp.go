package catalog

import (
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"
)

func (c *context) NewPDPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		respondProtoBuf(w, p)
	}
}

//func respondJson(w http.ResponseWriter, response interface{}) {
//	w.Header().Set("Content-Type", "application/json")
//
//	bytes, err := json.Marshal(response)
//	if err != nil {
//		log.Printf("failed to serialize: %v", err)
//	}
//
//	_, err = w.Write(bytes)
//	if err != nil {
//		log.Printf("failed to send result: %v", err)
//	}
//}

func respondProtoBuf(w http.ResponseWriter, response proto.Message) {
	w.Header().Set("Content-Type", "application/protobuf")

	bytes, err := proto.Marshal(response)
	if err != nil {
		log.Printf("failed to serialize: %v", err)
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to send result: %v", err)
	}
}
