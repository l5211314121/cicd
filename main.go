package main

import (
	"CICD/jenkins"
	"CICD/k8s"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"path"
	"time"
)

var jenkinscmd jenkins.JenkinsCmd


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

func SetRoute(e *gin.Engine, h *k8s.HelmClient) {
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//e.GET("/reconnect", jenkinscmd.IfReconnect)
	e.POST("/buildjob", jenkinscmd.BuildJob)
	e.POST("/getjob", jenkinscmd.GetJob)
	e.POST("/writerestodb", jenkinscmd.WriteResToDB)

	e.POST("/installchart", h.InstallChart)
	e.POST("/upgradechart", h.UpgradeChart)
	e.POST("/deletechart", h.DeleteChart)
	e.POST("/charthistory", h.ChartHistory)
	e.POST("/rollbackchart", h.RollbackChart)
}

func main() {
	//gin.SetMode(gin.ReleaseMode)
	jenkinscmd.JenkinsClient()
	e := gin.Default()
	helmClient := new(k8s.HelmClient)
	helmClient.Init()
	SetRoute(e, helmClient)
	//s := &http.Server{
	//  Addr:              ":80",
	//  Handler:           e,
	//  TLSConfig:         nil,
	//  ReadTimeout:       0,
	//  ReadHeaderTimeout: 0,
	//  WriteTimeout:      0,
	//  IdleTimeout:       0,
	//  MaxHeaderBytes:    0,
	//  TLSNextProto:      nil,
	//  ConnState:         nil,
	//  ErrorLog:          nil,
	//  BaseContext:       nil,
	//  ConnContext:       nil,
	//}
	//
	//gracehttp.Serve(s)
	e.Run()
}
