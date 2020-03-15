package main

import (
	"context"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"

	"./findURL"
	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*PORT http server port */
/*HOST http server host ip address */
/*MONGODBNAME mongo database name */
/*MONGOHOST mongo database host address */
/*MONGODBCOLLECTINNAME mongo database url collenction name */
const (
	PORT                 string = "8080"
	HOST                 string = "127.0.0.1"
	MONGODBNAME          string = "url-shortener"
	MONGOHOST            string = "localhost:27017"
	MONGODBCOLLECTINNAME string = "urls"
	SHORTURLPREFIX       string = "/p"
)

func main() {

	fmt.Println("Hello web server" + HOST + ":" + PORT)
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("this is web page")
		fmt.Println(w, "Hello, %q", html.EscapeString(r.URL.Path))
		tmpl, _ := template.ParseFiles("view/pages/home.html")
		pageDetail := PageData{
			PageTitle: "Home",
			URL:       "",
		}
		tmpl.Execute(w, pageDetail)
	})

	r.HandleFunc("/create-url", createURL)
	r.HandleFunc("/p/{key}", findData)
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("view/css"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         HOST + ":" + PORT,
		WriteTimeout: 100 * time.Second,
		ReadTimeout:  100 * time.Second,
	}

	srv.ListenAndServe()
}

func createURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("invalid_http_method")
		return

	}
	r.ParseForm()
	log.Println(r.Form)
	var randomShortURL = uniuri.NewLen(6)
	fmt.Printf("Random string: %s", randomShortURL)
	fmt.Printf("URL: %s", r.Form.Get("url"))

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://<user>:<password>@" + MONGOHOST + "/" + MONGODBNAME + "?retryWrites=true&w=majority&authMechanism=SCRAM-SHA-256")
	clientOptions.SetConnectTimeout(20 * time.Second)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal("client error: ", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("ping error: ", err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database(MONGODBNAME).Collection(MONGODBCOLLECTINNAME)
	insertResult, err := collection.InsertOne(context.TODO(), bson.M{"url": r.Form.Get("url"), "value": randomShortURL})
	if err != nil {
		log.Fatal("insert error: ", err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	tmpl, _ := template.ParseFiles("view/pages/home.html")
	pageDetail := PageData{
		URL:       HOST + ":" + string(PORT) + SHORTURLPREFIX + "/" + randomShortURL,
		PageTitle: "Home",
	}

	tmpl.Execute(w, pageDetail)

}

func findData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tmpl, _ := template.ParseFiles("view/pages/home.html")

	foundedURL := findURL.Find(vars["key"])
	var message string
	if foundedURL == "" {
		message = "INVALID_KEY"
	}

	u, err := url.ParseRequestURI(foundedURL)
	if err != nil || u == nil {
		message = "INVALID_URL"
	}

	log.Println("VALID_URL: ", u)

	if message == "" {
		http.Redirect(w, r, foundedURL, http.StatusSeeOther)
	}

	pageDetail := PageData{
		URL:       "",
		PageTitle: "Home",
		Message:   message,
	}
	tmpl.Execute(w, pageDetail)

}

/*PageData this struct on page details*/
type PageData struct {
	PageTitle string
	URL       string
	Message   string
}
