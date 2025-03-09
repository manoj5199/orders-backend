package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order/database"
	"order/types"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		h.ServeHTTP(w, r)
	})
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var customer types.Customer
	_id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		fmt.Println("Error in converting to object id")
	}
	err = database.CustomerCollection.FindOne(context.TODO(), bson.M{"_id": _id}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func GetTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Test")
}

func GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cursor, err := database.CustomerCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var customers []types.Customer
	for cursor.Next(context.TODO()) {
		var customer types.Customer
		if err := cursor.Decode(&customer); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		customers = append(customers, customer)
	}

	json.NewEncoder(w).Encode(customers)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var product types.Product

	_id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		fmt.Println("Error in converting to object id")
	}
	err = database.ProductCollection.FindOne(context.TODO(), bson.M{"_id": _id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(product)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cursor, err := database.ProductCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var products []types.Product
	for cursor.Next(context.TODO()) {
		var product types.Product
		if err := cursor.Decode(&product); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	json.NewEncoder(w).Encode(products)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var order types.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order.ID = primitive.NewObjectID() // Generate a unique ID
	var total float64 = 0
	for _, product := range order.Products {
		total += product.Price * float64(product.Quantity)
	}
	order.Total = total
	_, err = database.OrderCollection.InsertOne(context.TODO(), order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(order)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	var order types.Order

	_id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		fmt.Println("Error in converting to object id")
	}
	err = database.OrderCollection.FindOne(context.TODO(), bson.M{"_id": _id}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(order)
}

func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cursor, err := database.OrderCollection.Find(context.TODO(), bson.M{}) // options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var orders []types.Order
	for cursor.Next(context.TODO()) {
		var order types.Order
		if err := cursor.Decode(&order); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}

	json.NewEncoder(w).Encode(orders)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var updatedOrder types.Order
	err := json.NewDecoder(r.Body).Decode(&updatedOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_id, err := primitive.ObjectIDFromHex(params["id"])

	if err != nil {
		fmt.Println("Error in converting to object id")
	}
	updatedOrder.ID = _id
	var total float64 = 0
	for _, product := range updatedOrder.Products {
		total += product.Price * float64(product.Quantity)
	}
	updatedOrder.Total = total

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updatedOrder}

	_, err = database.OrderCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedOrder)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	_id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		fmt.Println("Error in converting to object id")
	}
	filter := bson.M{"_id": _id}

	_, err = database.OrderCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
