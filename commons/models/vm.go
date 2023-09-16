package models

type Page struct{
	IsPrimaryKey *bool
	IntValue *int 
	StringValue *string
}

type Pager struct{
	Columns []string 
	Pages [][]Page // b-tree
}

type Table struct {
	Name string 
	NumRows int 
	Pager *Pager 
}

type VMResult struct {
	MSG *string 
	Cursor *Pager 
	Err error
}