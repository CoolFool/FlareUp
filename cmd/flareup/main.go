package main

import (
	"flareup/internal/cloudflare"
	"flareup/internal/logging"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/net/publicsuffix"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var username string
var password string
var users = map[string]string{username: password}
var domains []string
var proxied = false
var log = logging.Init()

func handleRequest(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(`{"message": "No basic auth present"}`))
		if err != nil {
			log.Error("Error while writing to response")
		}
		return
	}
	if !isAuthorized(username, password) {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(`{"message": "Invalid username or password"}`))
		if err != nil {
			log.Error("Error while writing to response")
		}
		return
	}
	ip := getIP(r)
	go updateEntry(ip)
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, `{"Ok":"%s"}`, ip)
	if err != nil {
		log.Error("Error while writing to response")
	}
	return
}

func updateEntry(ip string) {
	for _, domain := range domains {
		TLD, err := publicsuffix.EffectiveTLDPlusOne(domain)
		if err != nil {
			log.Error(err)
		} else {
			TLD = "." + TLD
			hostname := strings.ReplaceAll(domain, TLD, "")
			TLD = strings.TrimPrefix(TLD, ".")
			if TLD == hostname {
				hostname = ""
			}
			go cloudflare.UpdateRecord(hostname, TLD, ip, proxied)
		}
	}
}

func getIP(r *http.Request) (ip string) {
	forwarded := r.Header.Get("X-FORWARDED-For")
	if forwarded != "" {
		ip = forwarded
	} else {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Error("Error while splitting remote address")
		}
		return ip
	}
	return
}

func isAuthorized(username, password string) bool {
	pass, ok := users[username]
	if !ok {
		return false
	}
	return password == pass
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Error(err)
	}
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")
	domain := os.Getenv("DOMAINS")
	port := os.Getenv("PORT")
	isProxy := os.Getenv("PROXIED")
	if username == "" || password == "" || domain == "" {
		log.Fatal("Environment variables not set")
	}
	if isProxy != "" {
		proxied, _ = strconv.ParseBool(isProxy)
	}
	if port == "" {
		port = ":5335"
	} else {
		port = ":" + port
	}
	users[username] = password
	domains = strings.Split(domain, ",")
	for i, j := range domains {
		domains[i] = strings.TrimSpace(j)
	}

	http.HandleFunc("/update", handleRequest)
	log.Info(fmt.Sprintf("flareup and Running on 0.0.0.0%s", port))
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Cant Start Server")
	}
}
