// paste
package models

import (
	//	"fmt"
	"conf"
	"html/template"

	"gopkg.in/mgo.v2/bson"
)

/*
func main() {
	fmt.Println("Hello World!")
}
*/
var TplFuncs = template.FuncMap{
	"Publicpastes":  PublicPastes,
	"Privatepastes": PrivatePastes,
}

type Paste struct {
	Id       bson.ObjectId
	UserId   string
	Title    string
	Content  string
	Language string
	IsPublic bool
}

func PublicPastes() []Paste {
	/*	session, e := mgo.Dial("localhost")
		if e != nil {
			panic(e)
		}
		collection := session.DB("gopastebin3-3").C("pastes")
	*/
	session, collection, err := conf.GetCollection(conf.PASTES)
	defer session.Close()
	if err != error(nil) {
		panic(err)
	}
	var pastes []Paste
	err = collection.Find(bson.M{"ispublic": true}).All(&pastes)
	if err != error(nil) {
		panic(err)
	}
	return pastes
}

func PrivatePastes(id string) []Paste {
	/*	session, e := mgo.Dial("localhost")
		if e != nil {
			panic(e)
		}
		collection := session.DB("gopastebin3-3").C("pastes")
	*/
	session, collection, err := conf.GetCollection(conf.PASTES)
	defer session.Close()
	if err != error(nil) {
		panic(err)
	}
	var pastes []Paste
	err = collection.Find(bson.M{"userid": id, "ispublic": false}).All(&pastes)
	if err != error(nil) {
		panic(err)
	}
	return pastes
}
