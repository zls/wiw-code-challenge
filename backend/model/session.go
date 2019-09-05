package model

import (
	"github.com/gin-gonic/gin"
)

type MockSesssion struct {
	User    User
	Account Account
}

func NewMockSession() gin.HandlerFunc {
	session := MockSesssion{
		User: User{
			ID:        1,
			FirstName: "Foo",
			LastName:  "Bar",
			Email:     "foo.bar@example.com",
		},
		Account: Account{
			ID:   1,
			Name: "Foo Company",
		},
	}
	return func(c *gin.Context) {
		c.Set("session", &session)
		c.Next()
	}
}

func SessionFromContext(c *gin.Context) *MockSesssion {
	return c.MustGet("session").(*MockSesssion)
}
