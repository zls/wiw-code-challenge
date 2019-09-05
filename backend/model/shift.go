package model

import (
	"errors"
	"fmt"
	"log"
	"strconv"
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

func GetShifts(c *gin.Context, userID string, startTime string, endTime string) (*[]Shift, error) {
	session := SessionFromContext(c)
	ddb := DynamoDBFromContext(c)

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
		log.Printf("found endTime, %v", endTime)
		filterExpr = expression.Name("EndTime").Between(expression.Value(startTime), expression.Value(endTime))
		filterSet = true
	}
	if len(userID) > 0 {
		uid, err := strconv.Atoi(userID)
		if err != nil {
			c.Error(fmt.Errorf("failed to convert userid to a string, %v", err))
			return nil, err
		}
		log.Printf("found userID, %v", uid)
		filterExpr = expression.Name("EndTime").Between(expression.Value(startTime), expression.Value(endTime))
		if filterSet {
			filterExpr.And(expression.Name("UserID").Equal(expression.Value(uid)))
		} else {
			filterExpr = expression.Name("UserID").Equal(expression.Value(uid))
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

// Return a shift by its uuid
// TODO <zls>: restrict to account
func GetShiftByID(c *gin.Context, idBytes []byte) (*Shift, error) {
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

	if err != nil {
		c.Error(fmt.Errorf("failed to find shift in db, %v", err))
		return nil, err
	}
	err = dynamodbattribute.UnmarshalMap(results.Items[0], &shift)
	if err != nil {
		c.Error(fmt.Errorf("failed to unmarshal results to shift, %v", err))
	}
	return &shift, nil
}

func DeleteShiftByID(c *gin.Context, shiftId []byte, userID int, startTime string) error {
	ddb := DynamoDBFromContext(c)
	// session := SessionFromContext(c)

	// exprBuilder := expression.NewBuilder()
	// cond := expression.And(
	// 	expression.Name("AccountID").Equal(expression.Value(session.Account.ID)),
	// 	expression.Name("ID").Equal(expression.Value(shiftId)))
	// exprBuilder = exprBuilder.WithCondition(cond)
	// expr, err := exprBuilder.Build()
	// if err != nil {
	// 	c.Error(fmt.Errorf("failed to build expression, %v", err))
	// 	return err
	// }
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(ddbTableName),
		// ConditionExpression:       expr.Condition(),
		// ExpressionAttributeNames:  expr.Names(),
		// ExpressionAttributeValues: expr.Values(),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				N: aws.String(strconv.Itoa(userID)),
			},
			"StartTime": {
				S: aws.String(startTime),
			},
		},
	}
	log.Printf("%+v", input)

	output, err := ddb.DeleteItem(input)
	if err != nil {
		c.Error(fmt.Errorf("failed to delete shift %v", err))
		return err
	}
	log.Printf("%+v", output)
	return nil

}

type Shift struct {
	ID        uuid.UUID
	UserID    int `form:"userID" binding:"required"`
	AccountID int
	StartTime time.Time `form:"startTime" time_format:"2006-01-02T15:04:05" binding:"required"`
	EndTime   time.Time `form:"endTime" time_format:"2006-01-02T15:04:05" binding:"required"`
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
	if !shift.isLinearTime() {
		return nil, errors.New("Start time is after end time")
	}
	return shift, nil
}

// Check shift does not require a time traveller
func (s *Shift) isLinearTime() bool {
	return s.StartTime.Before(s.EndTime)
}

func (s *Shift) Put(c *gin.Context) (*dynamodb.PutItemOutput, error) {
	av, err := dynamodbattribute.MarshalMap(*s)
	if err != nil {
		c.Error(fmt.Errorf("failed to marshal attr map %v", err.Error()))
		return nil, err
	}
	exprBuilder := expression.NewBuilder()
	cond := expression.Not(expression.Name("UserID").Equal(expression.Value(s.UserID)))

	// Start time does not fall between another start and end time
	cond1 := expression.Not(expression.Name("StartTime").Between(expression.Value(s.StartTime), expression.Value(s.EndTime)))
	// End time does not fall between another start and end time
	cond2 := cond1.And(expression.Not(expression.Name("EndTime").Between(expression.Value(s.StartTime), expression.Value(s.EndTime))))

	cond3 := cond.And(cond2)

	// For user
	exprBuilder = exprBuilder.WithCondition(cond3)

	expr, err := exprBuilder.Build()
	if err != nil {
		c.Error(fmt.Errorf("failed to build expression, %v", err))
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:                      av,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 aws.String(ddbTableName),
	}
	log.Printf("%+v", input)

	ddb := DynamoDBFromContext(c)
	output, err := ddb.PutItem(input)
	if err != nil {
		c.Error(fmt.Errorf("failed to put item %v", err.Error()))
		return nil, err
	}
	log.Printf("%+v", output)
	return output, nil
}
