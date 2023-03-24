package mongo_cli

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type client struct {
	cli *mongo.Client
	db  *mongo.Database
}

var Client *client

func NewMongoClient() *client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:123456@mongo6"))
	if err != nil {
		logrus.Fatal("NewClient err: %s", err.Error())
	}
	db := cli.Database("name", options.Database())
	return &client{
		cli: cli,
		db:  db,
	}
}

func (c *client) Collection(name string) *mongo.Collection {
	return c.db.Collection(name)
}

func init() {
	Client = NewMongoClient()
}
