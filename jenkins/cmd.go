package jenkins

import (
	"CICD/lib"
	"CICD/mysql"
	"database/sql"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type ReqData struct {
	JobName string `json:"job_name"`
	GitUrl string `json:"git_url"`
	Username string `json:"user_name"`
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
	if err != nil {
		log.Error("Get Job err: ", err)
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
				"insert into coderepo(url,modified_time,status,modified_by) values(?,?,?,?)",
				reqData.GitUrl,
				time.Now(), 1,
				reqData.Username)

			if err != nil {
				log.Error("Build Job write to mysql failed: ", err)
				c.JSON(500, lib.RespErr("Build Job write to mysql failed: ", err.Error()))
				return
			}
		} else {
			log.Error(err)
			c.JSON(500, lib.RespErr("scan failed, err ", err.Error()))
			return
		}
	}
	id, err := j.Jenkins.BuildJob(reqData.JobName, map[string]string{"BUILD_USER": reqData.Username})
	if err != nil {
		c.JSON(500, lib.RespErr("Build job error: ", err.Error()))
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
	//data, _ := ioutil.ReadAll(c.Request.Body)
	//fmt.Printf("data", string(data))
	reqData := new(Application)
	err := c.ShouldBindJSON(reqData)
	fmt.Println("shouldbindjson err", err, reqData)
	log.Info(fmt.Sprintf("reqdata: %s, %s, %s", reqData.PackageName, reqData.CoderepoUrl, reqData.User))
	db := mysql.Getconn()


	codeRepoId, err := getCodeRepoId(reqData.CoderepoUrl)
	if err != nil {
		c.JSON(500, lib.RespErr("Failed to get code repo id.", err.Error()))
		return
	}

	// 循环应用名
	for k, v := range(getArchPathAndName(reqData.PackageName)) {
		// 查找应用名是否存在
		err := writetodb(v, k, codeRepoId, reqData, db)
		if err !=nil {
			c.JSON(500, lib.RespErr("write to db error: ", err.Error()))
			return
		}
	}
}

func getCodeRepoId(repoUrl string) (int, error) {
	codeRepo := new(mysql.Coderepo)
	db := mysql.Getconn()
	row := db.QueryRow("select id from coderepo where url=?", repoUrl)
	log.Info("getCodeRepoId: repoUrl", repoUrl)
	if err := row.Scan(&codeRepo.Id); err != nil {
		log.Error("Failed to get code repo id")
		return 0, errors.New("query error")
	} else {
		return codeRepo.Id, nil
	}
}

func getArchPathAndName(svcString string) (map[string]string) {
	log.Info("Svc String: ", svcString)
	if len(svcString) == 0 {
		return map[string]string{}
	}
	s1 := strings.Split(svcString, "|")
	res := map[string]string{}
	for _, packagePath := range(s1) {
		packagePathSlice := strings.Split(packagePath, "/")
		res[packagePath] = (strings.Split(packagePathSlice[len(packagePathSlice)-1], "."))[0]
	}
	return res
}

func writetodb(svcName, archivePath string, codeRepoId int, reqData *Application, db *sql.DB) error {
	log.Info("exec loop: ", svcName, archivePath)
	application := new(mysql.Application)
	codeRepo := new(mysql.Coderepo)

	row := db.QueryRow("select id, coderepo_id from application where svc_name=?", svcName)
	if err := row.Scan(&application.Id, &application.CoderepoId); err != nil {
		log.Info("select id, coderepo_id from application where svc_name=", svcName)

		//  如果应用名不存在
		if err.Error() == "sql: no rows in result set" {
			_, err := db.Exec("insert into application(svc_name, archivepath, packagename, modified_time, modified_by, coderepo_id) values(?,?,?,?,?,?)",
				svcName, archivePath, time.Now(), reqData.User, codeRepoId)
			if err != nil {
				log.Error("Failed to write to db", err.Error())
				return err
			}
		} else {
			return err
		}
	} else { // 如果应用名存在，判断应用名对应的repo是不是本repo的应用名
		if codeRepoId != application.CoderepoId {
			row = db.QueryRow("select url from coderepo where id=?", codeRepoId)
			if err := row.Scan(&codeRepo.Url); err != nil {
				log.Error("`WriteResToDB` Error to get repo url: ", err.Error())
				return err
			} else {
				log.Errorf("Service name `%s` has been used by repo `%s`", svcName, codeRepo.Url)
				return err
			}
		} else {
			log.Info("Write to db: update application ", svcName)
			_, err := db.Exec("update application set modified_time=?, modified_by=? where svc_name=?", time.Now(), reqData.User, svcName)
			if err != nil {
				log.Error("Update service error ", err)
				return err
			}
		}
	}
	return nil
}