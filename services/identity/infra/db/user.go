package db

import (
	"context"
	"errors"

	"github.com/edmarfelipe/next-u/services/identity/entity"
	"github.com/edmarfelipe/next-u/services/identity/infra/tracer"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDB interface {
	Create(ctx context.Context, model entity.User) error
	Update(ctx context.Context, model entity.User) error
	FindAll(ctx context.Context) (*[]entity.User, error)
	FindOne(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByUsername(ctx context.Context, user string) (*entity.User, error)
}

type userDB struct {
	collectionName string
	db             *qmgo.Database
}

func NewUser(mongoDB *qmgo.Database) UserDB {
	return &userDB{
		db:             mongoDB,
		collectionName: "users",
	}
}

func (rep *userDB) coll() *qmgo.Collection {
	return rep.db.Collection(rep.collectionName)
}

func (rep *userDB) Create(ctx context.Context, model entity.User) error {
	childCtx, span := tracer.StartSpan(ctx, "database", "UserDB.Create")
	defer span.End()

	_, err := rep.coll().InsertOne(childCtx, model)
	if err != nil {
		return err
	}

	return nil
}

func (rep *userDB) Update(ctx context.Context, model entity.User) error {
	childCtx, span := tracer.StartSpan(ctx, "database", "UserDB.Update")
	defer span.End()

	updated := bson.M{
		"$set": bson.M{
			"name":     model.Name,
			"email":    model.Email,
			"password": model.Password,
			"username": model.Username,
			"active":   model.Active,
		},
	}

	err := rep.coll().UpdateOne(childCtx, bson.M{"_id": model.ID}, updated)
	if err != nil {
		return err
	}

	return nil
}

func (rep *userDB) FindAll(ctx context.Context) (*[]entity.User, error) {
	childCtx, span := tracer.StartSpan(ctx, "database", "UserDB.FindAll")
	defer span.End()

	var result []entity.User
	rep.coll().
		Find(childCtx, bson.M{}).
		All(&result)

	return &result, nil
}

func (rep *userDB) FindOne(ctx context.Context, id string) (*entity.User, error) {
	childCtx, span := tracer.StartSpan(ctx, "database", "UserDB.FindOne")
	defer span.End()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := rep.coll().Find(childCtx, bson.M{"_id": objectId})

	var result entity.User
	err = query.One(&result)

	if errors.Is(err, qmgo.ErrNoSuchDocuments) {
		return nil, nil
	}

	return &result, nil
}

func (rep *userDB) FindByUsername(ctx context.Context, user string) (*entity.User, error) {
	childCtx, span := tracer.StartSpan(ctx, "database", "UserDB.FindByUsername")
	defer span.End()

	query := rep.coll().Find(childCtx, bson.M{"username": user})

	var result entity.User
	err := query.One(&result)

	if errors.Is(err, qmgo.ErrNoSuchDocuments) {
		return nil, nil
	}

	return &result, nil
}

func (rep *userDB) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	childCtx, span := tracer.StartSpan(ctx, "database", "UserDB.FindByEmail")
	defer span.End()

	query := rep.coll().
		Find(childCtx, bson.M{"email": email})

	var result entity.User
	err := query.One(&result)

	if errors.Is(err, qmgo.ErrNoSuchDocuments) {
		return nil, nil
	}

	return &result, nil
}
