	package models

	import (
		"github.com/gocql/gocql"
		"fmt"
	)

	type LoginData struct{
		Username string
		Password string
	}
	type User struct {
		Logins LoginData
		Contacts []Contact

	}

	func (user * User) CheckUsernameExists(db* gocql.Session) error{
		var databasePassword string

		err := db.Query("SELECT password FROM user_logins WHERE username=?", user.Logins.Username).Scan( &databasePassword)
		fmt.Println(err)
		return err
	}

	func (user * User) QueryUser(db * gocql.Session) (string,error){
		var databasePassword string

		err := db.Query("SELECT password FROM user_logins WHERE username=?", user.Logins.Username).Scan( &databasePassword)
		return databasePassword,err
	}
	///////////////////////////////////////////////////////////////////////////////////////////////////////
	func (user * User) InsertUser(hashedPassword []byte , db * gocql.Session) error{
		err :=db.Query("INSERT INTO user_logins(username, password) VALUES(?, ?)", user.Logins.Username, hashedPassword).Exec()
		return err

	}

	func (user * User) GetUserContacts( db *gocql.Session) error{
		var newcontact Contact
		var newnumber PhoneNum
		rows := db.Query("select contact_id,contact_email,contact_fname,contact_lname from user_data where username= ?" , user.Logins.Username)
		scanner :=rows.Iter().Scanner()
		for scanner.Next(){
			scanner.Scan(&newcontact.Id , &newcontact.Email, &newcontact.FirstName , &newcontact.LastName )
			res :=  db.Query("select number_id , number from phone_number where contact_id= ?" , newcontact.Id)
			numberScanner :=res.Iter().Scanner()
			for numberScanner.Next() {
				numberScanner.Scan(&newnumber.NumberId , &newnumber.Phonenumber)
				newnumber.ContactId= newcontact.Id
				newcontact.PhoneNumbers = append(newcontact.PhoneNumbers , newnumber)
			}


			user.Contacts = append(user.Contacts, newcontact)
		}
		err := rows.Iter().Close()
		return err
	}

	func (user * User) AddtoContacts(contact  Contact) {
		user.Contacts = append(user.Contacts, contact)
		return
	}