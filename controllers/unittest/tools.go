package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http/httptest"
	"net/http"
	"os"
	"testing"
	"time"
	
	"sam-api/common"
	"sam-api/models"
	"sam-api/resources"
	"sam-api/routers"
)

func initTestEnv(t *testing.T, user, role string, silent bool) (client *http.Client, server *httptest.Server, token string) {
	// mock server & client for testing the Login service
	if os.Getenv("DEBUG") == "1" {
		common.LogInit(false)
	} else {
		common.LogInit(silent)
	}
	common.EnvInit("test", "test", "test")
	common.StartUp()

	// build enviroment
	client = &http.Client{}
	router := routers.InitRoutes()
	server = httptest.NewServer(router)	
	token = userLogin(t, user, role)
	
	return client, server, token
}


func userLogin(t *testing.T, user, role string) (token string) {
	client := &http.Client{}
	server := httptest.NewServer(http.HandlerFunc(createLoginHandler(userFormatter)))
	defer server.Close()
	
	// prepare payload for Login
	body := []byte("{\"data\":{\"user\": \"" + user + "\", \"role\": \"" + role + "\"}}")
	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for Login: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to Login: %v", err)
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != 200 {
		t.Errorf("Expected response status 200, received %s", res.Status)
	}

	// get token
	dataResource := resources.AuthUserResource{}
	err = json.NewDecoder(res.Body).Decode(&dataResource)
	if err != nil {
		t.Errorf("Expected AuthUserResource json: " + err.Error())
	}

	// if exists in json
	if !(len(dataResource.Data.Token) > 0) {
		t.Errorf("Expected token in json")
	}
	
	return dataResource.Data.Token
}

func newAccountId() string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	testAccountId = fmt.Sprintf("%d", r.Int())	
	return testAccountId
}

func accountCreate(t *testing.T, s *httptest.Server, c *http.Client, token string, a *models.Account) *models.Account {
	var dataReplyResource = resources.AccountReplyResource{Data: *a}
	j, err := json.Marshal(dataReplyResource);
	if err != nil {
		t.Errorf("Error in creating request for AccountCreate: %v", err)
		return nil
	}
	
	req, err := http.NewRequest("POST", s.URL + "/api/account", bytes.NewBuffer(j))
	if err != nil {
		t.Errorf("Error in creating POST request for AccountCreate: %v", err)
		return nil
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	
	// send test case to server
	res, err := c.Do(req)
	if err != nil {
		t.Errorf("Error in POST to AccountCreate: %v", err)
		return nil
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected response status %d, received %d", http.StatusCreated, res.StatusCode)
		return nil
	}

	dataResource := resources.AccountReplyResource{}
	err = json.NewDecoder(res.Body).Decode(&dataResource)
	if err != nil {
		t.Errorf("Expected AccountReplyResource json: " + err.Error())
		return nil
	}

	return &dataResource.Data
}

func accountsRead(t *testing.T, s *httptest.Server, c *http.Client, token string, status string, release int) *[]models.Account {
	route := fmt.Sprintf("/api/account/%s/%d", status, release)
	req, err := http.NewRequest("GET", s.URL + route, nil)
	if err != nil {
		t.Errorf("Error in GET request for AccountReadSome: %v", err)
		return nil
	}	
	req.Header.Add("Authorization", token)
	
	// send test case to server
	res, err := c.Do(req)
	if err != nil {
		t.Errorf("Error in GET for AccountReadSome: %v", err)
		return nil
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response status %d, received %d", http.StatusOK, res.StatusCode)
		return nil
	}

	dataResource := resources.AccountsReplyResource{}
	err = json.NewDecoder(res.Body).Decode(&dataResource)
	if err != nil {
		t.Errorf("Expected AccountResource json: " + err.Error())
		return nil
	}

	return &dataResource.Data
}

func accountKeyEq(a, b *models.Account) bool {
	return a.Status == b.Status &&
		a.ReleaseId == b.ReleaseId &&
		a.BscsAccount == b.BscsAccount
}

func accountUpdateAttributeOne(t *testing.T, s *httptest.Server, c *http.Client, token string, status string, release int, account string, attribute string, value interface{}) bool {
	route := fmt.Sprintf("/api/account/%s/%d/%s", status, release, account)
	var body []byte
	switch value.(type) {
	case string:
		body = []byte("{\"data\":{\"" + attribute + "\":\"" + fmt.Sprintf("%s", value) + "\"}}")
	case int:
		body = []byte("{\"data\":{\"" + attribute + "\":" + fmt.Sprintf("%d", value) + "}}")
	case time.Time:
		body = []byte("{\"data\":{\"" + attribute + "\":\"" + time.Time(value.(time.Time)).Format("2006-01-02T15:04:05Z") + "\"}}")
	default:
		t.Errorf("Invalid value type: %T", value)
		return false
	}
	req, err := http.NewRequest("PATCH", s.URL + route, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in PATCH request for AccountUpdateAttributeOne: %v", err)
		return false
	}
	req.Header.Add("Authorization", token)
	
	// send test case to server
	res, err := c.Do(req)
	if err != nil {
		t.Errorf("Error in PATCH for AccountUpdateAttributeOne: %v", err)
		return false
	}
	defer res.Body.Close()

	// check result(s)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response status %d, received %d", http.StatusOK, res.StatusCode)
		return false
	}

	dataResource := resources.AccountsReplyResource{}
	err = json.NewDecoder(res.Body).Decode(&dataResource)
	if err != nil {
		t.Errorf("Expected AccountResource json: " + err.Error())
		return false
	}

	// exactly 1 record was upated
	if dataResource.Count != 1 {
		t.Errorf("Expected updated records 1, got: %d", dataResource.Count)
	}
	
	return true
}

