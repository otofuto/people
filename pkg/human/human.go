package human

import (
	"log"
	"github.com/otofuto/people/pkg/database"
)

type Human struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
	Tension int `json:"tension"`
}

type StringData struct {
	Human int `json:"human"`
	Data string `json:"data"`
	Feel int `json:"feel"`
	Equal string `json:"equal"`
	Group []string `json:"group"`
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

	sql := "insert into `string_data` (`human`, `data`, `feel`) select * from (?, ?, ?) as tmp where not exists (select * from `string_data` where `human` = ? and `data` = ?) limit 1"
	ins, err := db.Prepare(sql)
	if err != nil {
		log.Println("human.go (sd *StringData) Insert()")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer ins.Close()

	_, err = ins.Exec(&sd.Human, &sd.Data, &sd.Feel, &sd.Human, &sd.Data)
	if err != nil {
		log.Println("human.go (sd *StringData) Insert()")
		log.Println(err)
		log.Println(sql)
		return err
	}

	return nil
}

func (sd *StringData) AddParam(param, val string) error {
	db := database.Connect()
	defer db.Close()

	sql := "insert into `string_data_than` (`human`, `data`, `" + param + "`) select * from (?, ?, ?) as tmp where not exists (select * from `string_data_than` where `human` = ? and `data` = ? and `" + param + "` = ?) limit 1"
	ins, err := db.Prepare(sql)
	if err != nil {
		log.Println("human.go (sd *StringData) AddParam(param, val string)")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer ins.Close()

	_, err = ins.Exec(&sd.Human, &sd.Data, &val, &sd.Human, &sd.Data, &val)
	if err != nil {
		log.Println("human.go (sd *StringData) AddParam(param, val string)")
		log.Println(err)
		log.Println(sql)
		return err
	}

	return nil
}