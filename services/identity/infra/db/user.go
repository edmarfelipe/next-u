package db

import (
	"context"
	"errors"

	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/edmarfelipe/next-u/services/identity/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDB interface {
	Create(ctx context.Context, model entity.User) error
	Update(ctx context.Context, model entity.User) error
	FindAll(ctx context.Context) ([]entity.User, error)
	FindOne(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByTokenNotDone(ctx context.Context, token string) (*entity.User, error)
}

type userDB struct {
	collectionName string
	db             *qmgo.Database
	logger         logger.Logger
}

func NewUser(mongoDB *qmgo.Database, logger logger.Logger) UserDB {
	return &userDB{
		db:             mongoDB,
		logger:         logger,
		collectionName: "users",
	}
}

func (rep *userDB) coll() *qmgo.Collection {
	return rep.db.Collection(rep.collectionName)
}

func (rep *userDB) Create(ctx context.Context, model entity.User) error {
	rep.logger.Info(ctx, "Inserting user "+model.Email)

	_, err := rep.coll().InsertOne(ctx, model)
	if err != nil {
		rep.logger.Error(ctx, "Failed to insert user", "error", err)
		return err
	}

	return nil
}

func (rep *userDB) Update(ctx context.Context, model entity.User) error {
	rep.logger.Info(ctx, "Updating user "+model.ID.String())

	updated := bson.M{
		"$set": bson.M{
			"name":           model.Name,
			"email":          model.Email,
			"password":       model.Password,
			"active":         model.Active,
			"passwordResets": model.PasswordResets,
		},
	}

	err := rep.coll().UpdateOne(ctx, bson.M{"_id": model.ID}, updated)
	if err != nil {
		rep.logger.Error(ctx, "Failed to update user", "error", err)
		return err
	}

	return nil
}

func (rep *userDB) FindAll(ctx context.Context) ([]entity.User, error) {
	filter := bson.M{}
	rep.logger.Error(ctx, "Finding users", "filter", filter)

	var result []entity.User
	err := rep.coll().
		Find(ctx, filter).
		All(&result)

	if err != nil {
		rep.logger.Error(ctx, "Failed to find users", "error", err)
		return nil, err
	}

	return result, nil
}

func (rep *userDB) FindOne(ctx context.Context, id string) (*entity.User, error) {
	rep.logger.Info(ctx, "Finding user "+id)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		rep.logger.Error(ctx, "Failed to parse objectId", "error", err)
		return nil, err
	}

	query := rep.coll().Find(ctx, bson.M{"_id": objectId})

	var result entity.User
	err = query.One(&result)

	if errors.Is(err, qmgo.ErrNoSuchDocuments) {
		return nil, nil
	}

	if err != nil {
		rep.logger.Error(ctx, "Failed to find user by id "+id, "error", err)
		return nil, err
	}

	return &result, nil
}

func (rep *userDB) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	rep.logger.Info(ctx, "Finding user by email "+email)

	query := rep.coll().
		Find(ctx, bson.M{"email": email})

	var result entity.User
	err := query.One(&result)

	if errors.Is(err, qmgo.ErrNoSuchDocuments) {
		return nil, nil
	}

	if err != nil {
		rep.logger.Error(ctx, "Failed to find user by email", "error", err)
		return nil, err
	}

	return &result, nil
}

func (rep *userDB) FindByTokenNotDone(ctx context.Context, token string) (*entity.User, error) {
	rep.logger.Info(ctx, "Finding user by token "+token[:8])

	filter := bson.D{
		bson.E{"passwordResets", bson.D{
			{"$elemMatch", bson.D{
				{"token", token},
				{"done", false},
			}},
		}},
	}

	query := rep.coll().Find(ctx, filter)

	var result entity.User
	err := query.One(&result)
	if errors.Is(err, qmgo.ErrNoSuchDocuments) {
		return nil, nil
	}

	if err != nil {
		rep.logger.Error(ctx, "Failed to find user by token", "error", err)
		return nil, err
	}

	return &result, nil
}
