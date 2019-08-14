package products

import (
	"encoding/json"
	"log"
	"net/http"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
)

type read struct {
	p  *pb.Product
	ok bool
}

func (c *context) NewPDPHandler() http.HandlerFunc {
	return newPDPHandler(c.locklessModel)
}

func newPDPHandler(modelAccess locklessModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		w.Header().Set("Content-Type", "application/json")
		//	fmt.Fprintf(w, "offset: %d", offset)

		uuid := r.FormValue("uuid")
		if uuid == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("uuid is missing"))
			if err != nil {
				log.Printf("failed to send result: %v", err)
			}
			return
		}

		m, close := modelAccess.read()
		defer close()

		p, ok := m.products[uuid]
		if !ok {
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
}
