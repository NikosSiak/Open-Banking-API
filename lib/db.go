package lib

import (
	"context"
	"time"

	"github.com/NikosSiak/Open-Banking-API/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewDB(env Env) Database {
	uri := env.DatabaseURI
	dbName := env.DatabaseName

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	database := client.Database(dbName)

	return Database{client: client, database: database}
}

func (d Database) Close(ctx context.Context) {
	d.client.Disconnect(ctx)
}

func (d Database) InsertOne(ctx context.Context, document models.Model) (*mongo.InsertOneResult, error) {
	return d.database.Collection(document.CollectionName()).InsertOne(ctx, document.GetBSON())
}

func (d Database) FindOne(ctx context.Context, res models.Model, query bson.M, projection bson.M) error {
	return d.database.Collection((res).CollectionName()).FindOne(ctx, query, options.FindOne().SetProjection(projection)).Decode(res)
}

func (d Database) Find(ctx context.Context, res []models.Model, query bson.M, projection bson.M) error {
	filterCursor, err := d.database.Collection((res)[0].CollectionName()).Find(ctx, query)
	if err != nil {
		return err
	}

	if err = filterCursor.All(ctx, res); err != nil {
		return err
	}

	return nil
}

func (d Database) UpdateByID(ctx context.Context, id primitive.ObjectID, updatedModel models.Model) error {
	updateQuery := bson.M{
		"$set": updatedModel.GetBSON(),
	}

	_, err := d.database.Collection((updatedModel).CollectionName()).UpdateByID(ctx, id, updateQuery)

	return err
}
