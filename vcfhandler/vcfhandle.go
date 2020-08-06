package vcfhandler

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"bitbucket.org/llg/vcard"
)

//TODO: in Future add more fields into Fields such as address ...

//NOTE: program still stupid cannot figure out which column contains a phone number or name

type VcfElmnt struct {
	Data   [][]string
	Wg     *sync.WaitGroup
	Chan   chan VcfChanRes
	Writer *os.File
}

type VcfChanRes struct {
	Ok  bool
	Err error
}

func NewVcf(data [][]string, wg *sync.WaitGroup, c chan VcfChanRes, filename string) (*VcfElmnt, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &VcfElmnt{data, wg, c, f}, nil
}

func (v *VcfElmnt) ExtWrite() {
	vciw := vcard.NewDirectoryInfoWriter(v.Writer)
	defer v.Wg.Done()
	chR := new(VcfChanRes)
	count := 0
	for i, n := range v.Data {
		var vc vcard.VCard
		if i != 0 {
			//GET NAMES FIRST n[1]
			if n[1] != "" {
				name := strings.Split(n[1], " ")
				vc.GivenNames = append(vc.GivenNames, name[0])
				vc.FormattedName = n[1]
				if len(name) > 1 {
					vc.FamilyNames = append(vc.FamilyNames, name[1])
				}
			}
			//GET PHONE NUMBER n[0]
			if n[0] != "" {
				vc.Telephones = append(vc.Telephones, vcard.Telephone{Number: n[0]})
			}
		}
		vc.WriteTo(vciw)
		count++
		fmt.Printf("%v\n", count)
	}
	chR.Err, chR.Ok = nil, true
	v.Chan <- *chR
}
