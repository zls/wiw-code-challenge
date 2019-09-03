package model

import (
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Shift struct {
	ID        uuid.UUID
	UserID    int
	AccountID int
	StartTime time.Time
	EndTime   time.Time
}

// Create a new shift structure
func NewShift(uid int, pid int, aid int, start time.Time, end time.Time) (*Shift, error) {
	sid, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err.Error())
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
		log.Fatal(err.Error())
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("shifts"),
	}
	ddb := DynamoDBFromContext(c)
	output, err := ddb.PutItem(input)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	return output, nil
}
