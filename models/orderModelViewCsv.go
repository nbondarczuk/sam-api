package models

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func (o *Order) ToCsvRecord() ([]string) {
	var r = make([]string, 12)

	r[0] = o.Status
	r[1] = o.ReleaseId
	r[2] = o.BscsAccount
	r[3] = o.SegmentCode
	r[4] = o.OrderNumber
	r[5] = o.ValidFromDate.Format(time.RFC3339)
	r[6] = o.EntryDate.Format(time.RFC3339)
	r[7] = o.EntryOwner
	r[8] = o.UpdateDate.Format(time.RFC3339)
	r[9] = o.UpdateOwner
	r[10] = o.ReleaseDate.Format(time.RFC3339)
	r[11] = o.ReleaseOwner
	
	return r
}

func (orders *Orders) ToCsv() (rv []byte, err error) {
	// prepeare output format for csv writer
	var records = make([][]string, len(orders.Data))
	for i, order := range orders.Data {
		records[i] = order.ToCsvRecord()
	}
	
	// open temporary file
    tmp, err := ioutil.TempFile(os.TempDir(), "sam-api-")
    if err != nil {
        return nil, fmt.Errorf("Can't create temporary file: %s", err.Error())
    }
    defer os.Remove(tmp.Name())
	log.Println("Created temporary file: " + tmp.Name())

	// write to it all records
	w := csv.NewWriter(tmp)
	w.WriteAll(records)
	if err = w.Error(); err != nil {
		return
	}
	w.Flush()

	// get the contents of the produced file
	rv, err = ioutil.ReadAll(tmp)
	if err != nil {
		return
	}
	
	// finish by closing file
	if err = tmp.Close(); err != nil {
        return
    }

	return
}
