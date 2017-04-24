package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Stuff struct { //Define type Stuff data
	/*
		Inside each of the tags there is an omitempty parameter.
		This means that if the property is null, it will be excluded from the JSON data rather than showing up as an empty string or value.
	*/
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	/*
		Inside the Person struct there is an Address property that is a pointer.
		This will represent a nested JSON object and it must be a pointer otherwise the omitempty will fail to work.
	*/
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var person []Stuff

func GetStuffEndpoint(w http.ResponseWriter, req *http.Request) {
	/*
		In the above GetPersonEndpoint we are trying to get a single record.
		Using the mux library we can get any parameters that were passed in with the request.
		We then loop over our global slice and look for any ids that match the id found in the request parameters.
		If a match is found, use the JSON encoder to display it, otherwise create an empty JSON object.
	*/
	params := mux.Vars(req)
	for _, item := range person {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Stuff{})
}

func GetPersonEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(person)
}

func CreateStuffEndpoint(w http.ResponseWriter, req *http.Request) {

	/*
		In the above we decode the JSON data that was passed in and store it in a Person object.
		We assign the new object an id based on what mux found and then we append it to our global slice.
		In the end, our global array will be returned and it should include everything including our newly added piece of data.
	*/
	params := mux.Vars(req)
	var stuff Stuff
	_ = json.NewDecoder(req.Body).Decode(&stuff)
	stuff.ID = params["id"]
	person = append(person, stuff)
	json.NewEncoder(w).Encode(person)
}

func DeleteStuffEndpoint(w http.ResponseWriter, req *http.Request) {
	//When the id to be deleted has been found, we can recreate our slice with all data excluding that found at the index.
	params := mux.Vars(req)
	for index, item := range person {
		if item.ID == params["id"] {
			person = append(person[:index], person[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(person)
}

func main() {
	/*
		We first create our new router and add two objects to our slice to get things started.
		Next up we have to create each of the endpoints that will call our endpoint functions.
		Notice we are using GET, POST, and DELETE where appropriate.
		We are also defining parameters that can be passed in.
	*/
	router := mux.NewRouter()
	person = append(person, Stuff{ID: "1", Firstname: "Cenk", Lastname: "Tosun", Address: &Address{City: "Istanbul", State: "Turkey"}})
	person = append(person, Stuff{ID: "2", Firstname: "Gokhan", Lastname: "Gonul"})
	router.HandleFunc("/person", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", GetStuffEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", CreateStuffEndpoint).Methods("POST")
	router.HandleFunc("/person/{id}", DeleteStuffEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
