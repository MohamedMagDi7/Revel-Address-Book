package tests


import (
	"github.com/revel/revel/testing"

)

type UserTest struct {
	testing.TestSuite
}



func (t *UserTest) TestUserPage() {
	t.Get("/userpage")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")

}

func (t *UserTest) TestLogout() {
	t.Get("/logout")
	t.AssertOk()
	t.AssertContains("User Name:")

}




