package compiler

import (
	"github.com/hendry19901990/db_tutorial/commons/models"
)

type Compiler interface{
	Call() *models.ByteCode
}

type compiler struct {
	lexer Lexer
	parser Parser
	generator Generator
}

func NewCompiler(sqlText string) Compiler{
   
    lexer := NewLexer(sqlText)
	parser := NewParser(lexer)
	gene := NewGenerator()

	return &compiler{
		lexer: lexer,
		parser: parser,
		generator: gene,
	}

}

func (comp * compiler) Call() *models.ByteCode {
	 comp.parser.ParseStatement()
	 return comp.generator.GenerateCode(comp.parser.GetTokens())
}