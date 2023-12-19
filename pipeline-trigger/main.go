package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RefUpdate struct {
	OldRev  string `json:"oldRev"`
	NewRev  string `json:"newRev"`
	RefName string `json:"refName"`
	Project string `json:"project"`
}

type GerritTrigger struct {
	Data RefUpdate `json:"refUpdate"`
	Type string    `json:"type"`
}

var repoMap map[string]string = make(map[string]string)
var hookMap map[string]string = make(map[string]string)

func registerRepos(repo string, branch string, hook string) {
	key := repo + "," + branch
	repoMap[key] = ""
	hookMap[key] = hook
	fmt.Printf("Registered:%s %s %s\n", repo, branch, hook)
}

func postTrigger(c *gin.Context) {
	var trigger GerritTrigger
	if err := c.BindJSON(&trigger); err != nil {
		fmt.Println("err binding:", err)
	}
	if trigger.Type != "ref-updated" {
		fmt.Println("not a ref-update event, ignored.")
		return
	}
	key := trigger.Data.Project + "," + trigger.Data.RefName
	curRev, ok := repoMap[key]
	if ok {
		fmt.Println("curRev:", curRev, " newRev:", trigger.Data.NewRev)
		if curRev != trigger.Data.NewRev {
			fmt.Println("going to trigger a new build...")
			if resp, err := http.Get(hookMap[key]); err != nil {
				fmt.Println("Hook Jenkins failed, ", err)
			} else {
				fmt.Println(resp.StatusCode)
			}
		}
		repoMap[key] = trigger.Data.NewRev
	}
}

func main() {
	registerRepos("/ems", "refs/heads/master", "http://ems-jenkins-service:8080/ems-jenkins/job/ems-build/build?token=test")
	registerRepos("/ems-team", "refs/heads/master", "http://ems-jenkins-service:8080/ems-jenkins/job/ems-team/build?token=p3CGMYzyvaMcJPq9gBjxgquaF6nM5kPN")
	registerRepos("/ems-team-worker", "refs/heads/master", "http://ems-jenkins-service:8080/ems-jenkins/job/ems-team-worker/build?token=p3CGMYzyvaMcJPq9gBjxgquaF6nM5kPN")
	router := gin.Default()
	router.POST("/trigger", postTrigger)

	router.Run(":8080")
}
