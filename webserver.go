package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	db := newDatabase()
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database struct {
	items map[string]dollars
	mu    sync.RWMutex
}

func newDatabase() *database {
	return &database{
		items: map[string]dollars{"shoes": 50, "socks": 5},
	}
}

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for item, price := range db.items {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db *database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	db.mu.RLock()
	price, ok := db.items[item]
	db.mu.RUnlock()
	if !ok {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

func (db *database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		http.Error(w, "Invalid price value", http.StatusBadRequest)
		return
	}
	db.mu.Lock()
	db.items[item] = dollars(price)
	db.mu.Unlock()
	fmt.Fprintf(w, "Created %s with price %s\n", item, dollars(price))
}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil || price < 0 {
		http.Error(w, "Invalid price value", http.StatusBadRequest)
		return
	}
	db.mu.Lock()
	_, exists := db.items[item]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		db.mu.Unlock()
		return
	}
	db.items[item] = dollars(price)
	db.mu.Unlock()
	fmt.Fprintf(w, "Updated %s to price %s\n", item, dollars(price))
}

func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	db.mu.Lock()
	_, exists := db.items[item]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		db.mu.Unlock()
		return
	}
	delete(db.items, item)
	db.mu.Unlock()
	fmt.Fprintf(w, "Deleted %s\n", item)
}
