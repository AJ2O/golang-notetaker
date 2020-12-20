package main

import (
	"encoding/gob"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/AJ2O/golang-notetaker/notes"
	"github.com/AJ2O/golang-notetaker/registration"
)

type LoggedInUser struct {
	Username      string
	Authenticated bool
}

const appCookie = "myappcookies"

var cookies *sessions.CookieStore

var (
	// Pages
	pageTemplates  = "templates/"
	homePage       = pageTemplates + "home.html"
	registerPage   = pageTemplates + "register.html"
	loginPage      = pageTemplates + "login.html"
	createNotePage = pageTemplates + "createNote.html"
	viewNotePage   = pageTemplates + "viewNote.html"
)

// getUser returns a user from session s. On error, it returns an empty user
func getUser(s *sessions.Session) LoggedInUser {
	val := s.Values["user"]
	var user = LoggedInUser{}
	user, ok := val.(LoggedInUser)
	if !ok {
		return LoggedInUser{Authenticated: false}
	}
	return user
}

func isAuthenticated(w http.ResponseWriter, r *http.Request, session *sessions.Session) bool {
	// session needs to be passed in from calling function here
	user := getUser(session)

	// Check if the user is not authenticated
	if auth := user.Authenticated; !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return false
	}
	return true
}

func redirectToHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func goToHomePage(w http.ResponseWriter, r *http.Request) {
	allNotes, _ := notes.ReadAllNotes(w, r)
	tmpl := template.Must(template.ParseFiles(homePage))
	tmpl.Execute(w, notes.NotePage{Notes: allNotes})
}

func setupRoutes(r *mux.Router) {
	// Home
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, _ := cookies.Get(r, appCookie)
		user := getUser(session)

		// Go to notes list
		if user.Authenticated {
			createNoteValue := r.FormValue("createnote")
			logoutValue := r.FormValue("logout")

			if createNoteValue != "" {
				http.Redirect(w, r, "/create", http.StatusSeeOther)
			} else if logoutValue != "" {
				http.Redirect(w, r, "/logout", http.StatusSeeOther)
			} else {
				goToHomePage(w, r)
			}
		} else {
			// Go to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	})

	// Registration operations
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		session, _ := cookies.Get(r, appCookie)
		user := getUser(session)

		// Go to home page
		if user.Authenticated {
			goToHomePage(w, r)
		} else {
			// Check if input is given
			auth := true
			registerValue := r.FormValue("register")
			backValue := r.FormValue("back")

			if registerValue != "" {
				err := registration.Register(w, r)
				if err != nil {
					auth = false
					tmpl := template.Must(template.ParseFiles(registerPage))
					tmpl.Execute(w, struct {
						Fail    bool
						Message string
					}{true, err.Error()})
				} else {
					user := &LoggedInUser{
						Username:      r.FormValue("username"),
						Authenticated: true,
					}
					session.Values["user"] = user
					session.Save(r, w)
				}
			} else if backValue != "" {
				// Go to login page
				auth = false
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else {
				// Go to registration page
				auth = false
				tmpl := template.Must(template.ParseFiles(registerPage))
				tmpl.Execute(w, nil)
			}

			if auth {
				redirectToHome(w, r)
			}
		}
	})
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		session, _ := cookies.Get(r, appCookie)
		user := getUser(session)

		// Go to home page
		if user.Authenticated {
			goToHomePage(w, r)
		} else {
			// Check if input is given
			auth := true
			registerValue := r.FormValue("register")
			loginValue := r.FormValue("login")

			if loginValue != "" {
				err := registration.Login(w, r)
				if err != nil {
					auth = false
					tmpl := template.Must(template.ParseFiles(loginPage))
					tmpl.Execute(w, struct {
						Fail    bool
						Message string
					}{true, err.Error()})
				} else {
					user := &LoggedInUser{
						Username:      r.FormValue("username"),
						Authenticated: true,
					}
					session.Values["user"] = user
					session.Save(r, w)
				}
			} else if registerValue != "" {
				auth = false
				http.Redirect(w, r, "/register", http.StatusSeeOther)
			} else {
				// Go to registration page
				auth = false
				tmpl := template.Must(template.ParseFiles(loginPage))
				tmpl.Execute(w, nil)
			}

			if auth {
				redirectToHome(w, r)
			}
		}
	})
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		session, _ := cookies.Get(r, appCookie)

		// Log user out
		session.Values["user"] = LoggedInUser{}
		session.Options.MaxAge = -1
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Note operations
	r.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		session, _ := cookies.Get(r, appCookie)
		user := getUser(session)

		if user.Authenticated {
			createValue := r.FormValue("create")
			backValue := r.FormValue("back")

			if createValue != "" {
				err := notes.CreateNote(w, r)
				if err == nil {
					redirectToHome(w, r)
				}
			} else if backValue != "" {
				redirectToHome(w, r)
			} else {
				tmpl := template.Must(template.ParseFiles(createNotePage))
				tmpl.Execute(w, nil)
			}
		} else {
			redirectToHome(w, r)
		}
	})
	r.HandleFunc("/read/{note}", func(w http.ResponseWriter, r *http.Request) {
		session, _ := cookies.Get(r, appCookie)
		user := getUser(session)
		vars := mux.Vars(r)
		var noteID = vars["note"]

		if user.Authenticated {
			if r.FormValue("update") != "" {
				err := notes.UpdateNote(w, r, noteID)
				if err == nil {
					redirectToHome(w, r)
				}
			} else if r.FormValue("delete") != "" {
				err := notes.DeleteNote(w, r, noteID)
				if err == nil {
					redirectToHome(w, r)
				}
			} else if r.FormValue("back") != "" {
				redirectToHome(w, r)
			} else {
				tmpl := template.Must(template.ParseFiles(viewNotePage))
				note, _ := notes.ReadNote(w, r, noteID)
				tmpl.Execute(w, note)
			}
		} else {
			redirectToHome(w, r)
		}
	})
}

func main() {
	// Session Cookies
	cookies = sessions.NewCookieStore([]byte("mysuperdupersecret"))
	cookies.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
	gob.Register(LoggedInUser{})

	r := mux.NewRouter()

	setupRoutes(r)
	http.ListenAndServe(":80", r)
}
