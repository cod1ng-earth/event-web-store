package checkout

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/golang/protobuf/proto"
)

type positions []*Position

func (a positions) Len() int           { return len(a) }
func (a positions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a positions) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (c *context) NewCartHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Set-Cookie, *")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			return
		}

		cookie, err := r.Cookie("cart")
		if err != nil {
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie = &http.Cookie{
				Name:    "cart",
				Value:   randomString(32),
				Expires: expiration,
			}
			http.SetCookie(w, cookie)
		}

		cartID := cookie.Value

		if r.Method == "POST" {
			err := addToCart(c, w, r, cartID)
			if err != nil {
				log.Printf("failed to add a product to the cart: %v", err)
				return
			}
		}

		m, free := c.read()
		defer free()

		pp := positions{}
		for uuid, count := range m.carts[cartID] {
			if _, ok := m.products[uuid]; ok {
				pp = append(pp, &Position{
					ProductID:     uuid,
					Price:         m.products[uuid].Price,
					Name:          m.products[uuid].Name,
					SmallImageURL: m.products[uuid].SmallImageURL,

					Quantity:    count,
					MoreInStock: m.stock[uuid] > count,
					InStock:     m.stock[uuid] >= count,
				})
			}
		}
		sort.Sort(pp)

		cart := &Cart{
			Positions: pp,
		}

		bytes, err := proto.Marshal(cart)
		if err != nil {
			log.Printf("failed to serialize cart: %v", err)
		}

		_, err = w.Write(bytes)
		if err != nil {
			log.Printf("failed to send result: %v", err)
		}
	}
}

func addToCart(c *context, w http.ResponseWriter, r *http.Request, cartID string) error {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("failed to read from http body: %v", err)
	}

	ccr := &ChangeProductQuantityRequest{}
	err = proto.Unmarshal(bytes, ccr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("failed to decode cart change: %v", err)
	}

	cc := &ChangeProductQuantity{
		CartID:    cartID,
		ProductID: ccr.ProductID,
		Quantity:  ccr.Quantity,
	}
	_, msgOffset, err := c.logChangeProductQuantity(cc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("failed to send cart change to kafka: %v", err)
	}

	c.await(msgOffset)

	return nil
}
