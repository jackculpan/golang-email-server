package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"io"
	"log"
	"net/http"
	"bytes"
	"time"
	"strings"
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


//at regular intervals, the MTA uses the MSA to read and delete a message from the user's outbox, then
//sending the emails from outbox to another email server, whose MTA uses its MSA to add the message to another user's inbox

// gets unsent email(s) from MSA, get network address(s) from BBS, sends email(s), deletes email(s) from MSA outbox
func send(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "Checking for new emails to send\n")
	//MSA network address required for initialization (otherwise it would not know how to connect to it to find the emails)
	res, err := http.Get("http://192.168.1.7:8888/outbox/emails")
	if err != nil {
	  fmt.Println("Err with GET in send method: %v \n", err)
		fmt.Fprintf(w, "No new emails found\n")
	  w.WriteHeader(http.StatusInternalServerError)
	}
	defer res.Body.Close()

	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
	  fmt.Println("Err with reqBody in send method: %v \n", err)
	  w.WriteHeader(http.StatusInternalServerError)
	}
	//decodes new emails found in the MSA's outbox
	dec := json.NewDecoder(strings.NewReader(string(reqBody)))
	for {
		var retrievedEmail email
		if err := dec.Decode(&retrievedEmail); err == io.EOF {
			return
		} else if err != nil {
			fmt.Fprintf(w, "No new emails found\n")
		  return
		}
		//contactBBS is used to find the MSA network address, used to delete the sent email
		MsaAddress := contactBBS(retrievedEmail.From, "Sender")

		req, err := http.NewRequest("DELETE", MsaAddress + "emails/" + retrievedEmail.ID, nil)
		if err != nil {
	  	fmt.Println("Err with DELETE req in send method: %v \n", err)
		  w.WriteHeader(http.StatusInternalServerError)
		}
		res, err = http.DefaultClient.Do(req)
		if err != nil {
	  	fmt.Println("Err with DELETE res in send method: %v \n", err)
		  w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Printf("Email with ID: %s deleted from outbox\n", retrievedEmail.ID)
		defer res.Body.Close()

		//contactBBS used again to find the MTA's network address for the outgoing email
		MtaAddress := contactBBS(retrievedEmail.To, "Outgoing")

		sendEmail(MtaAddress, retrievedEmail)
	}
}

//function contacts the BBS to find the network address for the email to be sent to
func contactBBS(retrievedEmailAddress string, direction string) string {
	var networkAddress string

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(retrievedEmailAddress)

	req, err := http.NewRequest("POST", "http://192.168.1.9:8888/servers", buf)
	if err != nil {
	  fmt.Println("Err with POST newRequest in contactbbs method: %v \n", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
	  fmt.Println("Err with DoRequest in contactbbs method: %v \n", err)
	}
	defer res.Body.Close()
	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
	  fmt.Println("Err with ReadAll in contactbbs method: %v \n", err)
	}

	err = json.Unmarshal(reqBody, &networkAddress)
	if err != nil {
	  fmt.Println("Err with unmashalling in contactbbs method: %v \n", err)
	}
	//section here breaks up the addresses depending on where the email is going
	networkAddresses := strings.Split(networkAddress, "|")
	if direction == "Outgoing" {
		return networkAddresses[1]+"incoming"
	} else if direction == "Incoming" {
		return networkAddresses[0]+"incoming"
	} else if direction == "Sender" {
		return networkAddresses[0]
	} else {
		return networkAddresses[0]
	}
}

//function sends email to network address (MTA) of the recipient
func sendEmail(MtaAddress string, retrievedEmail email) {
	fmt.Printf("Email with ID: %v sent by MTA\n", retrievedEmail.ID)
	buf := new(bytes.Buffer)
	retrievedEmail.Location = "Sending"
	json.NewEncoder(buf).Encode(retrievedEmail)

	req, err := http.NewRequest("POST", MtaAddress, buf)
	if err != nil {
	  fmt.Println(err)
	  return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
	  fmt.Println(err)
	  return
	}
	defer res.Body.Close()
}


//the receive method is used to accept POST requests with new emails from other MTA servers.
func receive(w http.ResponseWriter, r *http.Request) {
	var incomingEmail email

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  fmt.Println(err)
	  w.WriteHeader(http.StatusInternalServerError)
	  return
	}
	json.Unmarshal(reqBody, &incomingEmail)
	incomingEmail.Location = "Received"
	//Contacts the BBS to find the MSA's network address for the server
	MsaAddress := contactBBS(incomingEmail.To, "Incoming")

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(incomingEmail)
	//Makes a POST request with the new email that it has received (done instantly)
	req, err := http.NewRequest("POST", MsaAddress, buf)
	if err != nil {
	  fmt.Println("Err with NewRequest in receive method: %v \n", err)
	  return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
	  fmt.Println("Err with DoRequest in receive method: %v \n", err)
	  return
	}
	defer res.Body.Close()
}

// Home path
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the MTAs home!")
}

// handles all routing via mux
func handleRequests() {
	router := mux.NewRouter().StrictSlash( true )
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/send", send).Methods("GET")
	router.HandleFunc("/incoming", receive).Methods("POST")
	log.Fatal(http.ListenAndServe(":8888", router))
}

func main() {
	handleRequests()
}

