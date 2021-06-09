package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/unrolled/render"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"

	"sam-api/common"
	"sam-api/controllers"
	"sam-api/repository"	
	"sam-api/resources"
)

var (
	userFormatter = render.New(render.Options{
		IndentJSON: true,
	})
)

//
// scenario: simple login as USER/Booker
//

func createLoginHandler(formatter *render.Render) http.HandlerFunc {
	return controllers.UserLogin
}

func TestLogin(t *testing.T) {
	// mock server & client for testing the Login service
	common.LogInit(true)
	common.EnvInit("test", "test", "test")
	common.StartUp()
	client := &http.Client{}
	server := httptest.NewServer(http.HandlerFunc(createLoginHandler(userFormatter)))
	defer server.Close()

	user := "USER" 
	role := "Booker"
		
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

	// correct status?
	if res.StatusCode != 200 {
		t.Errorf("Expected response status 200, received %s", res.Status)
	}
	
	// get token
	dataResource := resources.AuthUserResource{}
	err = json.NewDecoder(res.Body).Decode(&dataResource)
	if err != nil {
		t.Errorf("Expected AuthUserResource json: " + err.Error())
	}

	// check if token exists in json
	if !(len(dataResource.Data.Token) > 0) {
		t.Errorf("Expected token in json")
	}

	// decode and check the validity of the token
	tokenStr := dataResource.Data.Token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			key := common.GetPublicKey()
			return key, nil
		})
	
	if token == nil {
		t.Errorf("No token in response payload")
	} else if !token.Valid {
		t.Errorf("Invalid token: " + err.Error())
	}

	// check claims if original user and role are there encoded
	if role, ok := claims["role"].(string); !ok {
		t.Errorf("Expected role in claims, not found")
	} else if role != "Booker" {
		t.Errorf("Expected role: " + role)
	}
	if user, ok := claims["user"].(string); !ok {
		t.Errorf("Expected user: " + user)
	} else if user != "USER" {
		t.Errorf("Expected user: " + user)
	}

	// check if the user was added to repository
	if !repository.UserExist(user) {
		t.Errorf("User not removed from repository: %s", user)
	}
}

//
// scenario: login successfull, log off attempt with token set, status ok
//

func createLogoffHandler(formatter *render.Render) http.HandlerFunc {
	return controllers.UserLogoff
}

func TestLogoffWithTokenSet(t *testing.T) {
	// mock server & client for testing the Login service
	common.LogInit(true)
	common.EnvInit("test", "test", "test")
	common.StartUp()
	client := &http.Client{}
	loginServer := httptest.NewServer(http.HandlerFunc(createLoginHandler(userFormatter)))
	defer loginServer.Close()

	user := "USER"
	role := "Booker"	
	token := userLogin(t, user, role)

	logoffServer := httptest.NewServer(http.HandlerFunc(createLogoffHandler(userFormatter)))
	defer logoffServer.Close()
	
	// prepare payload for Logoff with empty body
	body := []byte("{}")
	req, err := http.NewRequest("POST", logoffServer.URL, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for Logoff: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + token)

	// to simulate work of decorator WithAuthorize 
	req.Header.Set("user", user)
	// to simulate state enforced by previous login
	repository.UserLogin(user)
	
	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to Logoff: %v", err)
		return
	}
	defer res.Body.Close()
	
	// check result(s)

	// correct status?
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response status %d, received %d", http.StatusOK, res.StatusCode)
		return
	}

	// check if user was removed from repository 
	if repository.UserExist(user) {
		t.Errorf("User not removed from repository: %s", user)
		return
	}
}

//
// scenario: invalid role during login
//

func TestLoginInvalidRole(t *testing.T) {
	// mock server & client for testing the Login service
	common.LogInit(true)
	common.EnvInit("test", "test", "test")
	common.StartUp()
	client := &http.Client{}
	server := httptest.NewServer(http.HandlerFunc(createLoginHandler(userFormatter)))
	defer server.Close()

	user := "Nobody" 
	role := "Nobody"
		
	// prepare payload for Login
	body := []byte("{\"data\":{\"user\": \"" + user + "\", \"role\": \"" + role + "\"}}")
	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for Login: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	// send test case to server
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to Login: %v", err)
		return
	}
	defer res.Body.Close()

	// check result(s)

	// correct status?
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected response: %d, received %s", http.StatusUnauthorized, res.Status)
		return
	}
}
