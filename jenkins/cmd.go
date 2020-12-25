package jenkins

import (
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/gin-gonic/gin"
)

type params struct {
	jobName string
}

func (j *JenkinsCmd)JenkinsClient() (*JenkinsCmd, error) {
	jenkinsClient, err := gojenkins.CreateJenkins(nil, JENKINS_BASE_URL, JENKINS_ADMIN_USER, JENKINS_ADMIN_PASS).Init()
	j.Jenkins = jenkinsClient
	return j, err
}

func (j *JenkinsCmd)IfReconnect(c *gin.Context)  {
	_, err := j.Jenkins.GetJob("test")
	err = fmt.Errorf("error")
	if err != nil {
		fmt.Printf("abc\n")
		j, err = j.JenkinsClient()
	}
}

func (j *JenkinsCmd) BuildJob(c *gin.Context) {
	req := new(struct {
		JobName string `json:"job_name"`
	})
	//data, _ := ioutil.ReadAll(c.Request.Body)
	c.ShouldBindJSON(req)
	id, err := j.Jenkins.BuildJob(req.JobName)
	if err != nil {
		c.JSON(500, map[string]interface{}{"error": fmt.Sprintf("%s", err)})
		return
	}
	c.JSON(200, map[string]interface{}{"Success": "True", "ID": id})
}

func (j *JenkinsCmd) GetJob(c *gin.Context){
	req := new(struct {
		JobName string `json:"job_name"`
	})
	c.ShouldBindJSON(req)
	job, _ := j.Jenkins.GetJob(req.JobName)
	fmt.Printf("%s", job.Raw.LastBuild)
}