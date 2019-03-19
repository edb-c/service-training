package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/ardanlabs/service-training/15-validation/internal/platform/web"
	"github.com/ardanlabs/service-training/15-validation/internal/products"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	db  *sqlx.DB
	log *log.Logger

	http.Handler
}

// NewProducts creates a product handler with multiple routes defined.
func NewProducts(db *sqlx.DB, log *log.Logger) *Products {
	p := Products{db: db, log: log}

	r := chi.NewRouter()
	r.Post("/v1/products", web.Run(p.Create))
	r.Get("/v1/products", web.Run(p.List))
	r.Get("/v1/products/{id}", web.Run(p.Get))

	r.Post("/v1/products/{id}/sales", web.Run(p.AddSale))
	r.Get("/v1/products/{id}/sales", web.Run(p.ListSales))
	p.Handler = r

	return &p
}

// Create decodes the body of a request to create a new product. The full
// product with generated fields is sent back in the response.
func (s *Products) Create(w http.ResponseWriter, r *http.Request) error {
	var p products.Product
	if err := web.Decode(r, &p); err != nil {
		return errors.Wrap(err, "decoding new product")
	}

	if err := products.Create(r.Context(), s.db, &p); err != nil {
		return errors.Wrap(err, "creating new product")
	}

	return web.Encode(w, &p, http.StatusCreated)
}

// List gets all products from the service layer.
func (s *Products) List(w http.ResponseWriter, r *http.Request) error {
	list, err := products.List(r.Context(), s.db)
	if err != nil {
		return errors.Wrap(err, "getting product list")
	}

	return web.Encode(w, list, http.StatusOK)
}

// Get finds a single product identified by an ID in the request URL.
func (s *Products) Get(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	p, err := products.Get(r.Context(), s.db, id)
	if err != nil {
		return errors.Wrapf(err, "getting product %q", id)
	}

	return web.Encode(w, p, http.StatusOK)
}

// AddSale creates a new Sale for a particular product. It looks for a JSON
// object in the request body. The full model is returned to the caller.
func (s *Products) AddSale(w http.ResponseWriter, r *http.Request) error {
	var sale products.Sale
	if err := web.Decode(r, &s); err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	sale.ProductID = chi.URLParam(r, "id")

	if err := products.AddSale(r.Context(), s.db, &sale); err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Encode(w, sale, http.StatusCreated)
}

// ListSales gets all sales for a particular product.
func (s *Products) ListSales(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := products.ListSales(r.Context(), s.db, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Encode(w, list, http.StatusOK)
}