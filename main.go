package main

// title Orders API
// @version 1.0
// @description This is a first service
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email dbairamkulow@mail.ru
// @host localhost:8081
// @BasePath /

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	_ "fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"io/ioutil"
	"log"
	_ "log"
	"net/http"
	_ "net/http"
	"net/http/httputil"
	"os"
	_ "os"
	"time"
	_ "time"
)

type Metabolomic struct {
	Id        string `json:"Id"`
	Structure string `json:"SMILES"`
	Class     string `json:"Class chemical molecules"`
}

type Options struct {
	homePage            string
	returnAllProteins   string
	createNewProtein    string
	deleteProtein       string
	returnSingleProtein string
}

type Option func(*Options)

var meta []Metabolomic

// returnAllArticles godoc
// @SAccept json
// @Router /all

func logHandler(fn func(w http.ResponseWriter, r *http.Request) Option) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		select {
		case <-time.After(10 * time.Second):
			log.Println(fmt.Sprintf("%q", x))
			f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = fmt.Fprintln(f, fmt.Sprintf("%q", x))
			if err != nil {
				fmt.Println(err)
			}

		}
	}
}

func returnAllProteins(w http.ResponseWriter, r *http.Request) Option {
	return func(args *Options) {

		w.WriteHeader(http.StatusOK)
		fmt.Println("Endpoint Hit: returnAllArticles")
		select {
		case <-time.After(10 * time.Second):
			fmt.Println("Save")

		}
		json.NewEncoder(w).Encode(meta)

	}
}

// homePage godoc
// @SAccept json
// @Router /

func homePage(w http.ResponseWriter, r *http.Request) Option {
	return func(args *Options) {
		fmt.Fprintf(w, "Welcom to the HomePage")
		fmt.Println("Endpoint Hit: homePage")
	}
}

// createNewArticle godoc
// @SAccept json
// @Router /article [post]
func createNewProtein(w http.ResponseWriter, r *http.Request) Option {
	return func(args *Options) {
		reqBody, _ := ioutil.ReadAll(r.Body)
		var met Metabolomic
		json.Unmarshal(reqBody, &met)
		meta = append(meta, met)
		json.NewEncoder(w).Encode(met)
	}
}

// deleteArticle godoc
// @SAccept json
// @Router /article/{id} [delete]
func deleteProtein(w http.ResponseWriter, r *http.Request) Option {
	return func(args *Options) {
		vars := mux.Vars(r)
		id := vars["id"]

		for index, article := range meta {
			if article.Id == id {
				meta = append(meta[:index], meta[index+1:]...)
			}
		}
	}
}

// returnSingleArticle godoc
// @SAccept json
// @Router /article/{id}
func returnSingleProtein(w http.ResponseWriter, r *http.Request) Option {
	return func(args *Options) {
		vars := mux.Vars(r)
		key := vars["id"]
		for _, seq := range meta {
			if seq.Id == key {
				json.NewEncoder(w).Encode(seq)
			}
		}
	}
}

func handleRequests() {
	mainRouter := mux.NewRouter().StrictSlash(true)

	mainRouter.HandleFunc("/", logHandler(homePage))
	mainRouter.HandleFunc("/seq", logHandler(returnAllProteins))
	mainRouter.HandleFunc("/data/{id}", logHandler(returnSingleProtein))
	mainRouter.HandleFunc("/data", logHandler(createNewProtein)).Methods("POST")
	mainRouter.HandleFunc("/data/{id}", logHandler(deleteProtein)).Methods("Delete")

	log.Fatal(http.ListenAndServe(":8083", mainRouter))

}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	meta = []Metabolomic{
		Metabolomic{Id: "1", Structure: "c(6)c(5)o(16)", Class: "acid"},
		Metabolomic{Id: "2", Structure: "c(10)h(15)", Class: "alkens"},
		Metabolomic{Id: "3", Structure: "c(18)h(36)", Class: "alkans"},
	}

	handleRequests()
}
