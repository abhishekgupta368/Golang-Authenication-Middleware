package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// reference : https://github.com/FusionAuth/fusionauth-example-go-jwt-microservices
var (
	MySigningKey = []byte(os.Getenv("SECRET_KEY"))
	authorized   = true
	client       = "Krissanawat"
	aud          = "billing.jwtgo.io"
	iss          = "jwtgo.io"
	token        = ""
)

func GetJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = authorized
	claims["client"] = client
	claims["aud"] = aud
	claims["iss"] = iss
	claims["exp"] = time.Now().Add(time.Second * 5).Unix()

	tokenString, err := token.SignedString(MySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf(("Invalid Signing Method"))
				}
				aud := "billing.jwtgo.io"
				checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
				if !checkAudience {
					return nil, fmt.Errorf(("invalid aud"))
				}
				// verify iss claim
				iss := "jwtgo.io"
				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !checkIss {
					return nil, fmt.Errorf(("invalid iss"))
				}
				return MySigningKey, nil
			})
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			if token.Valid {
				endpoint(w, r)
			}
		} else if token != "" {
			token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf(("Invalid Signing Method"))
				}
				aud := "billing.jwtgo.io"
				checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
				if !checkAudience {
					return nil, fmt.Errorf(("invalid aud"))
				}
				// verify iss claim
				iss := "jwtgo.io"
				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !checkIss {
					return nil, fmt.Errorf(("invalid iss"))
				}
				return MySigningKey, nil
			})
			if err != nil {
				log.Println(err.Error())
				fmt.Fprintf(w, err.Error())
			}
			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "No Authorization Token provided")
		}
	})
}

func GenerateAuthToken(w http.ResponseWriter, r *http.Request) {
	validToken, err := GetJWT()
	if err != nil {
		fmt.Println("Failed to generate token")
	}
	token = validToken
	fmt.Println(token)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token Created"))
}

func ClearAuthToken(w http.ResponseWriter, r *http.Request) {
	token = ""
	w.Write([]byte("Token Cleared"))
}
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Super Secret Information")
}

func main() {
	//Generate token for the system
	http.HandleFunc("/genToken", GenerateAuthToken)
	// Clear token form the system
	http.HandleFunc("/clearToken", ClearAuthToken)

	// Try to access the token from the system
	http.Handle("/homepage", isAuthorized(http.HandlerFunc(homePage)))

	log.Println("Starting the service")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
