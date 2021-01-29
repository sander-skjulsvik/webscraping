package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

//type RealestateKey struct {
//	Address string
//	FinnId int
//}


func insertManyRealestate(c *mongo.Collection, realestates []interface{}) {
	// Find all distinct keys.
	// Key address and finnid
	c.Find(context.TODO(), bson.D{"", bson.D{{"id"}}})

	// compare with realestates, and filter out dups

	// insert new once

	// add updates to existing once
	insertRes, e := c.InsertMany(context.TODO(), realestates)
	logIfErr(e, "insertRes, e := c.InsertMany(context.TODO(), realestates), Failed.")
	log.Println("Inserted multiple documents:", insertRes.InsertedIDs)
}

//func getCollection(collectionName string) *mongo.Collection {
//
//}
//
//func getRealestateCollection()  *mongo.Collection {
//	return getCollection("finn_realestate")
//}

func getFinnRealestateCollection() *mongo.Collection{
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, e := mongo.Connect(context.TODO(), clientOptions)
	logIfErr(e, "client, e := mongo.Connect(context.TODO(), clientOptions), Failed.")

	e = client.Ping(context.TODO(), nil)
	logIfErr(e, "e = client.Ping(context.TODO(), nil), Failed")
	fmt.Println("Connected to mongoDB!")
	collection := client.Database("Finn").Collection("Realestate")
	fmt.Println("Collection:", collection)
	return collection
}