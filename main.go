package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/", index)
	http.HandleFunc("/robots.txt", robots)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	tweets, _, _ := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: 20,
	})
	for i := 0; i < len(tweets); i++ {
		fmt.Fprintf(w, "%+v", tweets[i])
	}
	renderTemplate(w, "index", r)
}

func robots(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("robots.txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%s", content)
}

func renderTemplate(w http.ResponseWriter, fname string, r *http.Request) {
	t, err := template.ParseFiles("templates/base.html", "templates/"+fname+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "base", "")
}
