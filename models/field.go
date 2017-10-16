// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package models

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// ErrInvalidDataType indicates that the field was created with an incorrect
// data type
var ErrInvalidDataType = errors.New("Invalid data type for field")

// DataType is a string which indicates the type of data in a given field.
type DataType string

// These are the available data types for fields on Tickets.
const (
	FloatField  DataType = "FLOAT"
	StringField          = "STRING"
	IntField             = "INT"
	DateField            = "DATE"
	OptionField          = "OPTION"
)

// DataTypes holds the available data types
var DataTypes = []DataType{
	FloatField,
	StringField,
	IntField,
	DateField,
	OptionField,
}

// Field is a ticket field
type Field struct {
	Name     string   `json:"name"`
	DataType DataType `json:"dataType"`

	// Options is only relevant for Fields of DataType OPTION
	Options []string `json:"options,omitempty" bson:"options,omitempty"`

	// Value holds the value of the given field
	Value interface{} `json:"value,omitempty" bson:"value,omitempty"`
}

// IsValidDataType is used to verify that the field has a data type we can
// support
func (f Field) IsValidDataType() bool {
	for _, t := range DataTypes {
		if t == f.DataType {
			return true
		}
	}

	return false
}

func (f Field) String() string {
	return jsonString(f)
}

// FieldScheme assigns fields to a ticke type.
type FieldScheme struct {
	ID   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name string        `json:"name" required:"true"`

	// Map ticket type to fields
	Fields map[string][]Field `json:"fields"`
}

// ValidateTicket verifies that all fields on t are valid
func (fs FieldScheme) ValidateTicket(t Ticket) error {
	fields, ok := fs.Fields[t.Type]
	if !ok {
		fields, ok = fs.Fields[""]
		if !ok {
			return errors.New("no fields set for this ticket type and default not set")
		}
	}

	for _, f := range t.Fields {
		if !validField(fields, f) {
			return fmt.Errorf("%s is not a valid field for type %s", f.Name, t.Type)
		}
	}

	return nil
}

func validField(fields []Field, field Field) bool {
	for _, f := range fields {
		if field.Name == f.Name {
			return true
		}
	}

	return false
}
