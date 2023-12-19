package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type Commit struct {
	Id                     uint   `json:"id"`
	ChangeId               string `json:"change_id"`
	Project                string `json:"project"`
	Branch                 string `json:"branch"`
	Commiter               string `json:"commiter"`
	Insertions             uint   `json:"insertions"`
	Deletions              uint   `json:"deletions"`
	TotalCommentCount      uint   `json:"total_comment"`
	UnresolvedCommentCount uint   `json:"unresolved_comment"`
	ScanCommentCount       uint   `json:"ci_comment"`
	JiraNumber             string `json:"jira_number"`
	ReviewBlockCount       uint   `json:"effective_review"`
	VerifyBlockCount       uint   `json:"effective_verify"`
	TeamName               string `json:"team_name"`
	DepartmentName         string `json:"department_name"`
}

type Item struct {
	Id              uint   `json:"id"`
	ChangeId        string `json:"change_id"`
	Project         string `json:"project"`
	Branch          string `json:"branch"`
	Commiter        string `json:"commiter"`
	ModifiedLines   uint   `json:"modified_lines"`
	TotalComment    uint   `json:"total_comment"`
	EffectiveReview uint   `json:"effective_review"`
	TeamName        string `json:"team_name"`
	DepartmentName  string `json:"department_name"`
}

func transform(c Commit) Item {
	var res Item
	res.Id = c.Id
	res.ChangeId = c.ChangeId
	res.Project = c.Project
	res.Branch = c.Branch
	res.Commiter = c.Commiter
	res.ModifiedLines = c.Insertions + c.Deletions
	res.TotalComment = c.TotalCommentCount - c.ScanCommentCount
	res.EffectiveReview = 0
	if c.ReviewBlockCount > 0 {
		res.EffectiveReview = 1
	}
	res.TeamName = c.TeamName
	res.DepartmentName = c.DepartmentName
	return res
}

func retrieveData() []Item {
	client := &http.Client{}
	url := getConfigString("ems.url") + "?from=" + strconv.FormatUint(getStatus().Cursor, 10) +
		"&pageSize=" + strconv.FormatUint(getConfigUInt("ems.pageSize"), 10)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	var commits []Commit
	items := make([]Item, 0)

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("cannot connect to ems: ", err)
		return items
	}
	if res.StatusCode != http.StatusOK {
		panic(errors.New("status=" + res.Status))
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(body, &commits)
	for _, c := range commits {
		items = append(items, transform(c))
	}
	return items
}
