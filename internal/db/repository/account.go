package repository

import (
	"context"
	"errors"
	"github.com/dw-account-service/internal/db"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/utilities"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"time"
)

type AccountRepository struct {
	Entity *entity.AccountBalance
}

func NewAccountRepository() AccountRepository {
	return AccountRepository{Entity: new(entity.AccountBalance)}
}

func (a *AccountRepository) Create() (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), 1500*time.Millisecond)
	defer cancel()

	result, err := db.Mongo.Collection.Account.InsertOne(ctx, a.Entity)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

// FindByID : id args accept interface{} or primitive.ObjectID make sure to convert it first
func (a *AccountRepository) FindByID(id interface{}) (*entity.AccountBalance, error) {
	filter := bson.D{{"_id", id}}
	ctx, cancel := context.WithTimeout(context.TODO(), 1500*time.Millisecond)
	defer cancel()

	var account = new(entity.AccountBalance)
	err := db.Mongo.Collection.Account.FindOne(ctx, filter).Decode(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (a *AccountRepository) FindOne() (*entity.AccountBalance, error) {
	var account entity.AccountBalance

	ctx, cancel := context.WithTimeout(context.TODO(), 1500*time.Millisecond)
	defer cancel()

	err := db.Mongo.Collection.Account.FindOne(ctx, GetDefaultAccountFilter(a.Entity)).Decode(&account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

/*
FindAllPaginated
function args:

	*request.PaginatedAccountRequest

return:

	HttpStatusCode 	int,
	Account 			interfaces{},
	TotalDocument 	int64,
	TotalPages 		int,
	err 			error
*/
func (a *AccountRepository) FindAllPaginated(request *entity.PaginatedAccountRequest) (interface{}, int64, int64, error) {
	filter := GetDefaultAccountStatusFilter(request.Status)

	if request.Type > 0 {
		filter = append(filter, bson.D{{"type", request.Type}}...)
	}

	skipValue := (request.Page - 1) * request.Size

	ctx, cancel := context.WithTimeout(context.TODO(), 1500*time.Millisecond)
	defer cancel()

	cursor, err := db.Mongo.Collection.Account.Find(
		ctx,
		filter,
		options.Find().
			SetProjection(bson.D{{"secretKey", 0}, {"lastBalance", 0}}).
			SetSkip(skipValue).
			SetLimit(request.Size),
	)

	if err != nil {
		return nil, 0, 0, err
	}

	totalDocs, _ := db.Mongo.Collection.Account.CountDocuments(ctx, filter)
	var accounts []entity.AccountBalance
	if err = cursor.All(ctx, &accounts); err != nil {
		return nil, 0, 0, err
	}

	if len(accounts) == 0 {
		return nil, 0, 0, errors.New("empty results or last pages has been reached")
	}

	totalPages := math.Ceil(float64(totalDocs) / float64(request.Size))
	return &accounts, totalDocs, int64(totalPages), nil
}

func (a *AccountRepository) DeactivateAccount(payload *entity.UnregisterAccount) error {

	// update field
	update := bson.D{
		{"$set", bson.D{
			{"active", false},
			{"updatedAt", time.Now().UnixMilli()},
		}},
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 1500*time.Millisecond)
	defer cancel()

	result, err := db.Mongo.Collection.Account.UpdateOne(
		ctx,
		GetDefaultAccountFilter(&entity.AccountBalance{
			PartnerID:  payload.PartnerID,
			MerchantID: payload.MerchantID,
			TerminalID: payload.TerminalID,
		}), update)

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("update failed, cannot find account with current uniqueId")
	}

	return nil
}

func (a *AccountRepository) InsertDeactivatedAccount(account *entity.UnregisterAccount) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), 1500*time.Millisecond)
	defer cancel()

	result, err := db.Mongo.Collection.UnregisterAccount.InsertOne(ctx, account)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil

}

func (a *AccountRepository) RemoveDeactivatedAccount(acc *entity.UnregisterAccount) (int, error) {
	filter := bson.D{{"uniqueId", acc.UniqueID}}

	ctx, cancel := context.WithTimeout(context.TODO(), 1500*time.Millisecond)
	defer cancel()

	result, err := db.Mongo.Collection.UnregisterAccount.DeleteOne(ctx, filter)

	if err != nil {
		return fiber.StatusInternalServerError, err
	}

	if result.DeletedCount > 0 {
		return fiber.StatusNoContent, nil
	}

	return fiber.StatusInternalServerError, errors.New("remove deactivated account failed, no document found")
}

// ----------------- MERCHANTS ----------------

func (a *AccountRepository) FindMembersPaginated(request *entity.PaginatedAccountRequest, isPeriod bool) (interface{}, int64, int64, error) {
	filter := GetDefaultAccountStatusFilter(request.Status)

	filter = append(filter, bson.D{
		{"partnerId", request.PartnerID},
		{"merchantId", request.MerchantID},
		{"type", request.Type},
	}...)

	if isPeriod {
		filter = append(filter,
			bson.D{
				{"createdAt", bson.D{
					{"$gte", request.Periods.StartDate.UnixMilli()},
					{"$lte", request.Periods.EndDate.UnixMilli()},
				}},
			}...)
	}

	skipValue := (request.Page - 1) * request.Size

	ctx, cancel := context.WithTimeout(context.TODO(), 1500*time.Millisecond)
	defer cancel()

	cursor, err := db.Mongo.Collection.Account.Find(
		ctx,
		filter,
		options.Find().
			SetProjection(bson.D{{"secretKey", 0}, {"lastBalance", 0}}).
			SetSkip(skipValue).
			SetLimit(request.Size),
	)

	if err != nil {
		return nil, 0, 0, err
	}

	totalDocs, _ := db.Mongo.Collection.Account.CountDocuments(ctx, filter)
	var accounts []entity.AccountBalance
	if err = cursor.All(ctx, &accounts); err != nil {
		return nil, 0, 0, err
	}

	if accounts == nil {
		return nil, 0, 0, errors.New("empty results or last pages has been reached")
	}

	totalPages := math.Ceil(float64(totalDocs) / float64(request.Size))
	return &accounts, totalDocs, int64(totalPages), nil
}

func (a *AccountRepository) FindMembers(request *entity.PaginatedAccountRequest) ([]entity.AccountBalance, error) {
	filter := GetDefaultAccountStatusFilter(request.Status)

	filter = append(filter, bson.D{
		{"partnerId", request.PartnerID},
		{"merchantId", request.MerchantID},
		{"type", request.Type},
	}...)

	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	cursor, err := db.Mongo.Collection.Account.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var accounts []entity.AccountBalance
	if err = cursor.All(ctx, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (a *AccountRepository) CountMembers() (int64, error) {
	filter := bson.D{
		{"partnerId", a.Entity.PartnerID},
		{"merchantId", a.Entity.MerchantID},
		{"type", utilities.AccountTypeRegular},
		{"active", true},
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 1500*time.Millisecond)
	defer cancel()

	totalDocs, err := db.Mongo.Collection.Account.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return totalDocs, nil
}
