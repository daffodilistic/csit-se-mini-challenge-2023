package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connection_string = "mongodb+srv://userReadOnly:7ZT817O8ejDfhnBM@minichallenge.q4nve1r.mongodb.net/"
const db_name = "minichallenge"

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/flight", func(c echo.Context) error {
		departureDate := c.QueryParam("departureDate")
		//returnDate := c.QueryParam("returnDate")
		destination := c.QueryParam("destination")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection_string))
		collection := client.Database(db_name).Collection("flights") // or hotels
		filter := bson.M{
			"destcity": bson.M{
				"$regex":   fmt.Sprintf("^%s", destination),
				"$options": "i",
			},
			"srccity": "Singapore",
		}

		if departureDate != "" {
			parsedDate, err := time.Parse(time.DateOnly, departureDate)
			if err == nil {
				filter["date"] = primitive.NewDateTimeFromTime(parsedDate)
			} else {
				fmt.Println(err)
			}
		}

		temp, err := json.Marshal(filter)
		fmt.Println(string(temp))

		cur, err := collection.Find(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(ctx)

		results := []bson.M{}
		for cur.Next(ctx) {
			var result bson.M
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			// do something with result....
			// log.Println(result)
			results = append(results, result)
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()

		// Sort by price ascending
		sort.Slice(results, func(i, j int) bool {
			return results[i]["price"].(int32) < results[j]["price"].(int32)
		})

		// res, _ := json.Marshal(results)
		return c.JSON(http.StatusOK, results)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
