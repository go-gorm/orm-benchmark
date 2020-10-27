package benchs

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
)

var sqlxDB *sqlx.DB

const (
	sqlxInsertBaseSQL   = `INSERT INTO models (name, title, fax, web, age, "right", counter) VALUES `
	sqlxInsertValuesSQL = `($1, $2, $3, $4, $5, $6, $7)`
	sqlxInsertSQL       = sqlxInsertBaseSQL + sqlxInsertValuesSQL
	sqlxUpdateSQL       = `UPDATE models SET name = $1, title = $2, fax = $3, web = $4, age = $5, "right" = $6, counter = $7 WHERE id = $8`
	sqlxSelectSQL       = `SELECT id, name, title, fax, web, age, "right", counter FROM models WHERE id = $1`
	sqlxSelectMultiSQL  = `SELECT id, name, title, fax, web, age, "right", counter FROM models WHERE id > 0 LIMIT 100`
)

func init() {
	st := NewSuite("sqlx")
	st.InitF = func() {
		st.AddBenchmark("Insert", 200*ORM_MULTI, SQLXInsert)
		st.AddBenchmark("MultiInsert 100 row", 200*ORM_MULTI, SQLXInsertMulti)
		st.AddBenchmark("Update", 200*ORM_MULTI, SQLXUpdate)
		st.AddBenchmark("Read", 200*ORM_MULTI, SQLXRead)
		st.AddBenchmark("MultiRead limit 100", 200*ORM_MULTI, SQLXReadSlice)

		sqlxDB, _ = sqlx.Open("postgres", ORM_SOURCE)
	}
}

func SQLXInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		var err error
		initDB()
		m = NewModel()
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		// pq dose not support the LastInsertId method.
		_, err := sqlxDB.Exec(sqlxInsertSQL, m.Name, m.Title, m.Fax, m.Web, m.Age, m.Right, m.Counter)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func sqlxInsert(m *Model) error {
	// pq dose not support the LastInsertId method.
	_, err := sqlxDB.Exec(sqlxInsertSQL, m.Name, m.Title, m.Fax, m.Web, m.Age, m.Right, m.Counter)
	if err != nil {
		return err
	}
	return nil
}

func SQLXInsertMulti(b *B) {
	var ms []*Model
	wrapExecute(b, func() {
		initDB()

		ms = make([]*Model, 0, 100)
		for i := 0; i < 100; i++ {
			ms = append(ms, NewModel())
		}
	})

	var valuesSQL string
	counter := 1
	for i := 0; i < 100; i++ {
		hoge := ""
		for j := 0; j < 7; j++ {
			if j != 6 {
				hoge += "$" + strconv.Itoa(counter) + ","
			} else {
				hoge += "$" + strconv.Itoa(counter)
			}
			counter++

		}
		if i != 99 {
			valuesSQL += "(" + hoge + "),"
		} else {
			valuesSQL += "(" + hoge + ")"
		}
	}

	for i := 0; i < b.N; i++ {
		nFields := 7
		query := sqlxInsertBaseSQL + valuesSQL
		args := make([]interface{}, len(ms)*nFields)
		for j := range ms {
			offset := j * nFields
			args[offset+0] = ms[j].Name
			args[offset+1] = ms[j].Title
			args[offset+2] = ms[j].Fax
			args[offset+3] = ms[j].Web
			args[offset+4] = ms[j].Age
			args[offset+5] = ms[j].Right
			args[offset+6] = ms[j].Counter
		}
		// pq dose not support the LastInsertId method.
		_, err := sqlxDB.Exec(query, args...)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func SQLXUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		var err error
		initDB()
		m = NewModel()
		sqlxInsert(m)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		_, err := sqlxDB.Exec(sqlxUpdateSQL, m.Name, m.Title, m.Fax, m.Web, m.Age, m.Right, m.Counter, m.Id)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func SQLXRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		var err error
		initDB()
		m = NewModel()
		sqlxInsert(m)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		var mout Model
		err := sqlxDB.Get(&mout, sqlxSelectSQL, 1)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func SQLXReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		var err error
		initDB()
		m = NewModel()
		for i := 0; i < 100; i++ {
			err = sqlxInsert(m)
			if err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		models := make([]Model, 100)
		err := sqlxDB.Select(&models, sqlxSelectMultiSQL)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}
