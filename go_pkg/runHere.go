package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func findDistinct(c *mongo.Collection, fieldName string) {

	res, err := c.Distinct(context.TODO(), fieldName, bson.D{})

	logIfErr(err, "Getting distinct : ")
	for ind, r := range res {
		fmt.Printf("%d: %s\n", ind, r)
	}
	file, _ := os.Create("unique.txt")
	for _, row := range res {
		file.Write([]byte(fmt.Sprintf("%v", row)))
	}
}

func groupBy(c *mongo.Collection) {
	pipeline := make([]bson.M, 0)
	groupStage := bson.M{
		"$group": bson.M{
			"_id":   bson.M{"id": "$id", "title": "$title", "address": "$address"},
			"count": bson.M{"$sum": 1},
		},
	}
	pipeline = append(pipeline, groupStage)
	matchStage := bson.M{
		"$match": bson.M{"count": bson.M{"$gt": 4}},
	}
	pipeline = append(pipeline, matchStage)

	data, err := c.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatalf("Aggregate failed: %s", err)
	}
	res := make([]interface{}, 0)
	err = data.All(context.TODO(), &res)
	logIfFatal(err, "Failed to gather all data")
	for ind, val := range res {
		fmt.Printf("%d, %v\n", ind, val)
	}

	fmt.Print(groupStage)
}

//func getAllActiveRealestKeys(c *mongo.Collection) []RealestKey {
//	pipeline := make([]bson.M, 0)
//	// //  Match active
//	// pipeline = append(pipeline, bson.M{
//	// 	"$match": bson.M{"Active": true},
//	// })
//	// Group by, unesessary after adding unique key to db
//	//	Title, Address, Id
//	pipeline = append(pipeline, bson.M{
//		"$group": bson.M{
//			"_id": bson.M{"id": "$id", "title": "$title", "address": "$address"},
//			// "count": bson.M{"$sum": 1},
//		},
//	})
//	// Select the key
//	pipeline = append(pipeline, bson.M{
//		"$project": bson.M{
//			"ID":      "$_id.id",
//			"Title":   "$_id.title",
//			"Address": "$_id.address",
//		},
//	})
//	// Aggregate
//	cursor, err := c.Aggregate(context.TODO(), pipeline)
//	if err != nil {
//		log.Fatalf("Aggragtion error, not able to get unique keys: %s", err)
//	}
//
//	// printCursorValues(cursor)
//
//	res := make([]RealestKey, 0)
//
//	for cursor.Next(context.TODO()) {
//		var elem RealestKey
//		err := cursor.Decode(&elem)
//		logIfFatal(err, "Failed to decode in getAllActiverealestateKeys: ")
//		res = append(res, elem)
//
//	}
//	fmt.Println(res)
//
//	return res
//}

type Eg struct {
	a int
	b string
}

func main() {
	UpdateFinnDB()

}
