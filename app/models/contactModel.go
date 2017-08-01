	package models

	import (
		"github.com/gocql/gocql"
		"fmt"

	)

	type Contact struct{
		Id gocql.UUID
		FirstName string
		LastName string
		Email string
		PhoneNumbers []PhoneNum


	}



	func (contact  * Contact) DeleteContact(id string ,username string , db *gocql.Session) error{
		batch := gocql.NewBatch(gocql.LoggedBatch)
		stmt1 := "delete from user_data where username = ? and contact_id = ?"
		stmt2 := "delete from phone_number where contact_id = ? "
		batch.Query(stmt1 , username , id)
		batch.Query(stmt2 , id)
		err  := db.ExecuteBatch(batch)
		fmt.Println(err)
		return err

	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////
	func (c  * Contact) InsertNewContact( username string , db * gocql.Session) (error){
		contactId ,err := gocql.RandomUUID()
		var numberId gocql.UUID
		if err!= nil { return err }
		batch := gocql.NewBatch(gocql.LoggedBatch)
		stmt1 := "insert into user_data (username ,contact_id , contact_email , contact_fname , contact_lname ) values(? , ? , ? , ? , ? ) "
		stmt2 := "insert into phone_number (contact_id , number_id , number) values(? , ? , ?)"
		batch.Query(stmt1 , username, contactId, c.Email , c.FirstName, c.LastName)
		for index , value := range c.PhoneNumbers{
			numberId , err = gocql.RandomUUID()
			c.PhoneNumbers[index].ContactId= contactId
			c.PhoneNumbers[index].NumberId = numberId
			batch.Query(stmt2 , contactId , numberId , value.Phonenumber)
		}
		err = db.ExecuteBatch(batch)
		if err !=nil {
			fmt.Println(err)
			return  err
		}
		//user.Contacts = append(user.Contacts, c)
		return  nil

	}