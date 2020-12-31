package k8s

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/time"
)

type HelmClient struct {
	settings     *cli.EnvSettings
	actionConfig *action.Configuration
	chartRoot    string
	chartPath  	 string
}

type ReqData struct {
	ServiceName string `json:"service_name"`
	ImageTag   string `json:"image_tag"`
}

type ResponseData struct {
	ChartHistory []*History `json:"chart_history"`
}

type History struct {
	Revision int `json:"revision"`
	Updated time.Time `json:"updated"`
	Status string `json:"status"`
	Description string 	`json:"description"`
}

