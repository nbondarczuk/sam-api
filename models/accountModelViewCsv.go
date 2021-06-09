package models

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func (a *Account) ToCsvRecord() ([]string) {
	var r = make([]string, 15)

	r[0] = a.Status
	r[1] = a.ReleaseId
	r[3] = a.BscsAccount
	r[4] = a.OfiSapAccount
	r[5] = a.ValidFromDate.Format(time.RFC3339)
	r[6] = a.VatCodeInd
	r[7] = a.OfiSapWbsCode
	r[8] = fmt.Sprintf("%d", a.CitMarkerVatFlag)
	r[9]  = a.EntryDate.Format(time.RFC3339)
	r[10] = a.EntryOwner
	r[11] = a.UpdateDate.Format(time.RFC3339)
	r[12] = a.UpdateOwner
	r[13] = a.ReleaseDate.Format(time.RFC3339)
	r[14] = a.ReleaseOwner
	
	return r
}

func (accounts *Accounts) ToCsv() (rv []byte, err error) {
	// prepeare output format for csv writer
	var records = make([][]string, len(accounts.Data))
	for i, account := range accounts.Data {
		records[i] = account.ToCsvRecord()
		log.Printf("Produced record: %#v", records[i])
	}
	
	// open temporary file
    tmp, err := ioutil.TempFile(os.TempDir(), "sam-api-")
    if err != nil {
        return nil, fmt.Errorf("Can't create temporary file: %s", err.Error())
    }
    //defer os.Remove(tmp.Name())
	log.Println("Created temporary file: " + tmp.Name())

	// write to it all records
	w := csv.NewWriter(tmp)
	w.WriteAll(records)
	if err = w.Error(); err != nil {
		return
	}
	w.Flush()

	var offset int64
	offset, err = tmp.Seek(offset, 0)
	if err != nil {
		return 
	}
	
	// get the contents of the produced file
	rv, err = ioutil.ReadAll(tmp)
	if err != nil {
		return
	}
	log.Printf("CSV: %s", string(rv))
	
	// finish by closing file
	if err = tmp.Close(); err != nil {
        return
    }

	return
}
