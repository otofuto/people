package human

import (
	"database/sql"
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
	Got      bool
}

type IntData struct {
	Human    int      `json:"human"`
	Data     int      `json:"data"`
	Then     []string `json:"then"`
	InCount  int      `json:"in_count"`
	OutCount int      `json:"out_count"`
}

type Talk struct {
	Opponent   string `json:"opponent"`
	Human      int    `json:"human"`
	Heard      string `json:"heard"`
	Talked     string `json:"talked"`
	IsQuestion bool   `json:"is_question"`
	TalkedAt   string `json:"talked_at"`
}

func (sd *StringData) Insert(db *sql.DB) error {
	if sd.Human == 0 {
		log.Println("human.go (sd *StringData) Insert()")
		return errors.New("sd.Human is not set")
	}
	if len([]rune(sd.Data)) > 255 {
		log.Println("human.go (sd *StringData) Insert()")
		log.Println("datas length is over 255")
		return nil
	}

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

func (sd *StringData) Get(db *sql.DB) error {
	if sd.Human == 0 {
		log.Println("human.go (sd *StringData) Get()")
		return errors.New("sd.Human is not set")
	}
	if sd.Data == "" {
		log.Println("human.go (sd *StringData) Get()")
		return errors.New("sd.Data is empty")
	}

	if sd.Got {
		return nil
	}

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
	sd.Got = true
	return nil
}

func (sd *StringData) Update(db *sql.DB) error {
	if sd.Human == 0 {
		log.Println("human.go (sd *StringData) Update()")
		return errors.New("sd.Human is not set")
	}
	if sd.Data == "" {
		log.Println("human.go (sd *StringData) Update()")
		return errors.New("sd.Data is empty")
	}

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

func (sd *StringData) AddParam(db *sql.DB, param, val string) error {
	if sd.Human == 0 {
		log.Println("human.go (sd *StringData) AddParam(param, val string)")
		return errors.New("sd.Human is not set")
	}
	if len([]rune(val)) > 255 {
		log.Println("human.go (sd *StringData) AddParam(param, val string)")
		log.Println("datas length is over 255")
		return nil
	}

	if param != "equal" && param != "group" && param != "big" && param != "not" && param != "then" && param != "next" && param != "action" {
		return errors.New("Invalid param name")
	}

	if param == "next" {
		return sd.AddNext(db, val)
	}
	if param == "action" {
		valint, err := strconv.Atoi(val)
		if err != nil {
			log.Println("human.go (sd *StringData) AddParam(param, val string)")
			log.Println(err)
			return err
		}
		return sd.AddAction(db, valint)
	}
	if param == "equal" {
		if err := sd.Get(db); err != nil {
			return err
		}
		sd.Equal = val
		return sd.Update(db)
	}
	if param == "group" && sd.Data == val {
		log.Println("human.go (sd *StringData) AddParam(param, val string)")
		log.Println("cannot add group '" + sd.Data + "' in '" + val + "'")
		return nil
	}

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

func (sd *StringData) AddNext(db *sql.DB, val string) error {
	if sd.Human == 0 {
		log.Println("human.go (sd *StringData) AddNext(val string)")
		return errors.New("sd.Human is not set")
	}
	if len([]rune(val)) > 255 {
		log.Println("human.go (sd *StringData) AddNext(val string)")
		log.Println("datas length is over 255")
		return nil
	}

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

func (sd *StringData) AddAction(db *sql.DB, val int) error {
	if sd.Human == 0 {
		log.Println("human.go (sd *StringData) AddAction(val int)")
		return errors.New("sd.Human is not set")
	}

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

func (id *IntData) Insert(db *sql.DB) error {
	if id.Human == 0 {
		log.Println("human.go (id *IntData) Insert()")
		return errors.New("sd.Human is not set")
	}
	if id.Data > 2147483647 {
		log.Println("human.go (id *IntData) Insert()")
		log.Println("data size is over max of int")
		return nil
	}

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

func (id *IntData) AddThen(db *sql.DB, then string) error {
	if id.Human == 0 {
		log.Println("human.go (id *IntData) AddThen(then string)")
		return errors.New("id.Human is not set")
	}
	if id.Data > 2147483647 {
		log.Println("human.go (id *IntData) AddThen(then string)")
		log.Println("data size is over max of int")
		return nil
	}

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

func LangSplit(db *sql.DB, text string) []StringData {
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
	ret2 = make([]string, 0)
	for _, word = range ret1 {
		added := ""
		for i := len([]rune(word)); i > 0; i-- {
			wd = string([]rune(word)[:i])
			if len([]rune(wd)) == 1 && !isKanji([]rune(wd)[0]) {
				break
			}
			sd1 := StartsWith(db, wd)
			if sd1.Data != "" {
				added = wd
				ret2 = append(ret2, wd)
				ret2 = append(ret2, word[len(wd):])
				break
			}
		}
		if added != word {
			ret2 = append(ret2, word)
		}
	}
	sds := make([]StringData, 0)
	ml2 := []string{"\n", "、", "。", "!", "！", "…", "?", "？"}
	for _, word = range ret2 {
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

func StartsWith(db *sql.DB, str string) StringData {
	var ret StringData
	if strings.TrimSpace(str) == "" {
		return ret
	}
	if _, err := strconv.Atoi(str); err == nil {
		return ret
	}

	sql := "select `string_data`.`data`, `feel`, `equal`, `in_count`, `out_count`, `param`, `val` " +
		"from `string_data` left outer join `string_data_than` on `string_data`.`data` = `string_data_than`.`data` " +
		"where `string_data`.`data` != '" + database.Escape(str) + "' and `string_data`.`data` collate utf8mb4_bin like '" + database.Escape(str) + "%'"
	rows, err := db.Query(sql)
	if err != nil {
		log.Println("human.go StartsWith(str string)")
		log.Println(err)
		log.Println(sql)
		return ret
	}
	defer rows.Close()

	for rows.Next() {
		var param string
		var val string
		err = rows.Scan(&ret.Data, &ret.Feel, &ret.Equal, &ret.InCount, &ret.OutCount, &param, &val)
		if err != nil {
			log.Println("human.go StartsWith(str string)")
			log.Println(err)
			return ret
		}
		if param == "group" {
			ret.Group = append(ret.Group, val)
		}
		if param == "big" {
			ret.Big = append(ret.Big, val)
		}
		if param == "not" {
			ret.Not = append(ret.Not, val)
		}
		if param == "then" {
			ret.Then = append(ret.Then, val)
		}
		if param == "next" {
			ret.Next = append(ret.Next, val)
		}
		if param == "action" {
			act, err := strconv.Atoi(val)
			if err == nil {
				ret.Action = append(ret.Action, act)
			}
		}
	}
	ret.Got = true
	return ret
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

func SaveWords(db *sql.DB, hid int, sds []StringData) error {
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
			if err = sd.Insert(db); err != nil {
				log.Println("human.go SaveWords(hid int, sds []StringData)")
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
				if err = id.Insert(db); err != nil {
					log.Println("human.go SaveWords(hid int, sds []StringData)")
					return err
				}
				if err = id.AddThen(db, sd.Data[len(numstr):]); err != nil {
					log.Println("human.go SaveWords(hid int, sds []StringData)")
					return err
				}
			}
		} else {
			id := IntData{
				Human: hid,
				Data:  num,
			}
			if err := id.Insert(db); err != nil {
				log.Println("human.go SaveWords(hid int, sds []StringData)")
				return err
			}
		}
		for i := 0; i < len(sd.Action); i++ {
			if err := sd.AddAction(db, sd.Action[i]); err != nil {
				log.Println("human.go SaveWords(hid int, sds []StringData)")
				return err
			}
		}
		if last.Data != "" {
			if err := last.AddNext(db, sd.Data); err != nil {
				log.Println("human.go SaveWords(hid int, sds []StringData)")
				return err
			}
		}
		if last.Data == "は" && ago2.Data != "" && notQuestion(last) {
			if err := ago2.AddParam(db, "equal", sd.Data); err != nil {
				log.Println("human.go SaveWords(hid int, sds []StringData)")
			}
			if err := sd.AddParam(db, "equal", ago2.Data); err != nil {
				log.Println("human.go SaveWords(hid int, sds []StringData)")
			}
			if err := ago2.AddParam(db, "group", sd.Data); err != nil {
				log.Println("human.go SaveWords(hid int, sds []StringData)")
			}
		}
		if last.Data == "の" && ago2.Data != "" {
			if err := ago2.AddParam(db, "group", sd.Data); err != nil {
				log.Println("human.go SaveWords(hid int, sds []StringData)")
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
				if _, err := strconv.Atoi(str); err != nil {
					if str != sd.Data {
						newword := StringData{
							Human:  hid,
							Data:   str,
							Action: last.Action,
						}
						if err := newword.Insert(db); err != nil {
							log.Println("human.go SaveWords(hid int, sds []StringData)")
							return err
						}
						for i := 0; i < len(newword.Action); i++ {
							if err := newword.AddAction(db, newword.Action[i]); err != nil {
								log.Println("human.go SaveWords(hid int, sds []StringData)")
								return err
							}
						}
						if err := newword.AddParam(db, "then", sd.Data); err != nil {
							log.Println("human.go SaveWords(hid int, sds []StringData)")
						}
					}
				}
			} else {
				if ago2.Data+last.Data != sd.Data && ago2.Data != sd.Data {
					if err := ago2.AddParam(db, "then", sd.Data); err != nil {
						log.Println("human.go SaveWords(hid int, sds []StringData)")
					}
				}
			}
		}
		if strings.HasSuffix(sd.Data, "ない") && notQuestion(sd) && ago2.Data != "" {
			str := sd.Data[:strings.LastIndex(sd.Data, "ない")]
			if str != "" {
				if _, err := strconv.Atoi(str); err != nil {
					newword := StringData{
						Human:  hid,
						Data:   str,
						Action: sd.Action,
					}
					if err := newword.Insert(db); err != nil {
						log.Println("human.go SaveWords(hid int, sds []StringData)")
						return err
					}
					for i := 0; i < len(newword.Action); i++ {
						if err := newword.AddAction(db, newword.Action[i]); err != nil {
							log.Println("human.go SaveWords(hid int, sds []StringData)")
							return err
						}
					}
					if err := last.AddParam(db, "not", str); err != nil {
						log.Println("human.go SaveWords(hid int, sds []StringData)")
					}
				}
			} else {
				if err := ago2.AddParam(db, "not", last.Data); err != nil {
					log.Println("human.go SaveWords(hid int, sds []StringData)")
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

func (tk *Talk) Insert(db *sql.DB) error {
	if tk.Opponent == "" {
		log.Println("human.go (tk *Talk) Insert()")
		return errors.New("tk.Opponent is empty")
	}
	if tk.Human == 0 {
		log.Println("human.go (tk *Talk) Insert()")
		return errors.New("tk.Human is not set")
	}
	if tk.Heard == "" {
		log.Println("human.go (tk *Talk) Insert()")
		return errors.New("tk.Heard is empty")
	}
	if tk.Talked == "" {
		log.Println("human.go (tk *Talk) Insert()")
		return errors.New("tk.Talked is empty")
	}

	sql := "insert into `talk` (`opponent`, `human`, `heard`, `talked`, `is_question`) values (?, ?, ?, ?, ?)"
	ins, err := db.Prepare(sql)
	if err != nil {
		log.Println("human.go (tk *Talk) Insert()")
		log.Println(err)
		log.Println(sql)
		return err
	}
	defer ins.Close()
	_, err = ins.Exec(&tk.Opponent, &tk.Human, &tk.Heard, &tk.Talked, &tk.IsQuestion)
	if err != nil {
		log.Println("human.go (tk *Talk) Insert()")
		log.Println(err)
		log.Println(sql)
		return err
	}
	return nil
}

func ResponseString(db *sql.DB, sd StringData) (Talk, error) {
	var ret Talk
	if sd.Human == 0 {
		log.Println("human.go ResponseString(sd StringData)")
		return ret, errors.New("sd.Human is not set")
	}
	if sd.Data == "" {
		log.Println("human.go ResponseString(sd StringData)")
		return ret, errors.New("sd.Data is empty")
	}
	_, err := strconv.Atoi(sd.Data)
	if err == nil {
		return ret, nil
	}

	sql := "select `data`, `param`, `val` from `string_data_than` where `param` in ('next', 'then') and `data` like '%" + database.Escape(sd.Data) + "%' and `human` = " + strconv.Itoa(sd.Human)
	if sdTopAct := most(sd.Action); sdTopAct == 65311 || sdTopAct == 63 {
		sql += " and `data` not in (select `data` from `string_data_than` where `param` = 'action' and `val` in (65311, 63))"
	}
	rows, err := db.Query(sql)
	if err != nil {
		log.Println("human.go ResponseString(sd StringData)")
		log.Println(err)
		log.Println(sql)
		return ret, err
	}
	defer rows.Close()
	exist := false
	for rows.Next() {
		exist = true
		//見つかった場合
		var data string
		var param string
		var val string
		err = rows.Scan(&data, &param, &val)
		if err != nil {
			log.Println("human.go ResponseString(sd StringData)")
			return ret, err
		}
		ret.Human = sd.Human
		ret.IsQuestion = true
		if param == "then" {
			ret.Talked = val
		} else {
			ret.Talked = data + val
		}
	}
	if !exist {
		//似た単語が見つからなかった場合
		ret.Human = sd.Human
		ret.IsQuestion = true
		ret.Talked = sd.Data
	}

	return ret, nil
}

func most(arr []int) int {
	type Act struct {
		Action int
		Count  int
	}
	acts := make([]Act, 0)
	for _, i := range arr {
		exists := false
		var a Act
		for _, act := range acts {
			if act.Action == i {
				exists = true
				a = act
				break
			}
		}
		if exists {
			a.Count++
		} else {
			a = Act{
				Action: i,
				Count:  1,
			}
			acts = append(acts, a)
		}
	}
	var ret Act
	for _, act := range acts {
		if ret.Count < act.Count {
			ret = act
		}
	}
	return ret.Action
}
