package tests

import (
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/client"
)

func TestCreateOnlyDockerCompose(t *testing.T) {
	dockerComposePath := "assets/only_docker_compose/docker-compose.yml"
	env, err := createEnvironment("onlyDCompose", dockerComposePath, "")
	if err != nil {
		log.Error("Error creating environment, err = ", err)
		t.FailNow()
	}
	newEnv := &client.Environment{}
	for {
		err = apiClient.Reload(&(env.Resource), newEnv)
		if err != nil {
			log.Error("Error updating environment, err = ", err)
			t.FailNow()
		}
		if newEnv.Transitioning != "yes" {
			break
		}
		<-time.After(2 * time.Second)
	}
	servicesMap := map[string]interface{}{}
	err = apiClient.GetLink(env.Resource, "services", &servicesMap)
	if err != nil {
		log.Error("Error getting services, err = ", err)
		t.FailNow()
	}
	services := servicesMap["data"].([]interface{})
	if len(services) == 0 {
		log.Error("Expected No. of Services = 1, but obtained %d", len(services))
		t.FailNow()
	}
	err = apiClient.Environment.Delete(newEnv)
	if err != nil {
		log.Error("Error while deleting stack ", newEnv.Name)
	}
}

func TestCreateDockerComposeAndRancherCompose(t *testing.T) {
	dockerComposePath := "assets/docker_compose_and_rancher_compose/docker-compose.yml"
	rancherComposePath := "assets/docker_compose_and_rancher_compose/rancher-compose.yml"
	env, err := createEnvironment("DockerRancherCompose", dockerComposePath, rancherComposePath)
	if err != nil {
		log.Error("Error creating environment, err = ", err)
		t.FailNow()
	}
	newEnv := &client.Environment{}
	for {
		err = apiClient.Reload(&(env.Resource), newEnv)
		if err != nil {
			log.Error("Error updating environment, err = ", err)
			t.FailNow()
		}
		if newEnv.Transitioning != "yes" {
			break
		}
		<-time.After(2 * time.Second)
	}
	servicesMap := map[string]interface{}{}
	err = apiClient.GetLink(env.Resource, "services", &servicesMap)
	if err != nil {
		log.Error("Error getting services, err = ", err)
		t.FailNow()
	}
	services := servicesMap["data"].([]interface{})
	if len(services) == 0 {
		log.Errorf("Expected No. of Services = 1, but obtained %d", len(services))
		t.FailNow()
	}
	err = apiClient.Environment.Delete(newEnv)
	if err != nil {
		log.Error("Error while deleting stack ", newEnv.Name)
	}
}

func TestCreateDockerMultipleServices(t *testing.T) {
	dockerComposePath := "assets/multiple_services/docker-compose.yml"
	rancherComposePath := "assets/multiple_services/rancher-compose.yml"
	env, err := createEnvironment("MultipleServices", dockerComposePath, rancherComposePath)
	if err != nil {
		log.Error("Error creating environment, err = ", err)
		t.FailNow()
	}
	newEnv := &client.Environment{}
	for {
		err = apiClient.Reload(&(env.Resource), newEnv)
		if err != nil {
			log.Error("Error updating environment, err = ", err)
			t.FailNow()
		}
		if newEnv.Transitioning != "yes" {
			break
		}
		<-time.After(2 * time.Second)
	}
	servicesMap := map[string]interface{}{}
	err = apiClient.GetLink(env.Resource, "services", &servicesMap)
	if err != nil {
		log.Error("Error getting services, err = ", err)
		t.FailNow()
	}
	services := servicesMap["data"].([]interface{})
	if len(services) != 2 {
		log.Errorf("Expected No. of Services = 2, but obtained %d", len(services))
		t.FailNow()
	}
	err = apiClient.Environment.Delete(newEnv)
	if err != nil {
		log.Error("Error while deleting stack ", newEnv.Name)
	}
}

func TestCreateMultipleEnvs(t *testing.T) {
	dockerComposePath := "assets/multiple_environments/docker-compose-1a5.yml"
	rancherComposePath := "assets/multiple_environments/rancher-compose-1a5.yml"
	env1, err := createEnvironment("MultipleEnvironments1", dockerComposePath, rancherComposePath)
	if err != nil {
		log.Error("Error creating environment, err = ", err)
		t.FailNow()
	}
	dockerComposePath2 := "assets/multiple_environments/docker-compose-1a6.yml"
	rancherComposePath2 := "assets/multiple_environments/rancher-compose-1a6.yml"
	env2, err := createEnvironment2("MultipleEnvironments2", dockerComposePath2, rancherComposePath2)
	if err != nil {
		log.Error("Error creating environment2, err = ", err)
		t.FailNow()
	}
	newEnv1 := &client.Environment{}
	newEnv2 := &client.Environment{}
	for {
		err = apiClient.Reload(&(env1.Resource), newEnv1)
		if err != nil {
			log.Error("Error updating environment, err = ", err)
			t.FailNow()
		}
		if newEnv1.Transitioning != "yes" {
			break
		}
		<-time.After(2 * time.Second)
	}
	for {
		err = apiClient.Reload(&(env2.Resource), newEnv2)
		if err != nil {
			log.Error("Error updating environment, err = ", err)
			t.FailNow()
		}
		if newEnv2.Transitioning != "yes" {
			break
		}
		<-time.After(2 * time.Second)
	}
	servicesMap := map[string]interface{}{}
	err = apiClient.GetLink(env1.Resource, "services", &servicesMap)
	if err != nil {
		log.Error("Error getting services, err = ", err)
		t.FailNow()
	}
	servicesMap2 := map[string]interface{}{}
	err = apiClient.GetLink(env2.Resource, "services", &servicesMap2)
	if err != nil {
		log.Error("Error getting services, err = ", err)
		t.FailNow()
	}
	services := servicesMap["data"].([]interface{})
	if len(services) != 1 {
		log.Errorf("Expected No. of Services = 1, but obtained %d", len(services))
		t.FailNow()
	}
	services2 := servicesMap2["data"].([]interface{})
	if len(services2) != 1 {
		log.Errorf("Expected No. of Services = 1, but obtained %d", len(services2))
		t.FailNow()
	}
	err = apiClient.Environment.Delete(newEnv1)
	if err != nil {
		log.Error("Error while deleting stack ", newEnv1.Name)
	}
	err = apiClient.Environment.Delete(newEnv2)
	if err != nil {
		log.Error("Error while deleting stack ", newEnv2.Name)
	}
}

func TestCreateWithBuildNoFile(t *testing.T) {
	dockerComposePath := "assets/build_image_from_url/docker-compose.yml"
	rancherComposePath := "assets/build_image_from_url/rancher-compose.yml"
	env, err := createEnvironment("BuildNoFile", dockerComposePath, rancherComposePath)
	if err != nil {
		log.Error("Error creating environment, err = ", err)
		t.FailNow()
	}
	newEnv := &client.Environment{}
	for {
		err = apiClient.Reload(&(env.Resource), newEnv)
		if err != nil {
			log.Error("Error updating environment, err = ", err)
			t.FailNow()
		}
		if newEnv.Transitioning != "yes" {
			break
		}
		<-time.After(2 * time.Second)
	}
	servicesMap := map[string]interface{}{}
	err = apiClient.GetLink(env.Resource, "services", &servicesMap)
	if err != nil {
		log.Error("Error getting services, err = ", err)
		t.FailNow()
	}
	services := servicesMap["data"].([]interface{})
	if len(services) != 1 {
		log.Errorf("Expected No. of Services = 1, but obtained %d", len(services))
		t.FailNow()
	}
	err = apiClient.Environment.Delete(newEnv)
	if err != nil {
		log.Error("Error while deleting stack ", newEnv.Name)
	}
}

func TestCreateWithSidekicks(t *testing.T) {
	dockerComposePath := "assets/sidekick/docker-compose.yml"
	rancherComposePath := "assets/sidekick/rancher-compose.yml"
	env, err := createEnvironment("Sidekick", dockerComposePath, rancherComposePath)
	if err != nil {
		log.Error("Error creating environment, err = ", err)
		t.FailNow()
	}
	newEnv := &client.Environment{}
	for {
		err = apiClient.Reload(&(env.Resource), newEnv)
		if err != nil {
			log.Error("Error updating environment, err = ", err)
			t.FailNow()
		}
		if newEnv.Transitioning != "yes" {
			break
		}
		<-time.After(2 * time.Second)
	}
	servicesMap := map[string]interface{}{}
	err = apiClient.GetLink(env.Resource, "services", &servicesMap)
	if err != nil {
		log.Error("Error getting services, err = ", err)
		t.FailNow()
	}
	services := servicesMap["data"].([]interface{})
	if len(services) != 1 {
		log.Errorf("Expected No. of Services = 1, but obtained %d", len(services))
		t.FailNow()
	}
	err = apiClient.Environment.Delete(newEnv)
	if err != nil {
		log.Error("Error while deleting stack ", newEnv.Name)
	}
}
