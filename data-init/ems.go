package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type deptRequestBody struct {
	Name string `json:"name"`
}

type teamRequestBody struct {
	Name         string `json:"name"`
	DepartmentId int    `json:"departmentId"`
}

type responseBody struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type repoRequestBody struct {
	Project string `json:"project"`
	Branch  string `json:"branch"`
	TeamId  int    `json:"team_id"`
}

func createDepartment(name string) int {
	url := getConfigString("ems.url") + "/departments"
	requestBody := deptRequestBody{
		Name: name,
	}
	payload, err := json.Marshal(requestBody)
	if err != nil {
		log.Panic(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		log.Panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()

	var respBody responseBody
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}
	json.Unmarshal(body, &respBody)
	return respBody.Id
}

func createTeam(name string, dept int) int {
	if len(name) == 0 {
		return 0
	}
	url := getConfigString("ems.url") + "/teams"
	requestBody := teamRequestBody{
		Name:         name,
		DepartmentId: dept,
	}
	payload, err := json.Marshal(requestBody)
	if err != nil {
		log.Panic(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		log.Panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()

	var respBody responseBody
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}
	json.Unmarshal(body, &respBody)
	return respBody.Id
}

func createRepo(proj string, branch string, team int) int {
	url := getConfigString("ems.url") + "/repos"
	requestBody := repoRequestBody{
		Project: proj,
		Branch:  branch,
		TeamId:  team,
	}
	payload, err := json.Marshal(requestBody)
	if err != nil {
		log.Panic(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		log.Panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()

	var respBody responseBody
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}
	json.Unmarshal(body, &respBody)
	return respBody.Id
}
