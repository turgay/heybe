package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseGlob("./tmpl/*.html"))

type ReadItem struct {
	Name  string
	Link  string
	Descr string
	Tag   string
}

type User struct {
	UserName string
	Password string
	Email    string
}

func (user *User) Match(userName string, passwd string) bool {
	return user.UserName == userName && user.Password == passwd
}

var repository = Repository{}

func renderTemplate(w http.ResponseWriter, tmpl string, params interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, params)
	//err := templates.ExecuteTemplate(w, tmpl+".html", params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func listHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	fmt.Println("Logged user = ", userName)
	if userName != "" {
		items, err := repository.LoadItems()
		fmt.Println(items, err)
		if err != nil {
			http.Redirect(response, request, "/list", http.StatusFound)
			return
		}

		params := map[string]interface{}{"authuser": userName, "items": items}
		renderTemplate(response, "list", params)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func aboutHandler(response http.ResponseWriter, request *http.Request) {
	renderTemplate(response, "about", nil)

}

type Errors map[string]string

func validateUser(userName string, email string) Errors {
	errs := make(map[string]string)

	re := regexp.MustCompile(".+@.+\\..+")
	matched := re.Match([]byte(email))
	if matched == false {
		errs["email"] = "Please enter a valid email address."
	}

	if strings.TrimSpace(userName) == "" {
		errs["userName"] = "User Name cannot be empty."
	}

	return errs
}

func registerHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		renderTemplate(response, "register", nil)
	} else {
		request.ParseForm()
		form := request.Form

		userName := form.Get("userName")
		email := form.Get("email")
		password := form.Get("passwd")

		errs := validateUser(userName, email)

		if len(errs) > 0 {

			form.Set("passwd", "")
			form.Set("conpasswd", "")
			params := map[string]interface{}{"userName": userName, "email": email, "Errors": errs}
			renderTemplate(response, "register", params)
		} else {

			//fmt.Printf("%+v", form)

			//TODO decode password
			newUser := User{UserName: userName, Email: email, Password: password}
			repository.AddUser(newUser)
			setSession(userName, response)

			http.Redirect(response, request, "/list", http.StatusFound)
		}
	}
}

func newItemHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		renderTemplate(response, "add", nil)
	} else {
		name := request.FormValue("name")
		link := request.FormValue("url")
		descr := request.FormValue("descr")
		tag := request.FormValue("tag")
		ri := ReadItem{Name: name, Link: link, Descr: descr, Tag: tag}

		repository.AddItem(ri)

		http.Redirect(response, request, "/list", http.StatusFound)
	}
}

//--------------

// login handler

func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("userName")
	pass := request.FormValue("password")
	redirectTarget := "/"
	fmt.Println("N=", name, " P=", pass)
	fmt.Printf("%+v\n", request)
	if name != "" && pass != "" {
		err := authUser(name, pass)
		if err == nil {
			setSession(name, response)
			redirectTarget = "list"
		} else {
			//TODO show error message!!!!
			if err != nil {
				http.Error(response, err.Error(), http.StatusInternalServerError)
				return

			}
		}
	}
	http.Redirect(response, request, redirectTarget, 302)
}

// logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	renderTemplate(response, "login", nil)
}

// server main method

var router = mux.NewRouter()

func main() {
	repository.Init()

	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/list", listHandler)
	http.HandleFunc("/add", newItemHandler)

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")
	router.HandleFunc("/about", aboutHandler)
	router.HandleFunc("/register", registerHandler)

	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.Handle("/static/", fileServer)

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)

}
