package model

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const ddbTableName = "shifts"

func ScanShifts(c *gin.Context) ([]Shift, error) {
	ddb := DynamoDBFromContext(c)
	results, err := ddb.Scan(&dynamodb.ScanInput{
		TableName: aws.String(ddbTableName),
	})
	if err != nil {
		c.Error(err)
		return nil, err
	}
	shifts := []Shift{}
	err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &shifts)
	if err != nil {
		c.Error(err)
		return nil, err
	}

	log.Printf("%+v", shifts)

	return shifts, nil
}

func GetShifts(c *gin.Context) (*[]Shift, error) {
	session := SessionFromContext(c)
	ddb := DynamoDBFromContext(c)

	userID := c.Query("userID")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	filterSet := false

	exprBuilder := expression.NewBuilder()

	// key condition
	keyCond := expression.Key("AccountID").Equal(expression.Value(session.Account.ID))
	if len(startTime) > 0 {
		log.Printf("found startTime, %v", startTime)
		if len(endTime) > 0 {
			keyCond = keyCond.And(expression.Key("StartTime").Between(expression.Value(startTime), expression.Value(endTime)))
		} else {
			keyCond = keyCond.And(expression.Key("StartTime").GreaterThanEqual(expression.Value(startTime)))
		}
	}
	exprBuilder = exprBuilder.WithKeyCondition(keyCond)

	// filter expression
	var filterExpr expression.ConditionBuilder
	if len(endTime) > 0 {
		log.Printf("found endTime, %v", startTime)
		filterExpr = expression.Name("EndTime").Between(expression.Value(startTime), expression.Value(endTime))
		filterSet = true
	}
	if len(userID) > 0 {
		log.Printf("found userID, %v", userID)
		filterExpr = expression.Name("EndTime").Between(expression.Value(startTime), expression.Value(endTime))
		if filterSet {
			filterExpr.And(expression.Name("UserID").Equal(expression.Value(userID)))
		} else {
			filterExpr = expression.Name("UserID").Equal(expression.Value(userID))
		}
		filterSet = true
	}
	if filterSet == true {
		log.Printf("filter expression set")
		exprBuilder = exprBuilder.WithFilter(filterExpr)
	}

	expr, err := exprBuilder.Build()
	if err != nil {
		c.Error(fmt.Errorf("failed to build query, %v", err))
		return nil, err
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(ddbTableName),
		IndexName:                 aws.String("ByAccount"),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
	}
	log.Printf("%+v", queryInput)
	results, err := ddb.Query(queryInput)
	if err != nil {
		c.Error(fmt.Errorf("failed to find shift in db, %v", err))
		return nil, err
	}

	shifts := []Shift{}
	err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &shifts)
	if err != nil {
		c.Error(fmt.Errorf("failed to unmarshal results to shift, %v", err))
	}
	return &shifts, nil
}

func GetShiftByID(c *gin.Context, id string) (*Shift, error) {
	idUUID, err := uuid.Parse(id)
	idBytes, err := idUUID.MarshalBinary()
	if err != nil {
		c.Error(fmt.Errorf("failed to marshal uuid to bytes, %v", err))
		return nil, err
	}
	shift := Shift{}
	ddb := DynamoDBFromContext(c)
	keyCond := expression.Key("ID").Equal(expression.Value(idBytes))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		c.Error(fmt.Errorf("failed to build keycondition, %v", err))
		return nil, err
	}
	results, err := ddb.Query(&dynamodb.QueryInput{
		TableName:                 aws.String(ddbTableName),
		IndexName:                 aws.String("ByID"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})

	// TODO <zls>: see if there is a constraint on ID to ensure uniqueness
	// TODO <zls>: split this out to more easily catch
	if err != nil || len(results.Items) > 1 {
		c.Error(fmt.Errorf("failed to find shift in db, %v", err))
		return nil, err
	}
	err = dynamodbattribute.UnmarshalMap(results.Items[0], &shift)
	if err != nil {
		c.Error(fmt.Errorf("failed to unmarshal results to shift, %v", err))
	}
	return &shift, nil
}

type Shift struct {
	ID        uuid.UUID `json:"id"`
	UserID    int       `form:"uid"`
	AccountID int       `form:"aid"`
	StartTime time.Time `form:"startTime" time_format:"2006-01-02T15:04:05"`
	EndTime   time.Time `form:"endTime" time_format:"2006-01-02T15:04:05"`
}

// Create a new shift structure
func NewShift(uid int, aid int, start time.Time, end time.Time) (*Shift, error) {
	sid, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.New("Failed to get new Shift UUID")
	}
	shift := &Shift{
		ID:        sid,
		UserID:    uid,
		AccountID: aid,
		StartTime: start,
		EndTime:   end,
	}
	return shift, nil
}

func (s *Shift) Overlaps(otherShift *Shift) bool {
	if s.UserID != otherShift.UserID {
		return false
	}
	if (s.StartTime.Unix() >= otherShift.StartTime.Unix() && s.StartTime.Unix() <= otherShift.EndTime.Unix()) ||
		(s.EndTime.Unix() <= otherShift.EndTime.Unix() && s.EndTime.Unix() >= otherShift.StartTime.Unix()) {
		return true
	}
	return false
}

func (s *Shift) Put(c *gin.Context) (*dynamodb.PutItemOutput, error) {
	av, err := dynamodbattribute.MarshalMap(*s)
	if err != nil {
		c.Error(fmt.Errorf("failed to marshal attr map %v", err.Error()))
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(ddbTableName),
	}
	ddb := DynamoDBFromContext(c)
	output, err := ddb.PutItem(input)
	if err != nil {
		c.Error(fmt.Errorf("failed to put item %v", err.Error()))
		return nil, err
	}
	return output, nil
}
