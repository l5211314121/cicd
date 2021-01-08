package mysql

type Application struct {
	Id int `json:"id"`
	SvcName string `json:"svc_name"`
	SvcDesc string `json:"svc_desc"`
	Archivepath string `json:"archivepath"`
	PackageName string `json:"packagename"`
	Modified_time string `json:"modified_time"`
	ModifiedBy string `json:"modified_by"`
	CoderepoId int `json:"coderepo_id"`
	Status int `json:"status"`
	DockerService string `json:"docker_service"`
}

type Coderepo struct {
	Id int `json:"id"`
	Url string `json:"url"`
	Modified_time string `json:"modified_time"`
	ModifiedBy string `json:"modified_by"`
	Status int `json:"status"`
	DockerBuild string `json:"docker_build"`
	ArchiveFile string `json:"arhive_file"`
	PostScript string `json:"post_script"`
}
