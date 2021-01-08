package jenkins

import "github.com/bndr/gojenkins"

const (
	JENKINS_BASE_URL   = "http://47.92.148.31:32008/"
	JENKINS_ADMIN_USER = "xxx"
	JENKINS_ADMIN_PASS = "xxx"
)

type JenkinsCmd struct {
	*gojenkins.Jenkins
	//Reconnect func() *JenkinsCmd
}


type Application struct {
	Id int `json:"id"`
	User string `json:"username"`
	SvcName string `json:"svc_name"`
	SvcDesc string `json:"svc_desc"`
	Archivepath string `json:"archivepath"`
	PackageName string `json:"packagename"`
	Modeifed_time string `json:"modeifed_time"`
	ModifiedBy string `json:"modified_by"`
	CoderepoUrl string `json:"coderepo_url"`
	Status int `json:"status"`
	DockerService string `json:"docker_service"`
}