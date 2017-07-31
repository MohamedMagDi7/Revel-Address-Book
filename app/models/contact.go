package models

import (
	"github.com/gocql/gocql"
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