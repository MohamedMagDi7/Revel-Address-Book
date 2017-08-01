package models

import (
	"github.com/gocql/gocql"
	"fmt"
)

type PhoneNum struct {
	NumberId gocql.UUID
	ContactId gocql.UUID
	Phonenumber string

}


func (number * PhoneNum) DeletePhoneNumber( numberId string ,contactId string , db *gocql.Session) error{
	err := db.Query("delete from phone_number where contact_id = ? and number_id = ?",contactId, numberId).Exec()
	return err

}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (number * PhoneNum) AddPhoneNumber(contactId string , db *gocql.Session) error{
	var err error
	number.NumberId , err = gocql.RandomUUID()
	if err != nil { return err}
	err = db.Query("insert into phone_number(contact_id , number_id , number) values(? , ? , ?) ", contactId ,number.NumberId , number.Phonenumber  ).Exec()
	fmt.Println(err)
	return err

}