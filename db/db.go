package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dipankarupd/go-graphql-crud/graph/model"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// load the .env file

// connection string for mongodb
var connString = ""
var dbName = "bookstore"
var collName = "book"

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	connString = os.Getenv("MONGODB_URL")
	if connString == "" {
		panic("MONGODB_URL environment variable is not set")
	}
}

type DB struct {
	client *mongo.Client
}

func Connect() *DB {

	clientOptions := options.Client().ApplyURI(connString)

	// connect to the mongodb:
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to MongoDB")

	return &DB{
		client: client,
	}

}

func (db *DB) AddBook(bookInfo model.AddBookInput) *model.Books {

	bookColl := db.client.Database(dbName).Collection(collName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	insert, err := bookColl.InsertOne(ctx, bson.M{"title": bookInfo.Title, "author": bookInfo.Author, "genre": bookInfo.Genre, "price": bookInfo.Price})
	if err != nil {
		panic(err)
	}
	bookId := insert.InsertedID.(primitive.ObjectID).Hex()
	resp := model.Books{
		ID:     bookId,
		Title:  bookInfo.Title,
		Author: bookInfo.Author,
		Genre:  bookInfo.Genre,
		Price:  bookInfo.Price,
	}
	return &resp
}

func (db *DB) UpdateBook(id string, bookInfo model.UpdateBookInput) *model.Books {
	bookColl := db.client.Database(dbName).Collection(collName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updateBookInfo := bson.M{}

	if bookInfo.Title != nil {
		updateBookInfo["title"] = bookInfo.Title
	}

	if bookInfo.Author != nil {
		updateBookInfo["author"] = bookInfo.Author
	}

	if bookInfo.Genre != nil {
		updateBookInfo["genre"] = bookInfo.Genre
	}

	if bookInfo.Price != nil {
		updateBookInfo["price"] = bookInfo.Price
	}

	_id, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateBookInfo}

	res := bookColl.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedBook model.Books

	if err := res.Decode(&updatedBook); err != nil {
		panic(err)
	}
	return &updatedBook
}

func (db *DB) GetAllBooks() []*model.Books {

	bookColl := db.client.Database(dbName).Collection(collName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var books []*model.Books

	res, err := bookColl.Find(ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	if err := res.All(context.TODO(), &books); err != nil {
		panic(err)
	}
	return books
}

func (db *DB) GetBook(id string) *model.Books {
	bookColl := db.client.Database(dbName).Collection(collName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": _id}

	var book *model.Books

	if err := bookColl.FindOne(ctx, filter).Decode(&book); err != nil {
		panic(err)
	}

	return book

}

func (db *DB) RemoveBook(id string) *model.RemoveBookResponse {
	bookColl := db.client.Database(dbName).Collection(collName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)

	_, err := bookColl.DeleteOne(ctx, bson.M{"_id": _id})

	if err != nil {
		panic(err)
	}

	return &model.RemoveBookResponse{DeletedBookID: id}
}
