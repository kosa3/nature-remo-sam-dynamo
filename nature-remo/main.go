package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/tenntenn/natureremo"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type NatureRemo struct {
	Type     string `json:"type"`
	DateTime string `json:"datetime"`
	Value    string `json:"value"`
}

func getValueFromNewsEventType(sensorValue natureremo.SensorValue) interface{} {
	bytes, err := json.Marshal(sensorValue)
	if err != nil {
		log.Fatal(err)
	}

	var m map[string]interface{}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		log.Fatal(err)
	}

	return m["val"]
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Environment variables
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	cli := natureremo.NewClient(os.Getenv("NATURE_REMO_ACCESS_TOKEN"))
	ctx := context.Background()

	sess := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion("ap-northeast-1")
	if len(endpoint) > 0 {
		config = config.WithEndpoint(endpoint)
	}

	db := dynamo.New(sess, config)
	table := db.Table(tableName)

	switch request.HTTPMethod {
	case "GET":
		var results []NatureRemo
		requestType := string(request.PathParameters["Type"])
		err := table.Get("Type", requestType).All(&results)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		bytes, err := json.Marshal(results)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{
			Body:       string(bytes),
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "origin,Accept,Authorization,Content-Type",
				"Content-Type":                 "application/json",
			},
		}, nil
	case "POST":
		// デバイスデータを取得する
		devices, err := cli.DeviceService.GetAll(ctx)
		if err != nil {
			log.Fatal(err)
		}

		var DateTime = strconv.FormatInt(time.Now().Unix(), 10)
		// DynamoDBにデバイスデータを送信する
		for _, device := range devices {
			_, err := dynamodb.New(sess, config).TransactWriteItems(&dynamodb.TransactWriteItemsInput{
				TransactItems: []*dynamodb.TransactWriteItem{
					{
						Put: &dynamodb.Put{
							TableName: aws.String(tableName),
							Item: map[string]*dynamodb.AttributeValue{
								"Type": {
									S: aws.String("hu"),
								},
								"DateTime": {
									S: aws.String(DateTime),
								},
								"Value": {
									S: aws.String(fmt.Sprintf("%v", getValueFromNewsEventType(device.NewestEvents["hu"]))),
								},
							},
						},
					},
					{
						Put: &dynamodb.Put{
							TableName: aws.String(tableName),
							Item: map[string]*dynamodb.AttributeValue{
								"Type": {
									S: aws.String("il"),
								},
								"DateTime": {
									S: aws.String(DateTime),
								},
								"Value": {
									S: aws.String(fmt.Sprintf("%v", getValueFromNewsEventType(device.NewestEvents["il"]))),
								},
							},
						},
					},
					{
						Put: &dynamodb.Put{
							TableName: aws.String(tableName),
							Item: map[string]*dynamodb.AttributeValue{
								"Type": {
									S: aws.String("mo"),
								},
								"DateTime": {
									S: aws.String(DateTime),
								},
								"Value": {
									S: aws.String(fmt.Sprintf("%v", getValueFromNewsEventType(device.NewestEvents["mo"]))),
								},
							},
						},
					},
					{
						Put: &dynamodb.Put{
							TableName: aws.String(tableName),
							Item: map[string]*dynamodb.AttributeValue{
								"Type": {
									S: aws.String("te"),
								},
								"DateTime": {
									S: aws.String(DateTime),
								},
								"Value": {
									S: aws.String(fmt.Sprintf("%v", getValueFromNewsEventType(device.NewestEvents["te"]))),
								},
							},
						},
					},
				},
			})

			if err != nil {
				log.Fatal(err)
			}
		}

		return events.APIGatewayProxyResponse{
			Body:       string(200),
			StatusCode: http.StatusOK,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string("Not Support Method"),
		StatusCode: http.StatusBadRequest,
	}, nil
}

func main() {
	lambda.Start(handler)
}
