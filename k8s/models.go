package k8s

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

type HelmClient struct {
	settings     *cli.EnvSettings
	actionConfig *action.Configuration
	chartRoot    string
	RequestData  *ReqData
}

type ReqData struct {
	ServiceName string `json:"service_name"`
	ImageName   string `json:"imagename"`
}
