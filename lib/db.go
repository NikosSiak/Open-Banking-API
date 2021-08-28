package lib

import (
  "context"
  "time"

  "github.com/NikosSiak/Open-Banking-API/models"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
  client *mongo.Client
  database *mongo.Database
}

func NewDB(env Env) Database {
  uri := env.DatabaseURI
  dbName := env.DatabaseName

  ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
  defer cancel()

  client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
  if err != nil {
    panic(err)
  }

  database := client.Database(dbName)

  return Database{ client: client, database: database }
}

func (d Database) Close(ctx context.Context) {
  d.client.Disconnect(ctx)
}

func (d Database) InsertOne(ctx context.Context, document models.Model) (*mongo.InsertOneResult, error) {
  return d.database.Collection(document.CollectionName()).InsertOne(ctx, document)
}
