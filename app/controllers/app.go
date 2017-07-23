	package controllers

import (
	"github.com/revel/revel"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"github.com/gocql/gocql"
)

type Logins struct {
	username string
	password string
}

	type App struct {
	*revel.Controller
	db *gocql.Session

}

func (c App) Index() revel.Result {


	return c.Render()
}


func (logins * Logins) CheckUsernameExists(db* gocql.Session) error{
	var databasePassword string

	err := db.Query("SELECT password FROM user_logins WHERE username=?", logins.username).Scan( &databasePassword)
	fmt.Println(err)
	return err
}

func (logins * Logins) QueryUser(db * gocql.Session) (string,error){
	var databasePassword string

	err := db.Query("SELECT password FROM user_logins WHERE username=?", logins.username).Scan( &databasePassword)
	return databasePassword,err
}
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (logins * Logins) InsertUser(hashedPassword []byte , db * gocql.Session) error{
	err :=db.Query("INSERT INTO user_logins(username, password) VALUES(?, ?)", logins.username, hashedPassword).Exec()
	return err

}

func (c App) SignIn(logins Logins) revel.Result{
	var err error
	var databasePassword string

	databasePassword, err = logins.QueryUser(c.db)
	if err == gocql.ErrNotFound {
		//no such user
		c.Flash.Error("Username doesn't exist")
		return c.Redirect( "/")

	} else if  err != nil {
		c.Flash.Error("Internal Server Error please try again")
		return c.Redirect( "/" )
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(logins.password))
	// If wrong password redirect to the login
	if  err != nil {
		//Wrong Password
		c.Flash.Error("wrong password")
		return c.Redirect( "/" )
	} else {
		// If the login succeeded
		c.Session["user"]= logins.username
		return c.Redirect( "/userpage" )
	}
}

func (c App) Register(logins Logins) revel.Result{
	err :=logins.CheckUsernameExists(c.db)
	switch {
	case err == nil:
		c.Flash.Error( "Please choose a different username")
		return c.Redirect( "/" )

	case err == gocql.ErrNotFound :
		// Username is available
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(logins.password), bcrypt.DefaultCost)
		if err != nil {
			c.Flash.Error("This Password is Not premitted")
			return c.Redirect( "/" )
		}

		err = logins.InsertUser(hashedPassword , c.db)
		if err != nil {
			return c.RenderError(err)
		}
		c.Session["user"]= logins.username
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
	logins := Logins{username:username , password:password}
	if c.Params.Get("register")!="" {
		return c.Register(logins)
	}



	if c.Params.Get("login")!="" {
		return c.SignIn(logins)
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