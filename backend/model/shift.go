package model

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
