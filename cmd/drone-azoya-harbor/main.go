package main

import (
	"os"
	"os/exec"
	"path"
)

func main() {
	var (
		registry = "hub.azoyagroup.com"
		username = os.Getenv("PLUGIN_USERNAME")
		password = os.Getenv("PLUGIN_PASSWORD")
	)

	os.Setenv("PLUGIN_REGISTRY", registry)
	os.Setenv("PLUGIN_REPO", path.Join(registry, "azoya", os.Getenv("DRONE_REPO_NAME")))

	// 如果未配置，默认使用运维账号权限
	if username == "" {
		username = "ops@haituncun.com"
	}
	if password == "" {
		password = "Azoya#0422"
	}

	os.Setenv("DOCKER_USERNAME", username)
	os.Setenv("DOCKER_PASSWORD", password)

	cmd := exec.Command("drone-docker")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}
