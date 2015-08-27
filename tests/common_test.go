package tests

import (
	"io/ioutil"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/client"
)

var apiClient *client.RancherClient
var apiClient2 *client.RancherClient

func TestMain(m *testing.M) {
	var err error

	apiUrl := "http://localhost:8080/v1/projects/1a5/schema"
	apiUrl2 := "http://localhost:8080/v1/projects/1a6/schema"
	accessKey := ""
	secretKey := ""

	apiClient, err = client.NewRancherClient(&client.ClientOpts{
		Url:       apiUrl,
		AccessKey: accessKey,
		SecretKey: secretKey,
	})
	if err != nil {
		log.Fatal("Error while initializing rancher client, err = ", err)
	}
	apiClient2, err = client.NewRancherClient(&client.ClientOpts{
		Url:       apiUrl2,
		AccessKey: accessKey,
		SecretKey: secretKey,
	})
	if err != nil {
		log.Fatal("Error while initializing rancher client, err = ", err)
	}
	os.Exit(m.Run())
}

func createEnvironmentWithClient(currClient *client.RancherClient, name, dockerComposePath, rancherComposePath string) (*client.Environment, error) {
	dockerComposeBytes, err := ioutil.ReadFile(dockerComposePath)
	if err != nil {
		return nil, err
	}
	dockerComposeString := string(dockerComposeBytes)
	rancherComposeString := ""
	if rancherComposePath != "" {
		rancherComposeBytes, err := ioutil.ReadFile(rancherComposePath)
		if err != nil {
			return nil, err
		}
		rancherComposeString = string(rancherComposeBytes)
	}
	return currClient.Environment.Create(&client.Environment{
		Name:           name,
		DockerCompose:  dockerComposeString,
		RancherCompose: rancherComposeString,
	})
}

func createEnvironment(name, dockerComposePath, rancherComposePath string) (*client.Environment, error) {
	return createEnvironmentWithClient(apiClient, name, dockerComposePath, rancherComposePath)
}

func createEnvironment2(name, dockerComposePath, rancherComposePath string) (*client.Environment, error) {
	return createEnvironmentWithClient(apiClient2, name, dockerComposePath, rancherComposePath)
}
