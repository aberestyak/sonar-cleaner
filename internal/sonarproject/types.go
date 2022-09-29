package sonarProject

type SonarProjects struct {
	Projects []SonarProject `json:"components"`
}

type SonarProject struct {
	Name     string `json:"key"`
	Branches ProjectBranches
	Analyses ProjectAnalyses
}

type ProjectBranches struct {
	ProjectBranches []ProjectBranch `json:"branches"`
}

type ProjectAnalyses struct {
	ProjectAnalyses []ProjectAnalysis `json:"analyses"`
}

type ProjectBranch struct {
	Name   string `json:"name"`
	IsMain bool   `json:"isMain"`
}

type ProjectAnalysis struct {
	Id   string `json:"key"`
	Date string `json:"date"`
}

type SonarPagination struct {
	Spec SonarPaginationSpec `json:"paging"`
}

type SonarPaginationSpec struct {
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
	Total     int `json:"total"`
}
