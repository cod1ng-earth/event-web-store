package products

import (
	"encoding/json"
	"log"
	"net/http"
)

func PDPHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
	w.Header().Set("Content-Type", "application/json")
	//	fmt.Fprintf(w, "offset: %d", offset)

	mut.RLock()
	defer mut.RUnlock()

	uuid := r.FormValue("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("uuid is missing"))
		if err != nil {
			log.Printf("failed to send result: %v", err)
		}
		return
	}

	p, found := products[uuid]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(p)
	if err != nil {
		log.Printf("failed to serialize: %v", err)
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to send result: %v", err)
	}
}
