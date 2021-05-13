package main

import "fmt"

func connectHTML(url string) {
	fmt.Println(url)
}

func getIssues(user, repo, label string) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?labels=%s&page=1&per_page=5", user, repo, label)
	connectHTML(url)
}

func main() {
	fmt.Println("Hola")
	getIssues("golang", "go", "Go2")
}
