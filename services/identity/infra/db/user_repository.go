package db

import (
	"context"
	"errors"

	"github.com/edmarfelipe/next-u/services/identity/entities"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepositorier interface {
	Create(ctx context.Context, model entities.User) error
	Update(ctx context.Context, id string, model entities.User) error
	FindAll(ctx context.Context) (*[]entities.User, error)
	FindOne(ctx context.Context, id string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByUser(ctx context.Context, user string) (*entities.User, error)
	FindByUserAndPass(ctx context.Context, user string, pass string) (*entities.User, error)
}

type UserRepository struct {
	collectionName string
}

func NewUserRepository() UserRepositorier {
	return UserRepository{
		collectionName: "users",
	}
}

func (repo UserRepository) Create(ctx context.Context, model entities.User) error {
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

func (repo UserRepository) Update(ctx context.Context, id string, model entities.User) error {
	coll, err := connect(ctx, repo.collectionName)
	if err != nil {
		return err
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updated := bson.M{
		"$set": bson.M{
			"name":   model.Name,
			"email":  model.Email,
			"pass":   model.Password,
			"user":   model.Username,
			"active": model.Active,
		},
	}

	err = coll.UpdateOne(ctx, bson.M{"_id": objectId}, updated)
	if err != nil {
		return err
	}

	return nil
}

func (repo UserRepository) FindAll(ctx context.Context) (*[]entities.User, error) {
	coll, err := connect(ctx, repo.collectionName)
	if err != nil {
		return nil, err
	}

	var result []entities.User
	coll.Find(ctx, bson.M{}).
		All(&result)

	return &result, nil
}

func (repo UserRepository) FindOne(ctx context.Context, id string) (*entities.User, error) {
	coll, err := connect(ctx, repo.collectionName)
	if err != nil {
		return nil, err
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := coll.Find(ctx, bson.M{"_id": objectId})

	var result entities.User
	err = query.One(&result)

	if errors.Is(err, qmgo.ErrNoSuchDocuments) {
		return nil, nil
	}

	return &result, nil
}

func (repo UserRepository) FindByUser(ctx context.Context, user string) (*entities.User, error) {
	coll, err := connect(ctx, repo.collectionName)
	if err != nil {
		return nil, err
	}

	query := coll.Find(ctx, bson.M{"user": user})

	var result entities.User
	err = query.One(&result)

	if errors.Is(err, qmgo.ErrNoSuchDocuments) {
		return nil, nil
	}

	return &result, nil
}

func (repo UserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	coll, err := connect(ctx, repo.collectionName)
	if err != nil {
		return nil, err
	}

	query := coll.Find(ctx, bson.M{"email": email})

	var result entities.User
	err = query.One(&result)

	if errors.Is(err, qmgo.ErrNoSuchDocuments) {
		return nil, nil
	}

	return &result, nil
}

func (repo UserRepository) FindByUserAndPass(ctx context.Context, user string, pass string) (*entities.User, error) {
	coll, err := connect(ctx, repo.collectionName)
	if err != nil {
		return nil, err
	}

	query := coll.Find(ctx, bson.M{
		"user":   user,
		"pass":   pass,
		"active": true,
	})

	var result entities.User
	err = query.One(&result)

	if errors.Is(err, qmgo.ErrNoSuchDocuments) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}
