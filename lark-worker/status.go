package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
)

type Status struct {
	Cursor        uint64 `json:"-"`
	CursorChanged bool   `json:"-"`
	BitableName   string `json:"bitable_name"`
	BitableId     string `json:"bitable_id"`
	TableId       string `json:"table_id"`
	Changed       bool   `json:"-"`
}

var curStatus Status

func loadStatus() {
	dat, err := os.ReadFile("./data/status.json")
	if err != nil {
		panic(err)
	}

	curStatus.BitableName = ""
	json.Unmarshal(dat, &curStatus)

	numdat, err := os.ReadFile("./data/cursor")
	if err != nil {
		panic(err)
	}

	num, err := strconv.ParseUint(strings.TrimSpace(string(numdat)), 0, 0)
	if err != nil {
		panic(err)
	}
	curStatus.Changed = false
	curStatus.CursorChanged = false
	curStatus.Cursor = num
}

func saveStatus() {
	if curStatus.Changed {
		dat, err := json.MarshalIndent(curStatus, "", " ")
		if err != nil {
			panic(err)
		}
		log.Println(string(dat))
		err = os.WriteFile("./data/status.json", dat, 0644)
		if err != nil {
			panic(err)
		}
	}
	if curStatus.CursorChanged {
		log.Println("newCursor=", curStatus.Cursor)
		err := os.WriteFile("./data/cursor", []byte(strconv.FormatUint(curStatus.Cursor, 10)), 0644)
		if err != nil {
			panic(err)
		}
	}
}

func getStatus() Status {
	return curStatus
}

func setCursor(c uint64) {
	if curStatus.Cursor == c {
		return
	}
	curStatus.Cursor = c
	curStatus.CursorChanged = true
}

func setBiTable(name string, id string, tableId string) {
	if curStatus.BitableId == id {
		return
	}
	curStatus.BitableName = name
	curStatus.BitableId = id
	curStatus.TableId = tableId
	curStatus.Changed = true
}
