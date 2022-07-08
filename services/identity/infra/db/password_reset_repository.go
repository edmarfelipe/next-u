package db

import (
	"context"

	"github.com/edmarfelipe/next-u/services/identity/entities"
	"go.mongodb.org/mongo-driver/bson"
)

type PasswordResetRepositorier interface {
	Create(ctx context.Context, model entities.PasswordReset) error
	Update(ctx context.Context, id string, model entities.PasswordReset) error
	FindAll(ctx context.Context) (*[]entities.PasswordReset, error)
	FindOne(ctx context.Context, id string) (*entities.PasswordReset, error)
}

type PasswordResetRepository struct {
	collectionName string
}

func NewPasswordResetRepository() PasswordResetRepositorier {
	return PasswordResetRepository{
		collectionName: "password-resets",
	}
}

func (repo PasswordResetRepository) Create(ctx context.Context, model entities.PasswordReset) error {
	coll, err := connect(ctx, repo.collectionName)
	if err != nil {
		return err
	}

	_, err = coll.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

func (repo PasswordResetRepository) Update(ctx context.Context, id string, model entities.PasswordReset) error {
	coll, err := connect(ctx, repo.collectionName)
	if err != nil {
		return err
	}

	err = coll.UpdateOne(ctx, bson.M{"id": id}, model)
	if err != nil {
		return err
	}

	return nil
}

func (repo PasswordResetRepository) FindAll(ctx context.Context) (*[]entities.PasswordReset, error) {
	coll, err := connect(ctx, repo.collectionName)
	if err != nil {
		return nil, err
	}

	var result []entities.PasswordReset
	coll.Find(ctx, bson.M{}).
		All(&result)

	return &result, nil
}

func (repo PasswordResetRepository) FindOne(ctx context.Context, id string) (*entities.PasswordReset, error) {
	coll, err := connect(ctx, repo.collectionName)
	if err != nil {
		return nil, err
	}

	var result *entities.PasswordReset
	coll.Find(ctx, bson.M{"_id": id}).
		One(result)

	return result, nil
}
