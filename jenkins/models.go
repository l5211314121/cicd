package jenkins

import "github.com/bndr/gojenkins"

const (
	JENKINS_BASE_URL   = "http://47.92.148.31:32008/"
	JENKINS_ADMIN_USER = "xxxx"
	JENKINS_ADMIN_PASS = "xxxx"
)

type JenkinsCmd struct {
	*gojenkins.Jenkins
	//Reconnect func() *JenkinsCmd
}
