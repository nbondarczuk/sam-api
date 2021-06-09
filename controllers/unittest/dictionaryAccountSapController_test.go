package controllers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"sam-api/common"
)

//
// scenario: create one entry
//
func TestDictionaryAccountSapCreate(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	// prepare payload for Account Create
	body := []byte("{\"data\":{\"sapOfiAccount\":\"123\",\"name\":\"Whatever\",\"status\":\"C\"}}")
	req, err := http.NewRequest("POST", server.URL + "/api/dictionary/account/sap", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for DictionaryAccountSapCreate: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to DictionaryAccountSapCreate: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected response status %d, received %d", http.StatusCreated, res.StatusCode)
		return
	}
}

//
// scenario: read the whole contnt of SA_OPI_ACCOUNTS
//

func TestDictionaryAccountSapReadSome(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Order Create
	req, err := http.NewRequest("GET", server.URL + "/api/dictionary/account/sap", nil)
	if err != nil {
		t.Errorf("Error in GET: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in GET: %v", err)
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
// scenario: clean the resource
//
func TestDictionaryAccountSapDeleteAll(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()

	// prepare payload for Account Create
	req, err := http.NewRequest("DELETE", server.URL + "/api/dictionary/account/sap", nil)
	if err != nil {
		t.Errorf("Error in creating DELETE request for DictionaryAccountSapDeleteAll: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + token)

	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in DELETE on DictionaryAccountSapDeleteAll: %v", err)
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
// scenario: create Excel based config
//
func TestDictionaryAccountSapCreateExcel(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()

	// prepare payload for Account Create
	ffn := common.AppConfig.RunPath + "/examples/" + "SAP_SEGMENT_ACCOUNTS_SHORT.xlsx"
	body, err := ioutil.ReadFile(ffn)
	if err != nil {
		t.Errorf("Error reading file: %s - %s", ffn, err)
		return		
	}
	
	req, err := http.NewRequest("POST", server.URL + "/api/dictionary/account/sap", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for DictionaryAccountSapCreate: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/xlsx")
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to DictionaryAccountSapCreate: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected response status %d, received %d", http.StatusCreated, res.StatusCode)
		return
	}
}
