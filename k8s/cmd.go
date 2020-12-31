package k8s

import (
	"fmt"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/strvals"
	log "github.com/sirupsen/logrus"
	"os"
)

func (h *HelmClient) installChart(name, chartPath string, args map[string]string) error {
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
		log.Error(err)
		return err
	}

	if err := strvals.ParseInto(args["set"], vals); err != nil {
		log.Error("failed parsing --set data")
		return err
	}

	chartRequested, err := loader.Load(chartPath)
	log.Debugf("Chart path is: %s", chartPath)
	if err != nil {
		log.Error(err)
		return err
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		log.Error(err)
		return err
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
					log.Error(err)
					return err
				}
			} else {
				log.Error(err)
				return err
			}
		}
	}

	client.Namespace = h.settings.Namespace()
	log.Debugf("Namespace is: %s", client.Namespace)
	release, err := client.Run(chartRequested, vals)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info(release.Manifest)
	return nil
}

func (h *HelmClient) upgradeChart(name, chartPath string, args map[string]string) error {
	log.Debugf("Chart path is: %s", chartPath)
	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		log.Error(err)
		return err
	}

	client := action.NewUpgrade(h.actionConfig)
	client.Recreate = true

	p := getter.All(h.settings)
	valueOpts := &values.Options{}
	vals, err := valueOpts.MergeValues(p)
	if err := strvals.ParseInto(args["set"], vals); err != nil {
		log.Fatal(errors.Wrap(err, "failed parsing --set data"))
	}

	release, err := client.Run(name, chartRequested, vals)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug(release.Manifest)
	return nil
}

func (h *HelmClient) history(name string) (*ResponseData, error) {
	res := new(ResponseData)
	client := action.NewHistory(h.actionConfig)
	release_list, err := client.Run(name)
	if err != nil {
		log.Error(err)
		return res, err
	}
	for release, i := range release_list {
		data := History{
			Revision:    0,
			Updated: 	 i.Info.LastDeployed,
			Status:      i.Info.Status.String(),
			Description: i.Info.Description,
		}
		res.ChartHistory = append(res.ChartHistory, &data)
		log.Infof("%s %s %s\n", release, i.Info.Description, i.Info.Status)
	}

	return res, err
}

func (h *HelmClient) rollbackChart(name string) error {
	client := action.NewRollback(h.actionConfig)
	err := client.Run(name)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (h *HelmClient) deleteChart(name string) error {
	client := action.NewUninstall(h.actionConfig)
	release, err := client.Run(name)
	if err != nil {
		log.Error(err)
	}
	log.Info(release.Release)
	return nil

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
	log.Debug(2, fmt.Sprintf(format, v...))
}

