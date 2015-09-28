package main

import (
	"fmt"
	"os"
	//	"os/exec"
	"strconv"
	"strings"

	"github.com/fsouza/go-dockerclient"
	//	"github.com/bitly/go-simplejson"
)

var ImageName = os.Getenv("DOCKER_IMAGE")
var DockerPort = os.Getenv("DOCKER_PORT")

//TODO(alex): add type Container struct {}

type ListContainersReq struct {
	Username string `json:"username" binding:"required"`
	Host     string //below is not part of the JSON input
}

type DeleteContainerReq struct {
	Id       string `json:"containerId" binding:"required"`
	Username string `json:"username" binding:"required"`
	Host     string //below is not part of the JSON input
}

type CreateContainerReq struct {
	Cores    string `json:"cores" binding:"required"`
	Memory   string `json:"memory" binding:"required"`
	Username string `json:"username" binding:"required"`
	Host     string //below is not part of the JSON input
	Port     string
}

type DockerManager interface {
	Create() CreateContainerReq
	Delete() DeleteContainerReq
	List() []map[string]string
}

//TODO: impelment a cache: user -> [containrs]
type Docker map[string]string

//creates+starts a container for cont.Username
func (contCache Docker) Create(cont CreateContainerReq,
	portList []string,
	envVars map[string]string,
	folderEnvVars map[string]string,
	containerName string) (map[string]string, error) {

	spark_driver_port_tcp := docker.Port(portList[0] + "/tcp")
	spark_fileserver_port_tcp := docker.Port(portList[1] + "/tcp")
	spark_broadcast_port_tcp := docker.Port(portList[2] + "/tcp")
	spark_replClassServer_port_tcp := docker.Port(portList[3] + "/tcp")
	//workaround for https://github.com/NFLabs/zeppelin-manager/issues/114
	portList3, _ := strconv.Atoi(portList[3])
	spark_replClassServe2_port_tcp := docker.Port(strconv.Itoa(portList3+1) + "/tcp")
	spark_blockManager_port_tcp := docker.Port(portList[4] + "/tcp")
	spark_ui_port_tcp := docker.Port(portList[5] + "/tcp")

	exposedPorts := map[docker.Port]struct{}{"8080/tcp": {},
		spark_driver_port_tcp:          {},
		spark_fileserver_port_tcp:      {},
		spark_broadcast_port_tcp:       {},
		spark_replClassServer_port_tcp: {},
		spark_replClassServe2_port_tcp: {},
		spark_blockManager_port_tcp:    {},
		spark_ui_port_tcp:              {},
	}
	//log.Debugf("Exposed Ports:  " + fmt.Sprintln(exposedPorts))

	portBindings := map[docker.Port][]docker.PortBinding{
		"8080/tcp":                     {{HostIP: "0.0.0.0", HostPort: cont.Port}},
		spark_driver_port_tcp:          {{HostIP: "0.0.0.0", HostPort: portList[0]}},
		spark_fileserver_port_tcp:      {{HostIP: "0.0.0.0", HostPort: portList[1]}},
		spark_broadcast_port_tcp:       {{HostIP: "0.0.0.0", HostPort: portList[2]}},
		spark_replClassServer_port_tcp: {{HostIP: "0.0.0.0", HostPort: portList[3]}},
		spark_replClassServe2_port_tcp: {{HostIP: "0.0.0.0", HostPort: strconv.Itoa(portList3 + 1)}},
		spark_blockManager_port_tcp:    {{HostIP: "0.0.0.0", HostPort: portList[4]}},
		spark_ui_port_tcp:              {{HostIP: "0.0.0.0", HostPort: portList[5]}},
	}
	//log.Debugf("Port Bindings:  " + fmt.Sprintln(portBindings))

	//mounts FS: expose volumes + bind them
	volumes := map[string]struct{}{"/usr/lib/zeppelin/conf": {}, "/usr/lib/spark": {}, "/zeppelin/notebook": {}}

	client, _ := docker.NewClientFromEnvWithHost(cont.Host, DockerPort)
	log.Debugf("User: %s, start container %v, from image %v on %v", cont.Username, containerName, ImageName, client.Endpoint)

	//try pulling the image first, will take a while to download ~1Gb.....
	er := client.PullImage(docker.PullImageOptions{Repository: ImageName},
		docker.AuthConfiguration{})
	if er != nil {
		log.Errorf("Can not pull the image %v, %v", ImageName, er)
	}

	//create
	hostConfig := &docker.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/usr/lib/zeppelin/conf", folderEnvVars[volumeZeppelinConfig]),
			fmt.Sprintf("%s:/zeppelin/notebook", folderEnvVars[volumeZeppelinNotebooks]),
			"/usr/lib/spark:/usr/lib/spark"},
		PortBindings:    portBindings,
		PublishAllPorts: true,
		ExtraHosts: []string{"bhsmaster1:211.43.191.6", "bhsmaster2:211.43.191.7",
			"bhsmaster3:211.43.191.8"},
	}

	envDockerZeppelinConfDir := fmt.Sprintf("ZEPPELIN_CONF_DIR=%s", "/usr/lib/zeppelin/conf")
	envDockerZeppelinNoteDir := fmt.Sprintf("ZEPPELIN_NOTEBOOK_DIR=%s", "/zeppelin/notebook")

	container, err := client.CreateContainer(
		docker.CreateContainerOptions{
			Name: containerName,
			Config: &docker.Config{
				AttachStdout: true,
				AttachStdin:  false,
				Cmd:          []string{"/bin/bash", "-c", "/usr/lib/zeppelin/bin/zeppelin.sh"},
				Env: []string{envVars[zeppelinConfDir],
					envVars[sparkPublicDNS],
					envVars[sparkLocalHostname],
					envDockerZeppelinConfDir,
					envDockerZeppelinNoteDir,
					envVars[zeppelinIntpJavaOpts]},
				Hostname:     cont.Host,
				Image:        ImageName,
				ExposedPorts: exposedPorts,
				Volumes:      volumes,
			},
			HostConfig: hostConfig,
		},
	)

	if err != nil {
		log.Errorf("Can not create a container %v, %v", containerName, err)
		return map[string]string{
			"error": "Can not create a container " + containerName + ". Reason: " + err.Error(),
		}, err
	}

	err = client.StartContainer(container.ID, hostConfig)
	if err != nil {
		log.Errorf("Can not start a container %v, %v", containerName, err)
		//TODO(bzz): docker rm containerName
		return map[string]string{
			"error": "Can not start a container " + containerName + ". Reason: " + err.Error(),
		}, err
	}

	return map[string]string{
		"containerId": container.ID,
		"port":        cont.Port,
	}, nil
}

//kill+rm contianer.ID for container.Username
func (contCache Docker) Delete(cont DeleteContainerReq) error {
	client, _ := docker.NewClientFromEnvWithHost(cont.Host, DockerPort)

	log.Debugf("User: %s, kill+rm container %v on %v", cont.Username, cont.Id, client.Endpoint)
	err := client.KillContainer(docker.KillContainerOptions{ID: cont.Id})
	if err != nil {
		log.Errorf("Can not kill a continer, %v", err)
		return err
	}

	err = client.RemoveContainer(docker.RemoveContainerOptions{ID: cont.Id})
	if err != nil {
		log.Errorf("Can not remove a continer, %v", err)
		return err
	}
	return nil
}

//lists all contianers for containerz.Username
func (contCache Docker) List(containerz ListContainersReq) map[string][]map[string]string {
	client, _ := docker.NewClientFromEnvWithHost(containerz.Host, DockerPort)

	log.Debugf("User: %s, listing containers on %v", containerz.Username, client.Endpoint)
	containers, _ := client.ListContainers(docker.ListContainersOptions{})

	result := []map[string]string{}
	for _, img := range containers {
		log.Tracef("%+v", img)

		//TODO(alex): move to Container struct: .ToString(c CreateContainerReq)
		name := img.Names[0]
		args := strings.Split(strings.TrimLeft(name, "/"), "-")
		if len(args) != 4 { //validate container name: user-memory-cores-port-contId
			continue
		}
		c := map[string]string{
			"username":    args[0],
			"cores":       args[1],
			"memory":      args[2],
			"port":        args[3],
			"containerId": img.ID,
		}
		result = append(result, c)
	}

	return toJSON(result)
}

func (contCache Docker) ListImages() map[string][]map[string]string {
	//this requeires at least DOCKER_HOST set
	client, _ := docker.NewClientFromEnv()
	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
	result := map[string][]map[string]string{} //{images: [...]}
	for _, img := range imgs {
		result["images"] = append(
			result["images"],
			map[string]string{
				"id":   img.ID,
				"tag":  img.RepoTags[0],
				"Size": strconv.FormatInt(img.Size, 10),
			})
	}
	return result
}

//Re-shapes given Datastructure to JSON-compatibe map
func toJSON(containersJson []map[string]string) map[string][]map[string]string {
	result := map[string][]map[string]string{
		"containers": containersJson,
	}
	return result
}

//Raw curl json call to create a container
func curlCreateContainer() {
	/*
		     //hostZeppelinConfDir := fmt.Sprintf("/data/users/%s/conf/", cont.Username)
		   	 //hostZeppelinNoteDir := fmt.Sprintf("/data/users/%s/notebooks/", cont.Username

					command := "curl"
					arg1 := "-H"
					contentType := "Content-Type: application/json"
					arg2 := "-X"
					post := "POST"
					arg3 := "-d"
					ccreateJson := fmt.Sprintf(`{   "AttachStdout": "%v",
																					"AttachStdin": "false",
																					"Cmd": [ "/bin/bash", "-c", "/usr/lib/zeppelin/bin/zeppelin.sh"],
																					"Image": "nflabs/zeppelin-bhs-aug-5-spark-1.4-hadoop-2.0.0-cdh-4.7.0:latest",
																					"ExposedPorts": { "8080/tcp": {} },
																					"Volumes": { "/usr/lib/zeppelin/conf": {}, "/usr/lib/spark": {}, "/zeppelin/notebook": {} },
																					"HostConfig": {
																							"Binds": [ "%s:/usr/lib/zeppelin/conf",
																											 "%s:/usr/lib/zeppelin/noteboks"],
																							"PortBindings": { "8080/tcp": [{ "HostIP": "0.0.0.0", "HostPort": %v }] },
																							"PublishAllPorts": "false"
																						 }
																			}`, "true", hostZeppelinConfDir, hostZeppelinNoteDir, cont.Port)
					 //log.Debugf(ccreateJson)

					//urlFull := fullHost + fmt.Sprintf("/containers/create?name=%v", containerName)
					url := "http://" + fullHost + ":DockerPort" + fmt.Sprintf("/containers/create?name=%v", containerName)
					cmd := exec.Command(command, arg1, contentType, arg2, post, arg3, ccreateJson, url)
					stdout, errExec := cmd.Output()
					 log.Debugf("URL IS: " + url)
					if errExec != nil {
						log.Errorf("Can not create a container (curl failed) %v, %v", containerName, errExec)
						return map[string]string{
							"error": "Can not create a container " + containerName + ". Reason: CURL" + errExec.Error(),
						}, errExec
					}
					createReplyJson, errJson := simplejson.NewJson(stdout)
					if errJson != nil {
						log.Debugf("Couldn't parse json reply: " + fmt.Sprintln(createReplyJson))
						return map[string]string{
							"error": "Can not create a container " + containerName + ". Reason: " + errJson.Error(),
						}, errJson
					}

					log.Debugf("Successfully got json from reply")
					containerMap, errMap := createReplyJson.Map()
					if errMap != nil {
						log.Errorf("Can not create a container %v, %v", containerName, errMap)
						return map[string]string{
							"error": "Can not create a container " + containerName + ". Reason: " + errMap.Error(),
						}, errJson
					}
					log.Debugf("CURL container create:  " + fmt.Sprintln(string(stdout)))
					log.Debugf("CURL container create error:  " + fmt.Sprintln(errJson))


					containerID := containerMap["Id"].(string)
					log.Debugf("\nContainer ID is :  " + containerID + "\n")
	*/
}
