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
	User UserContancts
	Db *gocql.Session
}

func StampContactId (contact  Contact) Contact{
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

	return contact
}

func (user * User) DeleteContact(id string) error{
	user.User.UserName = user.Session["user"]
	err := user.Db.Query("delete from user_data where username = ? and contact_id = ?",user.User.UserName , id).Exec()
	fmt.Println(err)
	return err

}
///////////////////////////////////////////////////////////////////////////////////////////////////////
/*func (user * User) GetUserId() (error){
	var id int
	err := user.Db.QueryRow("select id from users where username = ?",user.User.UserName).Scan(&id)
	user.User.Id = strconv.Itoa(id)
	return err

}*/
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user * User) DeleteContactNumber(id string , contactid string) error{
	user.User.UserName = user.Session["user"]
	err := user.Db.Query("delete contact_phonenumbers[?] from user_data where username = ? and contact_id = ?",id ,user.User.UserName , contactid ).Exec()
	return err

}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user * User) GetUserContacts() error{
	var newcontact Contact
	rows := user.Db.Query("select contact_id,contact_email,contact_fname,contact_lname,contact_phonenumbers from user_data where username= ?" , user.User.UserName)
	scanner :=rows.Iter().Scanner()
	for scanner.Next(){
		scanner.Scan(&newcontact.Id , &newcontact.Email, &newcontact.FirstName , &newcontact.LastName , &newcontact.PhoneNumbers)
		newcontact = StampContactId(newcontact)

		user.User.Contacts = append(user.User.Contacts, newcontact)
	}
	err := rows.Iter().Close()
	return err
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user * User) InsertNewContact() (Contact,error){
	//Start Transaction

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

	fmt.Println("before query")
	err := user.Db.Query("insert into user_data (username ,contact_id , contact_email , contact_fname , contact_lname , contact_phonenumbers ) values(? , uuid() , ? , ? , ? , ? ) ", user.User.UserName, user.Params.Get("email") , user.Params.Get("first-name"), user.Params.Get("last-name"), phonenumbers  ).Exec()
	if err !=nil {
		fmt.Println(err)
		return Contact{} , err
	}
	c = StampContactId(c)
	fmt.Println(c.PhoneNumbersStamped)
	user.User.Contacts = append(user.User.Contacts, c)
	return c , nil

}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user *User) Userpage() revel.Result {


	user.User.Contacts = nil
	Username := user.Session["user"]
	if Username == "" {
		return user.Redirect("/")

	}

	user.User.UserName = Username


	err :=user.GetUserContacts()
	if err != nil {
		fmt.Println("DB error")
		user.RenderError(err)

	}

	Myuser :=user.User
	return user.Render(Myuser)
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user *User) AddContact() revel.Result{
	user.User.UserName=user.Session["user"]
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
		fmt.Println("error")
		return user.RenderTemplate("User/userpage.html")
	}else {
		fmt.Println("No error")
		c , err := user.InsertNewContact()
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
		err := user.DeleteContact(user.Params.Get("id"))
		if err != nil {
			fmt.Println("DB error")
			return user.RenderError(err)
		}
		fmt.Println("row Deleted")
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
		err := user.DeleteContactNumber(user.Params.Get("id"),user.Params.Get("ID"))
		if err != nil {
			fmt.Println("DB error")
			return user.RenderError(err)
		}
	}
	fmt.Println("number Deleted")
	return user.Result
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user User) Logout() revel.Result{
	user.User.UserName=""
	user.User.Password=""
	user.User.Contacts=[] Contact{}
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