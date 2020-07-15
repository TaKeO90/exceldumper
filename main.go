package main

import (
	"fmt"
	"log"
	"sync"

	"./csvhandler"
	"./xlsxhandler"
)

//TODO: add flags to the program
const filepath string = "./1.xlsx"

var WG sync.WaitGroup

func main() {
	var x xlsxhandler.XlsxData
	var csF csvhandler.CsvFile
	c := make(chan xlsxhandler.ChanResult)
	cc := make(chan csvhandler.ChanRes)
	n := xlsxhandler.New(filepath, "Sheet1", c)
	x = n
	err := x.Open()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[*] Extracting Data \n")
	WG.Add(2)
	go x.Dump()
	D := <-c
	if D.Err != nil {
		log.Fatal(D.Err)
	}
	WG.Done()
	fmt.Printf("[*] Saving Data into csv file \n")
	csvFD, err := csvhandler.New("file.csv", D.CsvData, cc)
	if err != nil {
		log.Fatal(err)
	}
	csF = csvFD
	go csF.WriteData()
	V := <-cc
	if V.Err != nil {
		log.Fatal(V.Err)
	}
	if V.IsOk {
		fmt.Printf("Finished\n")
		WG.Done()
	}
	WG.Wait()
}
