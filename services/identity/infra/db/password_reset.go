package db

import (
	"context"
	"errors"

	entity "github.com/edmarfelipe/next-u/services/identity/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PasswordResetDB interface {
	Create(ctx context.Context, model entity.PasswordReset) error
	Update(ctx context.Context, model entity.PasswordReset) error
	FindAll(ctx context.Context) (*[]entity.PasswordReset, error)
	FindOne(ctx context.Context, id string) (*entity.PasswordReset, error)
	FindByTokenNotDone(ctx context.Context, token string) (*entity.PasswordReset, error)
}

type passwordResetDB struct {
	collectionName string
	db             *qmgo.Database
}

func NewPasswordReset(mongoDB *qmgo.Database) PasswordResetDB {
	return &passwordResetDB{
		db:             mongoDB,
		collectionName: "password-resets",
	}
}

func (rep *passwordResetDB) coll() *qmgo.Collection {
	return rep.db.Collection(rep.collectionName)
}

func (rep *passwordResetDB) Create(ctx context.Context, model entity.PasswordReset) error {
	_, err := rep.coll().InsertOne(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

func (rep *passwordResetDB) Update(ctx context.Context, model entity.PasswordReset) error {
	updated := bson.M{
		"$set": bson.M{
			"userId":   model.UserID,
			"token":    model.Token,
			"createAt": model.CreateAt,
			"done":     model.Done,
		},
	}

	err := rep.coll().UpdateOne(ctx, bson.M{"_id": model.ID}, updated)
	if err != nil {
		return err
	}

	return nil
}

func (rep *passwordResetDB) FindAll(ctx context.Context) (*[]entity.PasswordReset, error) {
	var result []entity.PasswordReset
	rep.coll().Find(ctx, bson.M{}).
		All(&result)

	return &result, nil
}

func (rep *passwordResetDB) FindOne(ctx context.Context, id string) (*entity.PasswordReset, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := rep.coll().Find(ctx, bson.M{"_id": objectId})

	var result entity.PasswordReset
	err = query.One(&result)

	if errors.Is(err, qmgo.ErrNoSuchDocuments) {
		return nil, nil
	}

	return &result, nil
}

func (rep *passwordResetDB) FindByTokenNotDone(ctx context.Context, token string) (*entity.PasswordReset, error) {
	query := rep.coll().Find(ctx, bson.M{"token": token, "done": false})

	var result entity.PasswordReset
	err := query.One(&result)

	if errors.Is(err, qmgo.ErrNoSuchDocuments) {
		return nil, nil
	}

	return &result, nil
}
