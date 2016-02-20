package main

import (
	"encoding/json"
	"fmt"
	"github.com/tonnerre/go-ldap"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

const MAX_FAILURE = 4

var LDAP_SERVER = os.Getenv("LDAP_SERVER")

var failedAddr = struct {
	sync.Mutex
	addrs map[string]int
}{addrs: make(map[string]int)}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	laddr := ":" + port
	http.HandleFunc("/authenticate", pwCheck)
	log.Println("Waiting for requests on:", laddr, "using LDAP server:", LDAP_SERVER)
	err := http.ListenAndServe(laddr, nil) // if balancer terminates ssl
	// err := http.ListenAndServeTLS(":8080", certFile, keyFile, handler) if app terminates ssl
	if err != nil {
		log.Fatalln("Error starting http server:", err)
	}
}

type Request struct {
	Username string
	Password string
}

func pwCheck(w http.ResponseWriter, r *http.Request) {
	failedAddr.Lock() // 1 PW check at a time
	defer failedAddr.Unlock()
	if remote, _, err := net.SplitHostPort(r.RemoteAddr); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Failed to parse host and port from:", remote)
	} else if failedAddr.addrs[remote] > MAX_FAILURE {
		// block blacklisted address
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Attempt to logic from blacklisted ip:", remote)
	} else {
		defer r.Body.Close() // don't care if close gets an error

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // bad json
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("Failed to parse request body:", err)
		} else if req.Password == "" || req.Username == "" { // missing fields
			http.Error(w, "Required { username:'', password: '' }", http.StatusBadRequest)
			log.Println("Received invalid request from:", remote)
		} else if con, err := ldap.Dial("tcp", LDAP_SERVER); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("Failed to connect to ldap server:", LDAP_SERVER, "reason:", err)
		} else {
			defer con.Close() // don't care if close gets an error
			if err := con.Bind(req.Username, req.Password); err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				failedAddr.addrs[remote]++
				log.Println("Failed attempt:", failedAddr.addrs[remote], "for:", req.Username, "from:", remote)
			} else {
				fmt.Fprintln(w, "Authenticated")
				delete(failedAddr.addrs, remote) // resets counter on successful login
				log.Println("Authenticated:", req.Username, "from:", remote)
			}
		}
	}
}
