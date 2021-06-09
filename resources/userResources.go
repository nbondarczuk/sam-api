package resources

import (
	"sam-api/models"
)

//Models for JSON resources
//Models for logical model resources envelopes
type (
	//
	// User resource
	//

	//For Post - /user/login
	LoginResource struct {
		Data LoginModel `json:"data"`
	}

	//Response for authorized user Post - /api/user/login
	AuthUserResource struct {
		Data AuthUserModel `json:"data"`
	}

	//Model for authentication
	LoginModel struct {
		User         string    `json:"user"`
		Role         string    `json:"role"`
		Password     string    `json:"password,omitempty"`
	}

	//Model for authorized user with access token
	AuthUserModel struct {
		User  models.User `json:"user"`
		Token string      `json:"token"`
	}
	
	//Response for authorized user Post - /api/user/info
	AuthUserInfoResource struct {
		Data AuthUserInfoModel `json:"data"`
	}

	//Model for authorized user without access token
	AuthUserInfoModel struct {
		User  models.User `json:"user"`
	}	
)
