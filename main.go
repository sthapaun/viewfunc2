// main
package main

import (
	"fmt"
	"net/http"
	"controllers"
)

func main() {
	//	fmt.Println("Hello World!")
	host := "127.0.0.1:8888"
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/oauth2callback", controllers.Callback)
	http.HandleFunc("/main", controllers.Pastebin)
	http.HandleFunc("/new", controllers.Create)
	http.HandleFunc("/paste/", controllers.Show)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/", controllers.Index)
	fmt.Println("Opening", host)
	http.ListenAndServe(host, (http.Handler)(nil))
}
