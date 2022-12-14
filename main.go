package main

import (
	"net/http"

	"drone-webhook/plugin"

	"github.com/drone/drone-go/plugin/webhook"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type spec struct {
	Bind   string `envconfig:"DRONE_BIND"`
	Debug  bool   `envconfig:"DRONE_DEBUG"`
	Secret string `envconfig:"DRONE_SECRET"`
	Bearer string `envconfig:"DRONE_BEARER"`
	URL string `envconfig:"DRONE_URL"`
	Master_branch string `envconfig:"MASTER_BRANCH"`
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Bearer == "" {
		logrus.Fatalln("missing Bearer key")
	}
	if spec.URL == "" {
		logrus.Fatalln("missing URL key")
	}
	if spec.Master_branch == "" {
		logrus.Fatalln("missing Master Branch key")
	}
	if spec.Bind == "" {
		spec.Bind = ":3000"
	}

	handler := webhook.Handler(
		plugin.New(spec.Bearer, spec.URL, spec.Master_branch),
		spec.Secret,
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Bind)

	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}
