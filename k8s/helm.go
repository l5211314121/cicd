package k8s

import (
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	_ "github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

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

func (h *HelmClient) Init() {
	h.settings = cli.New()
	h.actionConfig = new(action.Configuration)
	if err := h.actionConfig.Init(h.settings.RESTClientGetter(), h.settings.Namespace(),
		os.Getenv("HELM_DRIVER"), debug); err != nil {
		log.Error(err)
	}
	h.chartRoot = "/root/helm/"
}

func ParseData(c *gin.Context) *ReqData{
	reqData := new(ReqData)
	c.ShouldBindJSON(reqData)
	return reqData
}

func (h *HelmClient) InstallChart(c *gin.Context) {
	reqData := ParseData(c)
	c.ShouldBindJSON(reqData)
	fmt.Println(reqData.ServiceName)
	args := map[string]string{
		"set": "image.tag=" + reqData.ImageTag,
	}
	h.chartPath = h.chartRoot + reqData.ServiceName
	if err := h.installChart(reqData.ServiceName, h.chartPath, args); err != nil {
		c.JSON(500, map[string]string{"error": fmt.Sprintf("%s", err)})
	}
}

func (h *HelmClient) UpgradeChart(c *gin.Context) {
	reqData := ParseData(c)
	c.ShouldBindJSON(reqData)
	args := map[string]string{
		"set": "image.tag=" + reqData.ImageTag,
	}
	h.chartPath = h.chartRoot + reqData.ServiceName
	if err := h.upgradeChart(reqData.ServiceName, h.chartPath, args); err != nil {
		c.JSON(500, map[string]string {"error": fmt.Sprintf("%s", err)})
	}
}

func (h *HelmClient) DeleteChart(c *gin.Context) {
	// data, _ := ioutil.ReadAll(c.Request.Body)
	// fmt.Printf("0000: ", string(data))
	reqData := ParseData(c)
	c.ShouldBindJSON(reqData)
	if reqData.ServiceName == "" {
		c.JSON(500, map[string]string{"error": "ServiceName is needed!"})
		return
	}
	if err := h.deleteChart(reqData.ServiceName); err != nil {
		c.JSON(500, map[string]string {"error": fmt.Sprintf("%s", err)})
	}
}

func (h *HelmClient) ChartHistory(c *gin.Context) {
	reqData := ParseData(c)
	c.ShouldBindJSON(reqData)
	if reqData.ServiceName == "" {
		c.JSON(500, map[string]string{"error": "ServiceName is needed!"})
		return
	}
	if historyData, err := h.history(reqData.ServiceName); err != nil {
		c.JSON(500, map[string]string {"error": fmt.Sprintf("%s", err)})
	} else {
		c.JSON(200, historyData)
	}
}

func (h *HelmClient) RollbackChart(c *gin.Context) {
	reqData := ParseData(c)
	c.ShouldBindJSON(reqData)
	if reqData.ServiceName == "" {
		c.JSON(500, map[string]string{"error": "ServiceName is needed!"})
	}
	if err := h.rollbackChart(reqData.ServiceName); err != nil {
		c.JSON(500, map[string]string {"error": fmt.Sprintf("%s", err)})
	}
}

