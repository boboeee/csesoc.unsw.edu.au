package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllCategories retrieves a list of all post categories
func GetAllCategories(collection *mongo.Collection, count int, token string) ([]*Category, error) {
	findOptions := options.Find()
	if count == 0 || count < 50 {
		findOptions.SetLimit(int64(count))
	} else {
		findOptions.SetLimit(50)
	}

	var categories []*Category
	var cur *mongo.Cursor
	var err error

	cur, err = collection.Find(context.TODO(), bson.D{{}}, findOptions)

	if err != nil {
		return nil, err
	}

	// Iterate through all results
	for cur.Next(context.TODO()) {
		var elem Category
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &elem)
	}

	return categories, nil
}

// GetCategories retrieves a post category
func GetCategories(collection *mongo.Collection, id int, token string) (Category, error) {
	var result Category
	filter := bson.D{{Key: "categoryid", Value: id}}

	// Find a category
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

// NewCategories adds a post category
func NewCategories(collection *mongo.Collection, catID int, index int, name string, token string) error {
	category := Category{
		CategoryID:   catID,
		CategoryName: name,
		Index:        index,
	}

	// Insert a category
	_, err := collection.InsertOne(context.TODO(), category)
	return err
}

// PatchCategories updates some details in a post category
func PatchCategories(collection *mongo.Collection, catID int, name string, index int, token string) error {
	filter := bson.D{{Key: "categoryid", Value: catID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "categoryname", Value: name},
			{Key: "index", Value: index},
		}},
	}

	// Find a category by id and update it
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// DeleteCategories deletes a post category
func DeleteCategories(collection *mongo.Collection, id int, token string) error {
	filter := bson.D{{Key: "categoryid", Value: id}}

	// Find a category by id and delete it
	_, err := collection.DeleteOne(context.TODO(), filter)
	return err
}
