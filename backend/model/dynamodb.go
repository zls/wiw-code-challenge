package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
)

// Global DybamoDB connection
// var DDB *dynamodb.DynamoDB

// Get a new DynamoDB connection
// todo<zls>: read config from env and/or config files
func NewDynamoDB() gin.HandlerFunc {
	config := aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:8001"),
	}
	session := session.Must(session.NewSession(&config))
	ddb := dynamodb.New(session)
	return func(c *gin.Context) {
		c.Set("DDB", ddb)
		c.Next()
	}
}

func DynamoDBFromContext(c *gin.Context) *dynamodb.DynamoDB {
	return c.MustGet("DDB").(*dynamodb.DynamoDB)
}
