package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/read", read)
	http.HandleFunc("/crash", crash)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func read(rw http.ResponseWriter, req *http.Request) {
	key := req.FormValue("key")
	ret := fetchEtcdValue(key)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(ret)
}

func fetchEtcdValue(key string) []byte {
	ret, err := http.Get("http://127.0.0.1:4001/v2/keys/" + key)
	if err != nil {
		log.Printf("ERROR: Could not get key: %s, we got error: %s", key, err.Error())
		return nil
	}
	defer ret.Body.Close()
	body, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		log.Printf("ERROR: Could not read the body response, we got error: %s", err.Error())
		return nil
	}

	return body
}

func crash(rw http.ResponseWriter, req *http.Request) {
	log.Fatal("Good bye!")
}
