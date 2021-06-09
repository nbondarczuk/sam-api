package commontest

import (
	"testing"
	
	"sam-api/common"	
)

//
// scenario: check emptiness
//
func TestEmptyValue(t *testing.T) {
	ok := common.EmptyValue("")
	// check result(s)
	if !ok {
		t.Errorf("Expected empty value")
		return
	}
}

//
// scenario: check lack emptiness
//
func TestNotEmptyValue(t *testing.T) {
	ok := common.EmptyValue("xxx")
	// check result(s)
	if ok {
		t.Errorf("Expected empty value")
		return
	}
}


//
// scenario: check membership
//
func TestMemberOf(t *testing.T) {
	ok := common.MemberOf("x", "x", "y", "z")
	// check result(s)
	if !ok {
		t.Errorf("Expected x as member of x:y:z")
		return
	}
}

//
// scenario: check lack of membership
//
func TestNotMemberOf(t *testing.T) {
	ok := common.MemberOf("xx", "x", "y", "z")
	// check result(s)
	if ok {
		t.Errorf("Expected xx as not member of x:y:z")
		return
	}
}

