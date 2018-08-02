package main

import (
	"net/http"
	"html/template"
	"log"
)

type bianmaController struct {
}

func (this *bianmaController)IndexAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/html/tables/bianma.html")
	if (err != nil) {
		log.Println(err)
	}
	t.Execute(w, nil)
}
