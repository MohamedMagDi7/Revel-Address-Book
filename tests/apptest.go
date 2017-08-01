package tests

import (
	"github.com/revel/revel/testing"

	"net/url"
)

type AppTest struct {
	testing.TestSuite
}

func (t *AppTest) Before() {
	println("Getting things ready")
}

func (t *AppTest) TestIndexPage() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *AppTest) TestLoginFuncValidUser() {

	t.PostForm("/login",url.Values{"username": {"mohamedmega21@gmail.com" } , "password":{"1819"} , "login":{"Log In"}})
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
	t.AssertContains("You are logged in as:")
}

func (t *AppTest) TestLoginFuncNonValidUser() {

	t.PostForm("/login",url.Values{"username": {"mohamedaa@gmail.com" } , "password":{"1819"} , "login":{"Log In"}})
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
	t.AssertContains("User Name:")
}

func (t *AppTest) TestLoginFuncWrongPassword() {

	t.PostForm("/login",url.Values{"username": {"mohamedmega21@gmail.com" } , "password":{"18195sa"} , "login":{"Log In"}})
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
	t.AssertContains("User Name:")
}

/*func (t *AppTest) TestRegisterFuncValidUserName() {

	t.PostForm("/login",url.Values{"username": {"mohamed@gmail.com" } , "password":{"1819"} , "register":{"Register"}, "login":{""}})
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
	t.AssertContains("You are logged in as:")
}*/

func (t *AppTest) TestRegisterFuncInvalidUserName() {

	t.PostForm("/login",url.Values{"username": {"mohamed@gmail.com" } , "password":{"1819"} , "register":{"Register"} , "login":{""}})
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
	t.AssertContains("User Name:")
}

func (t *AppTest) After() {
	println("Finished")
}
