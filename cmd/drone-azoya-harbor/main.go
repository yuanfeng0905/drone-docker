package main

import (
	"os"
	"os/exec"
	"path"
)

func main() {
	var (
		registry = "hub.azoyagroup.com"
		email    = getenv("PLUGIN_EMAIL", "HEROKU_EMAIL")
		key      = getenv("PLUGIN_API_KEY", "HEROKU_API_KEY")
	)

	os.Setenv("PLUGIN_REGISTRY", registry)
	os.Setenv("PLUGIN_REPO", path.Join(registry, "azoya", os.Getenv("DRONE_REPO_NAME")))

	os.Setenv("DOCKER_PASSWORD", key)
	os.Setenv("DOCKER_USERNAME", email)
	os.Setenv("DOCKER_EMAIL", email)

	cmd := exec.Command("drone-docker")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}

func getenv(key ...string) (s string) {
	for _, k := range key {
		s = os.Getenv(k)
		if s != "" {
			return
		}
	}
	return
}
