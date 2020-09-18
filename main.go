package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/TaKeO90/exceldumper/csvhandler"
	"github.com/TaKeO90/exceldumper/vcfhandler"
	"github.com/TaKeO90/exceldumper/xlsxhandler"
)

var (
	excel     string
	csvfile   string
	vcffile   string
	sheet     string
	cntNumber int
)

type chanres struct {
	ok  bool
	err error
}

func excelToVcf(nx xlsxhandler.XlsxData, ch chan chanres) {
	dumpedData, err := fileOpenWorker(nx)
	if err != nil {
		ch <- *(&chanres{false, err})
	}
	ok, err := vcfWorker(dumpedData)
	if err != nil {
		ch <- *(&chanres{false, err})
	}
	ch <- *(&chanres{ok, nil})
}

func main() {

	var (
		WG sync.WaitGroup
	)

	flagParser()

	if excel != "" && csvfile != "" && vcffile == "" {
		WG.Add(2)

		cc := make(chan csvhandler.ChanRes)

		nx, err := xlsxhandler.New(excel, sheet, nil, cntNumber)
		checkError(err)
		err = nx.Open()
		checkError(err)
		fmt.Printf("[*] Extracting Data \n")
		data, err := nx.Dump()
		checkError(err)

		fmt.Printf("[*] Saving Data into csv file \n")
		csvFD, err := csvhandler.New(csvfile, data, cc, &WG)
		checkError(err)
		go csvFD.WriteData()
		V := <-cc
		checkError(V.Err)
		if V.IsOk {
			fmt.Printf("Finished\n")
		} else {
			fmt.Printf("Failed\n")
		}
		WG.Wait()

	} else if excel != "" && vcffile != "" && csvfile == "" {
		ch := make(chan chanres)
		nx, err := xlsxhandler.New(excel, sheet, nil, cntNumber)
		checkError(err)
		go excelToVcf(nx, ch)
		results := <-ch
		checkError(results.err)
		if results.ok {
			fmt.Println("finished")
		}

	} else {
		flag.PrintDefaults()
	}
}

func fileOpenWorker(x xlsxhandler.XlsxData) ([][]string, error) {
	err := x.Open()
	checkError(err)
	fmt.Printf("[*] Extracting Data \n")
	data, err := x.Dump()
	if err != nil {
		return data, err
	}
	return data, nil
}

func vcfWorker(dumpData [][]string) (bool, error) {
	vcf, err := vcfhandler.NewVcf(dumpData, vcffile)
	if err != nil {
		return false, err
	}
	ok := vcf.ExtWrite()
	return ok, nil
}

func flagParser() {
	flag.StringVar(&excel, "excel", "", "Specify the excel file path that you want to dump it data into <csv or vcf>")
	flag.StringVar(&csvfile, "csvfile", "", "Specify the name of the outputed csv file")
	flag.StringVar(&sheet, "sheet", "Sheet1", "Specify the name of the Sheet Default is Set to <Sheet1>")
	flag.StringVar(&vcffile, "vcffile", "", "Specify the name of the outputed vcf file")
	flag.IntVar(&cntNumber, "cntnumber", 0, "number of contacts to write into the vcf file")
	flag.Parse()

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
