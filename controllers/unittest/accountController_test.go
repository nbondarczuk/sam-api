package controllers

import (
	"bytes"
	"fmt"
	"github.com/unrolled/render"
	"net/http"
	"testing"
	"time"
	
	"sam-api/models"
	"sam-api/common"
)

var (
	accountFormatter = render.New(render.Options{
		IndentJSON: true,
	})
)

//
// scenario: clean the resource
//
func TestAccountDeleteAll(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()

	// prepare payload for Account Create
	req, err := http.NewRequest("DELETE", server.URL + "/api/account", nil)
	if err != nil {
		t.Errorf("Error in creating DELETE request for AccountDeleteAll: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + token)

	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in DELETE on AccountDeleteAll: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response status %d, received %d", http.StatusOK, res.StatusCode)
		return
	}
}

//
// scenario: account create on empty table
//
var testAccountId string
func TestAccountCreate(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()

	// prepare payload for Account Create
	testAccountId := newAccountId()
	body := []byte("{\"data\":{\"status\": \"W\", \"releaseId\": \"0\", \"bscsAccount\": \"" + testAccountId + "\"}}")
	req, err := http.NewRequest("POST", server.URL + "/api/account", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for AccountCreate: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to AccountCreate: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected response status %d, received %d", http.StatusCreated, res.StatusCode)
		return
	}
}

func TestAccountCreateWithDate(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()

	// prepare payload for Account Create
	testAccountId := newAccountId()
	var validFromDate string
	if ts, err := common.NextCutOffDate(); err != nil {
		t.Errorf("Error in creating cut off date: %v", err)
		return
	} else {
		validFromDate = ts.Format("2006-01-02")
	}
	
	body := []byte("{\"data\":{\"status\": \"W\", \"releaseId\": \"0\", \"bscsAccount\": \"" + testAccountId + "\", \"validFromDate\": \"" + validFromDate + "\"}}")
	req, err := http.NewRequest("POST", server.URL + "/api/account", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for AccountCreate: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to AccountCreate: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected response status %d, received %d", http.StatusCreated, res.StatusCode)
		return
	}
}

func TestAccountCreateWithDateInvalid(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()

	// prepare payload for Account Create
	testAccountId := newAccountId()
	var validFromDate string
	ts := time.Now()
	validFromDate = ts.Format("2006-01-02")
	body := []byte("{\"data\":{\"status\": \"W\", \"releaseId\": \"0\", \"bscsAccount\": \"" + testAccountId + "\", \"validFromDate\": \"" + validFromDate + "\"}}")
	req, err := http.NewRequest("POST", server.URL + "/api/account", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for AccountCreate: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to AccountCreate: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode == http.StatusCreated {
		t.Errorf("Expected response status %d, received %d", http.StatusCreated, res.StatusCode)
		return
	}
}

//
// scenario: an account previously created in test can be read with GET
//
func TestAccountReadSome(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Account Create
	req, err := http.NewRequest("GET", server.URL + "/api/account/W/0", nil)
	if err != nil {
		t.Errorf("Error in GET request for AccountReadSome: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in GET for AccountReadSome: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response status %d, received %d", http.StatusOK, res.StatusCode)
		return
	}
}

//
// scenario: an account previously created in test can be read with GET and last
//

func TestAccountReadSomeWithLast(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Account Create
	req, err := http.NewRequest("GET", server.URL + "/api/account/W/last", nil)
	if err != nil {
		t.Errorf("Error in GET request for AccountReadSome: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in GET for AccountReadSome: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response status %d, received %d", http.StatusOK, res.StatusCode)
		return
	}
}

//
// scenario: an account previously created in test can be read with GET and latest
//
func TestAccountReadSomeWithLatest(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Account Create
	req, err := http.NewRequest("GET", server.URL + "/api/account/W/latest", nil)
	if err != nil {
		t.Errorf("Error in GET request for AccountReadSome: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in GET for AccountReadSome: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode == http.StatusOK {
		t.Errorf("Expected not response status %d, received %d", http.StatusOK, res.StatusCode)
		return
	}
}

//
// scenario: an account previously created in test can be read with GET as active in default json format
//
func TestAccountReadAll(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Account Create
	req, err := http.NewRequest("GET", server.URL + "/api/account", nil)
	if err != nil {
		t.Errorf("Error in GET request for AccountReadSome: %v", err)
		return
	}
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in GET for AccountReadAll: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected not response status %d, received %d", http.StatusOK, res.StatusCode)
		return
	}
}

//
// scenario: an account previously created in test can be read with GET as active in csv format
//
func TestAccountReadAllAsJson(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Account Create
	req, err := http.NewRequest("GET", server.URL + "/api/account", nil)
	if err != nil {
		t.Errorf("Error in GET request for AccountReadSome: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in GET for AccountReadAll: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected not response status %d, received %d", http.StatusOK, res.StatusCode)
		return
	}
}

//
// scenario: an account previously created in test can be read with GET as active in csv format
//
func TestAccountReadAllAsCsv(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Account Create
	req, err := http.NewRequest("GET", server.URL + "/api/account", nil)
	if err != nil {
		t.Errorf("Error in GET request for AccountReadSome: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/csv")	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in GET for AccountReadAll: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected not response status %d, received %d", http.StatusOK, res.StatusCode)
		return
	}
}

//
// scenario: update previously created account
//
func TestAccountUpdateOne(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Account Update
	body := []byte("{\"data\":{\"status\": \"W\", \"releaseId\": \"0\", \"bscsAccount\": \"" + testAccountId + "\"}}")
	req, err := http.NewRequest("PUT", server.URL + "/api/account/W/0/"+ testAccountId, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in PUT request for AccountUpdateOne: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in PUT for AccountUpdateOne: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected response status %d, received %d", http.StatusNotFound, res.StatusCode)
		return
	}
}

//
// scenario: update one attribute of previously created account
//
func TestAccountUpdateAttributeOne(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Account Update
	body := []byte("{\"data\":{\"ofiSapAccount\":" + "\"XXX\"" + "}}")
	req, err := http.NewRequest("PATCH", server.URL + "/api/account/W/0/"+ testAccountId, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in PATCH request for AccountUpdateOne: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in PATCH for AccountUpdateOne: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected response status %d, received %d", http.StatusNotFound, res.StatusCode)
		return
	}
}

//
// scenario: delete one account
//
func TestAccountDeleteOne(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for AccountDeleteOne
	body := []byte("{\"data\":{\"status\": \"W\", \"releaseId\": \"0\", \"bscsAccount\": \"" + testAccountId + "\"}}")
	req, err := http.NewRequest("DELETE", server.URL + "/api/account/W/0/" + testAccountId, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in DELETE request for AccountDeleteOne: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in DELETE for AccountDeleteOne: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected response status %d, received %d", http.StatusNotFound, res.StatusCode)
		return
	}
}

//
// scenario: create N accounts and read them
//
func AccountCreateN(t *testing.T) {
	TestAccountDeleteAll(t)
	defer TestAccountDeleteAll(t)
	c, s, token := initTestEnv(t, "USER", "Booker", true)
	defer s.Close()
	const N int = 10
	
	var test = map[string]models.Account{}
	for i := 0; i < N; i++ {
		key := fmt.Sprintf("X%d", i)
		test[key] = models.Account{Status: "W", ReleaseId: "0", BscsAccount: key}
	}
	
	// create the test accounts
	for _, a := range test {
		aa := accountCreate(t, s, c, token, &a)
		if aa == nil {
			t.Errorf("No account created for: %#v", a)
			return 
		} else if !accountKeyEq(aa, &a) {
			t.Errorf("Difference detected: %#v != %#v", a, *aa)
			return 
		}
	}
}

//
// scenario: create N accounts, patch some of them, read all, check patched
//
func AccountCreate5Patch5(t *testing.T) {
	TestAccountDeleteAll(t)
	defer TestAccountDeleteAll(t)
	c, s, token := initTestEnv(t, "USER", "Booker", true)
	defer s.Close()
	const N int = 10
	
	var test = map[string]models.Account{}
	for i := 0; i < N; i++ {
		key := fmt.Sprintf("X%d", i)
		test[key] = models.Account{Status: "W", ReleaseId: "0", BscsAccount: key}
	}
	
	// create the test accounts
	for _, a := range test {
		aa := accountCreate(t, s, c, token, &a)
		if aa == nil {
			t.Errorf("No account created for: %#v", a)
			return
		} else if !accountKeyEq(aa, &a) {
			t.Errorf("Difference detected: %#v != %#v", a, *aa)
			return
		}
	}

	var ok bool
	
	// update one attrubute of one created account
	testBscsAccount1 := "X0"
	testOfiSapAccount1 := "XYZ"
	ok = accountUpdateAttributeOne(t, s, c, token, "W", 0, testBscsAccount1, "ofiSapAccount", testOfiSapAccount1) 
	if !ok {
		t.Errorf("Can not update account: " + testBscsAccount1)
		return
	}

	// update one attrubute of one created account
	testBscsAccount2 := "X1"
	testOfiSapAccount2 := "ABC"
	ok = accountUpdateAttributeOne(t, s, c, token, "W", 0, testBscsAccount2, "ofiSapAccount", testOfiSapAccount2)
	if !ok {
		t.Errorf("Can not update account: " + testBscsAccount2)
		return
	}

	// update one attrubute of one created account
	testBscsAccount3 := "X2"
	testVatCodeInd3 := "A"
	ok = accountUpdateAttributeOne(t, s, c, token, "W", 0, testBscsAccount3, "vatCodeInd", testVatCodeInd3)
	if !ok {
		t.Errorf("Can not update account: " + testBscsAccount3)
		return
	}

	// update one attrubute of one created account
	//testBscsAccount4 := "X3"
	//testOfiSapWbsCode4 := "B"
	//ok = accountUpdateAttributeOne(t, s, c, token, "W", 0, testBscsAccount4, "ofiSapWbsCode", testOfiSapWbsCode4)
	//if !ok {
	//	t.Errorf("Can not update account: " + testBscsAccount4)
	//	}

	// update one attrubute of one created account
	testBscsAccount5 := "X4"
	testCitMarkerVatFlag5 := 1
	ok = accountUpdateAttributeOne(t, s, c, token, "W", 0, testBscsAccount5, "citMarkerVatFlag", testCitMarkerVatFlag5)
	if !ok {
		t.Errorf("Can not update account: " + testBscsAccount5)
		return
	}
	
	// same number returned as created
	accounts := accountsRead(t, s, c, token, "W", 0)
	if len(*accounts) != len(test) {
		t.Errorf("Difference detected: %d (tested) != %d (created)", len(test), len(*accounts))
		return
	}

	// same accounts returned as in test
	var matches int = 0
	for _, a := range *accounts {
		if a.BscsAccount == testBscsAccount1 {
			if a.OfiSapAccount == testOfiSapAccount1 {
				matches++
			}
		}
		if a.BscsAccount == testBscsAccount2 {
			if a.OfiSapAccount == testOfiSapAccount2 {
				matches++
			}
		}
		if a.BscsAccount == testBscsAccount3 {
			if a.VatCodeInd == testVatCodeInd3 {
				matches++
			}
		}
		//if a.BscsAccount == testBscsAccount4 {
		//	if a.OfiSapWbsCode == testOfiSapWbsCode4 {
		//		matches++
		//	}
		//}
		if a.BscsAccount == testBscsAccount5 {
			if a.CitMarkerVatFlag == testCitMarkerVatFlag5 {
				matches++
			}
		}
	}

	if matches != 5 {
		t.Errorf("Can't find some of the updated records, expeced: 6, found: %d", matches)
		return
	}
}

