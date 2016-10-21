package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <stdint.h>
#include <unistd.h>
#include <sys/syscall.h>
#include <time.h>
#include <string.h>

char sessionid[256];

static char* generateSessionId()
{
	struct timespec tspec;
	memset(&tspec, 0, sizeof(tspec));
	// get monotonic time
	syscall(SYS_clock_gettime, CLOCK_MONOTONIC_RAW, &tspec);
	int sec = tspec.tv_sec;
	int msec = tspec.tv_nsec;
	sprintf(sessionid, "%x%x", sec,msec);
	return sessionid;
}
*/
import "C"
import "encoding/json"
import "fmt"
import "github.com/bitly/go-simplejson"
import "github.com/gorilla/mux"
import "io/ioutil"
import "log"
import "net/http"
import "strings"

//import "unsafe"
//import "strconv"
type HttpHeaderContent struct {
	ContentType  string
	ContentValue string
}

var HTTPHeader = HttpHeaderContent{"Content-Type", "application/json"}

// session num
const MAX_SESSION_NUM int = 1
const USERNAME string = "admin"
const USERPWD string = "admin"

// session hashmap
var smap = map[string]struct {
	session string
}{}

// resp result
type Result struct {
	Status bool `json:"status"`
}

func PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	// resp json header
	w.Header().Set(HTTPHeader.ContentType, HTTPHeader.ContentValue)
	var res = Result{false}
	if len(smap) < MAX_SESSION_NUM {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic("error")
		}
		js, err := simplejson.NewJson(body)
		// get username and userpwd
		username, _ := js.Get("username").String()
		userpwd, _ := js.Get("userpwd").String()
		if strings.Compare(username, USERNAME) == 0 &&
			strings.Compare(userpwd, USERPWD) == 0 {

			// generate session
			sessionid := C.GoString(C.generateSessionId())

			// put map
			var tmp = smap[sessionid]
			tmp.session = sessionid
			smap[sessionid] = tmp

			// set cookie
			cookie := http.Cookie{
				Name:  "session",
				Value: sessionid,
			}
			http.SetCookie(w, &cookie)
			// resp ok
			res = Result{true}
		}
	}
	j, _ := json.Marshal(res)
	w.Write(j)
}

func main() {
	fmt.Println("Starting webserver!")

	r := mux.NewRouter()

	r.HandleFunc("/login", PostLoginHandler).Methods("POST")

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("dist/"))))

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:80",
		// Good practice: enforce timeouts for servers you create!
		//WriteTimeout: 15 * time.Second,
		//ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
