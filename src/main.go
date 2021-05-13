package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type JsonIssue struct {
	Url    string
	Title  string
	Number int
}

func connectHTML(url string) []byte {
	// Make get request
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	// read body
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatal("Unexpected status code", res.StatusCode)
	}

	return body
}

type Milestone struct {
	name string
	desc string
}

type Datos struct {
	Url       string
	Title     string
	User      string
	Number    int
	Tags      []string
	Milestone Milestone
}

func getData(body []byte) []Datos {
	// Get url and title and number
	var issues []map[string]interface{}
	json.Unmarshal(body, &issues)

	totalIssues := len(issues)

	datos := []Datos{}

	for i := 0; i < totalIssues; i++ {
		var datoTemporal Datos

		url := issues[i]["url"]
		title := issues[i]["title"]
		number := issues[i]["number"]
		milestone := issues[i]["milestone"]

		datoTemporal.Url = fmt.Sprintf("%v", url)
		datoTemporal.Title = fmt.Sprintf("%v", title)
		datoTemporal.Number, _ = strconv.Atoi( fmt.Sprintf("%v", number) )

		usr := issues[i]["user"]
		tags := issues[i]["labels"]

		for key, value := range usr.(map[string]interface{}) {
			//fmt.Println(key, value)
			if key == "login" {
				datoTemporal.User = fmt.Sprintf("%v", value)
				break
			}
		}

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


		var mileTemporal Milestone
		for key, value := range milestone.(map[string]interface{}) {
			if key == "title" {
				mileTemporal.name = fmt.Sprintf("%v", value)
			}
			if key == "description" {
				mileTemporal.desc = fmt.Sprintf("%v", value)
			}
		}

		datoTemporal.Milestone = mileTemporal
		datos = append(datos, datoTemporal)

	}

	return datos
}

func saveData(dt []Datos) {
	records := len(dt)
	for record := 0; record < records; record++ {
		fmt.Println(dt[record].Url)

		
	}
}

func getIssues(user, repo, label string) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?labels=%s&page=1&per_page=1", user, repo, label)
	body := connectHTML(url)
	data := getData(body)
	saveData(data)
}

func main() {
	getIssues("golang", "go", "Go2")

}
