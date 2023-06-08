package mongodb

import (
	"context"
	"github.com/dw-account-service/configs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MgoCollection struct {
	Account           *mongo.Collection
	UnregisterAccount *mongo.Collection
	BalanceTopup      *mongo.Collection
}

type MgoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

const (
	AccountCollection           = "accountBalances"
	UnregisterAccountCollection = "accountDeactivated"
	BalanceTopupCollection      = "balanceTopup"
)

var Instance MgoInstance
var Collection MgoCollection

func (i *MgoInstance) Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(configs.MainConfig.Database.Mongo.Uri).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)))

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	// Send a ping to confirm a successful conn
	if err = client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		return err
	}

	Instance = MgoInstance{
		Client: client,
		DB:     client.Database(configs.MainConfig.Database.Mongo.DBName),
	}

	Collection = MgoCollection{
		Account:           Instance.DB.Collection(AccountCollection),
		UnregisterAccount: Instance.DB.Collection(UnregisterAccountCollection),
		BalanceTopup:      Instance.DB.Collection(BalanceTopupCollection),
	}

	log.Println("[INIT] database >> connected")

	return nil
}

func (i *MgoInstance) Disconnect() error {
	if Instance.Client == nil {
		return nil
	}

	err := Instance.Client.Disconnect(context.TODO())
	if err != nil {
		log.Println("error on closing mongodb connection: ", err.Error())
		return err
	}
	return nil
}
