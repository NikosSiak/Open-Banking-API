package services

import (
  "errors"
  "time"

  "github.com/dgrijalva/jwt-go"
  "github.com/google/uuid"
  "github.com/NikosSiak/Open-Banking-API/lib"
)

type AuthService struct {
  jwtSecret []byte
}

type TokenDetails struct {
  AccessToken string
  RefreshToken string
  AccessUuid string
  RefreshUuid string
  AtExpires int64
  RtExpires int64
}

func NewAuthService(env lib.Env) AuthService {
  jwtSecret := env.JWTSecret
  if jwtSecret == "" {
    panic("missing jwt secret")
  }

  return AuthService{ jwtSecret: []byte(jwtSecret) }
}

func (a AuthService) CreateTokens() (*TokenDetails, error) {
  td := &TokenDetails{}
  td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
  td.AccessUuid = uuid.NewString()

  td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
  td.RefreshUuid = uuid.NewString()

  var err error

  at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "access_uuid": td.AccessUuid,
    "exp": td.AtExpires,
  })
  td.AccessToken, err = at.SignedString(a.jwtSecret)
  if err != nil {
    return nil, err
  }

  rt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "refresh_uuid": td.RefreshUuid,
    "exp": td.RtExpires,
  })
  td.RefreshToken, err = rt.SignedString(a.jwtSecret)
  if err != nil {
    return nil, err
  }

  return td, nil
}

func (a AuthService) VerifyToken(tokenString string) (*jwt.Token, error) {
  token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
    return a.jwtSecret, nil
  })
  if err != nil {
    return nil, err
  }

  if !token.Valid {
    return nil, errors.New("invalid token")
  }

  return token, nil
}

func (a AuthService) GetAccessUuid(tokenString string) (string, error) {
  token, err := a.VerifyToken(tokenString)
  if err != nil {
    return "", err
  }

  claims, ok := token.Claims.(jwt.MapClaims)
  if !ok {
    return "", errors.New("invlaid token")
  }

  accessUuid, ok := claims["access_uuid"].(string)
  if !ok {
    return "", errors.New("no access uuid")
  }

  return accessUuid, nil
}
