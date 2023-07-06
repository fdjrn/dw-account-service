package repository

import (
	"context"
	"errors"
	"github.com/dw-account-service/internal/db"
	"github.com/dw-account-service/internal/db/entity"
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
	result, err := db.Mongo.Collection.Account.InsertOne(context.TODO(), a.Entity)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

// FindByID : id args accept interface{} or primitive.ObjectID make sure to convert it first
func (a *AccountRepository) FindByID(id interface{}) (*entity.AccountBalance, error) {
	filter := bson.D{{"_id", id}}
	var account = new(entity.AccountBalance)
	err := db.Mongo.Collection.Account.FindOne(context.TODO(), filter).Decode(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (a *AccountRepository) FindOne() (*entity.AccountBalance, error) {
	var account entity.AccountBalance
	err := db.Mongo.Collection.Account.FindOne(context.TODO(), GetDefaultAccountFilter(a.Entity)).Decode(&account)
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
	Data 			interfaces{},
	TotalDocument 	int64,
	TotalPages 		int,
	err 			error
*/
func (a *AccountRepository) FindAllPaginated(request *entity.PaginatedAccountRequest) (interface{}, int64, int64, error) {
	var filter = bson.D{}
	switch request.Status {
	case AccountStatusActive:
		filter = bson.D{{"active", true}}
	case AccountStatusDeactivated:
		filter = bson.D{{"active", false}}
	default:
	}

	if request.Type > 0 {
		filter = append(filter, bson.D{{"type", request.Type}}...)
	}

	skipValue := (request.Page - 1) * request.Size

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
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
	if err = cursor.All(context.TODO(), &accounts); err != nil {
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

	result, err := db.Mongo.Collection.Account.UpdateOne(
		context.TODO(),
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

	result, err := db.Mongo.Collection.UnregisterAccount.InsertOne(context.TODO(), account)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil

}

func (a *AccountRepository) RemoveDeactivatedAccount(acc *entity.UnregisterAccount) (int, error) {
	filter := bson.D{{"uniqueId", acc.UniqueID}}
	result, err := db.Mongo.Collection.UnregisterAccount.DeleteOne(context.TODO(), filter)

	if err != nil {
		return fiber.StatusInternalServerError, err
	}

	if result.DeletedCount > 0 {
		return fiber.StatusNoContent, nil
	}

	return fiber.StatusInternalServerError, errors.New("remove deactivated account failed, no document found")
}

// ----------------- MERCHANTS ----------------

func (a *AccountRepository) FindMembers(request *entity.PaginatedAccountRequest, isPeriod bool) (interface{}, int64, int64, error) {
	filter := bson.D{}
	switch request.Status {
	case AccountStatusActive:
		filter = bson.D{{"active", true}}
	case AccountStatusDeactivated:
		filter = bson.D{{"active", false}}
	default:
	}

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

	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
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
	if err = cursor.All(context.TODO(), &accounts); err != nil {
		return nil, 0, 0, err
	}

	if accounts == nil {
		return nil, 0, 0, errors.New("empty results or last pages has been reached")
	}

	totalPages := math.Ceil(float64(totalDocs) / float64(request.Size))
	return &accounts, totalDocs, int64(totalPages), nil
}

// -----------------  UNUSED  ----------------

//func (a *AccountRepository) FindAllMerchant(request *entity.PaginatedAccountRequest) (int, interface{}, int64, int64, error) {
//	//var filter interface{}
//
//	filter := bson.D{}
//	switch request.Status {
//	case AccountStatusActive:
//		filter = append(filter, bson.D{{"active", true}}...)
//	case AccountStatusDeactivated:
//		filter = append(filter, bson.D{{"active", false}}...)
//	default:
//	}
//
//	filter = append(filter, bson.D{{"type", 2}}...)
//
//	skipValue := (request.Page - 1) * request.Size
//
//	ctx, cancel := context.WithTimeout(context.TODO(), 500*time.Millisecond)
//	defer cancel()
//
//	cursor, err := db.Mongo.Collection.Account.Find(
//		ctx,
//		filter,
//		options.Find().
//			SetProjection(bson.D{{
//				"secretKey", 0},
//			//{"lastBalance", 0},
//			}).
//			SetSkip(skipValue).
//			SetLimit(request.Size),
//	)
//
//	if err != nil {
//		return fiber.StatusInternalServerError, nil, 0, 0, err
//	}
//
//	totalDocs, _ := db.Mongo.Collection.Account.CountDocuments(ctx, filter)
//	var accounts []entity.AccountBalance
//	if err = cursor.All(context.TODO(), &accounts); err != nil {
//		return fiber.StatusInternalServerError, nil, 0, 0, err
//	}
//
//	if len(accounts) == 0 {
//		return fiber.StatusInternalServerError, nil, 0, 0, errors.New("empty results or last pages has been reached")
//	}
//
//	totalPages := math.Ceil(float64(totalDocs) / float64(request.Size))
//	return fiber.StatusOK, &accounts, totalDocs, int64(totalPages), nil
//}

//func (a *AccountRepository) FindByMerchantID(ab *entity.AccountBalance) (int, interface{}, error) {
//
//	filter := bson.D{
//		{"merchantId", ab.MerchantID},
//		{"partnerId", ab.PartnerID},
//		{"terminalId", ab.TerminalID},
//		{"type", AccountTypeMerchant},
//	}
//
//	account, err := a.findByFilter(filter)
//
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return fiber.StatusNotFound, nil, err
//		}
//		return fiber.StatusInternalServerError, nil, err
//	}
//
//	return fiber.StatusOK, account, nil
//}

//func (a *AccountRepository) FindByMerchantStatus(ab *entity.AccountBalance, active bool) (int, interface{}, error) {
//
//	filter := bson.D{
//		{"merchantId", ab.MerchantID},
//		{"partnerId", ab.PartnerID},
//	}
//
//	if active {
//		filter = append(filter, bson.D{{"active", true}}...)
//	}
//
//	account, err := a.findByFilter(filter)
//
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return fiber.StatusNotFound, nil, err
//		}
//		return fiber.StatusInternalServerError, nil, err
//	}
//
//	return fiber.StatusOK, account, nil
//}

//func (a *AccountRepository) CountMerchantMember(account entity.AccountBalance) (int64, error) {
//	filter := bson.D{
//		{"partnerId", account.PartnerID},
//		{"merchantId", account.MerchantID},
//		{"active", account.Active},
//	}
//
//	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
//	defer cancel()
//
//	totalDocs, err := db.Mongo.Collection.Account.CountDocuments(ctx, filter)
//	if err != nil {
//		return 0, err
//	}
//
//	return totalDocs, nil
//}

//func (a *AccountRepository) FindMemberByMerchant(account entity.AccountBalance) (int64, interface{}, error) {
//	filter := bson.D{
//		{"partnerId", account.PartnerID},
//		{"merchantId", account.MerchantID},
//		{"active", account.Active},
//	}
//
//	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
//	defer cancel()
//
//	cursor, err := db.Mongo.Collection.Account.Find(ctx, filter)
//
//	if err != nil {
//		return 0, nil, err
//	}
//
//	var accounts []entity.AccountBalance
//	if err = cursor.All(context.TODO(), &accounts); err != nil {
//		return 0, nil, err
//	}
//
//	totalDocs, _ := db.Mongo.Collection.Account.CountDocuments(ctx, filter)
//
//	if len(accounts) == 0 {
//		return 0, nil, errors.New("empty results or last pages has been reached")
//	}
//
//	return totalDocs, &accounts, nil
//}

//func (a *AccountRepository) findByFilter(filter interface{}) (interface{}, error) {
//	account := new(entity.AccountBalance)
//
//	// filter condition
//	err := db.Mongo.Collection.Account.FindOne(context.TODO(), filter).Decode(&account)
//	if err != nil {
//		return nil, err
//	}
//
//	return account, nil
//}

//func (a *AccountRepository) FindOne(account *entity.AccountBalance) (interface{}, error) {
//
//	result, err := a.findByFilter(GetDefaultAccountFilter(account))
//
//	if err != nil {
//		//if err == mongo.ErrNoDocuments {
//		//	return nil, err
//		//}
//		return nil, err
//	}
//
//	return result, nil
//}

//func (a *AccountRepository) FindByUniqueID(id string, active bool) (int, interface{}, error) {
//
//	filter := bson.D{{"uniqueId", id}}
//	if active {
//		filter = bson.D{{"uniqueId", id}, {"active", true}}
//	}
//
//	account, err := a.findByFilter(filter)
//
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			//return fiber.StatusNotFound, nil, errors.New("account not found or it has been unregistered")
//			return fiber.StatusNotFound, nil, err
//		}
//		return fiber.StatusInternalServerError, nil, err
//	}
//
//	return fiber.StatusOK, account, nil
//}

//func (a *AccountRepository) FindByActiveStatus(id string, status bool) (int, interface{}, error) {
//	//id, _ := primitive.ObjectIDFromHex(accountId)
//
//	account, err := a.findByFilter(bson.D{
//		{Key: "uniqueId", Value: id},
//		{Key: "active", Value: status},
//	})
//
//	if err != nil {
//		return fiber.StatusInternalServerError, nil, err
//	}
//
//	return fiber.StatusOK, account, nil
//}

//func (a *AccountRepository) DeactivateMerchant(u *entity.UnregisterAccount) (int, error) {
//	//id, _ := primitive.ObjectIDFromHex(u.UniqueID)
//
//	// filter condition
//	filter := bson.D{{"merchantId", u.MerchantID}, {"partnerId", u.PartnerID}}
//
//	// update field
//	update := bson.D{
//		{"$set", bson.D{
//			{"active", false},
//			{"updatedAt", time.Now().UnixMilli()},
//		}},
//	}
//
//	result, err := db.Mongo.Collection.Account.UpdateOne(context.TODO(), filter, update)
//	if err != nil {
//		return fiber.StatusInternalServerError, err
//	}
//
//	if result.ModifiedCount == 0 {
//		return fiber.StatusBadRequest, errors.New(
//			"update failed, cannot find account with current uniqueId")
//	}
//
//	return fiber.StatusOK, nil
//}
