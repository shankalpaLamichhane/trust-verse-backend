package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

var Ctx context.Context

var Cancel context.CancelFunc

var Client *mongo.Client

var DB *mongo.Database

func Connect() {
	var err error

	Ctx, Cancel = context.WithTimeout(context.Background(), 30*time.Second)
	fmt.Print("THE DATABSE URL IS " + os.Getenv("MONGO_URI"))
	Client, err = mongo.Connect(Ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		panic(err)
	}
	DB = Client.Database("trustVerse")
}
