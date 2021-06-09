package controllers

import (
	"bytes"
	"github.com/unrolled/render"
	"net/http"
	"testing"
)

var (
	orderFormatter = render.New(render.Options{
		IndentJSON: true,
	})
)

var testSegmentCode string

//
// scenario: clean the resource
//
func TestOrderDeleteAll(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Control", true)
	defer server.Close()
	
	// prepare payload for Order Create
	req, err := http.NewRequest("DELETE", server.URL + "/api/order", nil)
	if err != nil {
		t.Errorf("Error in creating DELETE request for OrderDeleteAll: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + token)

	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in DELETE on OrderDeleteAll: %v", err)
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
// scenario: order create on empty table
//
func TestOrderCreate(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Control", true)
	defer server.Close()

	// prepare payload for Order Create
	testAccountId = newAccountId()
	testSegmentCode = "XXX"
	body := []byte("{\"data\":{\"status\": \"W\", \"releaseId\": \"0\", \"bscsAccount\": \"" + testAccountId + "\", \"segmentCode\": \"" + testSegmentCode + "\"}}")
	req, err := http.NewRequest("POST", server.URL + "/api/order", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for OrderCreate: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to OrderCreate: %v", err)
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
// scenario: an order previously created in test can be read with GET
//
func TestOrderReadSome(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Order Create
	req, err := http.NewRequest("GET", server.URL + "/api/order/W/0", nil)
	if err != nil {
		t.Errorf("Error in GET request for OrderReadSome: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in GET for OrderReadSome: %v", err)
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
// scenario: an order previously created in test can be read with GET last
//
func TestOrderReadSomeWithLast(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Order Create
	req, err := http.NewRequest("GET", server.URL + "/api/order/W/last", nil)
	if err != nil {
		t.Errorf("Error in GET request for OrderReadSome: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in GET for OrderReadSome: %v", err)
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
// scenario: an order previously created in test can be read with GET latest
//
func TestOrderReadSomeWithLatest(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Order Create
	req, err := http.NewRequest("GET", server.URL + "/api/order/W/latest", nil)
	if err != nil {
		t.Errorf("Error in GET request for OrderReadSome: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in GET for OrderReadSome: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode == http.StatusOK {
		t.Errorf("Expected response status not %d, received %d", http.StatusOK, res.StatusCode)
		return
	}
}

//
// scenario: an order previously created in test can be read with GET
//
func TestOrderReadAll(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Order Create
	req, err := http.NewRequest("GET", server.URL + "/api/order", nil)
	if err != nil {
		t.Errorf("Error in GET request for OrderReadSome: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in GET for OrderReadSome: %v", err)
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
// scenario: an order previously created in test can be read with GET
//
func TestOrderReadAllAsJson(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Order Create
	req, err := http.NewRequest("GET", server.URL + "/api/order", nil)
	if err != nil {
		t.Errorf("Error in GET request for OrderReadSome: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in GET for OrderReadSome: %v", err)
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
// scenario: an order previously created in test can be read with GET
//
func TestOrderReadAllAsCsv(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Order Create
	req, err := http.NewRequest("GET", server.URL + "/api/order", nil)
	if err != nil {
		t.Errorf("Error in GET request for OrderReadSome: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/csv")	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in GET for OrderReadSome: %v", err)
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
// scenario: update previously created order
//
func TestOrderUpdateOne(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Order Update
	testSegmentCode = "XXX"
	body := []byte("{\"data\":{\"status\": \"W\", \"releaseId\": \"0\", \"bscsAccount\": \"" + testAccountId + "\", \"segmentCode\": \"" + testSegmentCode + "\"}}")
	req, err := http.NewRequest("PUT", server.URL + "/api/order/W/0/"+ testAccountId + "/" + testSegmentCode, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in PUT request for OrderUpdateOne: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in PUT for OrderUpdateOne: %v", err)
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
// scenario: update one attribute of previously created order
//
func TestOrderUpdateAttributeOne(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Control", true)
	defer server.Close()
	
	// prepare payload for Order Update
	body := []byte("{\"data\":{\"orderNumber\":" + "\"XXX\"" + "}}")
	req, err := http.NewRequest("PATCH", server.URL + "/api/order/W/0/"+ testAccountId + "/" + testSegmentCode, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in PATCH request for OrderUpdateAttributeOne: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in PATCH for OrderUpdateAttributeOne: %v", err)
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
// scenario: delete one order
//
func TestOrderDeleteOne(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Control", true)
	defer server.Close()
	
	// prepare payload for OrderDeleteOne
	testSegmentCode = "XXX"
	body := []byte("{\"data\":{\"status\": \"W\", \"releaseId\": \"0\", \"bscsAccount\": \"" + testAccountId + "\", \"segmentCode\": \"" + testSegmentCode + "\"}}")
	req, err := http.NewRequest("DELETE", server.URL + "/api/order/W/0/" + testAccountId + "/" + testSegmentCode, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in DELETE request for OrderDeleteOne: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in DELETE for OrderDeleteOne: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response status %d, received %d", http.StatusOK, res.StatusCode)
		return
	}
}
