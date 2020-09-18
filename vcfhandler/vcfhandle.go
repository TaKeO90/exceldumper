package vcfhandler

import (
	"os"
	"strings"

	"bitbucket.org/llg/vcard"
)

//TODO: in Future add more fields into Fields such as address ...

//NOTE: program still stupid cannot figure out which column contains a phone number or name

type VcfElmnt struct {
	Data   [][]string
	Writer *os.File
}

func NewVcf(data [][]string, filename string) (*VcfElmnt, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &VcfElmnt{data, f}, nil
}

func (v *VcfElmnt) ExtWrite() bool {
	defer v.Writer.Close()
	vciw := vcard.NewDirectoryInfoWriter(v.Writer)
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
	}
	return true
}
