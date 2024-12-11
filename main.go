package main

import (
	"github.com/orlandofra/ghligh/ghligh"
	"fmt"
)

func main(){

	doc, err := ghligh.Open("test.pdf")
	defer doc.Close()
	if err != nil {
		fmt.Println("error")
		return
	}

	/* export annots */
	fmt.Printf(doc.Export())
}
