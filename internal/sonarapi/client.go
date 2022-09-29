package sonarclient

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
)

type QueryConfig struct {
	SonarqubeAddress string
	SonarqubeToken   string
	KeepDays         int
}

const (
	SonarqubeSearchProjectsEndpoint        = "/api/projects/search"
	SonarqubeSearchProjectAnalysisEndpoint = "/api/project_analyses/search"
	SonarqubeSearchProjectBranchesEndpoint = "/api/project_branches/list"
	SonarqubeDeleteProjectBranchesEndpoint = "/api/project_branches/delete"
	SonarqubeDeleteProjectAnalysesEndpoint = "/api/project_analyses/delete"
)

func NewSonarRequest(config QueryConfig) *http.Request {
	scheme := strings.Split(config.SonarqubeAddress, "://")[0]
	host := strings.Split(config.SonarqubeAddress, "://")[1]
	return &http.Request{
		URL: &url.URL{
			Scheme: scheme,
			Host:   host,
		},
		Header: map[string][]string{
			"Authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte(config.SonarqubeToken))},
		},
	}
}
