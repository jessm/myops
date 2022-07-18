package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const pingFile string = "ping.txt"

var fileLock sync.Mutex

func handler(w http.ResponseWriter, r *http.Request) {
	fileLock.Lock()
	defer fileLock.Unlock()

	data, err := ioutil.ReadFile(pingFile)
	if err != nil {
		http.Error(w, "Internal error 1: "+err.Error(), http.StatusInternalServerError)
		return
	}

	prevPingCount, err := strconv.Atoi(string(data))
	if err != nil {
		http.Error(w, "Internal error 2: "+err.Error(), http.StatusInternalServerError)
		return
	}

	pingCount := prevPingCount + 1

	err = ioutil.WriteFile(pingFile, []byte(strconv.Itoa(pingCount)), 0777)
	if err != nil {
		http.Error(w, "Internal error 3: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("ping %d\n", pingCount)))
}

func main() {
	_, err := os.Stat(pingFile)
	if errors.Is(err, os.ErrNotExist) {
		err := ioutil.WriteFile(pingFile, []byte("0"), 0777)
		if err != nil {
			panic(err)
		}
	}

	http.HandleFunc("/ping", handler)

	fmt.Println("Ping server now serving")
	http.ListenAndServe(":8123", nil)
}
