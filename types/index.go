package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Customer struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Name    string             `json:"name"`
	Email   string             `json:"email"`
	Address string             `json:"address"`
}

type Product struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"product_name"`
	Price    float64            `json:"price"`
	Quantity int                `json:"quantity"`
}

type Order struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Customer Customer           `json:"customer"`
	Products []Product          `json:"products"`
	Total    float64            `json:"total"`
}
