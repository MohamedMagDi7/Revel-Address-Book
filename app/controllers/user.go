package controllers

import (
	"github.com/revel/revel"
	"fmt"
	"strconv"
	"database/sql"
)

type PhoneNum struct {
	Id int64
	Phonenumber string

}

type Contact struct{
	Id int
	FirstName string
	LastName string
	Email string
	PhoneNumber []PhoneNum

}

type UserContancts struct {
	UserName string
	Id string
	Password string
	Contacts []Contact

}

type User struct {
	*revel.Controller
	User UserContancts
	Db *sql.DB
}

func (user * User) DeleteContact(id string) error{
	_ ,err := user.Db.Exec("delete from contact where contactID = ?",id)
	return err

}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user * User) GetUserId() (error){
	var id int
	err := user.Db.QueryRow("select id from users where username = ?",user.User.UserName).Scan(&id)
	user.User.Id = strconv.Itoa(id)
	return err

}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user * User) DeleteContactNumber(id string) error{
	_ ,err := user.Db.Exec("delete from phonenumbers where id = ?",id)
	return err

}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user * User) GetUserContacts() error{
	rows, err := user.Db.Query("select contactID,fname,lname,email,id,phonenumber from contact join`phonenumbers` on contact.contactID = phonenumbers.contact_id where userID= ?" , user.User.Id)
	var currentcontact Contact
	var newcontact Contact
	var phone PhoneNum

	for rows.Next() {

		rows.Scan(&newcontact.Id, &newcontact.FirstName, &newcontact.LastName , &newcontact.Email , &phone.Id , &phone.Phonenumber )

		if newcontact.Id!=currentcontact.Id && currentcontact.Id != 0{

			user.User.Contacts = append(user.User.Contacts, currentcontact)
			currentcontact = newcontact


		}else if currentcontact.Id == 0{

			currentcontact=newcontact

		}
		currentcontact.PhoneNumber = append(currentcontact.PhoneNumber, phone)
	}
	user.User.Contacts = append(user.User.Contacts, currentcontact)
	return err
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user * User) InsertNewContact() (Contact,error){
	//Start Transaction

	_ , err := user.Db.Exec("START TRANSACTION")
	if err!=nil {
		return Contact{},err
	}

	res, _err := user.Db.Exec("insert into contact values(? ,? ,? ,? ,? ) ", nil, user.Params.Get("first-name"), user.Params.Get("last-name"), user.Params.Get("email"), user.User.Id)
	if _err != nil {
		user.Db.Exec("ROLLBACK")
		return Contact{},err
	}
	id , _ := res.LastInsertId()

	c := Contact{
		FirstName:user.Params.Get("first-name"),
		LastName:user.Params.Get("last-name"),
		Email:user.Params.Get("email"),
		//PhoneNumber:r.FormValue("phone"),
	}
	i := 1
	for user.Params.Get("phone" + strconv.Itoa(i)) != "" {
		str := user.Params.Get("phone" + strconv.Itoa(i))
		res , err := user.Db.Exec("insert into phonenumbers values(?,?,?)", nil, str , id)
		if err != nil {
			user.Db.Exec("ROLLBACK")
			return Contact{},err
		}
		id , _ := res.LastInsertId()
		Phone := PhoneNum{Phonenumber:str , Id:id}
		c.PhoneNumber = append(c.PhoneNumber, Phone)
		i++
	}
	_ , err =user.Db.Exec("COMMIT")
	if err != nil {
		return Contact{}, err
	}
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
	err := user.GetUserId()
	if err != nil {
		fmt.Println("DB error")
		user.RenderError(err)

	}
	user.Session["userid"]=user.User.Id
	err =user.GetUserContacts()
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
	user.User.Id=user.Session["userid"]
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
		err := user.DeleteContactNumber(user.Params.Get("id"))
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
	user.User.Id=""
	user.User.Password=""
	user.User.Contacts=[] Contact{}
	user.Session["user"] = ""


	return user.Redirect("/")
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func startDatabase(user *User) revel.Result{
	var err error
	user.Db, err= sql.Open("mysql", "root:1819@tcp(127.0.0.1:3306)/my_add_bookDB")
	if err != nil {
		panic(err)
	}
	return nil
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func init(){
	revel.InterceptMethod(startDatabase , revel.BEFORE)
}