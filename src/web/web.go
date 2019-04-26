package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type MyMux struct {
}

func (p *MyMux)ServerHTTP(w http.ResponseWriter, r *http.Request){
	if r.URL.Path == "/"{
		Index(w, r)
		return
	}
	if r.URL.Path == "/login"{
		login(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func Index(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path: ", r.URL.Path)
	fmt.Println("scheme: ", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form{
		fmt.Println("key: ", k)
		fmt.Println("val: ", strings.Join(v, " "))
	}
	fmt.Fprintf(w, "hello chain!")
}

func login(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fmt.Println("method: ",r.Method)
	if r.Method == "GET"{
		t, _ := template.ParseFiles("./view/login.xml")
		t.Execute(w, nil)
	}else if r.Method == "POST"{
		fmt.Println("username: ", r.Form["username"])
		fmt.Println("password: ", r.Form["password"])
	}
}

func Start(){
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}


