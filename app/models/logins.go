package models

import (
	"github.com/gocql/gocql"
	"fmt"
)

type Logins struct {
	Username string
	Password string
}


func (logins * Logins) CheckUsernameExists(db* gocql.Session) error{
	var databasePassword string

	err := db.Query("SELECT password FROM user_logins WHERE username=?", logins.Username).Scan( &databasePassword)
	fmt.Println(err)
	return err
}

func (logins * Logins) QueryUser(db * gocql.Session) (string,error){
	var databasePassword string

	err := db.Query("SELECT password FROM user_logins WHERE username=?", logins.Username).Scan( &databasePassword)
	return databasePassword,err
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (logins * Logins) InsertUser(hashedPassword []byte , db * gocql.Session) error{
	err :=db.Query("INSERT INTO user_logins(username, password) VALUES(?, ?)", logins.Username, hashedPassword).Exec()
	return err

}