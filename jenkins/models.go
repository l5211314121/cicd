package jenkins

import "github.com/bndr/gojenkins"

const (
	JENKINS_BASE_URL = "http://47.92.148.31:32008/"
	JENKINS_ADMIN_USER = "yesheng"
	JENKINS_ADMIN_PASS = "chouxifu.521"
)

type JenkinsCmd struct {
	*gojenkins.Jenkins
	//Reconnect func() *JenkinsCmd
}

