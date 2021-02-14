package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//func getAllKeys(c *mongo.Collection) {
//
//}

func isInDb(c *mongo.Collection, realest Realest) (*mongo.SingleResult, bool) {
	// If realest in db, return the decoded realest and true, else nil and false.

	one := c.FindOne(context.TODO(), bson.M{
		"id":      realest.ID,
		"title":   realest.Title,
		"address": realest.Address,
	})
	return one, one.Err() == nil

}

func UpdateManyRealestate(c *mongo.Collection, realestates map[int]*Realest) {
	var old Realest
	var isUpdated bool
	for _, realest := range realestates {
		if oldBson, exist := isInDb(c, *realest); exist {
			// if realest exists in db, update if different
			oldBson.Decode(&old)
			if old.Updates[time.Now().String()], isUpdated = realest.RightUpdates(old); isUpdated  {
				// Insert update
				c.UpdateOne(context.TODO(), oldBson, old)

			} // Else skip

		} else {
			// Insert
			c.InsertOne(context.TODO(), *realest)
		}
	}
}


func getFinnRealestateCollection() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, e := mongo.Connect(context.TODO(), clientOptions)
	logIfErr(e, "client, e := mongo.Connect(context.TODO(), clientOptions), Failed.")

	e = client.Ping(context.TODO(), nil)
	logIfErr(e, "e = client.Ping(context.TODO(), nil), Failed")
	fmt.Println("Connected to mongoDB!")
	collection := client.Database("Finn").Collection("Realestate2.0")
	fmt.Println("Collection:", collection)
	return collection
}
