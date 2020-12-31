package jenkins

import (
	"CICD/lib"
	"CICD/mysql"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"path"
	"time"
)

type ReqData struct {
	JobName string `json:"job_name"`
	GitUrl string `json:"git_url"`
	Username string `json:"user_name"`
}

func init(){
	baseLogPath := path.Join("/root/Projects/src/CICD/", "cicd.log")
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour), // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)

	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	log.SetOutput(writer)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	//lfHook := lfshook.NewHook(lfshook.WriterMap{
	//	log.DebugLevel: writer, // 为不同级别设置不同的输出目的
	//	log.InfoLevel:  writer,
	//	log.WarnLevel:  writer,
	//	log.ErrorLevel: os.Stdout,
	//	log.FatalLevel: writer,
	//	log.PanicLevel: writer,
	//}, &log.JSONFormatter{})
	//log.AddHook(lfHook)

}

func ParseData(c *gin.Context) *ReqData{
	reqData := new(ReqData)
	c.ShouldBindJSON(reqData)
	return reqData
}

func (j *JenkinsCmd)JenkinsClient() (*JenkinsCmd, error) {
	jenkinsClient, err := gojenkins.CreateJenkins(nil, JENKINS_BASE_URL, JENKINS_ADMIN_USER, JENKINS_ADMIN_PASS).Init()
	j.Jenkins = jenkinsClient
	return j, err
}

func (j *JenkinsCmd)IfReconnect()  {
	_, err := j.Jenkins.GetJob("test")
	log.Error(err)
	if err != nil {
		fmt.Printf("abc\n")
		j, err = j.JenkinsClient()
	}
}

func (j *JenkinsCmd) BuildJob(c *gin.Context) {
	j.IfReconnect()
	reqData := ParseData(c)
	//data, _ := ioutil.ReadAll(c.Request.Body)
	coderepo := new(mysql.Coderepo)
	db := mysql.Getconn()
	row:= db.QueryRow("select id from coderepo where url=?",reqData.GitUrl)
	if err := row.Scan(&coderepo.Id); err != nil {
		log.Info("build job: ", err)

		if err.Error() == "sql: no rows in result set" {
			_, err := db.Exec(
				"insert into coderepo(url,modified,status,modified_by) values(?,?,?,?)",
				reqData.GitUrl,
				time.Now(), 1,
				reqData.Username)

			if err != nil {
				log.Error("Build Job write to mysql failed: ", err)
				c.JSON(500, lib.RespErr("Build Job write to mysql failed: ", err))
				return
			}
		} else {
			log.Error(err)
			c.JSON(500, lib.RespErr("scan failed, err ", err))
			return
		}
	}
	id, err := j.Jenkins.BuildJob(reqData.JobName)
	if err != nil {
		c.JSON(500, lib.RespErr("Build job error: ", err))
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

func (j *JenkinsCmd) WriteResToDB(c *gin.Context){
	reqData := new(mysql.Application)
	c.ShouldBindJSON(reqData)
	application := new(mysql.Application)
	db := mysql.Getconn()
	row := db.QueryRow("select id from application where svc_name=?", reqData.SvcName)
	if err := row.Scan(&application.Id); err != nil {
		if err.Error() == "sql: no rows in result set" {

		}
	}

}