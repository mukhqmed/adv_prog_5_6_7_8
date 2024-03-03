package auth

import (
	"database/sql"
	"html/template"
	"net/http"
	"net/smtp"
	"strconv"
)

var db *sql.DB

func Init(database *sql.DB) {
	db = database
}
func SendVerificationEmail(email string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	senderEmail := "gerardfrr@gmail.com"
	senderPassword := "talm rfiq pvqm vupp"

	recipientEmail := email

	message := []byte("Subject: Registration Confirmation\r\n" +
		"\r\n" +
		"Please click the following link to verify your email: http://localhost:8080/verify?email=" + email)

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	err := smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, senderEmail, []string{recipientEmail}, message)
	if err != nil {
		return err
	}

	return nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl, err := template.ParseFiles("register.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	case "POST":
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		ageStr := r.FormValue("age")
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			http.Error(w, "Invalid age", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO users (username, email, password, firstname, lastname, age) VALUES ($1, $2, $3, $4, $5, $6)", username, email, password, firstname, lastname, age)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send verification email
		err = SendVerificationEmail(email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if email == "" {
		http.Error(w, "Email parameter is missing", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var storedPassword string
	var role string
	err := db.QueryRow("SELECT password, role FROM users WHERE username = $1", username).Scan(&storedPassword, &role)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	if password != storedPassword {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Redirect based on role
	switch role {
	case "admin":
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	case "regular user":
		http.Redirect(w, r, "/user", http.StatusSeeOther)
	default:
		http.Error(w, "Unknown role", http.StatusInternalServerError)
	}
}
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
