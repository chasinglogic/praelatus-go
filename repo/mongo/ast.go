// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package mongo

import (
	"github.com/praelatus/praelatus/ql/ast"
	"gopkg.in/mgo.v2/bson"
)

func eval(exp ast.InfixExpression) bson.M {
	b := bson.M{}

	switch exp.Operator {
	case "AND":
		opName := "$and"
		if _, ok := b[opName]; ok {
			b[opName] = append(b[opName].([]bson.M),
				eval(exp.Left.(ast.InfixExpression)), eval(exp.Right.(ast.InfixExpression)))
		} else {
			b[opName] = []bson.M{
				eval(exp.Left.(ast.InfixExpression)),
				eval(exp.Right.(ast.InfixExpression)),
			}
		}
	case "OR":
		opName := "$or"
		if _, ok := b[opName]; ok {
			b[opName] = append(b[opName].([]bson.M),
				eval(exp.Left.(ast.InfixExpression)), eval(exp.Right.(ast.InfixExpression)))
		} else {
			b[opName] = []bson.M{
				eval(exp.Left.(ast.InfixExpression)),
				eval(exp.Right.(ast.InfixExpression)),
			}
		}
	case "=":
		val, ok := exp.Right.(ast.Literal)
		if !ok {
			break
		}
		makeFieldSearchDoc(exp, b, val.GetValue())
	case "~":
		val, ok := exp.Right.(ast.Literal)
		if !ok {
			break
		}
		makeFieldSearchDoc(exp, b, bson.M{"$regex": val.GetValue()})
	case "!=":
		val, ok := exp.Right.(ast.Literal)
		if !ok {
			break
		}
		makeFieldSearchDoc(exp, b, bson.M{"$ne": val.GetValue()})
	case ">":
		val, ok := exp.Right.(ast.Literal)
		if !ok {
			break
		}
		makeFieldSearchDoc(exp, b, bson.M{"$gt": val.GetValue()})
	case "<":
		val, ok := exp.Right.(ast.Literal)
		if !ok {
			break
		}
		makeFieldSearchDoc(exp, b, bson.M{"$lt": val.GetValue()})
	case ">=":
		val, ok := exp.Right.(ast.Literal)
		if !ok {
			break
		}
		makeFieldSearchDoc(exp, b, bson.M{"$gte": val.GetValue()})
	case "<=":
		val, ok := exp.Right.(ast.Literal)
		if !ok {
			break
		}
		makeFieldSearchDoc(exp, b, bson.M{"$lte": val.GetValue()})
	}

	return b
}

func evalAST(a ast.AST) bson.M {
	infix, ok := a.Query.Expression.(ast.InfixExpression)
	if !ok {
		return nil
	}

	return eval(infix)
}

func makeFieldSearchDoc(exp ast.InfixExpression, b bson.M, valDoc interface{}) {
	fn := exp.Left.(ast.FieldLiteral)
	if fn.Value == "status" {
		b["status.name"] = valDoc
	} else if fn.Value == "statusCategory" {
		b["status.type"] = valDoc
	} else if fn.Value == "key" {
		b["_id"] = valDoc
	} else if fn.IsCustomField() {
		customFieldDoc := bson.M{
			"name":  fn.Value,
			"value": valDoc,
		}

		if _, ok := b["fields"]; ok {
			b["fields"] = append(b["fields"].([]bson.M), customFieldDoc)
		} else {
			b["fields"] = []bson.M{customFieldDoc}
		}
	} else {
		b[fn.Value] = valDoc
	}
}
