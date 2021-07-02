package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	DEFAULT_REGISTRY = "hub.azoyagroup.com"
)

// tag 格式调整为：{版本号}-{构建流水号}-{分支名}，例如：v2.1.3-101-master。
func genTag() (ret string, err error) {
	curBranch := os.Getenv("DRONE_BRANCH") // 当前构建的分支
	gitTag := os.Getenv("DRONE_TAG")       // git仓库触发tag
	tag := os.Getenv("PLUGIN_TAGS")        // 用户自定义tag参数（兼容模式）
	if tag == "" {
		tag = os.Getenv("PLUGIN_TAG")
	}
	bn := os.Getenv("DRONE_BUILD_NUMBER") // 构建流水号

	if curBranch == "" && gitTag != "" {
		curBranch = "master" //tag 事件时，branch为空
	}
	if curBranch == "" {
		err = errors.New("[ERR] current branch is empty.")
		return
	}
	// 替换特殊字符
	curBranch = strings.ReplaceAll(curBranch, "/", "-")

	curTag := "latest" // default tag
	if tag != "" {
		curTag = tag
	}
	if gitTag != "" { // git tag 优先级最高
		curTag = gitTag
	}

	ret = fmt.Sprintf("%s-%s-%s", curTag, bn, curBranch)
	return
}

func main() {
	var (
		repo     = os.Getenv("PLUGIN_REPO")
		registry = os.Getenv("PLUGIN_REGISTRY")
		username = os.Getenv("PLUGIN_USERNAME")
		password = os.Getenv("PLUGIN_PASSWORD")
	)

	if registry == "" {
		registry = DEFAULT_REGISTRY
	}
	if repo == "" {
		repo = path.Join(registry, "azoya", os.Getenv("DRONE_REPO_NAME"))
	}

	os.Setenv("PLUGIN_REGISTRY", registry)
	os.Setenv("PLUGIN_REPO", repo)

	// 如果未配置，默认使用运维账号权限
	if username == "" {
		username = "ops@haituncun.com"
	}
	if password == "" {
		password = "Azoya#0422"
	}

	os.Setenv("DOCKER_USERNAME", username)
	os.Setenv("DOCKER_PASSWORD", password)

	// 生成自定义tag
	tag, err := genTag()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	os.Setenv("PLUGIN_TAG", tag)

	cmd := exec.Command("drone-docker")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
