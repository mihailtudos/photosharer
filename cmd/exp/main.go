package main

import (
	"os"
	"text/template"
)

type User struct {
	Name string
}

func main() {

	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, User{Name: "Mihail"})
	if err != nil {
		panic(err)
	}
}
