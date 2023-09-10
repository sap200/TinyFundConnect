package db

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/sap200/TinyFundConnect/secret"
	"github.com/sap200/TinyFundConnect/types"
	"google.golang.org/api/option"
)

type DB interface {
	init()
	Save(collectionName string, record interface{}) (bool, error)
	DoesRecordExists(collectionName string, key string) bool
	GetRecordDetails(collectionName string, key string) (interface{}, error)
}

type Repo struct {
	Client *firestore.Client
}

func New() *Repo {

	r := Repo{}
	return &r
}

func (r *Repo) init() {

	opt := option.WithCredentialsFile("./secret/firebase_secret.json")

	// Access the Database
	client, err := firestore.NewClient(context.Background(), secret.PROJECT_ID, opt)
	if err != nil {
		log.Fatalf("Error accessing Realtime Database: %v", err)
	}

	// Now you can use the 'client' to interact with the Realtime Database
	fmt.Println("Connected to Firebase Realtime Database")

	r.Client = client
}

func (r *Repo) Save(collectionName string, key string, m map[string]interface{}) {
	ctx := context.Background()
	// initialize that is open connection
	r.init()
	// Close the connection when done
	defer r.Client.Close()

	_, err := r.Client.Collection(collectionName).Doc(key).Set(ctx, m)
	if err != nil {
		fmt.Println(err)
	}

}

func (r *Repo) DoesRecordExists(collectionName, key string) (bool, error) {
	ctx := context.Background()
	r.init()
	defer r.Client.Close()

	collectionRef := r.Client.Collection(collectionName)
	docs, err := collectionRef.Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}

	// equate these keys
	for _, doc := range docs {
		if doc.Ref.ID == key {
			return true, nil
		}
	}

	return false, nil
}

func (r *Repo) GetRecordDetails(collectionName string, key string) (interface{}, error) {
	ctx := context.Background()
	r.init()
	defer r.Client.Close()

	docRef := r.Client.Collection(collectionName).Doc(key)
	_, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Get the data from the document snapshot
	var user types.User
	if err := docSnapshot.DataTo(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repo) GetPoolDetails(collectionName string, key string) (interface{}, error) {
	ctx := context.Background()
	r.init()
	defer r.Client.Close()

	docRef := r.Client.Collection(collectionName).Doc(key)
	_, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Get the data from the document snapshot
	var pool types.Pool
	if err := docSnapshot.DataTo(&pool); err != nil {
		return nil, err
	}

	return pool, nil
}

func (r *Repo) MarkUserVerified(collectionName string, key string, fieldName string, value bool) (bool, error) {
	ctx := context.Background()
	r.init()
	defer r.Client.Close()

	docRef := r.Client.Collection(collectionName).Doc(key)
	_, err := docRef.Get(ctx)
	if err != nil {
		return false, err
	}

	// Update a specific field in the fetched document
	updateField := fieldName
	newValue := value
	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: updateField, Value: newValue},
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *Repo) GetAllPools(collectionName string) (*[]types.Pool, error) {
	ctx := context.Background()
	r.init()
	defer r.Client.Close()

	collectionRef := r.Client.Collection(collectionName)
	docs, err := collectionRef.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var poolRecords = []types.Pool{}
	for i := 0; i < len(docs); i++ {
		var p types.Pool
		docs[i].DataTo(&p)
		poolRecords = append(poolRecords, p)
	}

	return &poolRecords, nil

}

func (r *Repo) GetChatDetailsByPool(collectionName string, key string) (*types.Chats, error) {
	ctx := context.Background()
	r.init()
	defer r.Client.Close()

	docRef := r.Client.Collection(collectionName).Doc(key)
	_, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Get the data from the document snapshot
	var chatLogs types.Chats
	if err := docSnapshot.DataTo(&chatLogs); err != nil {
		return nil, err
	}

	return &chatLogs, nil
}

func (r *Repo) GetPollDetailsByPollId(collectionName string, key string) (*types.Poll, error) {
	ctx := context.Background()
	r.init()
	defer r.Client.Close()

	docRef := r.Client.Collection(collectionName).Doc(key)
	_, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Get the data from the document snapshot
	var polls types.Poll
	if err := docSnapshot.DataTo(&polls); err != nil {
		return nil, err
	}

	return &polls, nil
}

func (r *Repo) GetAllPollsByPoolId(collectionName string, key string) (*[]types.Poll, error) {
	ctx := context.Background()
	r.init()
	defer r.Client.Close()

	collectionRef := r.Client.Collection(collectionName)
	docs, err := collectionRef.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var poolRecords = []types.Poll{}
	for i := 0; i < len(docs); i++ {
		var p types.Poll
		docs[i].DataTo(&p)
		if p.PoolId == key {
			poolRecords = append(poolRecords, p)
		}
	}

	return &poolRecords, nil

}

func (r *Repo) GetAllTxnLogsByEmailId(collectionName string, key string, txnType string) (*[]types.Transaction, error) {
	ctx := context.Background()
	r.init()
	defer r.Client.Close()

	collectionRef := r.Client.Collection(collectionName)
	docs, err := collectionRef.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var txnLogs = []types.Transaction{}
	for i := 0; i < len(docs); i++ {
		var p types.Transaction
		docs[i].DataTo(&p)
		if p.UserEmail == key && p.TxnType == txnType {
			txnLogs = append(txnLogs, p)
		}
	}

	return &txnLogs, nil
}

func (r *Repo) GetOrderEarnsByPoolId(collectionName string, key string) (*map[string]types.OrderEarns, error) {
	ctx := context.Background()
	r.init()
	defer r.Client.Close()

	docRef := r.Client.Collection(collectionName).Doc(key)
	_, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Get the data from the document snapshot
	var oe map[string]types.OrderEarns
	if err := docSnapshot.DataTo(&oe); err != nil {
		return nil, err
	}

	return &oe, nil
}

func (r *Repo) GetAllTxnLogsByPoolId(collectionName string, key string, txnType string) (*[]types.Transaction, error) {
	ctx := context.Background()
	r.init()
	defer r.Client.Close()

	collectionRef := r.Client.Collection(collectionName)
	docs, err := collectionRef.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var txnLogs = []types.Transaction{}
	for i := 0; i < len(docs); i++ {
		var p types.Transaction
		docs[i].DataTo(&p)
		if p.PoolId == key && p.TxnType == txnType {
			txnLogs = append(txnLogs, p)
		}
	}

	return &txnLogs, nil
}
