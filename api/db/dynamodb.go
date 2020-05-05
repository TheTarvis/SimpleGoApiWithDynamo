package db

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var (
	tableName = getEnv("DBNAME", "herbie")
	sess      = session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           "herbie",
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc = dynamodb.New(sess)
)

func PutItem(item interface{}) (*dynamodb.PutItemOutput, error) {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Printf("Got error marshalling item: %s", err.Error())
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	result, err := svc.PutItem(input)
	if err != nil {
		fmt.Printf("Got error calling PutItem: %s", err.Error())
		return nil, err
	}

	return result, nil
}

func GetItem(key map[string]*dynamodb.AttributeValue) (*dynamodb.GetItemOutput, error) {
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return result, nil
}

func UpdateItem(key map[string]*dynamodb.AttributeValue, expr expression.Expression) (*dynamodb.UpdateItemOutput, error) {
	input := &dynamodb.UpdateItemInput{
		Key: key,
		UpdateExpression: expr.Update(),
		TableName:        aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:     aws.String(dynamodb.ReturnValueAllNew),
	}

	result, err := svc.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return result, nil
}

