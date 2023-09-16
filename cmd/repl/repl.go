package repl

import(
	"bufio"
	"fmt"
	"os"
    "github.com/hendry19901990/db_tutorial/internal/compiler"
	"github.com/hendry19901990/db_tutorial/internal/vm"
)

type REPL interface{
	StartREPL()
}

type replStruct struct{

}

func New() REPL{
	return &replStruct{}
}

func (r *replStruct) StartREPL(){
	reader := bufio.NewReader(os.Stdin)

	for{
		fmt.Print("DB >> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v \n", err)
			continue
		}

		input = input[:len(input)-1]

		if input == "exit"{
			fmt.Println("Exiting ...")
			break
		}

		r.evaluate(input)
	}


}

func (r *replStruct) evaluate(str string){
    comp := compiler.NewCompiler(str)
	bt := comp.Call()
	/*for _, v := range bt.Instructions {
		fmt.Println("Output ", v)
	}*/

	machine := vm.NewVM("users", bt)
	res := machine.ExecuteByteCode()
	if res.Err != nil {
		fmt.Println(res.Err)
		return
	}

	if res.MSG != nil {
		fmt.Println(*res.MSG)
	}

	if res.Cursor != nil {
		for i := 0; i < len(res.Cursor.Columns); i++{
			fmt.Print("| ", res.Cursor.Columns[i], " ")
		}
		fmt.Println("\n-------------------------------------")

		for _, record := range res.Cursor.Pages{
			for _, currentP := range record {
				if currentP.IsPrimaryKey != nil && *currentP.IsPrimaryKey {
					fmt.Print("PK: ", *currentP.IntValue, " ")
					continue
				}

				if currentP.IntValue != nil  {
					fmt.Print("| ", *currentP.IntValue, " ")
					continue
				}

				if currentP.StringValue != nil  {
					fmt.Print("| ", *currentP.StringValue, " ")
					continue
				}
			}
			fmt.Println("\n-------------------------------------")
		}
	}
	
}