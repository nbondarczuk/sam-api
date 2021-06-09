package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"sam-api/common"
)

//
// scenario: release Account, Order status W -> C as Booker
//
func TestReleaseAsBooker(t *testing.T) {
	TestAccountDeleteAll(t)
	TestOrderDeleteAll(t)
	TestAccountCreate(t)
	TestOrderCreate(t)
	
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()

	// prepare payload for Release
	req, err := http.NewRequest("POST", server.URL + "/api/release/new", nil)
	if err != nil {
		t.Errorf("Error in creating POST request for Release: %v", err)
		return
	}
	req.Header.Add("Authorization", token)

	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST on Release: %v", err)
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
// scenario: release Account, Order status C -> P as Control
//

func TestReleaseAsControl(t *testing.T) {
	TestAccountDeleteAll(t)
	TestOrderDeleteAll(t)
	TestAccountCreate(t)
	TestOrderCreate(t)
	TestReleaseAsBooker(t)

	client, server, token := initTestEnv(t, "USER", "Control", true)
	defer server.Close()

	// prepare payload for Release
	req, err := http.NewRequest("POST", server.URL + "/api/release/new", nil)
	if err != nil {
		t.Errorf("Error in creating POST request for Release: %v", err)
		return
	}
	req.Header.Add("Authorization", token)

	// send test case to server	
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST on Release: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response status %d, received %d", http.StatusOK, res.StatusCode)
		return
	}

	TestAccountDeleteAll(t)
	TestOrderDeleteAll(t)
}

//
// scenario: release Account, Order status W -> C as Booker but mail server unavailable
//
func TestReleaseAsBookerWithNoEmail(t *testing.T) {
	TestAccountDeleteAll(t)
	TestOrderDeleteAll(t)
	TestAccountCreate(t)
	TestOrderCreate(t)
	
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()

	common.AppConfig.AlertMailServerAddress = "xxx:23123"
	
	// prepare payload for Release
	req, err := http.NewRequest("POST", server.URL + "/api/release/new", nil)
	if err != nil {
		t.Errorf("Error in creating POST request for Release: %v", err)
		return
	}
	req.Header.Add("Authorization", token)

	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST on Release: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode == http.StatusOK {
		t.Errorf("Expected response status not %d, received %d", http.StatusOK, res.StatusCode)
		return
	}

	dataResource := common.ErrorResource{}
	err = json.NewDecoder(res.Body).Decode(&dataResource)
	if err != nil {
		t.Errorf("Expected error json: " + err.Error())
		return
	}

	log.Printf("Got error: %#v", dataResource)
}
