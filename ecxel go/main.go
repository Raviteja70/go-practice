package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/tealeg/xlsx"
)

// Data represents the structure of data to be stored and exported
type Data struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Mobile   string `json:"mobile"`
	Location string `json:"location"`
}

func main() {
	// Connect to PostgreSQL
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=Ravi@123 dbname=Practice sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new router
	r := mux.NewRouter()

	// Route to insert data into PostgreSQL
	r.HandleFunc("/insert", InsertDataHandler(db)).Methods("POST")

	// Route to import data from Excel sheet
	r.HandleFunc("/import", ImportDataHandler(db)).Methods("POST")

	// Route to export data to Excel
	r.HandleFunc("/export", ExportDataHandler(db)).Methods("GET")

	// Route to retrieve data from database
	r.HandleFunc("/data", GetDataHandler(db)).Methods("GET")

	// Start the server
	http.Handle("/", r)
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// InsertDataHandler handles insertion of data into PostgreSQL
func InsertDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data Data
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		_, err = db.Exec("INSERT INTO data (name, email, role, mobile, location) VALUES ($1, $2, $3, $4, $5)", data.Name, data.Email, data.Role, data.Mobile, data.Location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Data inserted successfully")
	}
}

// ImportDataHandler handles importing data from an Excel sheet and inserting into PostgreSQL
func ImportDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get the file from the request
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Create a temporary file to save the uploaded Excel file
		tempFile, err := ioutil.TempFile("", "uploaded-*.xlsx")
		if err != nil {
			http.Error(w, "Error creating temporary file", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())

		// Copy the file content to the temporary file
		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Error(w, "Error copying file", http.StatusInternalServerError)
			return
		}

		// Open the temporary file with xlsx.OpenFile
		xlFile, err := xlsx.OpenFile(tempFile.Name())
		if err != nil {
			http.Error(w, "Error reading Excel file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Iterate over the sheets in the Excel file
		for _, sheet := range xlFile.Sheets {
			// Iterate over the rows in each sheet
			for _, row := range sheet.Rows {
				var data Data
				// Extract data from each row
				for idx, cell := range row.Cells {
					text := cell.String()
					switch idx {
					case 0:
						data.Name = text
					case 1:
						data.Email = text
					case 2:
						data.Role = text
					case 3:
						data.Mobile = text
					case 4:
						data.Location = text
					}
				}

				// Insert data into the database
				_, err := db.Exec("INSERT INTO data (name, email, role, mobile, location) VALUES ($1, $2, $3, $4, $5)", data.Name, data.Email, data.Role, data.Mobile, data.Location)
				if err != nil {
					http.Error(w, "Error inserting data into database: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Data imported successfully")
	}
}

// ExportDataHandler handles exporting data to an Excel file
func ExportDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query the database to retrieve data
		rows, err := db.Query("SELECT * FROM data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var data []Data
		for rows.Next() {
			var d Data
			err := rows.Scan(&d.ID, &d.Name, &d.Email, &d.Role, &d.Mobile, &d.Location)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data = append(data, d)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create a new Excel file
		file := xlsx.NewFile()
		sheet, err := file.AddSheet("Data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write data to Excel file
		for _, d := range data {
			row := sheet.AddRow()
			cell := row.AddCell()
			cell.Value = d.Name
			cell = row.AddCell()
			cell.Value = d.Email
			cell = row.AddCell()
			cell.Value = d.Role
			cell = row.AddCell()
			cell.Value = d.Mobile
			cell = row.AddCell()
			cell.Value = d.Location
		}

		// Save the Excel file
		filename := "data.xlsx"
		err = file.Save(filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Serve the file for download
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		http.ServeFile(w, r, filename)
	}
}

// GetDataHandler handles retrieving data from the database
func GetDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var data []Data
		for rows.Next() {
			var d Data
			err := rows.Scan(&d.ID, &d.Name, &d.Email, &d.Role, &d.Mobile, &d.Location)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data = append(data, d)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Encode data to JSON and write response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
