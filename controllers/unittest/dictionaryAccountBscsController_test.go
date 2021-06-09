package controllers

import (
	"net/http"
	"testing"
)

//
// scenario: read the whole contnt of GLACCOUNTS
//

func TestDictionaryAccountBscsReadSome(t *testing.T) {
	client, server, token := initTestEnv(t, "USER", "Booker", true)
	defer server.Close()
	
	// prepare payload for Order Create
	req, err := http.NewRequest("GET", server.URL + "/api/dictionary/account/bscs", nil)
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

