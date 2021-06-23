package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"
)

const (
	DEFAULT_REGISTRY = "hub.azoyagroup.com"
)

func fmtNow() string {
	return time.Now().Format("20060102150405")
}

// tag 格式调整为：{版本号}-{时间字符串}-{分支名}，例如：v2.1.3-20210608142201-master。
func genTag() (ret string, err error) {
	curBranch := os.Getenv("DRONE_BRANCH")
	curTag := os.Getenv("DRONE_TAG") // git仓库触发tag
	tag := os.Getenv("PLUGIN_TAG")   // 用户自定义tag参数（兼容模式）

	if curBranch == "" && curTag != "" {
		curBranch = "master" //tag 事件时，branch为空
	}
	if curBranch == "" {
		err = errors.New("[ERR] current branch is empty.")
		return
	}

	if curTag == "" {
		curTag = tag
		if curTag == "" {
			curTag = "latest"
		}
	}

	ret = fmt.Sprintf("%s-%s-%s", curTag, fmtNow(), curBranch)
	return
}

// 将当前tag写入临时目录
// 方便自动部署阶段使用该tag
func storeTag(tag string) error {
	bn := os.Getenv("DRONE_BUILD_NUMBER") // drone构建号
	if bn == "" {
		return errors.New("invalid drone build number")
	}
	return ioutil.WriteFile("/tmp/"+bn+".tag", []byte(tag), 0777)
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
	if err := storeTag(tag); err != nil {
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
