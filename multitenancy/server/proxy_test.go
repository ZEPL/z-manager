package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func setupBackend(t *testing.T, response string, status int) (*httptest.Server, string) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Came-From", r.URL.String())
		w.WriteHeader(status)
		w.Write([]byte(response))
	}))
	backendURL, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	return backend, backendURL.String()
}

func TestReverseProxyNoUrlRewrite(t *testing.T) {
	const backendResponse = "I am the backend"
	const backendStatus = 200
	backend, backendURL := setupBackend(t, backendResponse, backendStatus)
	defer backend.Close()

	//given
	HubURL = backendURL //this configres proxy dst below
	zeppelinHub := newHubProxy()

	frontend := httptest.NewServer(zeppelinHub)
	defer frontend.Close()

	getReqPath := "/api/v1/users/login"
	getReqURL := frontend.URL + getReqPath
	getReq, _ := http.NewRequest("GET", getReqURL, nil)
	getReq.Host = "some-name"
	getReq.Close = true

	//when
	res, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	//then
	if g, e := res.StatusCode, backendStatus; g != e {
		t.Errorf("got res.StatusCode %d; expected %d", g, e)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if g, e := string(bodyBytes), backendResponse; g != e {
		t.Errorf("got body %q; expected %q", g, e)
	}

	fmt.Printf("\tZeppelinHubProxy: re-writing from %v to %v",
		getReqPath, res.Header.Get("X-Came-From"))
	if g, e := res.Header.Get("X-Came-From"), getReqPath; g != e {
		t.Errorf(`got X-Came-From %q; expected %q
             as zeppelinHubProxy does not re-write URLs`, g, e)
	}
	fmt.Printf("\t OK!\n")
}

func TestReverseProxyWithUrlRewrite(t *testing.T) {
	const backendResponse = "I am the backend"
	const backendStatus = 200
	backend, backendURL := setupBackend(t, backendResponse, backendStatus)
	defer backend.Close()

	//given
	SparkURL = backendURL //this configres proxy dst below
	spark := newSparkMasterProxy()

	frontend := httptest.NewServer(spark)
	defer frontend.Close()

	getReqPath := "/api/v1/cluster"
	getReqURL := frontend.URL + getReqPath
	getReq, _ := http.NewRequest("GET", getReqURL, nil)
	getReq.Host = "some-name"
	getReq.Close = true

	//when
	res, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	//then
	if g, e := res.StatusCode, backendStatus; g != e {
		t.Errorf("got res.StatusCode %d; expected %d", g, e)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if g, e := string(bodyBytes), backendResponse; g != e {
		t.Errorf("got body %q; expected %q", g, e)
	}

	fmt.Printf("\tSparkMasterProxy: re-writing from %v to %v",
		getReqPath, res.Header.Get("X-Came-From"))
	if g, e := res.Header.Get("X-Came-From"), "/json/"; g != e {
		t.Errorf(`got X-Came-From %q; expected %q 
             as that's what sparkMasterProxy rewrites to`, g, e)
	}
	fmt.Printf("\t OK!\n")
}

// Runs backend echo server on given port
func runBackend(t *testing.T, port string) *http.ServeMux {
	mux := http.NewServeMux()
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true //allow all
		},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		s := port + " " + string(p[:])
		if err = conn.WriteMessage(messageType, []byte(s)); err != nil {
			return
		}
	})

	go func() {
		err := http.ListenAndServe("localhost:"+port, mux)
		if err != nil {
			t.Fatal("ListenAndServe: ", err)
		}
	}()

	return mux
}

func TestReverseProxyWebsocket(t *testing.T) {
	//given
	userA := "userA"
	userB := "userB"
	userAport := "9999"
	userBport := "10000"
	userHosts = map[string]string{userA: "localhost", userB: "localhost"}

	frontend := httptest.NewServer(newWebsocketProxy())
	defer frontend.Close()

	time.Sleep(time.Millisecond * 100)
	//different backends for users
	/*backendA := */ runBackend(t, userAport)
	/*backendB := */ runBackend(t, userBport)

	//when
	headerA := http.Header{"Cookie": {"username=" + userA, "port=" + userAport}}
	connA, respA, err := websocket.DefaultDialer.Dial("ws://"+frontend.URL[7:]+"/ws", headerA)
	if err != nil {
		t.Error(respA)
		t.Fatal(err)
	}
	connA.WriteMessage(websocket.TextMessage, []byte("somethingA"))

	headerB := http.Header{"Cookie": {"username=" + userB, "port=" + userBport}}
	connB, _, err := websocket.DefaultDialer.Dial("ws://"+frontend.URL[7:]+"/ws", headerB)
	if err != nil {
		t.Fatal(err)
	}
	connB.WriteMessage(websocket.TextMessage, []byte("somethingB"))

	//then
	//connA hits backendA
	connectionHitsPort(t, connA, userAport)
	//connB hits backendB
	connectionHitsPort(t, connB, userBport)
	fmt.Printf("\t OK!\n")
}

func connectionHitsPort(t *testing.T, conn *websocket.Conn, expectedPort string) {
	messageType, p, err := conn.ReadMessage()
	if err != nil {
		t.Error(err)
	}

	if messageType != websocket.TextMessage {
		t.Error("incoming message type is not Text")
	}

	if strings.HasPrefix(string(p), expectedPort) {
		t.Errorf("expecting port: %s, got: %s", expectedPort, string(p))
	}
}
