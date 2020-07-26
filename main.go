package main

import (
	"fmt"
	"log"
	"sync"

	"./csvhandler"
	"./xlsxhandler"
)

//TODO: add flags to the program
//TODO: progress bar or status
//TODO: support vcf format

const filepath string = "./1.xlsx"

func main() {
	var WG sync.WaitGroup
	var x xlsxhandler.XlsxData
	var csF csvhandler.CsvFile
	var csvfilename string = "file.csv"

	WG.Add(2)

	c := make(chan xlsxhandler.ChanResult)
	n := xlsxhandler.New(filepath, "Sheet2", c)
	x = n

	err := x.Open()
	checkError(err)

	excelData, err := workerI(&WG, &x, &c)
	checkError(err)

	workerII(&WG, csF, excelData, csvfilename)

	WG.Wait()
}

func workerI(wg *sync.WaitGroup, x *xlsxhandler.XlsxData, c *chan xlsxhandler.ChanResult) ([][]string, error) {
	fmt.Printf("[*] Extracting Data \n")
	defer wg.Done()
	xd := *x
	go xd.Dump()
	D := <-*c
	if D.Err != nil {
		return [][]string{}, D.Err
	}
	return D.CsvData, nil
}

func workerII(wg *sync.WaitGroup, csF csvhandler.CsvFile, DumpedData [][]string, csvfilename string) {
	cc := make(chan csvhandler.ChanRes)
	fmt.Printf("[*] Saving Data into csv file \n")
	csvFD, err := csvhandler.New(csvfilename, DumpedData, cc)
	checkError(err)
	defer wg.Done()
	csF = csvFD
	go csF.WriteData()
	V := <-cc
	checkError(V.Err)
	if V.IsOk {
		fmt.Printf("Finished\n")
	} else {
		fmt.Printf("Failed\n")
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
