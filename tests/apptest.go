package tests

import (
	"github.com/revel/revel/testing"
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

func (t *AppTest) TestUserPage() {
	t.Get("/userpage")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}


func (t *AppTest) After() {
	println("Finished")
}
