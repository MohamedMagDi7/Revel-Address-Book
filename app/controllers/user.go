package controllers

import (
	"github.com/revel/revel"
	"fmt"

	"strconv"

)

type Contact struct{
	Id int
	FirstName string
	LastName string
	Email string
	PhoneNumber []string

}
type User_Contancts struct {
	UserName string
	Id string
	Contacts []Contact

}

var MyUser =User_Contancts{}
type User struct {
	*revel.Controller
}

func (u User) Userpage() revel.Result {


	MyUser.Contacts = []Contact{}
	Username := u.Session["user"]
	if Username == "" {
		return u.Redirect("/")

	}

	MyUser.UserName = Username
	row := DB.QueryRow("select id from users where username= ?", Username)

	row.Scan(&MyUser.Id)
	rows, err := DB.Query("select contactID,fname,lname,email from contact where userID= ?", MyUser.Id)
	if err != nil {
		fmt.Println("DB error")
		u.RenderError(err)

	}

	for rows.Next() {
		var c Contact
		rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Email)

		res, err := DB.Query("select phonenumber from phonenumbers where contact_id= ?", c.Id)
		if err != nil {
			fmt.Println("DB error")
			u.RenderError(err)

		}

		for res.Next() {
			var N string
			res.Scan(&N)
			c.PhoneNumber = append(c.PhoneNumber, N)

		}
		MyUser.Contacts = append(MyUser.Contacts, c)

	}



	return u.Render(MyUser)
}

func (u User) AddContact() revel.Result{
	_, err := DB.Exec("insert into contact values(? ,? ,? ,? ,? ) ", nil, u.Params.Get("first-name"),
		u.Params.Get("last-name"), u.Params.Get("email"), MyUser.Id)
	if err != nil {
		return u.RenderError(err)
	}
	row := DB.QueryRow("select MAX(contactID) from contact")
	var id int
	row.Scan(&id)

	c := Contact{
		FirstName:u.Params.Get("first-name"),
		LastName:u.Params.Get("last-name"),
		Email:u.Params.Get("email"),
		//PhoneNumber:r.FormValue("phone"),
	}
	i := 1
	for u.Params.Get("phone" + strconv.Itoa(i)) != "" {
		str :=u.Params.Get("phone" + strconv.Itoa(i))
		c.PhoneNumber = append(c.PhoneNumber, str)
		_, err := DB.Exec("insert into phonenumbers values(?,?,?)", nil, str , id)
		if err != nil {
			return u.RenderError(err)
		}
		i++
	}

	MyUser.Contacts = append(MyUser.Contacts, c)
	return u.RenderJSON(c)
}

func (u User) Delete() revel.Result {
	fmt.Println(u.Params.Get("id"))
	u.Validation.Required(u.Params.Get("id"))
	_ ,err := DB.Exec("delete from contact where contactID = ?",u.Params.Get("id"))

	if err !=nil{
		fmt.Println("DB error")
		return u.RenderError(err)
	}
	fmt.Println("row Deleted")
	return u.Result
}

func (u User) Logout() revel.Result{

	u.Session["user"] = ""

	return u.Redirect("/")
}

func init(){
	revel.InterceptFunc(startDB , revel.BEFORE, &User{})
}