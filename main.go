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
	excel   string
	csvfile string
	vcffile string
	sheet   string
)

func main() {

	var (
		WG  sync.WaitGroup
		x   xlsxhandler.XlsxData
		csF csvhandler.CsvFile
	)

	flagParser()

	if excel != "" && csvfile != "" && vcffile == "" {
		WG.Add(2)

		cc := make(chan csvhandler.ChanRes)
		c := make(chan xlsxhandler.ChanResult)

		n, err := xlsxhandler.New(excel, sheet, c, &WG)
		checkError(err)
		x = n
		err = x.Open()
		checkError(err)
		fmt.Printf("[*] Extracting Data \n")
		go x.Dump()
		res := <-c
		checkError(res.Err)

		fmt.Printf("[*] Saving Data into csv file \n")
		csvFD, err := csvhandler.New(csvfile, res.DumpData, cc, &WG)
		checkError(err)
		csF = csvFD
		go csF.WriteData()
		V := <-cc
		checkError(V.Err)
		if V.IsOk {
			fmt.Printf("Finished\n")
		} else {
			fmt.Printf("Failed\n")
		}
		WG.Wait()

	} else if excel != "" && vcffile != "" && csvfile == "" {
		WG.Add(2)

		cB := make(chan vcfhandler.VcfChanRes)
		c := make(chan xlsxhandler.ChanResult)

		dumpedData, err := fileOpenWorker(&WG, c, x)
		checkError(err)
		ok, err := vcfWorker(&WG, cB, dumpedData)
		checkError(err)
		if ok {
			fmt.Println("Finished")
		}
		WG.Wait()

	} else {
		flag.PrintDefaults()
	}
}

func fileOpenWorker(wg *sync.WaitGroup, c chan xlsxhandler.ChanResult, x xlsxhandler.XlsxData) ([][]string, error) {
	n, err := xlsxhandler.New(excel, sheet, c, wg)
	checkError(err)
	x = n
	err = x.Open()
	checkError(err)
	fmt.Printf("[*] Extracting Data \n")
	go x.Dump()
	res := <-c
	if err != nil {
		return nil, err
	}
	return res.DumpData, nil
}

func vcfWorker(wg *sync.WaitGroup, c chan vcfhandler.VcfChanRes, dumpData [][]string) (bool, error) {
	vcf, err := vcfhandler.NewVcf(dumpData, wg, c, vcffile)
	if err != nil {
		return false, err
	}
	go vcf.ExtWrite()
	isFinished := <-c
	if isFinished.Err != nil {
		return false, err
	}
	if isFinished.Ok {
		return true, nil
	}
	return false, nil
}

func flagParser() {
	flag.StringVar(&excel, "excel", "", "Specify the excel file path that you want to dump it data into <csv or vcf>")
	flag.StringVar(&csvfile, "csvfile", "", "Specify the name of the outputed csv file")
	flag.StringVar(&sheet, "sheet", "Sheet1", "Specify the name of the Sheet Default is Set to <Sheet1>")
	flag.StringVar(&vcffile, "vcffile", "", "Specify the name of the outputed vcf file")
	flag.Parse()

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
