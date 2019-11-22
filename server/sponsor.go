package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllSponsors - Retrieve all sponsors
func GetAllSponsors(collection *mongo.Collection, count int) ([]*Sponsor, error) {
	findOptions := options.Find()
	if count < 50 {
		findOptions.SetLimit(int64(count))
	} else {
		findOptions.SetLimit(50)
	}

	var sponsors []*Sponsor
	var cur *mongo.Cursor
	var err error

	cur, err = collection.Find(context.TODO(), bson.D{{}}, findOptions)

	if err != nil {
		return nil, err
	}

	// Iterate through all results
	for cur.Next(context.TODO()) {
		var elem Sponsor
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		sponsors = append(sponsors, &elem)
	}

	return sponsors, nil
}

// GetSponsors - Retrieve a sponsor from the database
func GetSponsors(collection *mongo.Collection, id string, token string) (Sponsor, error) {
	parsedID := uuid.Must(uuid.Parse(id))

	var result Sponsor
	filter := bson.D{{Key: "sponsorid", Value: parsedID}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

// NewSponsors - Add a new sponsor
func NewSponsors(collection *mongo.Collection, expiryStr string, name string, logo string, tier string, link string, token string) error {
	// if !validToken(token) {
	// 	return
	// }

	expiryTime, _ := time.Parse(time.RFC3339, expiryStr)
	id := uuid.New()

	sponsor := Sponsor{
		SponsorID:   id,
		SponsorName: name,
		SponsorLogo: logo,
		SponsorTier: tier,
		SponsorLink: link,
		Expiry:      expiryTime.Unix(),
	}

	_, err := collection.InsertOne(context.TODO(), sponsor)
	return err
}

// DeleteSponsors - Delete a sponsor from the database
func DeleteSponsors(collection *mongo.Collection, id string, token string) error {
	// if !validToken(token) {
	// 	return
	// }

	parsedID := uuid.Must(uuid.Parse(id))

	// Find a sponsor by ID and delete it
	filter := bson.D{{Key: "sponsorid", Value: parsedID}}
	_, err := collection.DeleteOne(context.TODO(), filter)
	return err
}
