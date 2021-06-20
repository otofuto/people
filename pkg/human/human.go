package human

import (
	"log"
	"github.com/otofuto/human/pkg/database"
)

type StringData struct {
	Human int `json:"human"`
	Data string `json:"data"`
	Tension int `json:"tension"`
	Equal string `json:"equal"`
	Big []string `json:"big"`
	Not []string `json:"not"`
	Then []string `json:"then"`
	InCount int `json:"in_count"`
	OutCount int `json:"out_count"`
}

type IntData struct {
	Human int `json:"human"`
	Data int `json:"data"`
	Equal int `json:"equal"`
	Big []int `json:"big"`
	Not []int `json:"not"`
	Then []int `json:"then"`
	InCount int `json:"in_count"`
	OutCount int `json:"out_count"`
}

func (sd *StringData) Insert() error {
	db := database.Connect()
	defer db.Close()

	sql := "insert into `string_data` (`human`, `data`, `tension`) select * from (?, ?, ?) as tmp where not exists (select * from `string_data` where `human` = ? and `data` = ?) limit 1"
	ins, err := db.Prepare(sql)
	if err != nil {
		log.Println("human.go (sd *StringData) Insert()")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer ins.Close()

	_, err := ins.Exec(&sd.Human, &sd.Data, &sd.Tension, &sd.Human, &sd.Data)
	if err != nil {
		log.Println("human.go (sd *StringData) Insert()")
		log.Println(err)
		log.Println(sql)
		return err
	}

	return nil
}

func (sd *StringData) Big(big string) err {
	db := database.Connect()
	defer db.Close()

	sql := "insert into `string_data_than` (`data`, `big`) select * from (?, ?) as tmp where not exists (select * from `string_data_than` where `data` = ? and `big` = ?) limit 1"
	ins, err := db.Prepare(sql)
	if err != nil {
		log.Println("human.go (sd *StringData) Big(big string)")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer ins.Close()

	_, err := ins.Exec(&sd.Data, &big, &sd.Data, &big)
	if err != nil {
		log.Println("human.go (sd *StringData) Big(big string)")
		log.Println(err)
		log.Println(sql)
		return err
	}

	return nil
}