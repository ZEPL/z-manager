package main

//Right now we support 2 types of auth services:
// - using ZeppelinHub
// - using 'basic auth'

import (
	"net/http"
	"net/http/httputil"
)

type LoginService interface {
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Whoami(w http.ResponseWriter, r *http.Request)
}

//Auth API impl using reverse proxy to Hub
type HubLoginService struct {
	proxy *httputil.ReverseProxy
}

func (hub *HubLoginService) Login(w http.ResponseWriter, r *http.Request) {
	hub.forwardTo(w, r)
}

func (hub *HubLoginService) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Set-Cookie",
		"user_session=deleteMe; Path=/; Max-Age=0; Expires=Tue, 04-Aug-2015 03:29:18 GMT")
	hub.forwardTo(w, r)
}

func (hub *HubLoginService) Whoami(w http.ResponseWriter, r *http.Request) {
	hub.forwardTo(w, r)
}

func (hub *HubLoginService) forwardTo(w http.ResponseWriter, r *http.Request) {
	removeCORS(w) //Hub and us both set CORS headers on /api so we delete ours
	hub.proxy.ServeHTTP(w, r)
}

//Auth API impl using cookie-stored session
type BasicLoginService struct{}

func (hub *BasicLoginService) Login(w http.ResponseWriter, r *http.Request) {
	sessionKey := "user_session"
	sessionVal := ""
	cookies := r.Cookies()
	// assume a req wo/ cookies.user_session
	for _, cookie := range cookies {
		if cookie.Name == sessionKey {
			sessionVal = cookie.Value
		}
	}
	if sessionVal == "" {
		return
	}

	// read user.name and user.pass from request payload
	// if db does not have a user.name
	//    return 'Username and password do not match'
	// read username_pass_hash from db
	// if db.username_pass_hash == hash(user.name + user.pass)
	//    cookie.user_session = hash(user.name + user.pass)
	//    return true
}

func (hub *BasicLoginService) Logout(w http.ResponseWriter, r *http.Request) {
	// delete cookies.user_session
}

func (hub *BasicLoginService) Whoami(w http.ResponseWriter, r *http.Request) {
	// get username_pass_hash from cookies.user_session
	// read the db
	// if username_pass_hash == hash(db.username + db.hash)
	//   return User
	// else
	//   return guest
}
