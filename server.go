package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

var Count int = 0
var logger *log.Logger
var f *os.File
var countF *os.File
var message string = "Hello from Container Solutions." 
var version string
var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var ds Datastore

func main() {

	f, _ = os.OpenFile("/var/log/app.log", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
	countF, _ = os.OpenFile("/var/log/counter", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
	
	defer f.Close()
	defer countF.Close()

	version = os.Getenv("VERSION")
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	connection, useSQL := os.LookupEnv("SQL_DATASTORE_CONNECTION")
	if useSQL {
		ds = NewSQLDatastore()
		fmt.Printf("Initializing SQL datastore: %s", connection)
		err := ds.Init("postgres", connection)
		if err != nil {
			fmt.Printf("Cannot initialize SQL datastore: %v", err)
			os.Exit(1)
		}
	} else {
		ds = NewSliceDataStore()
		ds.Init(0)
	}

	http.HandleFunc("/", hello)
	http.HandleFunc("/readiness", readiness)
	http.HandleFunc("/liveness", liveness)
	http.HandleFunc("/counter", counter)
	http.HandleFunc("/avatar", avatarGen)
	http.HandleFunc("/mineBitcoin", mineBitcoin)
	http.HandleFunc("/store", storeHandler)
	
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		fmt.Println("Failed to start server: ", err)
	}

	fmt.Printf("Starting server..., listening on port %s\n", port)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received at %s", time.Now().String())
	f.WriteString("Request received at " + time.Now().String() + "\n")
	
	host, _ := os.Hostname()
	fmt.Fprint(w, message + "\nI'm running version " + version + " on ", host, "\n")
}

func readiness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Readiness Probe called at %s", time.Now().String())
	f.WriteString("Readiness Probe called at " + time.Now().String() + "\n")
	time.Sleep(500 * time.Millisecond)
	w.Write([]byte("I'm ready."))
}

func liveness(w http.ResponseWriter, r *http.Request) {
	f.WriteString("Liveness probe called at %s" + time.Now().String() + "\n")
	fmt.Println("Liveness probe called at %s", time.Now().String())
	time.Sleep(500 * time.Millisecond)
	w.Write([]byte("I am Alive!"))
}

func counter(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Counter called at %s", time.Now().String())
	f.WriteString("Counter called at %s" + time.Now().String() + "\n")
	Count = Count +1
	countF.WriteString("Count: " + strconv.Itoa(Count) + "\n")
	fmt.Fprintf(w,"Count: %d", Count)
}

func avatarGen(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Avatar called at %s", time.Now().String())
	f.WriteString("Avatar called at %s" + time.Now().String() + "\n")
	id := r.URL.Query().Get("id")
	// ??
	if len(id) < 1 {
		fmt.Printf("No id provided, will generate random avatar")
		id = randomString()
	} else {
		if len(id) < 3 {
			// Error
			fmt.Printf("Id must be at least 3 characters\n")
			w.WriteHeader(404)
			return
		}
	}

	img := FetchMonster(id, 96)
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		fmt.Println("Error encoding image.")
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buf.Bytes())))
	if _, err := w.Write(buf.Bytes()); err != nil {
		fmt.Println("Error sending image.")
	}
}

func mineBitcoin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get me some crypto currency.")
	seconds := r.URL.Query().Get("seconds")

	if len(seconds) < 1 {
		seconds = "30"
	}
	s, err := strconv.Atoi(seconds)
	if err != nil { panic(err) }

	fmt.Printf("Timeout set for %d seconds\n", s)

    f, err := os.Open(os.DevNull)
    if err != nil { panic(err) }
    defer f.Close()

    n := runtime.NumCPU()
    runtime.GOMAXPROCS(n)

	stopchan := make(chan struct{})
    for i := 0; i < n; i++ {
        go func() {
            for {
            	select {
      				default:
                		fmt.Fprintf(f, ".")
                	case <-stopchan:
        				// stop
        				return
    			}
            }
        }()
    }

    time.Sleep(time.Duration(s) * time.Second)
    close(stopchan)
	fmt.Println("timeout")
	w.Write([]byte("You are now rich."))
	return
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Store handler called...")
	status := http.StatusOK
	switch r.Method {

	case "GET":
		w.Header().Add("Content-type", "application/json")
		err := json.NewEncoder(w).Encode(ds.Get())
		if err != nil {
			status = http.StatusInternalServerError
			w.WriteHeader(status)
		}

	case "PUT":
		record := Record{}
		err := json.NewDecoder(r.Body).Decode(&record)
		if err != nil {
			status = http.StatusBadRequest
			w.WriteHeader(status)
			fmt.Fprintln(w, err)
		} else {
			ds.Add(record)
			status = http.StatusCreated
			w.WriteHeader(status)
			fmt.Fprintln(w, "OK")
		}

	case "DELETE":
		record := Record{}
		err := json.NewDecoder(r.Body).Decode(&record)
		if err != nil {
			status = http.StatusBadRequest
			w.WriteHeader(status)
			fmt.Fprintln(w, err)
		} else {
			ds.Rem(record)
			status = http.StatusOK
			w.WriteHeader(status)
			fmt.Fprintln(w, "OK")
		}

	default:
		status = http.StatusMethodNotAllowed
		w.WriteHeader(status)
	}

	fmt.Printf("%s - %s %d", r.Method, r.URL.Path, status)
}

func randomString() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}