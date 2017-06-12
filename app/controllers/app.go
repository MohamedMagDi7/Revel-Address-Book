package controllers

import (
	"github.com/revel/revel"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
)
var DB *sql.DB
var err error
type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	err := c.Flash.Data["Log"]
	fmt.Println(err)
	return c.Render(err)
}

func (c App) Login() revel.Result {

	if err = DB.Ping(); err !=nil {
		fmt.Println("Database is closed!")
		return c.RenderTemplate("errors/500.html")
	}
	username := c.Params.Get("username")
	password := c.Params.Get("password")


	if c.Params.Get("register")!="" {
		fmt.Println("registered")
		var user string

		err = DB.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

		switch {
		// Username is available
		case err == nil:
			//fmt.Println("Username is not available")
			c.Flash.Data["Log"]="UserName is Already in use!"
			c.FlashParams()
			return c.Redirect("/")

		case err == sql.ErrNoRows:
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				//fmt.Println("Couldn't Incrypt")
				c.Flash.Data["Log"]="Please choose another password!"
				return c.RenderError(err)
			}

			_, err = DB.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
			if err != nil {
				return c.RenderError(err)
			}
			c.Session["user"]=username

			return c.Redirect("/userpage")

		case err != nil:
			c.RenderError(err)
			c.Flash.Data["Log"]="DataBase error!"
			return c.Redirect("/")

		default:
			c.Flash.Data["Log"]=""
			c.Redirect("/")
		}

	}



	if c.Params.Get("login")!="" {

		fmt.Println("logged in")
		var databaseUsername string
		var databasePassword string

		err = DB.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)
		if err == sql.ErrNoRows {
			c.Flash.Data["Log"]="UserName Doesn't exist!"
			fmt.Println("no such user")
			return c.Redirect("/")


		} else if err != nil {
			c.Flash.Out["Log"]="DB error!"
			return c.RenderError(err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
		// If wrong password redirect to the login
		if err != nil {
			c.Flash.Data["Log"]="Wrong Password!"
			fmt.Println("wrong password")
			return c.Redirect("/")

		} else {
			fmt.Println("password match")
			// If the login succeeded
			c.Flash.Data["Log"]=""
			c.Session["user"]=username

			return c.Redirect("/userpage")
		}
	}


		// Grab from the database

	return c.Render()
}


func startDB(c *revel.Controller) revel.Result{
	DB, err = sql.Open("mysql", "root:1819@tcp(127.0.0.1:3306)/my_add_bookDB")
	if err != nil {
		panic(err)
	}
	return nil
}

func init(){
	revel.InterceptFunc(startDB , revel.BEFORE, &App{})
}