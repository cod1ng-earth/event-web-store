package checkout

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
)

func (c *context) NewOrderHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Set-Cookie, *")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			return
		}

		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cookie, err := r.Cookie("cart")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("failed to find cart cookie id: %v", err)
			return
		}

		cartID := cookie.Value

		err = orderCart(c, w, cartID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to create order request: %v", err)
			return
		}

		model, free := c.aggregator.read()
		defer free()

		_, isOrdered := model.orders[cartID]

		status := &OrderCartResonse{}
		status.Successful = isOrdered
		if isOrdered {
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie = &http.Cookie{
				Name:    "cart",
				Value:   randomString(32),
				Expires: expiration,
			}
			http.SetCookie(w, cookie)
		}

		bytes, err := proto.Marshal(status)
		if err != nil {
			log.Printf("failed to serialize order status: %v", err)
		}

		_, err = w.Write(bytes)
		if err != nil {
			log.Printf("failed to send result: %v", err)
		}
	}
}

func orderCart(c *context, w http.ResponseWriter, cartID string) error {

	_, msgOffset, err := c.internalTopic.logOrderCart(&OrderCart{
		CartID: cartID,
	})
	if err != nil {
		return fmt.Errorf("failed to send cart change to kafka: %v", err)
	}

	c.aggregator.await(msgOffset)

	return nil
}
