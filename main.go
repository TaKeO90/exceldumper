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

	c := make(chan xlsxhandler.ChanResult)

	if excel != "" && sheet != "" && csvfile != "" && vcffile == "" { //ADD HERE A CONDITION FOR CSV THEN ANOTHER ONE FOR VCF

		WG.Add(2)
		cc := make(chan csvhandler.ChanRes)

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

	} else if excel != "" && sheet != "" && vcffile != "" && csvfile == "" {
		WG.Add(2)
		cB := make(chan bool)
		n, err := xlsxhandler.New(excel, sheet, c, &WG)
		checkError(err)
		x = n
		err = x.Open()
		checkError(err)
		fmt.Printf("[*] Extracting Data \n")
		go x.Dump()
		res := <-c
		checkError(res.Err)

		vcf, err := vcfhandler.NewVcf(res.DumpData, &WG, cB, vcffile)
		checkError(err)

		go vcf.ExtWrite()
		isFinished := <-cB
		if isFinished {
			fmt.Println("Finished")
		}
	} else {
		flag.PrintDefaults()
	}
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
