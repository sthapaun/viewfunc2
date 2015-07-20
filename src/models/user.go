// user
package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"conf"
)
/*
func main() {
	fmt.Println("Hello World!")
}
*/

type User struct {
	Id         string `json: "id"`
	Name       string `json: "name"`
	GivenName  string `json: "given_name"`
	FamilyName string `json: "family_name"`
	Link       string `json: "link"`
	Picture    string `json: "picture"`
	Gender     string `json: "gender"`
	Locale     string `json: "locale"`
}

var defaultUser = User{
	Id:         "0",
	Name:       "Guest",
	GivenName:  "Guest",
	FamilyName: " ",
	Link:       " ",
	Picture:    " ",
	Gender:     " ",
	Locale:     "en",
}

func GetUser(r *http.Request) User {
	publicKey, e := ioutil.ReadFile("static/demo.rsa.pub")
	if e != error(nil) {
		fmt.Println("Failure to read public key: %v", e)
		os.Exit(1)
	}
	var user User
	cookie, err := r.Cookie(conf.CookieName)
	if err == error(nil) {

		tokenString := cookie.Value
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})
		obj := token.Claims["User"]

		j, _ := json.Marshal(obj)
		json.Unmarshal(j, &user)
	} else {
		user = defaultUser
	}
	return user
}

