package handlers

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/logger"
	"github.com/docker/libcompose/lookup"
	"github.com/docker/libcompose/project"
	"github.com/rancher/go-machine-service/events"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/rancher-compose/rancher"
	"gopkg.in/yaml.v2"
)

func CreateEnvironment(event *events.Event, apiClient *client.RancherClient) (err error) {
	log.WithFields(log.Fields{
		"resourceId": event.ResourceId,
		"eventId":    event.Id,
	}).Info("Environment Create Event Received")

	env, err := getEnvironment(event.ResourceId, apiClient)
	if err != nil {
		return handleByIdError(err, event, apiClient)
	}

	if env.DockerCompose == "" {
		reply := newReply(event)
		return publishReply(reply, apiClient)
	}

	composeUrl := os.Getenv("RANCHER_URL") + "/projects/" + env.AccountId + "/schema"
	log.Infof(composeUrl)
	projectName := env.Name
	composeBytes := []byte(env.DockerCompose)
	rancherComposeMap := map[string]rancher.RancherConfig{}
	if env.RancherCompose != "" {
		err := yaml.Unmarshal([]byte(env.RancherCompose), rancherComposeMap)
		if err != nil {
			return handleByIdError(err, event, apiClient)
		}
	}

	publishChan := make(chan string, 10)
	go republishTransitioningReply(publishChan, event, apiClient)

	publishChan <- "Starting rancher-compose"
	defer func() {
		close(publishChan)
	}()

	if err := createEnv(composeUrl, projectName, composeBytes, rancherComposeMap, publishChan); err != nil {
		return handleByIdError(err, event, apiClient)
	}

	publishChan <- "Completed environment create"
	reply := newReply(event)
	return publishReply(reply, apiClient)
}

func createEnv(rancherUrl, projectName string, composeBytes []byte, rancherComposeMap map[string]rancher.RancherConfig, publishChan chan<- string) error {
	context := rancher.Context{
		Url:           rancherUrl,
		RancherConfig: rancherComposeMap,
	}
	context.ProjectName = projectName
	context.ComposeBytes = composeBytes
	context.ConfigLookup = &lookup.FileConfigLookup{}
	context.EnvironmentLookup = &lookup.OsEnvLookup{}
	context.LoggerFactory = logger.NewColorLoggerFactory()
	context.ServiceFactory = &rancher.RancherServiceFactory{
		Context: &context,
	}

	p := project.NewProject(&context.Context)

	err := p.Parse()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("Error parsing docker-compose.yml")
		publishChan <- "Error parsing docker-compose.yml"
		return err
	}

	apiClient, err := client.NewRancherClient(&client.ClientOpts{
		Url:       rancherUrl,
		AccessKey: os.Getenv("RANCHER_ACCESS_KEY"),
		SecretKey: os.Getenv("RANCHER_SECRET_KEY"),
	})

	context.Client = apiClient

	c := &context

	envs, err := c.Client.Environment.List(&client.ListOpts{
		Filters: map[string]interface{}{
			"name":         c.ProjectName,
			"removed_null": nil,
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Error getting environment list from cattle")
		publishChan <- "Error getting environment list from cattle"
		return err
	}

	for _, env := range envs.Data {
		if strings.EqualFold(c.ProjectName, env.Name) {
			log.Debugf("Found stack: %s(%s)", env.Name, env.Id)
			c.Environment = &env
		}
	}

	/*log.Infof("Creating stack %s", c.ProjectName)
	env, err := c.Client.Environment.Create(&client.Environment{
		Name: c.ProjectName,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Error while creating stack")
		publishChan <- "Error while creating stack"
		return err
	}*/

	//c.Environment = env

	context.SidekickInfo = rancher.NewSidekickInfo(p)

	err = p.Create([]string{}...)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Error while creating project.")
		publishChan <- "Error while creating project"
		return err
	}
	return nil
}
