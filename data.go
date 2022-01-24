package main

import (
	"encoding/json"
	"gwentgo/gwent"
	"log"
	"os"
	"time"
)

//cookie -> data
var userDb = map[string]*gwent.PlayerData{}

func saveData() {
	f, err := os.OpenFile("user_db.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(&userDb)
	if err != nil {
		panic(err)
	}

	if _, err := f.Write(data); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func loadData() {
	f, err := os.Open("user_db.json")
	if err != nil {
		log.Printf("open error", err.Error())
		return
	}
	if err := json.NewDecoder(f).Decode(&userDb); err != nil {
		log.Println(err)
	}
}

func backupRoutine() {
	for {
		save()
		saveData()
		time.Sleep(time.Second * 10)
	}
}
