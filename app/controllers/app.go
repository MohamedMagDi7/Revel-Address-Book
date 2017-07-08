package controllers

import (
	"github.com/revel/revel"
	"fmt"
	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
	"github.com/gocql/gocql"
)

type App struct {
	*revel.Controller
	db *gocql.Session

}

func (c App) Index() revel.Result {

	return c.Render()
}

func (c App) CheckUsernameExists( username string) error{
	var databasePassword string

	err := c.db.Query("SELECT password FROM user_logins WHERE username=?", username).Scan( &databasePassword)
	fmt.Println(err)
	return err
}

func (c App) QueryUser(username string) (string,error){
	var databasePassword string

	err := c.db.Query("SELECT password FROM user_logins WHERE username=?", username).Scan( &databasePassword)
	return databasePassword,err
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (c App) InsertUser( username string, hashedPassword []byte) error{
	err :=c.db.Query("INSERT INTO user_logins(username, password) VALUES(?, ?)",username, hashedPassword).Exec()
	return err

}

func (c App) SignIn(username string , password string) revel.Result{
	var err error
	var databasePassword string
	databasePassword, err = c.QueryUser( username)
	if err == gocql.ErrNotFound {
		//no such user
		c.Flash.Error("Username doesn't exist")
		return c.Redirect( "/")

	} else if  err != nil {
		c.Flash.Error("Internal Server Error please try again")
		return c.Redirect( "/" )
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	// If wrong password redirect to the login
	if  err != nil {
		//Wrong Password
		c.Flash.Error("wrong password")
		return c.Redirect( "/" )
	} else {
		// If the login succeeded
		c.Session["user"]= username
		return c.Redirect( "/userpage" )
	}
}

func (c App) Register(username string , password string) revel.Result{
	err :=c.CheckUsernameExists(username)
	switch {
	case err == nil:
		c.Flash.Error( "Please choose a different username")
		return c.Redirect( "/" )

	case err == gocql.ErrNotFound :
		// Username is available
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.Flash.Error("This Password is Not premitted")
			return c.Redirect( "/" )
		}

		err = c.InsertUser(username,hashedPassword)
		if err != nil {
			return c.RenderError(err)
		}
		c.Session["user"]= username
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
	username := c.Params.Get("username")
	password := c.Params.Get("password")
	//Input Validation
	c.Validation.Required(username)
	c.Validation.Required(password)
	c.Validation.Length(username,50)
	c.Validation.Length(password,120)
	if !c.Validation.HasErrors() {
		fmt.Println("here i am")
		c.Validation.Keep()
		c.FlashParams()
		return  c.RenderTemplate("App/index.html")
	}

	if c.Params.Get("register")!="" {
		return c.Register(username,password)
	}



	if c.Params.Get("login")!="" {
		return c.SignIn(username , password)
	}


	return c.Render()
}

func startDB(c *App) revel.Result{
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "address_book"
	c.db, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	return nil
}

func init(){
	revel.InterceptMethod(startDB , revel.BEFORE)
}