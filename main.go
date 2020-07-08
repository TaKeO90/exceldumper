package main

import (
	"fmt"
	"log"

	"./xlsxhandler"
)

//TODO: get filepath from stdin
//TODO: Write data to a csv file
const filepath string = "./1.xlsx"

func main() {
	var x xlsxhandler.XlsxData
	c := make(chan xlsxhandler.ChanResult)
	n := xlsxhandler.New(filepath, "Sheet1", c)
	x = n
	err := x.Open()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Extracting Data please wait ....")
	go x.Dump()
	D := <-c
	if D.Err != nil {
		log.Fatal(D.Err)
	}
	fmt.Println(D.CsvData)
}
