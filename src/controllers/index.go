// index
package controllers

import (
	//	"fmt"
	"html/template"
	"net/http"
	"conf"
)

/*
func main() {
	fmt.Println("Hello World!")
}
*/

func Index(w http.ResponseWriter, r *http.Request) {
	view, err := template.ParseFiles("views/index.html")
	if err != error(nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	//Get the Google URL which shows the Authentication page to the user
	url := conf.OauthCfg.AuthCodeURL("")

	//redirect user to that page
	http.Redirect(w, r, url, http.StatusFound)
}
