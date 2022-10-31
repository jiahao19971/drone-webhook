package plugin

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/drone/drone-go/plugin/webhook"
	"github.com/sirupsen/logrus"
)

// New returns a new webhook extension.
func New(bearer string, url string) webhook.Plugin {
	return &plugin{
		Bearer: bearer,
		Url: url,
	}
}

type plugin struct {
	Bearer string
	Url string
}

type Builds struct{
	Id int
	Repo_id int
	Number int64
	Status string
	Event string
	Action string
	Link string
	Message string
	Before string
	After string
	Ref string
	Source_repo string
	Source string
	Target string
	Author_login string
	Author_name string
	Author_email string
	Author_avatar string
	Sender string
	Started int
	Finished int
	Created int
	Updated int
	Version int
}

func (p *plugin) Deliver(ctx context.Context, req *webhook.Request) error {
	if req.Action == "created" {
		url := p.Url

		repo := req.Repo.Slug
		
		build_client := http.Client{}
		build_url := url + "/api/repos/" + repo + "/builds"

		build_request, build_error := http.NewRequest("GET", build_url, nil)
		if build_error != nil {
			logrus.Errorf(build_error.Error())
		}

		authorization := "Bearer " + p.Bearer

		build_request.Header = http.Header{
			"Authorization": {authorization},
		}

		build_response, build_response_error := build_client.Do(build_request)
		if build_response_error != nil {
			logrus.Errorf(build_response_error.Error())
		}

		build_body, build_body_error := ioutil.ReadAll(build_response.Body)
		if build_body_error != nil {
			logrus.Errorf(build_body_error.Error())
		}

		string_build := string(build_body)

		var json_build []Builds
		json.Unmarshal([]byte(string_build), &json_build)

		current_build_number := []int64{}

		logrus.Infof("Current branch source %s", req.Build.Source)
		logrus.Infof("Current branch target %s", req.Build.Target)
		if req.Build.Source != req.Build.Target {
			logrus.Infof("Current branch is coming from a pull request")
		} else {
			logrus.Infof("Current branch is coming from a branch push")
		}

		logrus.Infof("Build events %s", req.Build.Event)
		logrus.Infof("Build action %s", req.Build.Action)
		logrus.Infof("Temporary apply only to chore/test-dronejsonnet")
		for _, m := range json_build {
			if m.Target == "chore/test-dronejsonnet" && 
			(m.Status == "running" || m.Status == "pending") {
				current_build_number = append(current_build_number, m.Number)
			}
		}

		should_terminate_build := []int64{}

		for _, i := range current_build_number {
			if i != req.Build.Number {
				should_terminate_build = append(should_terminate_build, i)
			}
		}

		if len(should_terminate_build) == 0 {
			logrus.Infof("No build will be kill")
		} else {
			logrus.Infof("Build that should be terminated %s", should_terminate_build)
			for _, i := range should_terminate_build {
				delete_client := http.Client{}
				delete_url := url + "/api/repos/" + repo + "/builds/" + strconv.Itoa(int(i))

				delete_request, delete_error := http.NewRequest("DELETE", delete_url, nil)
				if delete_error != nil {
					logrus.Errorf(delete_error.Error())
				}
	
				delete_request.Header = http.Header{
					"Authorization": {authorization},
				}
	
				delete_response, delete_response_err := delete_client.Do(delete_request)
				if delete_response_err != nil {
					logrus.Errorf(delete_response_err.Error())
				}
	
				if delete_response.StatusCode == 200 {
					logrus.Infof("Build %s was successfully killed", i)
				} else {
					delete_body, delete_body_error := ioutil.ReadAll(delete_response.Body)
					if delete_body_error != nil {
						logrus.Errorf(delete_body_error.Error())
					}
					delete_string := string(delete_body)
					logrus.Errorf(delete_string)
				}
			}
		}
	}

	return nil
}