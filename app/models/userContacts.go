package models

import (
	"github.com/gocql/gocql"
	"fmt"
)

type UserContancts struct {
	UserName string
	Password string
	Contacts []Contact

}

func (user * UserContancts) DeleteContact(id string , db *gocql.Session) error{
	err := db.Query("delete from user_data where username = ? and contact_id = ?",user.UserName , id).Exec()
	fmt.Println(err)
	return err

}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user * UserContancts) DeleteContactNumber(id string , contactid string , db *gocql.Session) error{
	err := db.Query("delete contact_phonenumbers[?] from user_data where username = ? and contact_id = ?",id ,user.UserName , contactid ).Exec()
	return err

}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user * UserContancts) GetUserContacts( db *gocql.Session) error{
	var newcontact Contact
	rows := db.Query("select contact_id,contact_email,contact_fname,contact_lname,contact_phonenumbers from user_data where username= ?" , user.UserName)
	scanner :=rows.Iter().Scanner()
	for scanner.Next(){
		scanner.Scan(&newcontact.Id , &newcontact.Email, &newcontact.FirstName , &newcontact.LastName , &newcontact.PhoneNumbers)
		newcontact.StampContactId()

		user.Contacts = append(user.Contacts, newcontact)
	}
	err := rows.Iter().Close()
	return err
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user * UserContancts) InsertNewContact(c Contact , db * gocql.Session) (Contact, error){

	err := db.Query("insert into user_data (username ,contact_id , contact_email , contact_fname , contact_lname , contact_phonenumbers ) values(? , uuid() , ? , ? , ? , ? ) ", user.UserName, c.Email , c.FirstName, c.LastName, c.PhoneNumbers).Exec()
	if err !=nil {
		fmt.Println(err)
		return Contact{} , err
	}
	c.StampContactId()
	user.Contacts = append(user.Contacts, c)
	return c , nil

}
