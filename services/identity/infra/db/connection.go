package db

import (
	"context"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	mongo "go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

func NewConnection(dbName string, dbURI string) (*qmgo.Database, error) {
	mongoConfig := options.ClientOptions{
		ClientOptions: &mongo.ClientOptions{
			Monitor: otelmongo.NewMonitor(),
		},
	}
	qmgoConfig := qmgo.Config{
		Uri: dbURI,
	}

	client, err := qmgo.NewClient(context.TODO(), &qmgoConfig, mongoConfig)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return db, nil
}
