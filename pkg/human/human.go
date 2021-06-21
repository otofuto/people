package human

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/otofuto/people/pkg/database"
)

type Human struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Tension int    `json:"tension"`
}

type StringData struct {
	Human    int      `json:"human"`
	Data     string   `json:"data"`
	Feel     int      `json:"feel"`
	Equal    string   `json:"equal"`
	Group    []string `json:"group"`
	Big      []string `json:"big"`
	Not      []string `json:"not"`
	Then     []string `json:"then"`
	Next     []string `json:"next"`
	Action   []int    `json:"action"`
	InCount  int      `json:"in_count"`
	OutCount int      `json:"out_count"`
}

type IntData struct {
	Human    int      `json:"human"`
	Data     int      `json:"data"`
	Then     []string `json:"then"`
	InCount  int      `json:"in_count"`
	OutCount int      `json:"out_count"`
}

func (sd *StringData) Insert() error {
	if sd.Human == 0 {
		log.Println("human.go (sd *StringData) Insert()")
		return errors.New("sd.Human is not set")
	}

	db := database.Connect()
	defer db.Close()

	sql := "select * from `string_data` where `human` = " + strconv.Itoa(sd.Human) + " and `data` = '" + database.Escape(sd.Data) + "'"
	rows, err := db.Query(sql)
	if err != nil {
		log.Println("human.go (sd *StringData) Insert()")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer rows.Close()
	if rows.Next() {
		sql = "update `string_data` set `in_count` = `in_count` + 1 where `human` = ? and `data` = ?"
		upd, err := db.Prepare(sql)
		if err != nil {
			log.Println("human.go (sd *StringData) Insert()")
			log.Println(err)
			log.Println(sql)
			return err
		}
		defer upd.Close()
		_, err = upd.Exec(&sd.Human, &sd.Data)
		if err != nil {
			log.Println("human.go (sd *StringData) Insert()")
			log.Println(err)
			log.Println(sql)
			return err
		}
	} else {
		sql = "insert into `string_data` (`human`, `data`, `feel`) values (?, ?, ?)"
		ins, err := db.Prepare(sql)
		if err != nil {
			log.Println("human.go (sd *StringData) Insert()")
			log.Println(err)
			log.Println(sql)
			return err
		}
		defer ins.Close()

		_, err = ins.Exec(&sd.Human, &sd.Data, &sd.Feel)
		if err != nil {
			log.Println("human.go (sd *StringData) Insert()")
			log.Println(err)
			log.Println(sql)
			return err
		}
	}

	return nil
}

func (sd *StringData) Get() error {
	if sd.Human == 0 {
		log.Println("human.go (sd *StringData) Get()")
		return errors.New("sd.Human is not set")
	}
	if sd.Data == "" {
		log.Println("human.go (sd *StringData) Get()")
		return errors.New("sd.Data is empty")
	}

	db := database.Connect()
	defer db.Close()

	sql := "select `feel`, `equal`, `in_count`, `out_count`, `param`, `val` from `string_data` left outer join `string_data_than` on `string_data`.`data` = `string_data_than`.`data` where `string_data`.`human` = " + strconv.Itoa(sd.Human) + " and `string_data`.`data` = '" + database.Escape(sd.Data) + "'"
	rows, err := db.Query(sql)
	if err != nil {
		log.Println("human.go (sd *StringData) Get()")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var param string
		var val string
		err = rows.Scan(&sd.Feel, &sd.Equal, &sd.InCount, &sd.OutCount, &param, &val)
		if err != nil {
			log.Println("human.go (sd *StringData) Get()")
			log.Println(err)
			return err
		}
		if param == "group" {
			sd.Group = append(sd.Group, val)
		}
		if param == "big" {
			sd.Big = append(sd.Big, val)
		}
		if param == "not" {
			sd.Not = append(sd.Not, val)
		}
		if param == "then" {
			sd.Then = append(sd.Then, val)
		}
		if param == "next" {
			sd.Next = append(sd.Next, val)
		}
		if param == "action" {
			act, err := strconv.Atoi(val)
			if err == nil {
				sd.Action = append(sd.Action, act)
			}
		}
	}
	return nil
}

func (sd *StringData) Update() error {
	if sd.Human == 0 {
		log.Println("human.go (sd *StringData) Update()")
		return errors.New("sd.Human is not set")
	}
	if sd.Data == "" {
		log.Println("human.go (sd *StringData) Update()")
		return errors.New("sd.Data is empty")
	}

	db := database.Connect()
	defer db.Close()

	sql := "update `string_data` set `feel` = ?, `equal` = ?, `in_count` = ?, `out_count` = ? where `human` = ? and `data` = ?"
	upd, err := db.Prepare(sql)
	if err != nil {
		log.Println("human.go (sd *StringData) Update()")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer upd.Close()
	_, err = upd.Exec(&sd.Feel, &sd.Equal, &sd.InCount, &sd.OutCount, &sd.Human, &sd.Data)
	if err != nil {
		log.Println("human.go (sd *StringData) Update()")
		log.Println(err)
		log.Println(sql)
		return err
	}
	return nil
}

func (sd *StringData) AddParam(param, val string) error {
	if sd.Human == 0 {
		log.Println("human.go (sd *StringData) AddParam(param, val string)")
		return errors.New("sd.Human is not set")
	}

	if param != "equal" && param != "group" && param != "big" && param != "not" && param != "then" && param != "next" && param != "action" {
		return errors.New("Invalid param name")
	}

	if param == "next" {
		return sd.AddNext(val)
	}
	if param == "action" {
		valint, err := strconv.Atoi(val)
		if err != nil {
			log.Println("human.go (sd *StringData) AddParam(param, val string)")
			log.Println(err)
			return err
		}
		return sd.AddAction(valint)
	}
	if param == "equal" {
		if err := sd.Get(); err != nil {
			return err
		}
		sd.Equal = val
		return sd.Update()
	}

	db := database.Connect()
	defer db.Close()

	sql := "insert into `string_data_than` (`human`, `data`, `param`, `val`) select * from (select " + strconv.Itoa(sd.Human) + ", '" + database.Escape(sd.Data) + "', '" + database.Escape(param) + "', '" + database.Escape(val) + "') as tmp where not exists (select * from `string_data_than` where `human` = ? and `data` = ? and `param` = ? and `val` = ?) limit 1"
	ins, err := db.Prepare(sql)
	if err != nil {
		log.Println("human.go (sd *StringData) AddParam(param, val string)")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer ins.Close()

	_, err = ins.Exec(&sd.Human, &sd.Data, &param, &val)
	if err != nil {
		log.Println("human.go (sd *StringData) AddParam(param, val string)")
		log.Println(err)
		log.Println(sql)
		return err
	}

	return nil
}

func (sd *StringData) AddNext(val string) error {
	if sd.Human == 0 {
		log.Println("human.go (sd *StringData) AddNext(val string)")
		return errors.New("sd.Human is not set")
	}

	db := database.Connect()
	defer db.Close()

	sql := "insert into `string_data_than` (`human`, `data`, `param`, `val`) values (?, ?, 'next', ?)"
	ins, err := db.Prepare(sql)
	if err != nil {
		log.Println("human.go (sd *StringData) AddNext(val string)")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer ins.Close()

	_, err = ins.Exec(&sd.Human, &sd.Data, &val)
	if err != nil {
		log.Println("human.go (sd *StringData) AddNext(val string)")
		log.Println(err)
		log.Println(sql)
		return err
	}

	return nil
}

func (sd *StringData) AddAction(val int) error {
	if sd.Human == 0 {
		log.Println("human.go (sd *StringData) AddAction(val int)")
		return errors.New("sd.Human is not set")
	}

	db := database.Connect()
	defer db.Close()

	sql := "insert into `string_data_than` (`human`, `data`, `param`, `val`) values (?, ?, 'action', ?)"
	ins, err := db.Prepare(sql)
	if err != nil {
		log.Println("human.go (sd *StringData) AddAction(val int)")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer ins.Close()

	_, err = ins.Exec(&sd.Human, &sd.Data, strconv.Itoa(val))
	if err != nil {
		log.Println("human.go (sd *StringData) AddAction(val int)")
		log.Println(err)
		log.Println(sql)
		return err
	}

	return nil
}

func (id *IntData) Insert() error {
	if id.Human == 0 {
		log.Println("human.go (id *IntData) Insert()")
		return errors.New("sd.Human is not set")
	}

	db := database.Connect()
	defer db.Close()

	sql := "select * from `int_data` where `human` = " + strconv.Itoa(id.Human) + " and `data` = " + strconv.Itoa(id.Data)
	rows, err := db.Query(sql)
	if err != nil {
		log.Println("human.go (id *IntData) Insert()")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return nil
	}
	sql = "insert into `int_data` (`human`, `data`) values (?, ?)"
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
	if id.Human == 0 {
		log.Println("human.go (id *IntData) AddThen(then string)")
		return errors.New("id.Human is not set")
	}

	db := database.Connect()
	defer db.Close()

	sql := "select * from `int_data_than` where `human` = " + strconv.Itoa(id.Human) + " and `data` = " + strconv.Itoa(id.Data) + " and `param` = 'then' and `val` = '" + database.Escape(then) + "'"
	rows, err := db.Query(sql)
	if err != nil {
		log.Println("human.go (id *IntData) AddThen(then string)")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return nil
	}
	sql = "insert into `int_data_than` (`human`, `data`, `param`, `val`) values (?, ?, 'then', ?)"
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

func LangSplit(text string) []StringData {
	ret1 := make([]string, 0)
	ml := []string{" ", "　", "\r", "\n", "、", "。", ":", "：", "=", "<", ">", "+", "?", "？", "#", "!", "！", "・", "…"}
	word := ""
	for _, s := range strings.Split(text, "") {
		if isMatch(s, ml) {
			ret1 = append(ret1, word)
			ret1 = append(ret1, s)
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
					lastKanji = false
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
					lastKanji = false
				}
			} else if isHiragana(r) || (lastHiragana && isHirakata(r)) {
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
					lastKanji = false
				}
			} else if isKatakana(r) || (lastKatakana && isHirakata(r)) {
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
					lastKanji = false
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
					lastKanji = true
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
		for count = 1; ret1[i-count] == "" && i > count; count++ {
		}
		last := []rune(ret1[i-count])[0]
		if ((isHiragana(current) || isHirakata(current)) && isKanji(last)) ||
			((isKatakana(current) || isHirakata(current)) && isKanji(last)) ||
			(isKanji(current) && isNumber(last)) ||
			(current <= 255 && last <= 255) {
			ret1[i-1] = ret1[i-1] + word
			ret1[i] = ""
		}
	}
	sds := make([]StringData, 0)
	ml2 := []string{"\n", "、", "。", "!", "！", "…", "?", "？"}
	for _, word = range ret1 {
		if strings.TrimSpace(word) != "" {
			if isMatch(word, ml2) {
				act := []rune(word)[0]
				for i := len(sds) - 1; i >= 0 && sds[i].Action[0] == 0; i-- {
					sds[i].Action[0] = int(act)
				}
			}
			if !isMatch(word, ml) {
				sds = append(sds, StringData{
					Data:   strings.TrimSpace(word),
					Action: []int{0},
				})
			}
		}
	}
	return sds
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

func isInt(r rune) bool {
	return (48 <= r && r <= 57)
}

func isHiragana(r rune) bool {
	return (12353 <= r && r < 12441) || (12444 < r && r <= 12446)
}

func isKatakana(r rune) bool {
	return 12449 <= r && r <= 12538
}

func isHirakata(r rune) bool {
	return r == 12540 || (12441 <= r && r <= 12444)
}

func isKanji(r rune) bool {
	return (19968 <= r && r <= 40879) || r == 12293
}

func SaveWords(hid int, sds []StringData) error {
	if hid == 0 {
		log.Println("human.go SaveWords(hid int, sds []StringData)")
		return errors.New("hid is not set")
	}

	last := StringData{
		Human: hid,
		Data:  "",
	}
	ago2 := StringData{
		Human: hid,
		Data:  "",
	}
	for _, sd := range sds {
		sd.Human = hid
		if num, err := strconv.Atoi(sd.Data); err != nil {
			if err = sd.Insert(); err != nil {
				return err
			}
			if isNumber([]rune(sd.Data)[0]) {
				numstr := ""
				for _, r := range sd.Data {
					if isInt(r) {
						numstr += string(r)
					} else {
						break
					}
				}
				num, _ = strconv.Atoi(numstr)
				id := IntData{
					Human: hid,
					Data:  num,
				}
				if err = id.Insert(); err != nil {
					return err
				}
				if err = id.AddThen(sd.Data[len(numstr):]); err != nil {
					return err
				}
			}
		} else {
			id := IntData{
				Human: hid,
				Data:  num,
			}
			if err := id.Insert(); err != nil {
				return err
			}
		}
		for i := 0; i < len(sd.Action); i++ {
			if err := sd.AddAction(sd.Action[i]); err != nil {
				return err
			}
		}
		if last.Data != "" {
			if err := last.AddNext(sd.Data); err != nil {
				return err
			}
		}
		if last.Data == "は" && ago2.Data != "" && notQuestion(last) {
			if err := ago2.AddParam("equal", sd.Data); err != nil {
				return err
			}
			if err := sd.AddParam("equal", ago2.Data); err != nil {
				return err
			}
			if err := ago2.AddParam("group", sd.Data); err != nil {
				return err
			}
		}
		if last.Data == "の" && ago2.Data != "" {
			if err := ago2.AddParam("group", sd.Data); err != nil {
				return err
			}
		}
		if (strings.HasSuffix(last.Data, "なら") ||
			strings.HasSuffix(last.Data, "たら") ||
			last.Data == "で") &&
			ago2.Data != "" && notQuestion(last) {
			str := ""
			if strings.HasSuffix(last.Data, "なら") {
				str = last.Data[:strings.LastIndex(last.Data, "なら")]
			} else if strings.HasSuffix(last.Data, "たら") {
				str = last.Data[:strings.LastIndex(last.Data, "たら")]
			}
			if str != "" {
				newword := StringData{
					Human:  hid,
					Data:   str,
					Action: last.Action,
				}
				if err := newword.Insert(); err != nil {
					return err
				}
				for i := 0; i < len(newword.Action); i++ {
					if err := newword.AddAction(newword.Action[i]); err != nil {
						return err
					}
				}
				if err := newword.AddParam("then", sd.Data); err != nil {
					return err
				}
			} else {
				if err := ago2.AddParam("then", sd.Data); err != nil {
					return err
				}
			}
		}
		if strings.HasSuffix(sd.Data, "ない") && notQuestion(sd) && ago2.Data != "" {
			str := sd.Data[:strings.LastIndex(sd.Data, "ない")]
			if str != "" {
				newword := StringData{
					Human:  hid,
					Data:   str,
					Action: sd.Action,
				}
				if err := newword.Insert(); err != nil {
					return err
				}
				for i := 0; i < len(newword.Action); i++ {
					if err := newword.AddAction(newword.Action[i]); err != nil {
						return err
					}
				}
				if err := last.AddParam("not", str); err != nil {
					return err
				}
			} else {
				if err := ago2.AddParam("not", last.Data); err != nil {
					return err
				}
			}
		}
		ago2 = last
		last = sd
	}
	return nil
}

func contains(arr []int, t int) bool {
	for _, a := range arr {
		if a == t {
			return true
		}
	}
	return false
}

func notQuestion(sd StringData) bool {
	for _, a := range sd.Action {
		if a == 65311 || a == 63 {
			return false
		}
	}
	return true
}
