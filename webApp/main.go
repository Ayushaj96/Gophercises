package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type PageVariables struct {
	Name string
	User string
	Key  string
}

type userData struct {
	Name  string `json:"name"`
	Bio   string `json:"bio"`
	login string
}

func getUserData(user string) (pv PageVariables) {
	url := "https://api.github.com/users/" + user
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
	}

	var p userData
	err = json.Unmarshal(body, &p)

	if err != nil {
		panic(err)
	}

	pv.Name = p.Name
	pv.User = p.Bio
	pv.Key = user
	return pv
}

func main() {
	http.HandleFunc("/user/", HomePage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HomePage(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["key"]

	if !ok || len(keys[0]) < 1 {
		fmt.Println(ok, keys)
		fmt.Println("Url Param 'key' is missing")
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	key := keys[0]

	pv := getUserData(key)

	t, err := template.ParseFiles("assets/index.html") //parse the html file homepage.html
	if err != nil {                                    // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, pv) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {        // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
