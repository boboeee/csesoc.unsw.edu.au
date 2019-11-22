package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// GetPosts - Retrieve a post from the database
func GetPosts(collection *mongo.Collection, id int, category int) (Post, error) {
	var result Post

	// Search for post by id and category
	filter := bson.D{{Key: "postid", Value: id}, {Key: "category", Value: category}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

// GetAllPosts - Retrieve all posts
func GetAllPosts(collection *mongo.Collection, count int, cat int) ([]*Post, error) {
	findOptions := options.Find()
	if count < 50 {
		findOptions.SetLimit(int64(count))
	} else {
		findOptions.SetLimit(50)
	}

	var posts []*Post
	var cur *mongo.Cursor
	var err error

	if cat == 0 { // No specified category
		cur, err = collection.Find(context.TODO(), bson.D{{}}, findOptions)
	} else {
		filter := bson.D{{Key: "postcategory", Value: cat}}
		cur, err = collection.Find(context.TODO(), filter, findOptions)
	}

	if err != nil {
		return nil, err
	}

	// Iterate through all results
	for cur.Next(context.TODO()) {
		var elem Post
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &elem)
	}

	return posts, nil
}

// NewPosts - Add a new post
func NewPosts(collection *mongo.Collection, id int, category int, showInMenu bool,
	title string, subtitle string, postType string, content string, image string, resource string, canonical string) error {
	currTime := time.Now()
	post := Post{
		PostID:        id,
		PostTitle:     title,
		PostSubtitle:  subtitle,
		PostType:      postType,
		PostCategory:  category,
		CreatedOn:     currTime.Unix(),
		PostContent:   content,
		ImageLink:     image,
		ResourceLink:  resource,
		CanonicalLink: canonical,
		ShowInMenu:    showInMenu,
	}

	_, err := collection.InsertOne(context.TODO(), post)
	return err
}

// UpdatePosts - Update a post with new information
func UpdatePosts(collection *mongo.Collection, id int, category int, showInMenu bool,
	title string, subtitle string, postType string, content string, image string, resource string, canonical string) error {
	filter := bson.D{{Key: "postid", Value: id}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "posttitle", Value: title},
			{Key: "postsubtitle", Value: subtitle},
			{Key: "posttype", Value: postType},
			{Key: "postcategory", Value: category},
			{Key: "lasteditedon", Value: time.Now()},
			{Key: "postcontent", Value: content},
			{Key: "resourcelink", Value: resource},
			{Key: "imagelink", Value: image},
			{Key: "canonicallink", Value: canonical},
		}},
	}

	// Find a post by id and update it
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// DeletePosts - Delete a post from the database
func DeletePosts(collection *mongo.Collection, id int) error {
	filter := bson.D{{Key: "postid", Value: id}}

	// Find a post by id and delete it
	_, err := collection.DeleteOne(context.TODO(), filter)
	return err
}
