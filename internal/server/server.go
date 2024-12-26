package server

import (
	"context"
	"fmt"
	"net/http"
	"rest_api_pks/internal/handlers"
	"rest_api_pks/internal/repository"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Server struct {
	router *http.ServeMux
	repo   *repository.ProductRepository
}

func NewServer() *Server {
	dbpool, err := pgxpool.Connect(context.Background(), "postgres://postgres:Googleapple123@localhost:5432/product_db")
	if err != nil {
		panic(err)
	}

	repo := repository.NewProductRepository(dbpool)
	handler := handlers.NewProductHandler(repo)

	router := http.NewServeMux()
	router.HandleFunc("/products", handler.GetProductsHandler)
	router.HandleFunc("/products/create", handler.CreateProductHandler)
	router.HandleFunc("/products/", handler.GetProductByIDHandler)
	router.HandleFunc("/products/update/", handler.UpdateProductHandler)
	router.HandleFunc("/products/delete/", handler.DeleteProductHandler)
	router.HandleFunc("/products/quantity/", handler.UpdateProductQuantityHandler)
	router.HandleFunc("/products/favorite/", handler.ToggleFavoriteHandler)
	router.HandleFunc("/products/cart/", handler.ToggleCartHandler)
	router.HandleFunc("/orders/create", handler.CreateOrderHandler)
	router.HandleFunc("/orders", handler.GetOrdersHandler)

	return &Server{
		router: router,
		repo:   repo,
	}
}

func (s *Server) Start(addr string) error {
	fmt.Printf("Server is running on http://localhost%s !\n", addr)
	return http.ListenAndServe(addr, s.router)
}
