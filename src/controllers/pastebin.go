// pastebin
package controllers

import (
	"conf"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"models"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"code.google.com/p/goauth2/oauth"
	"github.com/dgrijalva/jwt-go"
)

/*
func main() {
	fmt.Println("Hello World!")
}
*/
const profileInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"

func Callback(w http.ResponseWriter, r *http.Request) {
	//Get the code from the response
	code := r.FormValue("code")

	t := &oauth.Transport{Config: &conf.OauthCfg}

	// Exchange the received code for a token
	t.Exchange(code)

	//now get user data based on the Transport which has the token
	resp, _ := t.Client().Get(profileInfoURL)
	// strip null byte off the tail end of the resoponse body
	buf := make([]byte, 1024)
	resp.Body.Read(buf)
	if resp != (*http.Response)(nil) {
		resp.Body.Close()
	}
	str := string(buf)
	str = strings.Trim(str, "\x00")
	b := []byte(str)
	// convert the body to a Google user account
	var account models.User
	err := json.Unmarshal(b, &account)
	if err != error(nil) {
		fmt.Println("Unmarshal error:", err)
		os.Exit(1)
	}
	// convert the account data to a JSON web token
	var privateKey []byte
	privateKey, _ = ioutil.ReadFile("static/demo.rsa") // location of demo.rsa
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["User"] = account
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}
	cookie := http.Cookie{Name: conf.CookieName,
		Value:   tokenString,
		Expires: time.Now().AddDate(0, 0, 1),
	}
	http.SetCookie(w, &cookie)
	view, q := template.New("layout.html").Funcs(models.TplFuncs).ParseFiles("views/layout.html", "views/pastebin.html")
	//view, q := template.ParseFiles("views/pastebin.html", "views/layout.html")
	if q != error(nil) {
		http.Error(w, q.Error(), http.StatusInternalServerError)
	}
	if e := view.Execute(w, struct {
		Title     string
		User      models.User
		Languages map[string]string
	}{
		Title:     "Golang Pastebin",
		User:      account,
		Languages: models.Languages,
	}); e != error(nil) {
		fmt.Println(e)
		os.Exit(1)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	/*	TplFuncs := template.FuncMap{
			"Publicpastes":  models.PublicPastes,
			"Privatepastes": models.PrivatePastes,
		}
	*/
	user := models.GetUser(r)
	title := r.FormValue("title")
	content := r.FormValue("content")
	language := r.FormValue("language")
	isPublicString := r.FormValue("ispublic")
	var isPublic bool
	if isPublicString == "true" || user.Id == "0" {
		isPublic = true
	} else {
		isPublic = false
	}
	userId := user.Id
	paste := models.Paste{
		Id:       bson.NewObjectId(),
		UserId:   userId,
		Title:    title,
		Content:  content,
		Language: language,
		IsPublic: isPublic,
	}
	prism := models.Languages[language]
	//	session, _ := mgo.Dial("localhost")
	//	collection := session.DB("gopastebin3-3").C("pastes")
	session, collection, _ := conf.GetCollection(conf.PASTES)
	defer session.Close()
	collection.Insert(&paste)
	log.Println("New id:", paste.Id)
	log.Println("After insert, paste =", paste)
	//	t, err := template.ParseFiles("views/layout.html", "views/create.tpl")
	t, err := template.New("layout.html").Funcs(models.TplFuncs).ParseFiles("views/layout.html", "views/create.tpl")
	if err != error(nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//	t2 := t.Funcs(funcmap)
	t.Execute(w, struct {
		Title string
		User  models.User
		Paste models.Paste
		Prism string
	}{Title: "Verify Paste",
		User:  user,
		Paste: paste,
		Prism: prism,
	})
}

func Show(w http.ResponseWriter, r *http.Request) {
	/*	TplFuncs := template.FuncMap{
			"Publicpastes":  models.PublicPastes,
			"Privatepastes": models.PrivatePastes,
		}
	*/
	user := models.GetUser(r)
	url := r.URL.Path
	parts := strings.Split(url, "/")
	log.Println("parts[0]:", parts[0], "parts[1]:", parts[1], "parts[2]:", parts[2])
	rawId := parts[2]
	pasteId := strings.TrimLeft(rawId, "ObjectIdHex(")
	pasteId = strings.TrimRight(pasteId, ")")
	pasteId = strings.Trim(pasteId, "\"")
	realId := bson.ObjectIdHex(pasteId)
	session, collection, _ := conf.GetCollection(conf.PASTES)
	defer session.Close()
	var result models.Paste
	err := collection.Find(bson.M{"id": realId}).One(&result)
	if err != error(nil) {
		panic(err)
	}

	prism := models.Languages[result.Language]
	//	t, err2 := template.ParseFiles("views/layout.html", "views/create.tpl")
	t, err2 := template.New("layout.html").Funcs(models.TplFuncs).ParseFiles("views/layout.html", "views/create.tpl")
	if err2 != error(nil) {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
	}
	t.Execute(w, struct {
		Title string
		User  models.User
		Paste models.Paste
		Prism string
	}{
		Title: "Show Paste",
		User:  user,
		Paste: result,
		Prism: prism,
	})
}
func Pastebin(w http.ResponseWriter, r *http.Request) {
/*	TplFuncs := template.FuncMap{
		"Publicpastes":  models.PublicPastes,
		"Privatepastes": models.PrivatePastes,
	}
	*/
	user := models.GetUser(r)
	//	templ, _ := template.ParseFiles("views/layout.html", "views/create.tpl")
	templ, err := template.New("layout.html").Funcs(models.TplFuncs).ParseFiles("views/layout.html", "views/pastebin.html")
	if err != error(nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templ.Execute(w, struct {
		Title         string
		User          models.User
		Languages     map[string]string
		Publicpastes  []models.Paste
		Privatepastes []models.Paste
	}{
		Title:         "Add A Paste",
		User:          user,
		Languages:     models.Languages,
		Publicpastes:  models.PublicPastes(),
		Privatepastes: models.PrivatePastes(user.Id),
	})
}
