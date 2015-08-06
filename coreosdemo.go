package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

func increaseEtcdCnt() {

	data := fmt.Sprintf("value=%+v", rand.Int())
	log.Printf("here ut comes %s", data)
	req, _ := http.NewRequest("PUT", "http://192.168.122.125:2379/v2/keys/cnt", bytes.NewReader([]byte(data)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: Could not update key in etcd: %v", err)
		return
	}

	defer res.Body.Close()
	return

}

func main() {
	http.HandleFunc("/read", read)
	http.HandleFunc("/crash", crash)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func read(rw http.ResponseWriter, req *http.Request) {
	increaseEtcdCnt()
	key := req.FormValue("key")
	ret := fetchEtcdValue(key)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(ret)
}

func fetchEtcdValue(key string) []byte {
	ret, err := http.Get("http://192.168.122.125:2379/v2/keys/" + key)
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
