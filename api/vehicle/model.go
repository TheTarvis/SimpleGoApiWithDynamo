package vehicle

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"herbie/api/db"
	"time"
)

type Vehicle struct {
	AccountId    string     `json:"-" dynamodbav:"primary_key"`
	SK           string     `json:"-" dynamodbav:"sort_key"`
	ID           string     `json:"id" dynamodbav:"vehicle_id"`
	Make         *string    `json:"make" dynamodbav:"make,omitempty"`
	Model        *string    `json:"model"  dynamodbav:"model,omitempty"`
	Year         *int       `json:"year"  dynamodbav:"year,omitempty"`
	Color        *string    `json:"color"  dynamodbav:"color,omitempty"`
	LicensePlate *string    `json:"license_plate"  dynamodbav:"license_plate,omitempty"`
	CreatedOn    *time.Time `json:"created_on" dynamodbav:"created_on,omitempty"`
	UpdatedOn    *time.Time `json:"updated_on" dynamodbav:"updated_on,omitempty"`
}

func (v *Vehicle) NewVehicle(accountId string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	var now = time.Now()
	v.CreatedOn = &now
	v.UpdatedOn = &now
	v.AccountId = accountId
	v.ID = id.String()
	v.SK = getSortKey(v.ID)

	return nil
}

func (v *Vehicle) GetVehicleById(accountId string, id string) error {
	key := createVehicleDynamoKey(accountId, id)

	result, err := db.GetItem(key)
	if err != nil {
		log.Warnf("Error while getting item: %s", err)
		return err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &v)
	if err != nil {
		log.Warnf("Error while unmarshalling item: %s", err)
		return err
	}

	return nil
}

func (v *Vehicle) UpdateVehicle() error {
	key := createVehicleDynamoKey(v.AccountId, v.ID)
	update := expression.Set(
		expression.Name("make"),
		expression.Value(v.Make),
	).Set(
		expression.Name("model"),
		expression.Value(v.Model),
	).Set(
		expression.Name("year"),
		expression.Value(v.Year),
	).Set(
		expression.Name("color"),
		expression.Value(v.Color),
	).Set(
		expression.Name("license_plate"),
		expression.Value(v.LicensePlate),
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
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &v)
	if err != nil {
		log.Warnf("Error while unmarshalling updated item. %s", err)
		return err
	}

	return nil
}

func (v *Vehicle) SaveVehicle() error {
	result, err := db.PutItem(v)
	if err != nil {
		return err
	}

	err = dynamodbattribute.UnmarshalMap(result.Attributes, &v)
	if err != nil {
		log.Warnf("Error while unmarshalling updated item. %s", err)
		return err
	}

	return nil
}

func createVehicleDynamoKey(accountId string, vehicleId string) map[string]*dynamodb.AttributeValue {
	key := map[string]*dynamodb.AttributeValue{
		"primary_key": {
			S: aws.String(accountId),
		},
		"sort_key": {
			S: aws.String(getSortKey(vehicleId)),
		},
	}
	return key
}

func getSortKey(vehicleId string) string {
	return "vehicle#" + vehicleId
}