package main

import(

	"fmt"
	"os"
	"github.com/hendry19901990/db_tutorial/cmd/repl"
)

func main() {
	cmd := repl.New()
	cmd.StartREPL()

	fmt.Println("Goodbye ...")
	os.Exit(0)
}	