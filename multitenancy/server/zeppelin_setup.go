package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/bitly/go-simplejson"
)

const (
	sparkPorts              = 6
	zeppelinConfDir         = "ZEPPELIN_CONF_DIR"
	sparkPublicDNS          = "SPARK_PUBLIC_DNS"
	sparkLocalHostname      = "SPARK_LOCAL_HOSTNAME"
	zeppelinIntpJavaOpts    = "ZEPPELIN_INTP_JAVA_OPTS"
	volumeZeppelinConfig    = "VOLUME_ZEPPELIN_CONFIG"
	volumeZeppelinNotebooks = "VOLUME_ZEPPELIN_NOTEBOOKS"
)

var DefaultUsersFolder = os.Getenv("USERS_FOLDER_PATH")
var DefaultUser = os.Getenv("DEFAULT_USER_FOLDER_NAME")

// returns slice of the given size with open ports to use
func getFreePorts(num int) []string {
	ports := make([]string, num)
	for i := 0; i < num; i++ {
		ports[i] = getLocalFreePort()
	}
	return ports
}

// returns a free port from system resource, tries to bind locally
func getLocalFreePort() string {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Debugf("Couldn't allocate free port")
	}
	port := ln.Addr().String()[5:]
	ln.Close()
	return port
}

// finds N free ports of the given host using find_open_port.py service
func getRemoteFreePorts(N int, host string) []string {
	ports := make([]string, N)
	for i := 0; i < N; i++ {
		ports[i] = getRemoteFreePort(host)
	}
	return ports
}

// finds single free port on the given host, using find_open_port.py service
func getRemoteFreePort(host string) string {
	if host == "localhost" || host == "127.0.0.1" || host == "0.0.0.0" {
		return getLocalFreePort()
	}
	url := getFQDN("http://"+host, "7777")

	response, errGet := http.Get(url)
	if errGet != nil {
		port := strconv.Itoa(rand.Intn(64510) + 1024)
		//TODO(bzz): try dailing the port to check if it's free one
		log.Errorf("Can not pick a free port on %v, using %v instead. Reason: %v", url, port, errGet)
		return port
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	return string(body)
}

func getEnvVars(username string, hostName string, portList []string) map[string]string {
	envVars := make(map[string]string)
	envVars[sparkPublicDNS] = fmt.Sprintf("SPARK_PUBLIC_DNS=%s", hostName)
	envVars[sparkLocalHostname] = fmt.Sprintf("SPARK_LOCAL_HOSTNAME=%s", hostName)
	envVars[zeppelinIntpJavaOpts] = fmt.Sprintf("ZEPPELIN_INTP_JAVA_OPTS=-Dspark.replClassServer.port=%s", portList[3])
	return envVars
}

type fileError struct {
	Message string
	Path    string
}

func (e fileError) Error() string {
	return fmt.Sprintf("%v Path: %v", e.Message, e.Path)
}

/* returns map with the key declared as constants above (starting with 'volume') e.g. volumeZeppelinConfig
   the value is the path e.g. .../username/conf */
func setVolumes(username string) map[string]string {
	defaultConfPath := filepath.Join(DefaultUsersFolder, DefaultUser, "conf")
	userConf := filepath.Join(DefaultUsersFolder, username, path.Base(defaultConfPath))
	defConfExists, err1 := exists(defaultConfPath)
	userExists, err2 := exists(userConf)

	if err1 != nil {
		log.Debugf("Error while reading file system : " + fmt.Sprintln(err1))
		return nil
	}
	if err2 != nil {
		log.Debugf("Error while reading file system : " + fmt.Sprintln(err2))
		return nil
	}
	if !defConfExists {
		log.Debugf("Default Zeppelin configuration folder under " + defaultConfPath + " doesn't exist. \n Please create one")
		return nil
	}
	if !userExists {
		log.Debugf("Creating folder for user under " + userConf)
		err3 := copyDir(defaultConfPath, userConf)
		if err3 != nil {
			log.Debugf("Couldn't copy conf dir from " + defaultConfPath + " to " + userConf)
			return nil
		}
	}

	userNotebooksPath := filepath.Join(DefaultUsersFolder, username, "notebooks")
	userNotebooksExist, err4 := exists(userNotebooksPath)
	if err4 != nil {
		log.Debugf("Error while reading file system : " + fmt.Sprintln(err4))
		return nil
	}

	if !userNotebooksExist { // create notebook folder inside username
		fileInfo, err5 := os.Stat(DefaultUsersFolder)
		if err5 != nil {
			log.Debugf("Error while reading file system : " + fmt.Sprintln(err5))
			return nil
		}
		err6 := os.MkdirAll(userNotebooksPath, fileInfo.Mode())
		if err6 != nil {
			log.Debugf("Error while creating directory : " + fmt.Sprintln(err5))
			return nil
		}
	}
	volumeVars := make(map[string]string)
	volumeVars[volumeZeppelinConfig] = userConf
	volumeVars[volumeZeppelinNotebooks] = userNotebooksPath
	return volumeVars
}

func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else { // permissions, etc
			return true, err
		}
	}
	return true, nil
}

func copyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()
	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()
	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}
	return nil
}

func copyDir(source string, dest string) (err error) {
	// get properties of source dir
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !sourceInfo.IsDir() {
		return &fileError{"Source is not a directory", source}
	}

	_, err = os.Open(dest)
	if !os.IsNotExist(err) { // check whether dest dir exists
		return &fileError{"Destination already exists", dest}
	}

	// create dest dir
	err = os.MkdirAll(dest, sourceInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(source)
	for _, entry := range entries {
		sourceFilePointer := source + "/" + entry.Name()
		destFilePointer := dest + "/" + entry.Name()
		if entry.IsDir() {
			err = copyDir(sourceFilePointer, destFilePointer)
			if err != nil {
				log.Debugf(fmt.Sprintln(err))
			}
		} else { // perform copy
			err = copyFile(sourceFilePointer, destFilePointer)
			if err != nil {
				log.Debugf(fmt.Sprintln(err))
			}
		}
	}
	return nil
}

type interpreterError struct {
	Message string
}

func (e interpreterError) Error() string {
	return fmt.Sprintf("%v", e.Message)
}

func replaceInterpVars(filePath, cores, memory, containerName string, ports []string) error {
	if len(ports) != sparkPorts {
		return &interpreterError{"Wrong number of Spark ports"}
	}

	fd, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	interpreterJson, err2 := simplejson.NewJson(fd)
	if err2 != nil {
		return err2
	}

	interpreterSettings, isOk := interpreterJson.CheckGet("interpreterSettings")
	if !isOk {
		return &interpreterError{"Couldn't load interpreterSettings"}
	}

	interpreterSettingsMap, err3 := interpreterSettings.Map()
	if err3 != nil {
		return err3
	}

	for _, item := range interpreterSettingsMap {
		innerMap, innerMapOk := item.(map[string]interface{})
		if innerMapOk {
			if val, ok := innerMap["name"]; ok {

				if val == "spark-cluster" {
					sparkProperties, propsOk := innerMap["properties"].(map[string]interface{})
					if propsOk {
						//log.Debugf(fmt.Sprintln("\n\nCORES:   " + sparkProperties["spark.cores.max"].(string)))
						sparkProperties["spark.cores.max"] = cores
						sparkProperties["spark.executor.memory"] = memory
						sparkProperties["spark.app.name"] = containerName

						sparkProperties["spark.driver.port"] = ports[0]
						sparkProperties["spark.fileserver.port"] = ports[1]
						sparkProperties["spark.broadcast.port"] = ports[2]
						sparkProperties["spark.replClassServer.port"] = ports[3]
						sparkProperties["spark.blockManager.port"] = ports[4]
						sparkProperties["spark.ui.port"] = ports[5]
					} else {
						return interpreterError{"sparkProperties isn't of Map type"}
					}
				}
			}
		}
	}

	outFormat, outErr := interpreterJson.EncodePretty()
	if outErr != nil {
		return outErr
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	writeErr := ioutil.WriteFile(filePath, outFormat, fileInfo.Mode())

	if writeErr != nil {
		return writeErr
	}
	return nil
}
