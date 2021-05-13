package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Basic structr for a http github request to json
type JsonIssue struct {
	Url    string
	Title  string
	Number int
}

// Checks if there's an error, if so... panic!
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Makes an http request
// @param url: the link to make the request
// returns a []byte of the response
func connectHTML(url string) []byte {
	// Make get request
	res, err := http.Get(url)

	checkError(err)

	// read body
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	checkError(err)

	// 200 code is ok, else there's an error
	if res.StatusCode != 200 {
		log.Fatal("Unexpected status code", res.StatusCode)
	}

	return body
}

// Basic structure of the Milestone variable
type Milestone struct {
	name string
	desc string
}

// Basic structure of the data that an Issue will have
type Datos struct {
	Url       string
	Title     string
	User      string
	Number    int
	Tags      []string
	Milestone Milestone
}

// From an http response, it formats the data and filters it
// so we can use only what itis important
// @param body: a http response
// returns an array of Datos with all the issues' information
func getData(body []byte) []Datos {

	// Format all issues into a map
	var issues []map[string]interface{}

	// Formats issues to json
	json.Unmarshal(body, &issues)

	totalIssues := len(issues)

	// To save the filtered data, toReturn
	datos := []Datos{}

	// For each issue...
	for i := 0; i < totalIssues; i++ {
		var datoTemporal Datos

		// Get basic issue data
		url := issues[i]["html_url"]
		title := issues[i]["title"]
		number := issues[i]["number"]
		milestone := issues[i]["milestone"]

		// Save basic data to a temporal Datos struct
		datoTemporal.Url = fmt.Sprintf("%v", url)
		datoTemporal.Title = fmt.Sprintf("%v", title)
		datoTemporal.Number, _ = strconv.Atoi(fmt.Sprintf("%v", number))

		// These variables has more data inside them
		// we have to interate them in order to filter
		// what's important
		usr := issues[i]["user"]
		tags := issues[i]["labels"]

		// Get the name autor
		for key, value := range usr.(map[string]interface{}) {
			//fmt.Println(key, value)
			if key == "login" {
				datoTemporal.User = fmt.Sprintf("%v", value)
				break
			}
		}

		// Get all the name tags
		tagsTemporales := []string{}
		for j := 0; j < len(tags.([]interface{})); j++ {
			var tagTemp string

			for key, value := range tags.([]interface{})[j].(map[string]interface{}) {
				//fmt.Println(key, value)
				if key == "name" {
					tagTemp = fmt.Sprintf("%v", value)
				}
			}

			tagsTemporales = append(tagsTemporales, tagTemp)
		}

		datoTemporal.Tags = tagsTemporales

		// Get milestone name and description
		var mileTemporal Milestone

		// If milestone exists:
		if milestone != nil {
			for key, value := range milestone.(map[string]interface{}) {
				if key == "title" {
					mileTemporal.name = fmt.Sprintf("%v", value)
				}
				if key == "description" {
					mileTemporal.desc = fmt.Sprintf("%v", value)
				}
			}
		} else {
			// If it doesn't exist, save empty values
			mileTemporal.name = " "
			mileTemporal.desc = " "
		}

		// Save to our datos array
		datoTemporal.Milestone = mileTemporal
		datos = append(datos, datoTemporal)
	}

	return datos
}

// Opens a sqlite database and
// saves the date of Datos struct
// @param vl: Data to save
func insertQuery(vl Datos) {
	// Opens database
	db, err := sql.Open("sqlite3", "./src/goRepoDB.db")

	checkError(err)

	// Preparse the query
	stmt, err := db.Prepare("INSERT into issues (numero, url, nombre, autor, mile_nombre, mile_desc, tagArray) VALUES (?, ?, ?, ?, ?, ?, ?)")
	checkError(err)

	// Save tags as one string and execute query
	justTags := strings.Join(vl.Tags, " ")
	_, err = stmt.Exec(vl.Number, vl.Url, vl.Title, vl.User, vl.Milestone.name, vl.Milestone.desc, justTags)

	// Close
	checkError(err)
	db.Close()
}

// Saves all the info in Datos array
// into a database
// @param dt: The Data array with the issues information
func saveData(dt []Datos) {

	records := len(dt)

	// For each record, save into the database
	for record := 0; record < records; record++ {
		//TOFIX. Check if it has already been saved
		insertQuery(dt[record])
	}
}

// Prints to the user all the infromation saved in the database
// @param table: The name of the table to be printed with the correct fromat
func printDataTable(table string) {
	// Open database
	db, err := sql.Open("sqlite3", "./src/goRepoDB.db")
	checkError(err)

	defer db.Close()

	// Select everything from the table
	query := fmt.Sprintf("SELECT * FROM %s", table)
	rows, err := db.Query(query)
	checkError(err)

	// Data to be saved into
	var number int
	var url string
	var title string
	var user string
	var mile_name string
	var mile_desc string
	var tags string

	// For each row
	for rows.Next() {
		// Save the info into our variables
		err = rows.Scan(&number, &url, &title, &user, &mile_name, &mile_desc, &tags)
		checkError(err)

		// Print them
		fmt.Printf("\nIssue #%d\n", number)
		fmt.Printf("- URL: %s\n", url)
		fmt.Printf("- Nombre: %s\n", title)
		fmt.Printf("- Autor: %s\n", user)
		fmt.Printf("- Tags: %s\n", tags)
		fmt.Printf("- Milestone\n  - Nombre: %s\n", mile_name)
		fmt.Printf("  - Desc: %s\n", mile_desc)
	}

}

// Checks into a github repository
// and get's all the issued with a given tag
// @param user: the name of the user repo
// @param repo: the github repository name
// @param label: the name of the lable to be filter on
// return an array with the issues data
func getIssues(user, repo, label string) []Datos {
	// Format the http link
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?labels=%s&page=1&per_page=100", user, repo, label)

	// Get response
	body := connectHTML(url)

	// Filter data
	data := getData(body)

	return data

}

func main() {
	table := "issues"

	// Get issues data from go's repo with Go2 tag
	data := getIssues("golang", "go", "Go2")

	// Save it to a database
	saveData(data)

	// Print it to the user
	printDataTable(table)
}
