package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// Set connection info of PostgreSQL Server
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

const (
	host     = "phonebookdb"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// DB set up
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func setupDB() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)

	checkErr(err)

	return db
}

//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// Define Contact structure
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

type Contact struct {
	PhoneNumber string `json:"phonenumber"`
	FullName    string `json:"fullname"`
	Address     string `json:"address"`
	Email       string `json:"email"`
}

//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// Define JsonResponse structure
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

type JsonResponse struct {
	Type    string    `json:"type"`
	Data    []Contact `json:"data"`
	Message string    `json:"message"`
}

//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// Main function
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func main() {

	// Init the mux router
	router := mux.NewRouter()

	// Route handles & endpoints

	router.HandleFunc("/contacts/search/", SearchContact).Methods("POST")

	router.HandleFunc("/contacts/get-all/", GetAllContacts).Methods("GET")

	router.HandleFunc("/contacts/new/", CreateNewContact).Methods("POST")

	router.HandleFunc("/contacts/delete-by-number/{phonenumber}", DeleteContactByNumber).Methods("DELETE")

	router.HandleFunc("/contacts/delete-by-name/", DeleteContactByName).Methods("DELETE")

	router.HandleFunc("/contacts/delete-all/", DeleteAllContacts).Methods("DELETE")

	// serve the app
	fmt.Println("Server running at port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// handling messages
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// handling errors
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// GetAllContacts
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func GetAllContacts(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Getting All contacts...")

	rows, err := db.Query("SELECT * FROM contacts")

	checkErr(err)

	var contacts []Contact

	// Foreach movie
	for rows.Next() {
		var id int
		var phonenumber string
		var fullname string
		var address string
		var email string

		err = rows.Scan(&id, &phonenumber, &fullname, &address, &email)

		checkErr(err)

		contacts = append(contacts, Contact{PhoneNumber: phonenumber, FullName: fullname, Address: address, Email: email})
	}

	var response = JsonResponse{Type: "success", Data: contacts}

	json.NewEncoder(w).Encode(response)
}

//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// CreateNewContact
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func CreateNewContact(w http.ResponseWriter, r *http.Request) {
	PhoneNumber := r.FormValue("phonenumber")
	FullName := r.FormValue("fullname")
	Address := r.FormValue("address")
	Email := r.FormValue("email")

	var response = JsonResponse{}

	if PhoneNumber == "" || FullName == "" || Address == "" || Email == "" {
		response = JsonResponse{Type: "error", Message: "You are missing PhoneNumber or FullName or Address or Email parameter."}
	} else {
		db := setupDB()

		fmt.Println("Inserting new contact into DB")
		fmt.Println("Name:" + FullName)
		fmt.Println("PhoneNumber:" + PhoneNumber)
		fmt.Println("Email:" + Email)
		fmt.Println("Address:" + Address)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO contacts(PhoneNumber, FullName, Address, Email) VALUES($1, $2, $3, $4) returning id;", PhoneNumber, FullName, Address, Email).Scan(&lastInsertID)

		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The contact has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// DeleteContactByNumber
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func DeleteContactByNumber(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	PhoneNumber := params["phonenumber"]

	var response = JsonResponse{}

	if PhoneNumber == "" {
		response = JsonResponse{Type: "error", Message: "You are missing PhoneNumber parameter."}
	} else {
		db := setupDB()

		fmt.Println("Deleting a contact from DB")

		_, err := db.Exec("DELETE FROM contacts where PhoneNumber = $1", PhoneNumber)

		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The contact has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// DeleteAllContacts
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func DeleteAllContacts(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Deleting all contacts...")

	_, err := db.Exec("DELETE FROM contacts")

	checkErr(err)

	printMessage("All contacts have been deleted successfully!")

	var response = JsonResponse{Type: "success", Message: "All contacts have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}

//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// Delete Contact by Name
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func DeleteContactByName(w http.ResponseWriter, r *http.Request) {
	FullName := r.FormValue("fullname")

	var response = JsonResponse{}

	if FullName == "" {
		response = JsonResponse{Type: "error", Message: "You are missing FullName parameter."}
	} else {
		db := setupDB()

		fmt.Println("Deleting a contact from DB")

		_, err := db.Exec("DELETE FROM contacts where FullName = $1", FullName)

		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The contact has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// Search in Contact items
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func SearchContact(w http.ResponseWriter, r *http.Request) {
	Search := r.FormValue("search")

	var response = JsonResponse{}

	if Search == "" {
		response = JsonResponse{Type: "error", Message: "You are missing Search parameter."}
		fmt.Println("You are missing Search parameter.")
	} else {
		db := setupDB()

		fmt.Println("Searching for:" + Search)

		rows, err := db.Query("SELECT * FROM contacts WHERE fullname ~* $1 ", Search)

		checkErr(err)

		var contacts []Contact

		// Foreach movie
		for rows.Next() {
			var id int
			var phonenumber string
			var fullname string
			var address string
			var email string

			err = rows.Scan(&id, &phonenumber, &fullname, &address, &email)

			checkErr(err)

			contacts = append(contacts, Contact{PhoneNumber: phonenumber, FullName: fullname, Address: address, Email: email})

			response = JsonResponse{Type: "success", Data: contacts}
		}
	}
	json.NewEncoder(w).Encode(response)
}
