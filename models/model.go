package models

type Model interface {
	CollectionName() string
	GetBSON() interface{}
}
