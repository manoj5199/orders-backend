package main

import (
	"fmt"
	"log"
	"net/http"
	"order/controller"
	"order/database"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.ConnectDB("mongodb+srv://puser:N00Bmanoj@atlascluster.w2gn9mm.mongodb.net/?retryWrites=true&w=majority&appName=AtlasCluster", "order_management", "customers_collection", "products_collection", "orders_collection")

	router := mux.NewRouter()

	router.Handle("/api/customers/{id}", controller.Cors(http.HandlerFunc(controller.GetCustomer))).Methods("GET")
	router.Handle("/api/customers", controller.Cors(http.HandlerFunc(controller.GetAllCustomers))).Methods("GET")
	router.Handle("/api/products/{id}", controller.Cors(http.HandlerFunc(controller.GetProduct))).Methods("GET")
	router.Handle("/api/products", controller.Cors(http.HandlerFunc(controller.GetAllProducts))).Methods("GET")

	router.Handle("/api/orders", controller.Cors(http.HandlerFunc(controller.CreateOrder))).Methods("POST")
	router.Handle("/api/orders/{id}", controller.Cors(http.HandlerFunc(controller.GetOrder))).Methods("GET")
	router.Handle("/api/orders", controller.Cors(http.HandlerFunc(controller.GetAllOrders))).Methods("GET")
	router.Handle("/api/orders/{id}", controller.Cors(http.HandlerFunc(controller.UpdateOrder))).Methods("PUT")
	router.Handle("/api/orders/{id}", controller.Cors(http.HandlerFunc(controller.DeleteOrder))).Methods("DELETE")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	fmt.Println("Server listening on port 3010...")
	log.Fatal(http.ListenAndServe(":3010", c.Handler(router)))
}
