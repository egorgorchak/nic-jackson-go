package handlers

import (
	"log"
	"net/http"
	"product/api/data"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		p.l.Println("Handle PUT Products")
		rx := regexp.MustCompile("/([0-9]+)")

		g := rx.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("more than 1 ids")
			http.Error(rw, "Invalid URL", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("more than 1 capture groups")
			http.Error(rw, "Invalid URL", http.StatusBadRequest)
			return
		}

		idStr := g[0][1]
		id, err := strconv.Atoi(idStr)

		if err != nil {
			p.l.Println("unable convert number", idStr)
			http.Error(rw, "Invalid URL", http.StatusBadRequest)
		}

		p.l.Println("got id", id)
		p.updateProduct(id, rw, r)
		return
	}
	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, h *http.Request) {
	p.l.Println("Handle GET Products")

	lp := data.GetProducts()
	err := lp.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	p.l.Printf("Prod: %#v", prod)
	err2 := data.UpdateProduct(id, prod)

	if err2 == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err2 != nil {
		http.Error(rw, "Product not found - Internal error", http.StatusInternalServerError)
		return
	}
}
