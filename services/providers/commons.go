package providers

import (
	"context"

	"github.com/NikosSiak/Open-Banking-API/lib"
	"github.com/NikosSiak/Open-Banking-API/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAccountToUser(ctx context.Context, db lib.Database, provider Bank, code, _userId string) error {
	userId, err := primitive.ObjectIDFromHex(_userId)
	if err != nil {
		return err
	}

	user := models.User{}

	err = db.FindOne(ctx, &user, bson.M{"_id": userId}, bson.M{})
	if err != nil {
		return err
	}

	account, err := provider.GetUserTokens(code)
	if err != nil {
		return err
	}

	user.AddAccount(provider.Name(), &account)

	err = db.UpdateByID(ctx, user.ID, &user)

	return err
}
