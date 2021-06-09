package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"sam-api/common"
	"sam-api/models"
	"sam-api/repository"
	"sam-api/resources"
	"sam-api/valid"
)

//
// Handler for /api/user/login
//
func UserLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start UserLogin")

	var err error
	
	// Decode the incoming Login json
	var dataResource resources.LoginResource
	err = json.NewDecoder(r.Body).Decode(&dataResource);
	if err != nil {
		common.DisplayAppError(w, fmt.Errorf("Invalid Login data"), "Invalid Json Format", http.StatusInternalServerError)
		return
	}

	loginModel := dataResource.Data
	loginUser := models.User{
		User:     loginModel.User,
		Role:     loginModel.Role,
		//Password: loginModel.Password,
	}

	log.Printf("Login user:%s, role:%s", loginUser.User, loginUser.Role)

	if !valid.IsUserValid(loginUser.User) {
		common.DisplayAppError(w, fmt.Errorf("Invalid user: %s", loginUser.User), "Login Error", http.StatusUnauthorized)
		return
	}

	if !valid.IsRoleValid(loginUser.Role) {
		common.DisplayAppError(w, fmt.Errorf("Invalid role: %s", loginUser.Role), "Login Error", http.StatusUnauthorized)
		return
	}

	if common.AppConfig.Testing == "Y" && loginUser.User != "TEST" { // backdoor
		// Authenticate the login user with password
		err = common.LdapServerAuth(loginUser.User, loginUser.Role, loginUser.Password);
		if err != nil {
			common.DisplayAppError(w, fmt.Errorf("Invalid authentication in LDAP server for user: %s, %#v", loginUser.User, err), "LDAP Error", http.StatusUnauthorized)
			return
		}
	}

	// Internally register user as logged in
	err = repository.UserLogin(loginUser.User)
	if err != nil {
		common.DisplayAppError(w, fmt.Errorf("Error while registering user: %s", loginUser.User), "User Login Error", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, expiry, err := common.GenerateJWToken(loginUser.User, loginUser.Role)
	if err != nil {
		common.DisplayAppError(w, fmt.Errorf("Error while gererating JWT token for user: %s", loginUser.User), "JWT Error", http.StatusInternalServerError)
		return
	}

	// build payload of the response
	authUser := resources.AuthUserModel{
		User:  loginUser,
		Token: token,
	}
	j, err := json.Marshal(resources.AuthUserResource{Data: authUser})
	if err != nil {
		common.DisplayAppError(w, fmt.Errorf("An unexpected error has occurred"), "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Build headers with cors
	w.Header().Set("X-Expires-After", expiry.Format(common.TokenDateFormat))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")	
	common.SetupCorsResponse(&w, false)

	// flush payload
	w.WriteHeader(http.StatusOK)
	w.Write(j)

	log.Printf("Finished UserLogin, status: %d, response: %#v", http.StatusOK, w)
}

//
// Handler for /api/user/relogin
//
func UserRelogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start UserRelogin")

	// Attributes already set by Authorize check so we can reuse them
	user := r.Header.Get("user")
	role := r.Header.Get("role")
	log.Printf("Relogin of existing user:%s, role:%s", user, role)

	// Generate JWT token
	token, expiry, err := common.GenerateJWToken(user, role)
	if err != nil {
		common.DisplayAppError(w, err, "JWT Error", http.StatusInternalServerError)
		return
	}	
	log.Printf("New token granted user:%s, role:%s, token: %s", user, role, token)

	// Build json response with token and user id
	loginUser := models.User{
		User: user,
		Role: role,
	}
	authUser := resources.AuthUserModel{
		User:  loginUser,
		Token: token,
	}
	j, err := json.Marshal(resources.AuthUserResource{Data: authUser})
	if err != nil {
		common.DisplayAppError(w, err, "Json Encoding Error", http.StatusInternalServerError)
		return
	}

	// make headers
	w.Header().Set("X-Expires-After", expiry.Format(common.TokenDateFormat))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	common.SetupCorsResponse(&w, false)

	// flush data
	w.WriteHeader(http.StatusOK)
	w.Write(j)

	log.Printf("Finished UserRelogin, status: %d, response: %#v", http.StatusOK, w)
}

//
// Handler for /api/user/logoff
//
func UserLogoff(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start UserLogoff")

	// Internally deregister user
	user := r.Header.Get("user")
	log.Printf("Logoff user:%s", user)
	if err := repository.UserLogoff(user); err != nil {
		common.DisplayAppError(w, err, "Eror while deregistering user:"+user, http.StatusInternalServerError)
		return
	}

	common.SetupCorsResponse(&w, false)
	w.WriteHeader(http.StatusOK)

	log.Printf("Finished UserLogoff, status: %d, response: %#v", http.StatusOK, w)
}

//
// Handler for /api/user/info
//
func UserInfo(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start UserInfo")

	// Attributes already set by Authorize check so we can reuse them
	user := r.Header.Get("user")
	role := r.Header.Get("role")
	log.Printf("Info existing user:%s, role:%s", user, role)

	// Build json response with token and user id
	loginUser := models.User{
		User: user,
		Role: role,
	}
	authUserInfo := resources.AuthUserInfoModel{
		User:  loginUser,
	}
	j, err := json.Marshal(resources.AuthUserInfoResource{Data: authUserInfo})
	if err != nil {
		common.DisplayAppError(w, err, "An unexpected error has occurred", http.StatusInternalServerError)
		return
	}

	// make headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	common.SetupCorsResponse(&w, false)

	// flush data
	w.WriteHeader(http.StatusOK)
	w.Write(j)

	log.Printf("Finished UserInfo, status: %d, response: %#v", http.StatusOK, w)
}
