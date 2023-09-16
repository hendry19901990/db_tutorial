package vm

import (
	"fmt"
	"os"
	"encoding/json"
	"github.com/hendry19901990/db_tutorial/commons/models"
)

const FILE_NAME = "DB.txt"


type VM interface {
	ExecuteByteCode() *models.VMResult
}


type MiniVM struct {
	tableName string 
	byteCode *models.ByteCode
}

func NewVM(t string , byteCode *models.ByteCode) VM {
	return &MiniVM{tableName: t, byteCode: byteCode}
}

func (vm_ *MiniVM) ExecuteByteCode() *models.VMResult{
	tbl, err := vm_.dbOpen()
	res := &models.VMResult{}
	if err != nil{
		res.Err = err
        return  res	
	}

	typeB := vm_.byteCode.Instructions[0].Type

	if typeB == models.ByteCodeOperationTypeINSERT {
		str, err := vm_.executeInsert(tbl)
		res.MSG = str
		res.Err = err 
		res.Cursor = nil
		return res
	}

	if typeB == models.ByteCodeOperationTypeSELECT {
		cursor := vm_.executeSelect(tbl)
		res.Cursor = cursor
		return res
	} 

	res.Err = fmt.Errorf("%s", "Operation NOT_FOUND")
    return res
}

func (vm_ *MiniVM) executeSelect(t *models.Table) *models.Pager{
	n := vm_.byteCode.Instructions[2].Count
	cols := []string{}
	for i:= 3; i <n+3; i++{
		cols = append(cols, *vm_.byteCode.Instructions[i].Identifier)
	}
	t.Pager.Columns= cols 

	return t.Pager
}

func chedckFileExists(fPath string) bool{
  _, err := os.Open(fPath)
  return err == nil
}

func (vm_ *MiniVM) executeInsert(t *models.Table) (*string, error){
    n := vm_.byteCode.Instructions[2].Count
    
	if t.NumRows ==0 {
		cols := []string{}
		for i:= 3; i <n+3; i++{
			cols = append(cols, *vm_.byteCode.Instructions[i].Identifier)
		}
		t.Pager = &models.Pager{
			Columns: cols,
			Pages: [][]models.Page{},
		}
	}

	index := n+4
	b := true
	p := []models.Page {
		{
			IsPrimaryKey: &b,
			IntValue: vm_.byteCode.Instructions[index].IntValue,
		},
	}

	if t.NumRows > 0 {
		for _, currentP := range t.Pager.Pages {
			if currentP[0].IntValue != nil &&
			   *currentP[0].IntValue == *vm_.byteCode.Instructions[index].IntValue{
				return nil, fmt.Errorf("Duplicated key %d", *vm_.byteCode.Instructions[index].IntValue)
			}
		}
	}

	index++
	for index < len(vm_.byteCode.Instructions){
		pg := models.Page{
			IntValue: vm_.byteCode.Instructions[index].IntValue,
			StringValue: vm_.byteCode.Instructions[index].StringValue,
		}

		p = append(p, pg)
		index++
	}

	t.Pager.Pages = append(t.Pager.Pages, p)
	t.NumRows = len(t.Pager.Pages)

	if err := vm_.write(t.Pager); err != nil {
		return nil, err
	}

	str := "Record stored successfully"
	return &str, nil

}


func (vm_ *MiniVM) write(p *models.Pager) error {
	f, err := os.Create(FILE_NAME)
	if err != nil {
		return err
	}
    defer f.Close()

	bts, err := json.Marshal(p)
	if err != nil {
		return err
	}

	f.Write(bts)
	f.Sync()
	
    return nil
}


func (vm_ *MiniVM) dbOpen () (*models.Table, error){
	if !chedckFileExists(FILE_NAME) {
		f, _ := os.OpenFile(FILE_NAME, os.O_CREATE | os.O_APPEND, 0644)
		f.Close()
	}

    dat, err := os.ReadFile(FILE_NAME)
	if err != nil {
		return nil, err
	}

	var p models.Pager 
	if len(dat) >0 {
		if err := json.Unmarshal(dat, &p); err != nil {
			return nil, err
		}
	}

	t := models.Table{
		Name: vm_.tableName,
		NumRows: len(p.Pages),
		Pager: &p,
	}
    
	return &t, nil
}