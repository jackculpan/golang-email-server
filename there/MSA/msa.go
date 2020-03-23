package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/gorilla/mux"
)

type email struct {
	ID          string 		`json:"ID"`
	To          string 		`json:"To"`
	From        string 		`json:"From"`
	Subject     string 		`json:"Subject"`
	Body 				string 		`json:"Body"`
	Time				time.Time `json:"Time"`
	Status 		  string 		`json:"Status"`
	Location 		string 		`json:"Location"`
}

type allEmails []email

var emails = allEmails{
	{
		ID:          "1",
		To:					 "fred@here.com",
		From:				 "barney@there.com",
		Subject:     "EMAIL 1 FROM THERE",
		Body: 			 "you are receiving this email from the THERE.COM inbox",
		Time:				 time.Now(),
		Status:		 	 "Unread",
		Location:		 "Outbox",
	},
	{
		ID:          "2",
		To:					 "wilma@here.com",
		From:				 "betty@there.com",
		Subject:     "EMAIL 2 FROM THERE",
		Body: 			 "you are receiving this email 2nd email from the THERE.COM inbox",
		Time:				 time.Now(),
		Status:		 	 "Unread",
		Location:		 "Outbox",
	},
}


//msa lists messages in the user's outbox and inbox, can also mark as read/delete
//msa moves newly created email to outbox


//creates email using post request from user. Newly created email is added to the outbox ready for the MTA
func createEmail(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
	var newEmail email

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  fmt.Println(err)
		fmt.Fprintf(w, "Please enter data with the email to, from, subject and body only in order to create new email")
	  w.WriteHeader(http.StatusBadRequest)
	  return
	}
	newEmail.ID = strconv.Itoa(len(emails)+1)
	newEmail.Time = time.Now()
	newEmail.Location = "Outbox"
	newEmail.Status = "Unread"
	//creates new email with some preset information - once created it gets added to the outbox
	json.Unmarshal(reqBody, &newEmail)
	emails = append(emails, newEmail)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEmail)
}

// can be used by the user to find one email by ID
func getOneEmail(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
	emailID := mux.Vars(r)["id"]

	for _, singleEmail := range emails {
		if singleEmail.ID == emailID {
			json.NewEncoder(w).Encode(singleEmail)
			return
		}
	}
	fmt.Printf("No email with ID %v found - err: %v\n", emailID, http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
	return
}

// can be used by the user to find all the emails in their outbox (function used by the MTA)
func getAllOutbox(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  count := 0

	for _, singleEmail := range emails {
		if singleEmail.Location == "Outbox" {
			json.NewEncoder(w).Encode(singleEmail)
			count+=1
		}
	}
	if count == 0 {
		w.WriteHeader(http.StatusNoContent)
	}
}

// can be used by the user to find all the emails in their outbox
func getAllInbox(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  count := 0

	for _, singleEmail := range emails {
		if singleEmail.Location == "Inbox" {
			json.NewEncoder(w).Encode(singleEmail)
			count+=1
		}
	}
	if count == 0 {
		fmt.Printf("No emails in Inbox - err: %v\n", http.StatusNoContent)
		w.WriteHeader(http.StatusNoContent)
	}
}

// can be used by the user to find all the emails in their outbox and inbox
func getAllEmails(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
	if len(emails) == 0 {
		fmt.Fprintf(w, "No emails available")
		return
	}
	json.NewEncoder(w).Encode(emails)
}

//Method is used by the MTA to send across new emails that have been sent to the user - moves them to inbox
func addToInbox(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

	var incomingEmail email
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  fmt.Println(err)
	  w.WriteHeader(http.StatusInternalServerError)
	  return
	}
	json.Unmarshal(reqBody, &incomingEmail)
	incomingEmail.Location = "Inbox"
	incomingEmail.ID = strconv.Itoa(len(emails)+1)
	for _, singleEmail := range emails {
		if incomingEmail.ID == singleEmail.ID {
			incomingEmail.ID = strconv.Itoa(len(emails))
		}
	}
	json.NewEncoder(w).Encode(incomingEmail)
	emails = append(emails, incomingEmail)
	fmt.Printf("Email with ID: %v received by MSA from MTA\n", incomingEmail.ID)
}

// can be used by the user to mark an email as read in their inbox/outbox
func markRead(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
	emailID := mux.Vars(r)["id"]

	for i, singleEmail := range emails {
		if singleEmail.ID == emailID {
			//singleEmail.Status = "Read"
			emails[i].Status = "Read"
			fmt.Fprintf(w, "The email with ID %v has been marked as read successfully", emailID)
			return
		}
	}
	fmt.Printf("No email with ID %v found - err: %v\n", emailID, http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
	return
}
// can be used by the user to delete emails from their inbox/outbox
func deleteEmail(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
	emailID := mux.Vars(r)["id"]

	for i, singleEmail := range emails {
		if singleEmail.ID == emailID {
			emails = append(emails[:i], emails[i+1:]...)
			fmt.Fprintf(w, "The email with ID %v has been deleted successfully", emailID)
			return
		}
	}
	fmt.Printf("No email with ID %v found - err: %v\n", emailID, http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
	return
}

// home path
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the MSA's home!")
}

// handles all routing via mux
func handleRequests() {
	router := mux.NewRouter().StrictSlash( true )
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/email", createEmail).Methods("POST")
	router.HandleFunc("/inbox/emails", getAllInbox).Methods("GET")
	router.HandleFunc("/incoming", addToInbox).Methods("POST")
	router.HandleFunc("/outbox/emails", getAllOutbox).Methods("GET")
	router.HandleFunc("/emails", getAllEmails).Methods("GET")
	router.HandleFunc("/emails/{id}", getOneEmail).Methods("GET")
	router.HandleFunc("/emails/{id}", markRead).Methods("PATCH")
	router.HandleFunc("/emails/{id}", deleteEmail).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8888", router))
}

func main() {
	handleRequests()
}
