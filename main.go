
package main
//Import libraries
import (
  "fmt"
  "net/http"
  "html/template"

  "database/sql"
  _ "github.com/mattn/go-sqlite3"

  "encoding/json"
  "net/url"
  "io/ioutil"
  "encoding/xml"
)
//Struct Page
type Page struct {
  Name string
  DBStatus bool
}
//Struct for the search result
type SearchResult struct {
  Title string `xml:"title,attr"`
  Author string `xml:"author,attr"`
  Year string `xml:"hyr,attr"`
  ID string `xml:"owi,attr"`
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
    var results []SearchResult
    var err error
    //Call the function search
    if results, err = search(r.FormValue("search")); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
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
//Struct for the Classify the resutl request
type ClassifySearchResponse struct {
  Results []SearchResult `xml:"works>work"`
}
//Function Search
func search(query string) ([]SearchResult, error) {
  //Method from the HTTP package
  var resp *http.Response
  //var for the error
  var err error
  //Request using the method get
  if resp, err = http.Get("http://classify.oclc.org/classify2/Classify?&summary=true&title=" + url.QueryEscape(query)); err != nil {
    return []SearchResult{}, err
  }
  //Read the result
  defer resp.Body.Close()
  var body []byte
  //ioutil read the content xml and set the body with the xml
  if body, err = ioutil.ReadAll(resp.Body); err != nil {
    return []SearchResult{}, err
  }
  //c content the structure
  var c ClassifySearchResponse
  //set c the the xml
  err = xml.Unmarshal(body, &c)
  //return the result or the error
  return c.Results, err
}
