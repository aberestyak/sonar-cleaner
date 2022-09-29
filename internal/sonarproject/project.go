package sonarProject

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/kofoworola/godate"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	sonarapi "github.com/aberestyak/sonar-cleaner/internal/sonarapi"
)

const (
	paginationSize = "100"
)

func GetOutdatedProjects(days int, config sonarapi.QueryConfig) (*SonarProjects, error) {
	projects := SonarProjects{}
	analyzedBefore := godate.Now(time.UTC).Sub(days, godate.DAY).ToDateString()

	pages, err := getPaginationPages(analyzedBefore, config)
	if err != nil {
		return nil, err
	}

	for i := 1; i <= pages; i++ {
		request := sonarapi.NewSonarRequest(config)
		project := SonarProjects{}
		request.Method = "GET"
		request.URL.Path = sonarapi.SonarqubeSearchProjectsEndpoint

		query := request.URL.Query()
		query.Add("ps", paginationSize)
		query.Add("p", strconv.Itoa(i))
		query.Add("analyzedBefore", analyzedBefore)
		request.URL.RawQuery = query.Encode()

		body, err := sonarapi.DoHttpRequest(request)
		if err != nil {
			return nil, errors.Wrap(err, "Couldn't get projects list")
		}
		if err := json.Unmarshal(*body, &project); err != nil {
			return nil, errors.Wrap(err, "Couldn't parse response body")
		}
		projects.Projects = append(projects.Projects, project.Projects...)
	}

	return &projects, nil
}

func (project *SonarProject) GetProjectBranches(config sonarapi.QueryConfig) error {
	request := sonarapi.NewSonarRequest(config)

	request.Method = "GET"
	request.URL.Path = sonarapi.SonarqubeSearchProjectBranchesEndpoint

	branches := ProjectBranches{}
	query := request.URL.Query()
	query.Add("project", project.Name)
	request.URL.RawQuery = query.Encode()

	body, err := sonarapi.DoHttpRequest(request)
	if err != nil {
		return errors.Wrapf(err, "Couldn't get branches of project %s", project.Name)
	}

	if err := json.Unmarshal(*body, &branches); err != nil {
		return errors.Wrapf(err, "Couldn't parse response body %s", project.Name)
	}
	project.Branches = branches
	return nil
}

func (project *SonarProject) GetProjectAnalysis(config sonarapi.QueryConfig) error {
	request := sonarapi.NewSonarRequest(config)

	request.Method = "GET"
	request.URL.Path = sonarapi.SonarqubeSearchProjectAnalysisEndpoint

	analyses := ProjectAnalyses{}
	query := request.URL.Query()
	query.Add("project", project.Name)
	request.URL.RawQuery = query.Encode()

	body, err := sonarapi.DoHttpRequest(request)
	if err != nil {
		return errors.Wrapf(err, "Couldn't get analisys of project %s", project.Name)
	}

	if err := json.Unmarshal(*body, &analyses); err != nil {
		return errors.Wrapf(err, "Couldn't parse response body %s", project.Name)
	}
	project.Analyses = analyses
	return nil
}

func (project *SonarProject) CleanBranches(config sonarapi.QueryConfig) error {
	for _, branch := range project.Branches.ProjectBranches {
		request := sonarapi.NewSonarRequest(config)
		request.Method = "POST"
		request.URL.Path = sonarapi.SonarqubeDeleteProjectBranchesEndpoint

		query := request.URL.Query()
		query.Add("project", project.Name)

		// Do not master branch
		if branch.Name != "master" {
			query.Add("branch", branch.Name)
		} else {
			continue
		}
		request.URL.RawQuery = query.Encode()

		_, err := sonarapi.DoHttpRequest(request)
		if err != nil {
			return errors.Wrapf(err, "Couldn't delete branches of project %s", project.Name)
		}
		log.WithFields(log.Fields{"prefix": project.Name, "type": "branch"}).Debugf("deleted %s", branch.Name)
	}
	log.Infof("Cleaned up branches of project %s", project.Name)
	return nil
}

func (project *SonarProject) CleanAnalyses(config sonarapi.QueryConfig) error {
	for i, analysis := range project.Analyses.ProjectAnalyses {
		// Skip first analyses, because we can't delete all of them and it's the most fresh one.
		if i == 0 {
			continue
		}

		request := sonarapi.NewSonarRequest(config)
		request.Method = "POST"
		request.URL.Path = sonarapi.SonarqubeDeleteProjectAnalysesEndpoint

		query := request.URL.Query()
		query.Add("analysis", analysis.Id)
		request.URL.RawQuery = query.Encode()

		_, err := sonarapi.DoHttpRequest(request)
		if err != nil {
			return errors.Wrapf(err, "ouldn't delete analisys of project %s", project.Name)
		}
		log.WithFields(log.Fields{"prefix": project.Name, "type": "analysis"}).Debug("deleted %s", analysis.Id)
	}
	log.Infof("Cleaned up analisys of project %s", project.Name)
	return nil
}

func makeProjectsList(config sonarapi.QueryConfig) (*SonarProjects, error) {
	foundSonarProjects, err := GetOutdatedProjects(config.KeepDays, config)
	cleanableSonarProjects := SonarProjects{}
	if err != nil {
		return nil, err
	}

	for _, project := range foundSonarProjects.Projects {
		// Get project branches
		if err := project.GetProjectBranches(config); err != nil {
			return nil, err
		}
		// Get project analysis
		if err := project.GetProjectAnalysis(config); err != nil {
			return nil, err
		}
		// Process project if it has not only default branch and more then 1 analyses
		if len(project.Branches.ProjectBranches) > 1 || len(project.Analyses.ProjectAnalyses) > 1 {
			cleanableSonarProjects.Projects = append(cleanableSonarProjects.Projects, project)
			log.Debugf("%s project will be cleaned up", project.Name)
		}

	}
	return &cleanableSonarProjects, nil
}

func CleanProjects(config sonarapi.QueryConfig) error {
	sonarProjects, err := makeProjectsList(config)
	for _, project := range sonarProjects.Projects {
		if err := project.CleanAnalyses(config); err != nil {
			log.Error(err.Error())
		}
		if err := project.CleanBranches(config); err != nil {
			log.Error(err.Error())
		}
	}
	if err != nil {
		log.Fatalf(err.Error())
	}

	return nil
}

func ShowProjects(config sonarapi.QueryConfig) error {
	sonarProjects, err := makeProjectsList(config)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Infof("List of outdated projects: %s", sonarProjects.Projects)
	return nil
}
