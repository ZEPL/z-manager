// Package main provides main processes.
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"net/http"

	"github.com/getlantern/autoupdate"
	"github.com/getlantern/golog"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/staticbin"
)

const (
	internalVersion      = "0.0.3"
	versionPrintInterval = time.Second * 70
	updateCheckInterval  = time.Second * 80
	assetFsPath          = "../web/ui/public"
)

var log = golog.LoggerFor("z-manager")

//for load-balancing and sticky sessions
var userHosts = getUserHosts("DOCKER_HOSTS")
var serveAssetsFromFs = os.Getenv("SERVE_WEBAPP_FROM_FS")

var loginService = HubLoginService{proxy: newHubProxy()}

// main executes a webserver.
func main() {
	var m = martini.Classic()
	m.Use(render.Renderer())
	if serveAssetsFromFs != "" {
		log.Debugf("Dev mode: serving webapp from %s", assetFsPath)
		m.Use(martini.Static(assetFsPath))
		m.Use(cors.Allow(&cors.Options{
			AllowOrigins:     []string{"http://localhost:9090"},
			AllowMethods:     []string{"PUT", "GET", "POST"},
			AllowHeaders:     []string{"Origin", "Content-Type"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))
	} else {
		log.Debugf("Production mode: serving webapp from the binary")
		m.Use(staticbin.Static(assetFsPath, Asset))
	}

	//reverse proxy to Zeppelin-container
	zeppelin := newZeppelinProxy()
	m.Get("/zeppelin**", func(w http.ResponseWriter, r *http.Request) {
		zeppelin.ServeHTTP(w, r)
	})
	m.Post("/zeppelin**", func(w http.ResponseWriter, r *http.Request) {
		zeppelin.ServeHTTP(w, r)
	})
	m.Delete("/zeppelin**", func(w http.ResponseWriter, r *http.Request) {
		zeppelin.ServeHTTP(w, r)
	})

	m.Get("/assets/**", func(w http.ResponseWriter, r *http.Request) {
		zeppelin.ServeHTTP(w, r)
	})
	m.Get("/styles/**", func(w http.ResponseWriter, r *http.Request) {
		zeppelin.ServeHTTP(w, r)
	})
	m.Get("/scripts/**", func(w http.ResponseWriter, r *http.Request) {
		zeppelin.ServeHTTP(w, r)
	})
	m.Get("/components/**", func(w http.ResponseWriter, r *http.Request) {
		zeppelin.ServeHTTP(w, r)
	})
	m.Get("/app/**", func(w http.ResponseWriter, r *http.Request) {
		zeppelin.ServeHTTP(w, r)
	})
	m.Get("/fonts/**", func(w http.ResponseWriter, r *http.Request) {
		zeppelin.ServeHTTP(w, r)
	})

	//Websocket proxy
	websocketRP := newWebsocketProxy()
	m.Get("/ws**", func(w http.ResponseWriter, r *http.Request) {
		websocketRP.ServeHTTP(w, r)
	})

	//Cluster API - reverse proxy for Spark Master
	sparkMaster := newSparkMasterProxy()
	m.Get("/api/v1/cluster", func(w http.ResponseWriter, r *http.Request) {
		sparkMaster.ServeHTTP(w, r)
	})

	//Auth API
	m.Group("/api/v1/users", func(r martini.Router) {
		r.Put("/login", func(w http.ResponseWriter, r *http.Request) {
			loginService.Login(w, r)
		})

		r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
			loginService.Logout(w, r)
		})

		r.Get("/whoiam", func(w http.ResponseWriter, r *http.Request) {
			loginService.Whoami(w, r)
		})
	})

	//Containers API
	m.Map(&Docker{})
	m.Group("/api/v1/containers", func(r martini.Router) {
		r.Post("/list", binding.Bind(ListContainersReq{}), ListContainers)

		//IN: 'user-memory-cpu'
		r.Post("/create", binding.Bind(CreateContainerReq{}), CreateContainer)

		//IN: containerId
		r.Post("/delete", binding.Bind(DeleteContainerReq{}), DeleteContainer)

		r.Get("/images", func(r render.Render, docker *Docker) {
			r.JSON(200, docker.ListImages())
		})
	})

	m.Get("/**", func(w http.ResponseWriter, r *http.Request) {
		zeppelin.ServeHTTP(w, r)
	})

	go startAutoupdate()
	go func() { //debug output
		for {
			log.Debugf("Running program version: %v", internalVersion)
			time.Sleep(versionPrintInterval)
		}
	}()
	m.Run()
}

func ListContainers(containers ListContainersReq, r render.Render, docker *Docker) {
	containers.Host = userHosts[containers.Username]
	r.JSON(200, docker.List(containers)) //{ containers: [..,..]  }
}

func CreateContainer(container CreateContainerReq, r render.Render, docker *Docker) {
	container.Host = userHosts[container.Username]
	container.Port = getRemoteFreePort(container.Host)

	//TODO(alex): move to Container struct: .ToString(c CreateContainerReq)
	containerName := strings.Join([]string{container.Username, container.Cores, container.Memory, container.Port}, "-")

	portList := getRemoteFreePorts(sparkPorts, container.Host)
	envVars := getEnvVars(container.Username, container.Host, portList)
	volumes := setVolumes(container.Username)

	err := replaceInterpVars(volumes[volumeZeppelinConfig]+"/interpreter.json",
		container.Cores, container.Memory, containerName, portList)
	if err != nil {
		log.Debugf("Couldn't parse interpreter.json : " + fmt.Sprintln(err))
	}

	createdContainer, err := docker.Create(container, portList, envVars, volumes, containerName)
	if err != nil {
		r.JSON(404, createdContainer)
		return
	}
	r.JSON(200, createdContainer)
}

func DeleteContainer(container DeleteContainerReq, res http.ResponseWriter, docker *Docker) {
	container.Host = userHosts[container.Username]
	err := docker.Delete(container)
	if err != nil {
		res.WriteHeader(404)
		return
	}
	res.WriteHeader(200)
}

func removeCORS(w http.ResponseWriter) {
	w.Header().Del("Access-Control-Allow-Origin")
	w.Header().Del("Access-Control-Allow-Methods")
	w.Header().Del("Access-Control-Allow-Credentials")
}

// reads given environment variable and parses it to the
//  Map: username -> hostname
func getUserHosts(envVar string) map[string]string {
	result := map[string]string{}
	hosts := strings.Split(os.Getenv(envVar), "\n")

	log.Debugf("%s: %v hosts", envVar, len(hosts))
	for _, host := range hosts {
		hostUsers := strings.Split(host, ":")
		if len(hostUsers) < 2 {
			continue
		}
		host := strings.Trim(hostUsers[0], " ")
		users := strings.Split(hostUsers[1], ",")

		if len(users) == 1 && strings.Trim(users[0], " ") == "" { // filter empty
			users = append(users[:0], users[1:]...) // slices don't have 'delete'
		}
		log.Debugf("Host: %s, %v users", host, len(users))
		for _, user := range users {
			username := strings.Trim(user, " ")
			log.Debugf("\tuser: %s", username)
			result[username] = host
		}
	}
	return result
}

func startAutoupdate() {
	err := autoupdate.ApplyNext(&autoupdate.Config{
		CurrentVersion: internalVersion,
		URL:            "http://staging.nflabs.com:6868/update",
		CheckInterval:  updateCheckInterval,
		PublicKey:      []byte("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzSoibtACnqcp2uTGjCMJ\ntTOLDIMQ4oGPhGHT4Q/epum+H3hcbBNs9jRnMRWgX4z++xxuNJnhmoJw0eUXB7B4\nvj5DYpPajq6gPY8JuraF4ngfP5oxKj2BqpEUR9bx+3SjOSInrirM0JZO+aAW38BQ\nNJB+sS7JvbPjcwdjwKc5IKzc9kxxJNoZoFE9GMnYzaOrAlpCuAKWH8SCXYtCTxsX\nfKexdDxsI5Vzm5lQHJLMeqhLTQTUm9oQofwNAOGOkn6dD4ObMlmFTOsf1G03/Dl9\nsVgjaWaZ9bGjvJ9B85UxNeWwduy+uMrqFytxG6bbq0PbDEVu6ZQCPyiyCA7l945J\nOQIDAQAB\n-----END PUBLIC KEY-----\n"),
	})
	if err != nil {
		log.Debugf("Error getting update: %v", err)
		return
	}
	log.Debugf("Got update, please restart z-manager")
}
