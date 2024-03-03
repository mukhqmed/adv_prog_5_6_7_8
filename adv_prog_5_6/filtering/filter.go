// filtering/filter.go
package filtering

import (
	"database/sql"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Barber struct {
	ID         int
	Name       string
	BasicInfo  string
	Price      int
	Experience string
	Status     string
	ImagePath  string
}

var db *sql.DB

func Init(database *sql.DB) {
	db = database
}

func FilteredBarbersHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	statusFilter := r.URL.Query().Get("status")
	experienceFilter := r.URL.Query().Get("experience")
	sortBy := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	itemsPerPage := 3

	log.WithFields(log.Fields{
		"action":           "filter_barbers",
		"timestamp":        time.Now().Format(time.RFC3339),
		"statusFilter":     statusFilter,
		"experienceFilter": experienceFilter,
		"sortBy":           sortBy,
		"page":             pageStr,
	}).Info("Filtering and sorting barbers")

	// Convert page number to integer
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Calculate offset
	offset := (page - 1) * itemsPerPage

	// Construct SQL query based on filters, sorting, and pagination
	query := "SELECT id, name, basic_info, price, experience, status, image_path FROM barbers WHERE true"
	if statusFilter != "" {
		query += " AND status = '" + statusFilter + "'"
	}
	if experienceFilter != "" {
		query += " AND experience = '" + experienceFilter + "'"
	}
	switch sortBy {
	case "name":
		query += " ORDER BY name"
	case "price":
		query += " ORDER BY price"
	default:
		query += " ORDER BY id" // Default sorting
	}
	query += " LIMIT " + strconv.Itoa(itemsPerPage) + " OFFSET " + strconv.Itoa(offset)

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		log.WithFields(log.Fields{
			"action":    "filter_barbers",
			"timestamp": time.Now().Format(time.RFC3339),
			"error":     err,
		}).Error("Error occurred while executing query")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Fetch results
	var barbers []Barber
	for rows.Next() {
		var b Barber
		if err := rows.Scan(&b.ID, &b.Name, &b.BasicInfo, &b.Price, &b.Experience, &b.Status, &b.ImagePath); err != nil {
			log.WithFields(log.Fields{
				"action":    "filter_barbers",
				"timestamp": time.Now().Format(time.RFC3339),
				"error":     err,
			}).Error("Error occurred while scanning row")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		barbers = append(barbers, b)
	}
	if err := rows.Err(); err != nil {
		log.WithFields(log.Fields{
			"action":    "filter_barbers",
			"timestamp": time.Now().Format(time.RFC3339),
			"error":     err,
		}).Error("Error occurred while iterating rows")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the template with the data
	tmpl, err := template.ParseFiles("barbers.html")
	if err != nil {
		log.WithFields(log.Fields{
			"action":    "filter_barbers",
			"timestamp": time.Now().Format(time.RFC3339),
			"error":     err,
		}).Error("Error occurred while parsing template")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, gin.H{"Barbers": barbers}); err != nil {
		log.WithFields(log.Fields{
			"action":    "filter_barbers",
			"timestamp": time.Now().Format(time.RFC3339),
			"error":     err,
		}).Error("Error occurred while executing template")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func GetBarbersFromDB(db *sql.DB) ([]Barber, error) {
	// Query to select all barbers from the database
	rows, err := db.Query("SELECT id, name, basic_info, price, experience, status, image_path FROM barbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows and parse barber data
	var barbers []Barber
	for rows.Next() {
		var barber Barber
		err := rows.Scan(&barber.ID, &barber.Name, &barber.BasicInfo, &barber.Price, &barber.Experience, &barber.Status, &barber.ImagePath)
		if err != nil {
			return nil, err
		}
		barbers = append(barbers, barber)
	}

	return barbers, nil
}
