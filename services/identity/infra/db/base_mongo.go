package db

import (
	"context"

	"github.com/qiniu/qmgo"
)

var DB_NAME = "db-identity"

func connect(ctx context.Context, collectionName string) (*qmgo.Collection, error) {
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017"})
	if err != nil {
		return nil, err
	}

	db := client.Database(DB_NAME)
	coll := db.Collection(collectionName)

	return coll, nil
}
