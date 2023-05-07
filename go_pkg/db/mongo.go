package db

import (
	"context"
	"fmt"
	"time"

	utils "github.com/sander-skjulsvik/webscraping/go_pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func IsInDb(c *mongo.Collection, realest Realestate) (*mongo.SingleResult, bool) {
	// If realest in db, return the decoded realest and true, else nil and false.

	one := c.FindOne(context.TODO(), bson.M{
		"id":      realest.ID,
		"title":   realest.Title,
		"address": realest.Address,
	})
	return one, one.Err() == nil

}

func UpdateManyRealestate(c *mongo.Collection, realestates map[int]*Realestate) {
	var old Realestate
	var isUpdated bool
	for _, realest := range realestates {
		if oldBson, exist := IsInDb(c, *realest); exist {
			// if realest exists in db, update if different
			oldBson.Decode(&old)
			if old.Updates[time.Now().String()], isUpdated = realest.RightUpdates(old); isUpdated {
				// Insert update
				c.UpdateOne(context.TODO(), oldBson, old)

			} // Else skip

		} else {
			// Insert
			c.InsertOne(context.TODO(), *realest)
		}
	}
}

func UpdateRealestate(collection *mongo.Collection, realestate Realestate) {
	var old Realestate
	var isUpdated bool
	if oldBson, exist := IsInDb(collection, realestate); exist {
		// if realest exists in db, update if different
		oldBson.Decode(&old)
		if old.Updates[time.Now().String()], isUpdated = realestate.RightUpdates(old); isUpdated {
			// Insert update
			collection.UpdateOne(context.TODO(), oldBson, old)

		} // Else skip

	} else {
		// Insert
		collection.InsertOne(context.TODO(), realestate)
	}
}

func GetFinnRealestateCollection() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, e := mongo.Connect(context.TODO(), clientOptions)
	utils.LogIfErr(e, "client, e := mongo.Connect(context.TODO(), clientOptions), Failed.")

	e = client.Ping(context.TODO(), nil)
	utils.LogIfErr(e, "e = client.Ping(context.TODO(), nil), Failed")
	fmt.Println("Connected to mongoDB!")
	collection := client.Database("Finn").Collection("Realestate2.0")
	fmt.Println("Collection:", collection)
	return collection
}
