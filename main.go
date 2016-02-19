package main

import (
	"encoding/json"
	"fmt"
	"github.com/tonnerre/go-ldap"
	"log"
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
	http.HandleFunc("/authenticate", pwCheck)
	http.ListenAndServe(":8080", nil) // if balancer terminates ssl
	// http.ListenAndServeTLS(":8080", certFile, keyFile, handler) if app terminates ssl
}

type Request struct {
	Username string
	Password string
}

func pwCheck(w http.ResponseWriter, r *http.Request) {
	failedAddr.Lock() // 1 PW check at a time
	defer failedAddr.Unlock()
	if failedAddr.addrs[r.RemoteAddr] > MAX_FAILURE {
		// block blacklisted address
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else {
		defer r.Body.Close() // don't care if close gets an error
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // bad json
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else if req.Password == "" || req.Username == "" { // missing fields
			http.Error(w, "Required { username:'', password: '' }", http.StatusBadRequest)
		} else if con, err := ldap.Dial("tcp", LDAP_SERVER); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else {
			defer con.Close() // don't care if close gets an error
			if err := con.Bind(req.Username, req.Password); err != nil {
				failedAddr.addrs[r.RemoteAddr]++
				log.Println("Failed attempt:", failedAddr.addrs[r.RemoteAddr], "for:", req.Username, "from:", r.RemoteAddr)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				fmt.Fprintln(w, "Authenticated")
			}
		}
	}
}
