package sonarProject

import (
	"encoding/json"

	sonarapi "github.com/aberestyak/sonar-cleaner/internal/sonarapi"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func getPaginationPages(analyzedBefore string, config sonarapi.QueryConfig) (int, error) {
	pagination := SonarPagination{}

	request := sonarapi.NewSonarRequest(config)

	request.Method = "GET"
	request.URL.Path = sonarapi.SonarqubeSearchProjectsEndpoint

	query := request.URL.Query()
	query.Add("ps", paginationSize)
	query.Add("analyzedBefore", analyzedBefore)
	request.URL.RawQuery = query.Encode()

	body, err := sonarapi.DoHttpRequest(request)
	if err != nil {
		return 0, errors.Wrap(err, "Couldn't get projects list")
	}

	if err := json.Unmarshal(*body, &pagination); err != nil {
		return 0, errors.Wrap(err, "Couldn't parse response body with pagination")
	}

	log.Debugf("Found %d projects, which weren't analised sinse %s", pagination.Spec.Total, analyzedBefore)

	if pagination.Spec.Total > pagination.Spec.PageSize {
		return pagination.Spec.Total/pagination.Spec.PageSize + 1, nil
	}
	return 1, nil
}
