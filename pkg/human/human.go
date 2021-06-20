package human

import (
	"log"
	"errors"
	"strconv"
	"strings"
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
	Then []string `json:"then"`
	InCount int `json:"in_count"`
	OutCount int `json:"out_count"`
}

func (sd *StringData) Insert() error {
	db := database.Connect()
	defer db.Close()

	sql := "insert into `string_data` (`human`, `data`, `feel`) select * from (select " + strconv.Itoa(sd.Human) + ", '" + database.Escape(sd.Data) + "', " + strconv.Itoa(sd.Feel) + ") as tmp where not exists (select * from `string_data` where `human` = ? and `data` = ?) limit 1"
	ins, err := db.Prepare(sql)
	if err != nil {
		log.Println("human.go (sd *StringData) Insert()")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer ins.Close()

	_, err = ins.Exec(&sd.Human, &sd.Data)
	if err != nil {
		log.Println("human.go (sd *StringData) Insert()")
		log.Println(err)
		log.Println(sql)
		return err
	}

	return nil
}

func (sd *StringData) AddParam(param, val string) error {
	if param != "group" && param != "big" && param != "not" && param != "then" {
		return errors.New("Invalid param name")
	}

	db := database.Connect()
	defer db.Close()

	sql := "insert into `string_data_than` (`human`, `data`, `" + param + "`) select * from (select " + strconv.Itoa(sd.Human) + ", '" + database.Escape(sd.Data) + "', '" + database.Escape(val) + "') as tmp where not exists (select * from `string_data_than` where `human` = ? and `data` = ? and `" + param + "` = ?) limit 1"
	ins, err := db.Prepare(sql)
	if err != nil {
		log.Println("human.go (sd *StringData) AddParam(param, val string)")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer ins.Close()

	_, err = ins.Exec(&sd.Human, &sd.Data, &val)
	if err != nil {
		log.Println("human.go (sd *StringData) AddParam(param, val string)")
		log.Println(err)
		log.Println(sql)
		return err
	}

	return nil
}

func (id *IntData) Insert() error {
	db := database.Connect()
	defer db.Close()

	sql := "insert into `int_data` (`human`, `data`) select * from (select " + strconv.Itoa(id.Human) + ", " + strconv.Itoa(id.Data) + ") as tmp where not exists (select * from `int_data` where `human` = ? and `data` = ?) limit 1"
	ins, err := db.Prepare(sql)
	if err != nil {
		log.Println("human.go (id *IntData) Insert()")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer ins.Close()

	_, err = ins.Exec(&id.Human, &id.Data)
	if err != nil {
		log.Println("human.go (id *IntData) Insert()")
		log.Println(err)
		log.Println(sql)
		return err
	}

	return nil
}

func (id *IntData) AddThen(then string) error {
	db := database.Connect()
	defer db.Close()

	sql := "insert into `int_data_than` (`human`, `data`, `then`) select * from (select " + strconv.Itoa(id.Human) + ", " + strconv.Itoa(id.Data) + ", '" + database.Escape(then) + "') as tmp where not exists (select * from `int_data_than` where `human` = ? and `data` = ? and `then` = ?) limit 1"
	ins, err := db.Prepare(sql)
	if err != nil {
		log.Println("human.go (id *IntData) AddThen(then string)")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer ins.Close()

	_, err = ins.Exec(&id.Human, &id.Data, &then)
	if err != nil {
		log.Println("human.go (id *IntData) AddThen(then string)")
		log.Println(err)
		log.Println(sql)
		return err
	}

	return nil
}

func LangSplit(text string) []string {
	ret1 := make([]string, 0)
	ml := []string { " ", "　", "\n", "、", "。", ":", "：", "=", "<", ">", "+", "?", "？", "#", "!", "！", "・", "…" }
	word := ""
	for _, s := range strings.Split(text, "") {
		if isMatch(s, ml) {
			ret1 = append(ret1, word)
			word = ""
		} else {
			word += s
		}
	}
	if word != "" {
		ret1 = append(ret1, word)
	}
	ret2 := make([]string, 0)
	wd := ""
	for _, word = range ret1 {
		wd = ""
		lastInt := false
		lastHarf := false
		lastHiragana := false
		lastKatakana := false
		lastKanji := false
		for _, r := range word {
			if isNumber(r) {
				if lastInt {
					wd += string(r)
				} else {
					if wd != "" {
						ret2 = append(ret2, wd)
					}
					wd = string(r)
					lastInt = true
					lastHarf = false
					lastHiragana = false
					lastKatakana = false
					lastKanji= false
				}
			} else if r <= 255 {
				if lastHarf {
					wd += string(r)
				} else {
					if wd != "" {
						ret2 = append(ret2, wd)
					}
					wd = string(r)
					lastInt = false
					lastHarf = true
					lastHiragana = false
					lastKatakana = false
					lastKanji= false
				}
			} else if isHiragana(r) || (lastHiragana && r == 12540) {
				if lastHiragana {
					wd += string(r)
				} else {
					if wd != "" {
						ret2 = append(ret2, wd)
					}
					wd = string(r)
					lastInt = false
					lastHarf = false
					lastHiragana = true
					lastKatakana = false
					lastKanji= false
				}
			} else if isKatakana(r) || (lastKatakana && r == 12540) {
				if lastKatakana {
					wd += string(r)
				} else {
					if wd != "" {
						ret2 = append(ret2, wd)
					}
					wd = string(r)
					lastInt = false
					lastHarf = false
					lastHiragana = false
					lastKatakana = true
					lastKanji= false
				}
			} else if isKanji(r) {
				if lastKanji {
					wd += string(r)
				} else {
					if wd != "" {
						ret2 = append(ret2, wd)
					}
					wd = string(r)
					lastInt = false
					lastHarf = false
					lastHiragana = false
					lastKatakana = false
					lastKanji= true
				}
			} else {
				if wd != "" && (lastHarf || lastHiragana || lastKatakana || lastKanji) {
					ret2 = append(ret2, wd)
					wd = ""
				}
				wd += string(r)
				lastInt = false
				lastHarf = false
				lastHiragana = false
				lastKatakana = false
				lastKanji = false
			}
		}
		if wd != "" {
			ret2 = append(ret2, wd)
		}
	}
	ret1 = make([]string, 0)
	wd = ""
	for i, word := range ret2 {
		ret1 = append(ret1, word)
		current := []rune(word)[0]
		if i == 0 {
			continue
		}
		var count int
		for count = 1; ret1[i - count] == "" && i > count; count++ {}
		last := []rune(ret1[i - count])[0]
		if (isHiragana(current) && isKanji(last)) ||
			(isKatakana(current) && isKanji(last)) ||
			(isKanji(current) && isNumber(last)) ||
			(current <= 255 && last <= 255) {
			ret1[i - 1] = ret1[i - 1] + word
			ret1[i] = ""
		}
	}
	ret2 = make([]string, 0)
	for _, word = range ret1 {
		if word != "" {
			ret2 = append(ret2, word)
		}
	}
	return ret2
}

func isMatch(s string, match []string) bool {
	for i := 0; i < len(match); i++ {
		if match[i] == s {
			return true
		}
	}
	return false
}

func isNumber(r rune) bool {
	return (48 <= r && r <= 57) || r == 44 || r == 46
}

func isHiragana(r rune) bool {
	return 12353 <= r && r <= 12446
}

func isKatakana(r rune) bool {
	return 12449 <= r && r <= 12538
}

func isKanji(r rune) bool {
	return 19968 <= r && r <= 40879
}