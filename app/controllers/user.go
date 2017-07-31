package controllers

import (
	"github.com/revel/revel"
	"fmt"
	"strconv"
	. "MyRevelApp/app/models"
	"MyRevelApp/app"
)


type User struct {
	*revel.Controller
}

///////////////////////////////////////////////////////////////////////////////////////////////////////
func (user *User) Userpage() revel.Result {

	myuser := UserContancts{Contacts:nil}
	Username := user.Session["user"]
	if Username == "" {
		return user.Redirect("/")

	}

	myuser.UserName = Username


	err :=myuser.GetUserContacts(app.DB)
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
		c := Contact{
			FirstName:user.Params.Get("first-name"),
			LastName:user.Params.Get("last-name"),
			Email:user.Params.Get("email"),
			PhoneNumbers:phonenumbers,
		}
		c , err := myuser.InsertNewContact(c , app.DB)
		if err !=nil {
			return user.RenderError(err)
		}
		fmt.Println(c)
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
		err := myuser.DeleteContact(user.Params.Get("id") , app.DB)
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
		err := myuser.DeleteContactNumber(user.Params.Get("id") , user.Params.Get("ID") , app.DB)
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