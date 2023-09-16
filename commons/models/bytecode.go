package models

import(
	"encoding/json"
)

type ByteCodeOperationType int

const (
	ByteCodeOperationTypeINSERT = iota
	ByteCodeOperationTypeSELECT 
	ByteCodeOperationTypeTableName
	ByteCodeOperationTypeIdentifier
	ByteCodeOperationTypeValue
	ByteCodeOperationTypeCount
)

type ByteCodeValue struct{
	Type ByteCodeOperationType
	Identifier *string
	IntValue *int
	StringValue *string
	Count int
}

func (b ByteCodeValue) String() string{
	s, _ := json.Marshal(b)
	return string(s)
}

type ByteCode struct {
	Instructions []ByteCodeValue
}