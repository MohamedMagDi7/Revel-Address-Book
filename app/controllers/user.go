package controllers

import (
	"github.com/revel/revel"
	"strconv"
	"MyRevelApp/app/models"
	"MyRevelApp/app"
	"fmt"
)


type User struct {
	*revel.Controller
}

///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user *User) Userpage() revel.Result {

	myUser := &models.User{Contacts:nil}
	myUser.Logins.Username = user.Session["user"]
	if myUser.Logins.Username == "" {
		return user.Redirect("/")

	}

	err :=myUser.GetUserContacts(app.DB)
	if err != nil {
		user.RenderError(err)

	}
	fmt.Println(myUser)
	return user.Render(myUser)
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user *User) AddContact() revel.Result{
	myUser := models.User{}
	myUser.Logins.Username=user.Session["user"]

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
		contact := models.Contact{
			FirstName:user.Params.Get("first-name"),
			LastName:user.Params.Get("last-name"),
			Email:user.Params.Get("email"),
			PhoneNumbers:phonenumbers,
		}
		 err := contact.InsertNewContact(myUser.Logins.Username, app.DB)
		if err !=nil {
			return user.RenderError(err)
		}else {
			myUser.AddtoContacts(contact)
			return user.RenderJSON(contact)
		}
	}
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user User) Delete() revel.Result {
	user.Validation.Clear()
	user.Validation.Required(user.Params.Get("id"))
	if user.Validation.HasErrors() {
		return user.RenderTemplate("user/userpage.html")
	} else {
		username :=user.Session["user"]
		contact := models.Contact{}
		err := contact.DeleteContact(user.Params.Get("id") ,username, app.DB)
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
		username :=user.Session["user"]
		contact := models.Contact{}
		err := contact.DeleteContactNumber(user.Params.Get("id") , user.Params.Get("ID") ,username, app.DB)
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