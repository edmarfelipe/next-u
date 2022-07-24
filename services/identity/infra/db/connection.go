package db

import (
	"context"

	"github.com/qiniu/qmgo"
)

func NewConnection(dbName string, dbURI string) (*qmgo.Database, error) {
	client, err := qmgo.NewClient(context.TODO(), &qmgo.Config{Uri: dbURI})
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return db, nil
}
