package main
import "fmt"
import "net/http"
func main() {
  fmt.Println("Hello World!")
  http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
    fmt.Fprint(res, "hello")
} )
http.ListenAndServe(":3000", nil)
}
