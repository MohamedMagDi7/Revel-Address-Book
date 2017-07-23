package controllers

import (
	"github.com/revel/revel"
	"fmt"
	"strconv"
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

type UserContancts struct {
	UserName string
	Password string
	Contacts []Contact

}

type User struct {
	*revel.Controller
	Db *gocql.Session
}

func (contact  *Contact)StampContactId () {
	i :=0
	contact.PhoneNumbersStamped = []PhoneNum{}
	fmt.Println(len(contact.PhoneNumbers))
	contactid := contact.Id
	for i<len(contact.PhoneNumbers){
		numberid := i
		phonenumber := contact.PhoneNumbers[i]
		contact.PhoneNumbersStamped = append(contact.PhoneNumbersStamped , PhoneNum{ContactId:contactid , Id:numberid , Phonenumber:phonenumber})
		i++
	}
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
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user *User) Userpage() revel.Result {

	myuser := UserContancts{Contacts:nil}
	Username := user.Session["user"]
	if Username == "" {
		return user.Redirect("/")

	}

	myuser.UserName = Username


	err :=myuser.GetUserContacts(user.Db)
	if err != nil {
		user.RenderError(err)

	}
	return user.Render(myuser)
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user *User) AddContact() revel.Result{
	myuser := UserContancts{}
	myuser.UserName=user.Session["user"]
	user.Validation.Required(user.Params.Get("first-name"))
	user.Validation.Required(user.Params.Get("last-name"))
	user.Validation.Required(user.Params.Get("email"))
	user.Validation.MaxSize(user.Params.Get("first-name") ,50)
	user.Validation.MaxSize(user.Params.Get("last-name") ,50)
	user.Validation.MaxSize(user.Params.Get("email") ,50)
	user.Validation.MinSize(user.Params.Get("email") , 7)


	if user.Validation.HasErrors() {
		user.Validation.Keep()
		user.FlashParams()
		return user.RenderTemplate("User/userpage.html")
	}else {
		var phonenumbers [] string

		i := 1
		for user.Params.Get("phone" + strconv.Itoa(i)) != "" {
			str := user.Params.Get("phone" + strconv.Itoa(i))
			phonenumbers = append(phonenumbers,str)
			i++
		}
		fmt.Println(phonenumbers)
		c := Contact{
			FirstName:user.Params.Get("first-name"),
			LastName:user.Params.Get("last-name"),
			Email:user.Params.Get("email"),
			PhoneNumbers:phonenumbers,
		}
		c , err := myuser.InsertNewContact(c , user.Db)
		if err !=nil {
			return user.RenderError(err)
		}
		return user.RenderJSON(c)
	}
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user User) Delete() revel.Result {
	user.Validation.Clear()
	user.Validation.Required(user.Params.Get("id"))
	if user.Validation.HasErrors() {
		return user.RenderTemplate("user/userpage.html")
	} else {
		myuser :=UserContancts{UserName:user.Session["user"]}
		err := myuser.DeleteContact(user.Params.Get("id") , user.Db)
		if err != nil {
			return user.RenderError(err)
		}
		return user.Result
	}
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user User) DeleteNum() revel.Result{
	user.Validation.Clear()
	user.Validation.Required(user.Params.Get("id"))
	if user.Validation.HasErrors() {
		return user.RenderTemplate("user/userpage.html")
	} else {
		myuser :=UserContancts{UserName:user.Session["user"]}
		err := myuser.DeleteContactNumber(user.Params.Get("id") , user.Params.Get("ID") , user.Db)
		if err != nil {
			return user.RenderError(err)
		}
	}
	return user.Result
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user User) Logout() revel.Result{
	user.Session["user"] = ""


	return user.Redirect("/")
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func startDatabase(user *User) revel.Result{
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "address_book"
	user.Db, err= cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	return nil
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func init(){
	revel.InterceptMethod(startDatabase , revel.BEFORE)
}