package k8s

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/strvals"
	"helm.sh/helm/v3/pkg/chart/loader"
	"log"
	"os"
)

type HelmClient struct {
	settings *cli.EnvSettings
	actionConfig *action.Configuration
}

func (h *HelmClient) Init(){
	h.settings = cli.New()
	h.actionConfig = new(action.Configuration)
	if err := h.actionConfig.Init(h.settings.RESTClientGetter(), h.settings.Namespace(),
		os.Getenv("HELM_DRIVER"), debug); err != nil {
		log.Fatal(err)
	}
}

func (h *HelmClient) InstallChart(c * gin.Context) {
	args := map[string]string {
		"set": "image.tag=25",
	}
	h.installChart("testchart", "/root/helm/mychart", args)
}

func (h *HelmClient) UpgradeChart(c *gin.Context){
	args := map[string] string {
		"set": "image.tag=25",
	}
	h.upgradeChart("testchart", "/root/helm/mychart", args)
}

func (h *HelmClient) DeleteChart(c *gin.Context){
	h.deleteChart("testchart")
}

func (h*HelmClient) ChartHistory(c *gin.Context){
	h.history("testchart")
}

func (h *HelmClient) RollbackChart(c *gin.Context) {
	h.rollbackChart("testchart")
}

func (h *HelmClient) installChart(name, chartPath string, args map[string]string) {
	client := action.NewInstall(h.actionConfig)
	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}

	client.ReleaseName = name
	//cp, err := client.ChartPathOptions.LocateChart(repo, settings)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//debug("CHART PATH: %s\n", cp)

	p := getter.All(h.settings)
	valueOpts := &values.Options{}
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		log.Fatal(err)
	}

	if err := strvals.ParseInto(args["set"], vals); err != nil {
		log.Fatal(errors.Wrap(err, "failed parsing --set data"))
	}

	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		log.Fatal(err)
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		log.Fatal(err)
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        chartPath,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: h.settings.RepositoryConfig,
					RepositoryCache:  h.settings.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}
	}

	client.Namespace = h.settings.Namespace()
	release, err := client.Run(chartRequested, vals)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(release.Manifest)
}

func (h *HelmClient) upgradeChart(name, chartPath string, args map[string]string) {
	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		log.Fatal(err)
	}

	client := action.NewUpgrade(h.actionConfig)
	client.Recreate = true

	p := getter.All(h.settings)
	valueOpts := &values.Options{}
	vals, err := valueOpts.MergeValues(p)
	if err := strvals.ParseInto(args["set"], vals); err != nil {
		log.Fatal(errors.Wrap(err, "failed parsing --set data"))
	}
	fmt.Println("---------------", vals)

	release, err := client.Run(name, chartRequested, vals)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(release.Manifest)
}

func (h *HelmClient) history(name string)  {
	client := action.NewHistory(h.actionConfig)
	release_list, err := client.Run(name)
	if err != nil {
		log.Fatal(err)
	}
	for release, i := range(release_list) {
		fmt.Println("-------", release, i.Info.Description, i.Info.Status)
	}
}

func (h *HelmClient) rollbackChart(name string) {
	client := action.NewRollback(h.actionConfig)
	err := client.Run(name)
	if err != nil {
		log.Fatal(err)
	}

}

func (h *HelmClient) deleteChart(name string){
	client := action.NewUninstall(h.actionConfig)
	release, err := client.Run(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(release.Release)

}


func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v...))
}