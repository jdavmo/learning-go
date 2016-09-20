package main
//Import libraries
import (
  "fmt"
  "net/http"
  "html/template"

  "database/sql"
  _ "github.com/mattn/go-sqlite3"

  "encoding/json"
)
//Struct Page
type Page struct {
  Name string
  DBStatus bool
}
type SearchResult struct {
  Title string
  Author string
  Year string
  ID string
}
//Function main
func main() {
  //Initialize and Parse the template
  templates := template.Must(template.ParseFiles("templates/index.html"))
  //establish conection with the DB
  db, _ := sql.Open("sqlite3", "dev.db")
  //Stating the route /
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    //New instance of page
    p := Page{Name: "Gopher"}
    //Access to the form data with FormValue
    if name := r.FormValue("name"); name != "" {
      //Assing value from get
      p.Name = name
    }
    p.DBStatus = db.Ping() == nil
    //Execute the template
    if err := templates.ExecuteTemplate(w, "index.html", p); err != nil {
      //Alert the user something is wrong
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
  })
  //Stating the route /search
  http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
    //hard code results
    results := []SearchResult{
      SearchResult{"Moby-Dick", "Herman Melville", "1851", "222222"},
      SearchResult{"The Adventures of Huckleberry finn", "Mark Twain", "1884", "44444"},
      SearchResult{"The Catcher in the Rye", "JD Salinger", "1951", "333333"},
    }
    //return json
    encoder := json.NewEncoder(w)
    if err := encoder.Encode(results); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
  })

  //Start the web server using the default mux
  fmt.Println(http.ListenAndServe(":8080", nil))
}
