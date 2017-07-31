package controllers

import (
	"github.com/revel/revel"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"github.com/gocql/gocql"
	."MyRevelApp/app/models"
	"MyRevelApp/app"
)



type App struct {
	*revel.Controller
	logins * Logins
}

func (c App) Index() revel.Result {


	return c.Render()
}




func (c App) SignIn() revel.Result{
	var err error
	var databasePassword string

	databasePassword, err = c.logins.QueryUser(app.DB)
	if err == gocql.ErrNotFound {
		//no such user
		c.Flash.Error("Username doesn't exist")
		return c.Redirect( "/")

	} else if  err != nil {
		c.Flash.Error("Internal Server Error please try again")
		return c.Redirect( "/" )
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(c.logins.Password))
	// If wrong password redirect to the login
	if  err != nil {
		//Wrong Password
		c.Flash.Error("wrong password")
		return c.Redirect( "/" )
	} else {
		// If the login succeeded
		c.Session["user"]= c.logins.Username
		return c.Redirect( "/userpage" )
	}
}

func (c App) Register() revel.Result{
	err :=c.logins.CheckUsernameExists(app.DB)
	switch {
	case err == nil:
		c.Flash.Error( "Please choose a different username")
		return c.Redirect( "/" )

	case err == gocql.ErrNotFound :
		// Username is available
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.logins.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Flash.Error("This Password is Not premitted")
			return c.Redirect( "/" )
		}

		err = c.logins.InsertUser(hashedPassword , app.DB)
		if err != nil {
			return c.RenderError(err)
		}
		c.Session["user"]=c.logins.Username
		return c.Redirect("/userpage")

	case err != nil:
		//Database Error
		c.Flash.Error( "Internal Server Error please try again")
		return c.Redirect( "/" )

	default:
		return c.Redirect( "/" )
	}
}

func (c App) Login() revel.Result {
	c.Params.Bind(&c.logins , "logins")
	//Input Validation
	c.Validation.Required(c.logins.Username)
	c.Validation.Required(c.logins.Password)
	c.Validation.Length(c.logins.Username,50)
	c.Validation.Length(c.logins.Password,120)
	if !c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return  c.RenderTemplate("App/index.html")
	}
	if c.Params.Get("register")!="" {
		return c.Register()
	}

	if c.Params.Get("login")!="" {
		return c.SignIn()
	}


	return c.Render()
}
