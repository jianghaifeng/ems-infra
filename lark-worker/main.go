package main

import (
	"time"
)

func main() {
	loadConfig()

	for {
		loadStatus()
		dat := retrieveData()
		refreshToken()
		refreshCurBiTable()
		pushDataToLark(dat)
		saveStatus()
		time.Sleep(time.Second * time.Duration(getConfigUInt("ems.interval")))
	}
}

/*
1. config
	ems_url
	app_id
	app_secret
	template_id
	folder_id
	content_file_id
2. status
	cursor
	bitable_name
	bitable_id
	table_id
*/
