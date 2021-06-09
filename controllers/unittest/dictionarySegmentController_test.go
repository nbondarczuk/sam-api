package controllers

import (
	"bytes"
	"net/http"
	"testing"
)

//
// scenario: create one segment
//
func TestSegmentCreate(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	// prepare payload for Account Create
	body := []byte("{\"data\":{\"csTradeRef\":\"X\", \"segmCategory\": \"PRIV\"}}")
	req, err := http.NewRequest("POST", server.URL + "/api/dictionary/segment", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for SegmentCreate: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to SegmentCreate: %v", err)
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
// scenario: read the whole resource
//
func TestDictionarySegmentReadSome(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Order Create
	req, err := http.NewRequest("GET", server.URL + "/api/dictionary/segment", nil)
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
// scenario: update previously created item
//
func TestDictionarySegmentUpdateOne(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Order Update
	body := []byte("{\"data\":{\"csTradeRef\": \"X\", \"segmCategory\": \"PRIVATE\"}}")
	req, err := http.NewRequest("PUT", server.URL + "/api/dictionary/segment/X", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in PUT request for DictionarySegmentUpdateOne: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in PUT for DictionarySegmentUpdateOne: %v", err)
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
// scenario: update previously created item
//
func TestDictionarySegmentUpdateAttributeOne(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Order Update
	body := []byte("{\"data\":{\"segmCategory\": \"PRIV\"}}")
	req, err := http.NewRequest("PATCH", server.URL + "/api/dictionary/segment/X", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in PATCH request for DictionarySegmentUpdateOne: %v", err)
		return
	}	
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in PATCH for DictionarySegmentUpdateOne: %v", err)
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
// scenario: read the whole contnt of GLACCOUNTS
//
func TestDictionarySegmenDeleteOne(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Order Create
	req, err := http.NewRequest("DELETE", server.URL + "/api/dictionary/segment/X", nil)
	if err != nil {
		t.Errorf("Error in DELETE: %v", err)
		return
	}
	req.Header.Add("Authorization", "Bearer " + token)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in DELETE: %v", err)
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
func TestDictionarySegmentDeleteAll(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()

	// prepare payload for Account Create
	req, err := http.NewRequest("DELETE", server.URL + "/api/dictionary/segment", nil)
	if err != nil {
		t.Errorf("Error in creating DELETE request for SegmentDeleteAll: %v", err)
		return
	}
	req.Header.Add("Authorization", "Bearer " + token)

	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in DELETE on SegmentDeleteAll: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response status %d, received %d", http.StatusOK, res.StatusCode)
		return
	}
}
