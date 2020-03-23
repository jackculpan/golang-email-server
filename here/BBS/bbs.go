package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"github.com/gorilla/mux"
)


//returns the server server network address from the source or destination address of an server

type server struct {
	ID          			string `json:"ID"`
	Address 					string `json:"Address"`
	MsaAddress 				string `json:"MsaAddress"`
	MtaAddress 				string `json:"MtaAddress"`
}

type allServers []server

var servers = allServers{
	{
		ID:         	 "1",
		Address: 			 "here.com",
		MsaAddress:    "http://192.168.1.7:8888/",
		MtaAddress:    "http://192.168.1.8:8888/",
	},
	{
		ID:         	 "2",
		Address: 			 "there.com",
		MsaAddress:    "http://192.168.1.2:8888/",
		MtaAddress:    "http://192.168.1.3:8888/",
	},
}

//finds network address for MTA to send email based on the 'To' field of email
func findAddress(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

	var incomingAddress string
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  fmt.Println(err)
	  w.WriteHeader(http.StatusInternalServerError)
	  return
	}

	err = json.Unmarshal(reqBody, &incomingAddress)
	if err != nil {
	  fmt.Println("Err with unmarshalling request in findaddress method: %v \n", err)
	}
	address := strings.Split(incomingAddress, "@")
	//searches through servers to find one that matches, returns both the MSA and MTA address ready for use by the MTA
	for _, singleServer := range servers {
		if singleServer.Address == address[1] {
			json.NewEncoder(w).Encode(singleServer.MsaAddress+"|"+singleServer.MtaAddress)
			return
		}
	}
	fmt.Printf("No server found in BBS - err: %v", http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
	return
}

//returns all servers
func getAllServers(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servers)
}

//allows user to create new server
func createServer(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

	var newServer server
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  fmt.Println(err)
		fmt.Fprintf(w, "Kindly enter data with the server's address, MSA and MTA network addresses only in order to create new server")
	  w.WriteHeader(http.StatusInternalServerError)
	  return
	}
	newServer.ID = strconv.Itoa(len(servers)+1)

	json.Unmarshal(reqBody, &newServer)
	servers = append(servers, newServer)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newServer)
}

//returns one server based on ID
func getOneServer(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
	serverID := mux.Vars(r)["id"]

	for _, singleServer := range servers {
		if singleServer.ID == serverID {
			json.NewEncoder(w).Encode(singleServer)
			return
		}
	}
	fmt.Printf("No server ID found in BBS - err: %v\n", http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
	return
}

//allows updating of server
func updateServer(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
	serverID := mux.Vars(r)["id"]

	var updatedServer server
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  fmt.Println(err)
		fmt.Fprintf(w, "Kindly enter data with the server's address, MSA and MTA network addresses only in order to update server")
	  w.WriteHeader(http.StatusInternalServerError)
	  return
	}
	json.Unmarshal(reqBody, &updatedServer)

	for i, singleServer := range servers {
		if singleServer.ID == serverID {
			servers[i] = updatedServer
			fmt.Fprintf(w, "The server with ID %v has been updated successfully", serverID)
			return
		}
	}
	fmt.Printf("No server ID found in BBS - err: %v\n", http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
	return
}

//allows deletion of server
func deleteServer(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
	serverID := mux.Vars(r)["id"]

	for i, singleServer := range servers {
		if singleServer.ID == serverID {
			servers = append(servers[:i], servers[i+1:]...)
			fmt.Fprintf(w, "The server with ID %v has been deleted successfully", serverID)
			return
		}
	}
	fmt.Printf("No server ID found in BBS - err: %v\n", http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
	return
}

//home route
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the BBS's home!")
}

//handles requests via mux
func handleRequests() {
	router := mux.NewRouter().StrictSlash( true )
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/server", createServer).Methods("POST")
	router.HandleFunc("/servers/{id}", getOneServer).Methods("GET")
	router.HandleFunc("/servers/{id}", updateServer).Methods("PATCH")
	router.HandleFunc("/servers/{id}", deleteServer).Methods("DELETE")
	router.HandleFunc("/servers", getAllServers).Methods("GET")
	router.HandleFunc("/servers", findAddress).Methods("POST")
	log.Fatal(http.ListenAndServe(":8888", router))
}

func main() {
	handleRequests()
}
