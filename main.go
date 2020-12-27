package main

import (
	"CICD/jenkins"
	"CICD/k8s"

	"github.com/gin-gonic/gin"
)

var jenkinscmd jenkins.JenkinsCmd

func SetRoute(e *gin.Engine, h *k8s.HelmClient) {
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	e.GET("/reconnect", jenkinscmd.IfReconnect)
	e.POST("/buildjob", jenkinscmd.BuildJob)
	e.POST("/getjob", jenkinscmd.GetJob)

	e.GET("/installchart", h.InstallChart)
	e.GET("/upgradechart", h.UpgradeChart)
	e.POST("/deletechart", h.ParseData, h.DeleteChart)
	e.GET("/charthistory", h.ChartHistory)
	e.GET("/rollbackchart", h.RollbackChart)
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
