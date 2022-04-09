package services_test

import (
	"testing"
	"time"

	"github.com/NikosSiak/Open-Banking-API/lib"
	"github.com/NikosSiak/Open-Banking-API/services"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/suite"
)

type NewAuthServiceSuite struct {
	suite.Suite
}

func (suite *NewAuthServiceSuite) TestNewAuthService() {
	jwtSecret := "test_secret"
	env := lib.Env{JWTSecret: jwtSecret}

	suite.NotPanics(func() { services.NewAuthService(env) }, "new auth service did panic for no reason")
}

func (suite *NewAuthServiceSuite) TestNewAuthServiceWithoutJwtSecret() {
	env := lib.Env{}

	suite.Panics(func() { services.NewAuthService(env) }, "new auth service did not panic without jwt secret")
}

func TestNewAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(NewAuthServiceSuite))
}

type CreateTokensSuite struct {
	suite.Suite
	jwtSecret   []byte
	authService services.AuthService
}

func (suite *CreateTokensSuite) SetupTest() {
	secret := "test_secret"

	suite.jwtSecret = []byte(secret)
	suite.authService = services.NewAuthService(lib.Env{JWTSecret: secret})
}

func (suite *CreateTokensSuite) TestCreatesTokenDetails() {

}

func TestCreateTokensSuite(t *testing.T) {
	suite.Run(t, new(CreateTokensSuite))
}

type VerifyTokenSuite struct {
	suite.Suite
	jwtSecret   []byte
	authService services.AuthService
}

func (suite *VerifyTokenSuite) SetupTest() {
	secret := "test_secret"

	suite.jwtSecret = []byte(secret)
	suite.authService = services.NewAuthService(lib.Env{JWTSecret: secret})
}

func (suite *VerifyTokenSuite) TestAValidToken() {
	expiresAt := time.Now().Add(time.Minute * 15).Unix()
	textToEncode := "random token"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"access_uuid": textToEncode,
		"exp":         expiresAt,
	})
	signedToken, _ := token.SignedString(suite.jwtSecret)

	resToken, err := suite.authService.VerifyToken(signedToken)

	suite.NotNil(resToken)
	suite.Nil(err)

	claims, ok := resToken.Claims.(jwt.MapClaims)
	suite.True(ok, "token was incorrectly flagged as invalid")
	suite.Equal(textToEncode, claims["access_uuid"].(string), "access uuid did not match the original text")
}

func (suite *VerifyTokenSuite) TestAnExpiredToken() {
	expiresAt := time.Now().Unix()
	textToEncode := "random token"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"access_uuid": textToEncode,
		"exp":         expiresAt,
	})
	signedToken, _ := token.SignedString(suite.jwtSecret)

	resToken, err := suite.authService.VerifyToken(signedToken)

	suite.Nil(resToken)
	suite.EqualError(err, "Token is expired")
}

func TestVerifyTokenSuite(t *testing.T) {
	suite.Run(t, new(VerifyTokenSuite))
}

type GetAccessUuidSuite struct {
	suite.Suite
	jwtSecret   []byte
	authService services.AuthService
}

func (suite *GetAccessUuidSuite) SetupTest() {
	secret := "test_secret"

	suite.jwtSecret = []byte(secret)
	suite.authService = services.NewAuthService(lib.Env{JWTSecret: secret})
}

func (suite *GetAccessUuidSuite) TestAValidToken() {
	expiresAt := time.Now().Add(time.Minute * 15).Unix()
	textToEncode := "random token"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"access_uuid": textToEncode,
		"exp":         expiresAt,
	})
	signedToken, _ := token.SignedString(suite.jwtSecret)

	accessUuid, err := suite.authService.GetAccessUuid(signedToken)

	suite.Equal(textToEncode, accessUuid, "access uuid does not match original string")
	suite.Nil(err)
}

func (suite *GetAccessUuidSuite) TestATokenWithoutAccessUuid() {
	expiresAt := time.Now().Add(time.Minute * 15).Unix()
	textToEncode := "random token"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"random_claim": textToEncode,
		"exp":          expiresAt,
	})
	signedToken, _ := token.SignedString(suite.jwtSecret)

	accessUuid, err := suite.authService.GetAccessUuid(signedToken)

	suite.Equal("", accessUuid)
	suite.EqualError(err, "no access uuid")
}

func TestGetAccessUuidSuite(t *testing.T) {
	suite.Run(t, new(GetAccessUuidSuite))
}
