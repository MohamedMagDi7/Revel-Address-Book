package models

import (
	"github.com/gocql/gocql"
	"fmt"
)

type PhoneNum struct {
	Id int
	ContactId gocql.UUID
	Phonenumber string

}

type Contact struct{
	Id gocql.UUID
	FirstName string
	LastName string
	Email string
	PhoneNumbersStamped []PhoneNum
	PhoneNumbers []string

}


func (contact  * Contact)StampContactId () {
	i :=0
	contact.PhoneNumbersStamped = []PhoneNum{}
	contactid := contact.Id
	for i<len(contact.PhoneNumbers){
		numberid := i
		phonenumber := contact.PhoneNumbers[i]
		contact.PhoneNumbersStamped = append(contact.PhoneNumbersStamped , PhoneNum{ContactId:contactid , Id:numberid , Phonenumber:phonenumber})
		i++
	}
}

func (contact  * Contact) DeleteContact(id string ,username string , db *gocql.Session) error{
	err := db.Query("delete from user_data where username = ? and contact_id = ?", username, id).Exec()
	fmt.Println(err)
	return err

}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (contact  * Contact) DeleteContactNumber(id string , contactid string ,username string, db *gocql.Session) error{
	err := db.Query("delete contact_phonenumbers[?] from user_data where username = ? and contact_id = ?",id ,username, contactid ).Exec()
	return err

}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (c  * Contact) InsertNewContact( username string , db * gocql.Session) (error){

	err := db.Query("insert into user_data (username ,contact_id , contact_email , contact_fname , contact_lname , contact_phonenumbers ) values(? , uuid() , ? , ? , ? , ? ) ", username, c.Email , c.FirstName, c.LastName, c.PhoneNumbers).Exec()
	if err !=nil {
		fmt.Println(err)
		return  err
	}
	c.StampContactId()
	//user.Contacts = append(user.Contacts, c)
	return  nil

}