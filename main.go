package main

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"fmt"
)

type FirebaseContext struct {
	w   http.ResponseWriter
	r   http.Request
	c   http.Client      // client used to GET from ctftime.org
	ctx context.Context  // context used in connection to Firestore
	fb  firestore.Client // client used in connection to Firestore
}

func fetch(url string, fbc FirebaseContext) *http.Response {
	resp, err := fbc.c.Get(url)
	if err != nil {
		http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return resp
}

func checkCurrentRankingsHandler(w http.ResponseWriter, r *http.Request) {
	//highestNode := int(getHighestNode(fbc))
	token := generateToken()
	FbClient, ctx := connect(token)
	if FbClient != nil && ctx != nil {
		fbc := FirebaseContext{
			w, *r, http.Client{}, ctx, *FbClient,
		}
		for i := 1; i < 287; i++ {
			url := fmt.Sprintf("https://ctftime.org/stats/2017?page=%d", i)
			response := fetch(url, fbc)
			if response != nil {
				err := parseAndStoreRankings(response, i, fbc)
				if err != nil {
					http.Error(fbc.w, err.Error(), http.StatusInternalServerError)
				}
			}
		}
		fbc.w.Write([]byte("Finished storing contents"))
	} else {
		http.Error(w,
			"Failed to connect to Firestore",
			http.StatusInternalServerError)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {}

func main() {
	http.HandleFunc("/favicon.ico", defaultHandler)
	http.HandleFunc("/", checkCurrentRankingsHandler)
	http.ListenAndServe("localhost:8080", nil)
}
