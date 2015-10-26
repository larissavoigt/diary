package main

import (
	"fmt"
	"net/http"
)

func main() {
	var body string = `
		<html>
			<head>
				<title>My Diary</title>
			</head>
			<body>
				<h1>Hello World %s</h2>
			</body>
		</html>
	`
	fmt.Println("Hello World!")
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, body, req.URL.Query().Get("name"))
	})
	http.ListenAndServe(":3000", nil)
}
