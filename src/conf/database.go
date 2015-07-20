// database
package conf

import (
	//	"fmt"

	"gopkg.in/mgo.v2"
)
/*
func main() {
	fmt.Println("Hello World!")
}
*/
const PASTES = "pastes"

func GetCollection(collectionName string) (*mgo.Session, *mgo.Collection, error) {
	session, err := mgo.Dial("127.0.0.1")
	if err != error(nil) {
		return session, nil, err
	}
	col := session.DB("viewfunc1").C(collectionName)
	return session, col, nil
}
