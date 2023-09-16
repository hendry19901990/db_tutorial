package compiler

import (
	"github.com/hendry19901990/db_tutorial/commons/models"
	"strings"
	"strconv"
)

type Generator interface{
	GenerateCode(tokens []Token) *models.ByteCode
}

func NewGenerator() Generator{
	return &generator{}
}

type generator struct{

}

func (g *generator) GenerateCode(tokens []Token) *models.ByteCode{
   var bt  *models.ByteCode

   if tokens[0].Type == TokenKeyword && strings.ToUpper(tokens[0].value) == "INSERT"{
	 bt = g.generateCodeInsert(tokens)
   }

   if tokens[0].Type == TokenKeyword && strings.ToUpper(tokens[0].value) == "SELECT"{
	bt = g.generateCodeSelect(tokens)
    }

   return bt
}

func stringToStringPtr(s string) *string{
  return &s
}

func stringToIntPtr(s string) *int{
  v, _ := strconv.Atoi(s)
  return &v
}

func isInteger(s string) bool{
	_, e := strconv.Atoi(s)
    return e == nil
}

func (g *generator)  generateCodeInsert(tokens []Token)*models.ByteCode {
	 bt := &models.ByteCode{}

	ints := []models.ByteCodeValue{
		{
			Type: models.ByteCodeOperationTypeINSERT,
			Identifier: stringToStringPtr("INSERT"),
		},
	}

	tblName :=  models.ByteCodeValue{
		Type: models.ByteCodeOperationTypeTableName,
		Identifier: stringToStringPtr(tokens[2].value),
	}


	varNames := []models.ByteCodeValue{}
	i:= 3
	for i < len(tokens){
		if tokens[i].Type == TokenKeyword && strings.ToUpper(tokens[i].value) == "VALUES"{
			break
		}

		if tokens[i].Type == TokenIdentifier{
			v := models.ByteCodeValue{
				Type: models.ByteCodeOperationTypeIdentifier,
				Identifier: stringToStringPtr(tokens[i].value),
			}
			varNames = append(varNames, v)
		}
		
		i++
	}
	i++
	

	countVarNames :=  models.ByteCodeValue{
		Type: models.ByteCodeOperationTypeCount,
		Count: len(varNames),
	}


	varVals  := []models.ByteCodeValue{}
	for i < len(tokens){
		if tokens[i].Type == TokenSymbol && strings.ToUpper(tokens[i].value) == ";"{
			break
		}

		if tokens[i].Type == TokenIdentifier{
			v := models.ByteCodeValue{
				Type: models.ByteCodeOperationTypeIdentifier,
				//Identifier: stringToStringPtr(tokens[i].value),
			}
			if isInteger(tokens[i].value){
				v.IntValue = stringToIntPtr(tokens[i].value)
			}else{
				v.StringValue = stringToStringPtr(tokens[i].value)
			}
			varVals = append(varVals, v)
		}
		
		i++
	}

	countvarVals :=  models.ByteCodeValue{
		Type: models.ByteCodeOperationTypeCount,
		Count: len(varVals),
	}


	ints = append(ints, tblName)
	ints = append(ints, countVarNames)
	ints = append(ints, varNames...)
	ints = append(ints, countvarVals)
	ints = append(ints, varVals...)

    bt.Instructions = ints
	return bt
}

func (g *generator)  generateCodeSelect(tokens []Token)*models.ByteCode {
	bt := &models.ByteCode{}

	ints := []models.ByteCodeValue{
		{
			Type: models.ByteCodeOperationTypeSELECT,
			Identifier: stringToStringPtr("SELECT"),
		},
	}

	varNames := []models.ByteCodeValue{}
	i:= 1
	for i < len(tokens){
		if tokens[i].Type == TokenKeyword && strings.ToUpper(tokens[i].value) == "FROM"{
			break
		}

		if tokens[i].Type == TokenIdentifier{
			v := models.ByteCodeValue{
				Type: models.ByteCodeOperationTypeIdentifier,
				Identifier: stringToStringPtr(tokens[i].value),
			}
			varNames = append(varNames, v)
		}
		
		i++
	}
	i++
	tblName :=  models.ByteCodeValue{
		Type: models.ByteCodeOperationTypeTableName,
		Identifier: stringToStringPtr(tokens[i].value),
	}

	countVarNames :=  models.ByteCodeValue{
		Type: models.ByteCodeOperationTypeCount,
		Count: len(varNames),
	}

	ints = append(ints, tblName)
	ints = append(ints, countVarNames)
	ints = append(ints, varNames...)

    bt.Instructions = ints
	return bt
}