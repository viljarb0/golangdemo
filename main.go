package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/viljarb0/golangdemo/somemodule"
	"golang.org/x/crypto/bcrypt"
)

func htmlDoc(x int) string {
	var strr string = "11"
	for ii := 1; ii <= x; ii++ {
		strr += strconv.Itoa(ii)
	}
	return strr
}

type user struct {
	name string
	age  int
}

func testFunc() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	funcid := 5
	time.Sleep(time.Duration(funcid) * time.Second)
	fmt.Printf("done: %d\n", funcid)
}

func serveHTML(w http.ResponseWriter, filename string) {
	filepath := fmt.Sprintf("templates/%s", filename)
	content, err := os.ReadFile(filepath)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

func main() {
	u1 := user{
		name: "myname",
		age:  33,
	}
	var addr string = "127.0.0.1:8080"
	var password string = fmt.Sprintf("%s%s123", u1.name, strconv.Itoa(u1.age))
	fmt.Println(password)

	somemodule.Secondfunc()
	// for _ := 0;; {
	// nil
	// }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(hashedPassword))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHTML(w, "index.html")
	})
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		serveHTML(w, "about.html")
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		serveHTML(w, "register.html")
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		serveHTML(w, "login.html")
	})

	conn, err := InitDB("db.sqlite3")
	if err != nil {
		log.Fatal("error")
	}
	AddUser(conn, u1.name, "u1@email.com", password)
	log.Printf("listening on http://%s\n", addr)
	result, err := Login(conn, u1.name, password)
	if err != nil {
		log.Fatal(err)
	}
	if result {
		fmt.Println("login successful")
	} else {
		fmt.Println("login failed")
	}

	log.Fatal(http.ListenAndServe(addr, mux))
}
