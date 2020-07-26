package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/TaKeO90/exceldumper/csvhandler"
	"github.com/TaKeO90/exceldumper/xlsxhandler"
)

//TODO: progress bar or status
//TODO: support vcf format

var (
	excel   string
	csvfile string
	sheet   string
)

func main() {

	var (
		WG  sync.WaitGroup
		x   xlsxhandler.XlsxData
		csF csvhandler.CsvFile
	)

	flagParser()

	if excel == "" || sheet == "" || csvfile == "" {
		flag.PrintDefaults()
	} else {

		WG.Add(2)

		c := make(chan xlsxhandler.ChanResult)
		n, err := xlsxhandler.New(excel, sheet, c)
		checkError(err)
		x = n

		err = x.Open()
		checkError(err)

		excelData, err := workerI(&WG, &x, &c)
		checkError(err)

		workerII(&WG, csF, excelData, csvfile)

		WG.Wait()
	}
}

func flagParser() {
	flag.StringVar(&excel, "excel", "", "Specify the excel file path that you want to dump it data into <csv or vcf>")
	flag.StringVar(&csvfile, "csvfile", "", "Specify the name of the outputed csv file")
	flag.StringVar(&sheet, "sheet", "Sheet1", "Specify the name of the Sheet Default is Set to <Sheet1>")
	flag.Parse()
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
