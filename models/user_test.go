package models_test

import (
	"testing"
	"time"

	"github.com/NikosSiak/Open-Banking-API/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHashPassword(t *testing.T) {
	password := "test password"
	user := models.User{Password: password}

	err := user.HashPassword()
	assert.Nil(t, err)

	assert.NotEqual(t, password, user.Password, "user password did not change")
}

func TestCheckPasswordHash(t *testing.T) {
	password := "test password"
	user := models.User{Password: password}
	user.HashPassword()

	assert.True(t, user.CheckPasswordHash(password), "hashed password did not match with the original password")
}

type GetBSONSuite struct {
	suite.Suite
	user          models.User
	userCreatedAt primitive.DateTime
	userUpdatedAt primitive.DateTime
}

func (suite *GetBSONSuite) SetupTest() {
	user := models.User{}

	suite.user = user
	suite.userCreatedAt = user.CreatedAt
	suite.userUpdatedAt = user.UpdatedAt
}

func (suite *GetBSONSuite) TestCreatedAtGetsSet() {
	userBson := suite.user.GetBSON()

	suite.NotEqual(suite.userCreatedAt, userBson.(*models.User).CreatedAt, "user created at time did not change")
}

func (suite *GetBSONSuite) TestCreatedAtDoesNotChange() {
	createdTime := primitive.NewDateTimeFromTime(time.Now().UTC())
	suite.user.CreatedAt = createdTime
	suite.userCreatedAt = createdTime

	userBson := suite.user.GetBSON()

	suite.Equal(suite.userCreatedAt, userBson.(*models.User).CreatedAt, "user created at time did change while it was already set")
}

func (suite *GetBSONSuite) TestUpdatedAtChanges() {
	userBson := suite.user.GetBSON()

	suite.NotEqual(suite.userUpdatedAt, userBson.(*models.User).UpdatedAt, "user updated at did not change")
}

func TestGetBSONSuite(t *testing.T) {
	suite.Run(t, new(GetBSONSuite))
}

type AddAccountSuite struct {
	suite.Suite
	user        models.User
	accountName string
	account     models.AccountTokens
}

func (suite *AddAccountSuite) SetupTest() {
	user := models.User{}

	suite.user = user
	suite.accountName = "testAccount"
	suite.account = models.AccountTokens{AccessToken: "access_token", RefreshToken: "refresh_token"}
}

func (suite *AddAccountSuite) TestAddNewAccount() {
	suite.user.AddAccount(suite.accountName, &suite.account)
	suite.Equal(&suite.account, suite.user.Accounts[suite.accountName], "account was not added")
}

func (suite *AddAccountSuite) TestAddMoreAccounts() {
	suite.user.Accounts = map[string]*models.AccountTokens{suite.accountName: &suite.account}

	secondAccountName := "testAccount2"
	secondAccount := models.AccountTokens{}

	suite.user.AddAccount(secondAccountName, &secondAccount)

	suite.NotNil(suite.user.Accounts[suite.accountName], "previous account was overwritten")
	suite.Equal(&secondAccount, suite.user.Accounts[secondAccountName], "second account was not added")
}

func TestAddAccountSuite(t *testing.T) {
	suite.Run(t, new(AddAccountSuite))
}
