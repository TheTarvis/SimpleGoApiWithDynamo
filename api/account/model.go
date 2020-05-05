package account

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"herbie/api/db"
	"time"
)

type Account struct {
	ID        string     `json:"id" dynamodbav:"primary_key"`
	SK        string     `json:"-" dynamodbav:"sort_key,omitempty"`
	FirstName *string    `json:"first_name" dynamodbav:"first_name,omitempty"`
	LastName  *string    `json:"last_name"  dynamodbav:"last_name,omitempty"`
	CreatedOn *time.Time `json:"created_on" dynamodbav:"created_on,omitempty"`
	UpdatedOn *time.Time `json:"updated_on" dynamodbav:"updated_on,omitempty"`
}

func (a *Account) NewAccount() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	var now = time.Now()
	a.CreatedOn = &now
	a.UpdatedOn = &now
	a.ID = id.String()
	a.SK = "meta"

	return nil
}

func (a *Account) GetAccountById(id string) error {
	key := createAccountDynamoKey(id)

	result, err := db.GetItem(key)
	if err != nil {
		log.Warnf("Error while getting item: %s", err)
		return err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &a)
	if err != nil {
		log.Warnf("Error while unmarshalling item: %s", err)
		return err
	}

	return nil
}

func (a *Account) UpdateAccount() error {
	key := createAccountDynamoKey(a.ID)
	update := expression.Set(
		expression.Name("first_name"),
		expression.Value(a.FirstName),
	).Set(
		expression.Name("last_name"),
		expression.Value(a.LastName),
	).Set(
		expression.Name("updated_on"),
		expression.Value(time.Now().Format(time.RFC3339Nano)),
	)
	expr, err := expression.NewBuilder().
		WithUpdate(update).
		Build()

	result, err := db.UpdateItem(key, expr)

	if err != nil {
		log.Warnf("Error while updating item. %s", err)
		return err
	}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &a)
	if err != nil {
		log.Warnf("Error while unmarshalling updated item. %s", err)
		return err
	}

	return nil
}

func (a *Account) SaveAccount() error {
	result, err := db.PutItem(a)
	if err != nil {
		return err
	}

	err = dynamodbattribute.UnmarshalMap(result.Attributes, &a)
	if err != nil {
		log.Warnf("Error while unmarshalling updated item. %s", err)
		return err
	}

	return nil
}

func createAccountDynamoKey(id string) map[string]*dynamodb.AttributeValue {
	key := map[string]*dynamodb.AttributeValue{
		"primary_key": {
			S: aws.String(id),
		},
		"sort_key": {
			S: aws.String("meta"),
		},
	}
	return key
}
