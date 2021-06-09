package commontest

import (
	"testing"
	"time"
	
	"sam-api/common"	
)

//
// scenario: check date in string
//
func TestStringCutOffDate(t *testing.T) {
	var sdt string = "2019-11-01"
	ok := common.IsCutOffDate(sdt)
	// check result(s)
	if !ok {
		t.Errorf("Expected 1st day of the month: " + sdt)
		return
	}
}

//
// scenario: check date in string
//
func TestStringNotCutOffDate(t *testing.T) {
	var sdt string = "2019-11-02"
	ok := common.IsCutOffDate(sdt)
	// check result(s)
	if ok {
		t.Errorf("Expected not a cut off date: " + sdt)
		return
	}
}

//
// scenario: check date in Time
//
func TestTimeCutOffDate(t *testing.T) {
	var sdt string = "2019-11-01"
	var dt time.Time
	var err error
	// convert string to Time
	dt, err = time.Parse(common.CutOffDateFormat, sdt)
	if err != nil {
		t.Errorf("Wrong date format: " + sdt + " - not YYYY-MM-DD - " + err.Error())
		return		
	}
	ok := common.IsCutOffDate(dt)
	// check result(s)
	if !ok {
		t.Errorf("Expected 1st day of the month: " + sdt)
		return
	}
}

//
// scanario: produced cut off date is correct
//
func TestNextCutOffDate(t *testing.T) {
	ts, err := common.NextCutOffDate()
	if err != nil {
		t.Errorf("Can't produce cut off date")
		return		
	}
	ok := common.IsCutOffDate(ts)
	// check result(s)
	if !ok {
		t.Errorf("Expected 1st day of the month: %#v", ts)
		return
	}
}	
