package main

import (
	"log"
	"sync"
)

func processAll() {
	m := getConfig("lark.departments")
	items := m.([]interface{})
	for _, item := range items {
		m := item.(map[string]interface{})
		processOne(m["name"].(string), m["table"].(string))
	}
}

func getResult(bitable string, table string, pagetoken string, c chan recordResponseBody) {
	body := readLark(bitable, table, pagetoken)
	c <- body
}

func processOne(name string, table string) {
	log.Println("start processing:", name)
	deptId := createDepartment(name)

	body := recordResponseBody{Data: recordResponseData{Has_more: true, Page_token: ""}}
	page_num := int(1)
	bitable := getConfigString("lark.bitable")
	teamsMap := make(map[string]int)

	var wg sync.WaitGroup
	for {
		c := make(chan recordResponseBody)
		if body.Data.Has_more {
			go getResult(bitable, table, body.Data.Page_token, c)
		}

		for _, item := range body.Data.Items {
			t := item.Fields["agile Team"]
			if _, exist := teamsMap[t]; !exist {
				teamsMap[t] = createTeam(t, deptId)
			}
		}
		wg.Add(1)
		go func(items []recordResponseItem) {
			defer wg.Done()
			for _, item := range items {
				if len(item.Fields["agile Team"]) == 0 {
					continue
				}
				// fmt.Println(item.Fields["repository"], item.Fields["branch"], teamsMap[item.Fields["agile Team"]])
				createRepo(item.Fields["repository"], item.Fields["branch"], teamsMap[item.Fields["agile Team"]])
			}
		}(body.Data.Items)

		if !body.Data.Has_more {
			break
		}

		body = <-c
		log.Println("got page:", page_num)
		page_num++
	}

	wg.Wait()
	log.Println("done")
}

func main() {
	loadConfig()
	refreshToken()
	processAll()
}
