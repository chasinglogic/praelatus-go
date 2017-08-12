package models

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

// ErrInvalidDataType indicates that the field was created with an incorrect
// data type
var ErrInvalidDataType = errors.New("Invalid data type for field")

// DataTypes holds the available data types
var DataTypes = []string{
	"FLOAT",
	"STRING",
	"INT",
	"DATE",
	"OPT",
}

// Field is a ticket field
type Field struct {
	Name     string `json:"name"`
	DataType string `json:"dataType"`

	// Options is only relevant for Fields of DataType OPT
	Options []string `json:"options,omitempty" bson:"options,omitempty"`

	// Value holds the value of the given field
	Value interface{} `json:"value,omitempty" bson:"value,omitempty"`
}

// IsValidDataType is used to verify that the field has a data type we can
// support
func (f *Field) IsValidDataType() bool {
	for _, t := range DataTypes {
		if t == f.DataType {
			return true
		}
	}

	return false
}

func (f *Field) String() string {
	return jsonString(f)
}

// FieldScheme assigns fields to a ticke type.
type FieldScheme struct {
	ID   bson.ObjectId `json:"id" bson:"_id"`
	Name string

	// Map ticket type to fields
	Fields map[string][]Field
}
