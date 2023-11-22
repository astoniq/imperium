package objecttype

import (
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

type Model interface {
	GetId() int64
	GetTypeId() string
	GetDefinition() string
	SetDefinition(value string)
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetDeletedAt() *time.Time
	ToSpec() (*Spec, error)
}

type ObjectType struct {
	Id         int64      `postgres:"id" sqlite:"id"`
	TypeId     string     `postgres:"type_id" sqlite:"typeId"`
	Definition string     `postgres:"definition" sqlite:"definition"`
	CreatedAt  time.Time  `postgres:"created_at" sqlite:"createdAt"`
	UpdatedAt  time.Time  `postgres:"updated_at" sqlite:"updatedAt"`
	DeletedAt  *time.Time `postgres:"deleted_at" sqlite:"deletedAt"`
}

func (objectType *ObjectType) GetId() int64 {
	return objectType.Id
}

func (objectType *ObjectType) GetTypeId() string {
	return objectType.TypeId
}

func (objectType *ObjectType) GetDefinition() string {
	return objectType.Definition
}

func (objectType *ObjectType) SetDefinition(newDefinition string) {
	objectType.Definition = newDefinition
}

func (objectType *ObjectType) GetCreatedAt() time.Time {
	return objectType.CreatedAt
}

func (objectType *ObjectType) GetUpdatedAt() time.Time {
	return objectType.UpdatedAt
}

func (objectType *ObjectType) GetDeletedAt() *time.Time {
	return objectType.DeletedAt
}

func (objectType *ObjectType) ToSpec() (*Spec, error) {
	var objectTypeSpec Spec
	err := json.Unmarshal([]byte(objectType.Definition), &objectTypeSpec)
	if err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling object type %s", objectType.TypeId)
	}
	objectTypeSpec.CreatedAt = objectType.CreatedAt
	return &objectTypeSpec, nil
}
