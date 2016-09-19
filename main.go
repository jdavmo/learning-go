package main
//Import libraries
import (
  "fmt"
  "net/http"
  "html/template"

  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)
//Struct Page
type Page struct {
  Name string
  DBStatus bool
}
//Function main
func main() {
  //Initialize and Parse the template
  templates := template.Must(template.ParseFiles("templates/index.html"))
  //establish conection with the DB
  db, _ := sql.Open("sqlite3", "dev.db")
  //Stating the route
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
  //Start the web server using the default mux
  fmt.Println(http.ListenAndServe(":8080", nil))
}
