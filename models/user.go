package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	AccessToken  string `bson:"access_token"`
	RefreshToken string `bson:"refresh_token"`
}

type User struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	Email       string              `bson:"email,omitempty" json:"email,omitempty"`
	Password    string              `bson:"hashed_password,omitempty" json:"password,omitempty"`
	PhoneNumber string              `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	HasTwoFa    bool                `bson:"has_two_fa,omitempty" json:"has_two_fa,omitempty"`
	Accounts    map[string]*Account `bson:"accounts,omitempty" json:"accounts,omitempty"`
	CreatedAt   primitive.DateTime  `bson:"created_at"`
	UpdatedAt   primitive.DateTime  `bson:"updated_at"`
}

func (u User) CollectionName() string {
	return "users"
}

func (u *User) GetBSON() interface{} {
	now := primitive.NewDateTimeFromTime(time.Now().UTC())

	if u.CreatedAt == 0 {
		u.CreatedAt = now
	}

	u.UpdatedAt = now

	return u
}

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 13)
	u.Password = string(bytes)

	return err
}

func (u User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) AddAccount(name string, account *Account) {
	if u.Accounts == nil {
		u.Accounts = make(map[string]*Account)
	}

	u.Accounts[name] = account
}
