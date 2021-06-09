package common

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

func loadRSAPrivateKeyFromDisk(location string) *rsa.PrivateKey {
	keyData, e := ioutil.ReadFile(location)
	if e != nil {
		panic(e.Error())
	}
	key, e := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if e != nil {
		panic(e.Error())
	}
	return key
}

func loadRSAPublicKeyFromDisk(location string) *rsa.PublicKey {
	keyData, e := ioutil.ReadFile(location)
	if e != nil {
		panic(e.Error())
	}
	key, e := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if e != nil {
		panic(e.Error())
	}
	return key
}

func makeJWToken(c jwt.Claims, key interface{}) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	s, err := token.SignedString(key)
	if err != nil {
		panic(err.Error())
	}
	return s
}

//
// Private key for signing and public key for verification
//
var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	// openssl genrsa -out app.rsa 1024
	privateKeyPath string
	// openssl rsa -in app.rsa -pubout > app.rsa.pub
	publicKeyPath string
)

func GetPrivateKey() *rsa.PrivateKey { return privateKey }

func GetPublicKey() *rsa.PublicKey { return publicKey }

//
// Read the key files before starting http handlers, may panic as it is init phase
//
func initKeys() {
	log.Printf("Loading keys from: %s", AppConfig.KeyPath)

	// Load private sign key
	privateKeyPath = AppConfig.KeyPath + "/" + "app.rsa"
	privateKey = loadRSAPrivateKeyFromDisk(privateKeyPath)

	// Load public sign key
	publicKeyPath = AppConfig.KeyPath + "/" + "app.rsa.pub"
	publicKey = loadRSAPublicKeyFromDisk(publicKeyPath)
}

//
// Generate JWT token containing user name role and name used between sessions
//
func GenerateJWToken(user, role string) (token string, expiry time.Time, err error) {
	log.Printf("Start generate JWT token for: user:%s, role:%s", user, role)

	var validity int
	validity, err = strconv.Atoi(AppConfig.JWTTokenValidHours)
	if err != nil {
		err = fmt.Errorf("Invalid config parameter JWTTokenValidHours format: %s", AppConfig.JWTTokenValidHours)
		return
	}
	expiry = time.Now().Add(time.Hour * time.Duration(validity))
	var claims = &jwt.MapClaims{
		"user": user,
		"role": role,
		"exp":  expiry,
		"iat":  time.Now().Unix(),
	}
	token = makeJWToken(claims, privateKey)
	log.Printf("Produced JWT token: %s", token)

	return
}

//
// Middleware for validating JWT tokens with public key
// It loads role and name of the user as the side effect
//
func WithAuthorize(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Printf("Start JWT Authorize on metod: %s", r.Method)

	// Must be handlet first as it si always called
	if r.Method == "OPTIONS" {
		log.Printf("JWT Authorize skip: OPTIONS received, quiting")
		next(w, r)
		return
	}

	// no CORS preflight request so we get Autorization header with token
	extractor := request.AuthorizationHeaderExtractor
	keyfunc := func(*jwt.Token) (interface{}, error) {
		return publicKey, nil
	}

	log.Printf("Parsing JWT token")
	token, err := request.ParseFromRequestWithClaims(r, extractor, jwt.MapClaims{}, keyfunc)
	if token == nil {
		DisplayAppError(w, AuthorizationError, "Token not found", http.StatusInternalServerError)
		return
	} else {
		log.Printf("Received JWT Token: %v", token)
	}

	// JWT token errors caught
	if err != nil {
		log.Printf("JWT token error handling")
		switch err.(type) {
		case *jwt.ValidationError: // JWT validation error
			e := err.(*jwt.ValidationError)
			switch e.Errors {
			case jwt.ValidationErrorExpired: //JWT expired
				DisplayAppError(w, AuthorizationError, "Access Token is expired, get a new Token", http.StatusUnauthorized)
				return
			default:
				DisplayAppError(w, AuthorizationError, "Error while parsing the Access Token: " + e.Error(), http.StatusInternalServerError)
				return
			}
		default:
			DisplayAppError(w, AuthorizationError, "Error while parsing Access Token", http.StatusInternalServerError)
			return
		}
	}

	if token.Valid {
		// Valid token but wrong signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			DisplayAppError(w, AuthorizationError, fmt.Sprintf("Invalid token, unexpected signing method: %v", token.Header["alg"]), http.StatusInternalServerError)
		}
		
		log.Printf("JWT token is valid")

		// needed in the controller to open the context
		var claims jwt.MapClaims = token.Claims.(jwt.MapClaims)

		// Get role from claims
		if role, ok := claims["role"].(string); !ok {
			DisplayAppError(w, AuthorizationError, "Invalid token: no role found in token claims", http.StatusUnauthorized)
			return
		} else {
			log.Printf("JWT Claim role: %s", role)
			r.Header.Set("role", role)
		}

		var tester bool = false
		
		// Get user from claims
		if user, ok := claims["user"].(string); !ok {
			DisplayAppError(w, AuthorizationError, "Invalid token: no user found in token claims", http.StatusUnauthorized)
			return
		} else {
			log.Printf("JWT Claim user: %s", user)
			r.Header.Set("user", user)
			if user == "TEST" {
				tester = true
			}
		}
		
		// expiry date from claims
		if exp, ok := claims["exp"].(string); !ok {
			DisplayAppError(w, AuthorizationError, "Invalid token: no exp field found in token claims", http.StatusUnauthorized)
			return
		} else if !(AppConfig.Testing == "Y" && tester) { // backdoor
			log.Printf("JWT Claim exp date: %s", exp)
			now := time.Now()
			if expiry, err := time.Parse(TokenDateFormat, exp); err != nil {
				DisplayAppError(w, AuthorizationError, fmt.Sprintf("Invalid token: Can't parse date format: %s", exp), http.StatusUnauthorized)
				return
			} else if !IsBefore(now, expiry) {
				DisplayAppError(w, AuthorizationError, fmt.Sprintf("Invalid token: expired"), http.StatusUnauthorized)
				return
			}
		}		
	} else {
		DisplayAppError(w, AuthorizationError, "Invalid Access Token", http.StatusUnauthorized)
		return
	}

	log.Printf("JWT Authorize successful")
	
	next(w, r)
}
