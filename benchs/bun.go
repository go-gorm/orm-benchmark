package benchs

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var bundb *bun.DB

func init() {
	st := NewSuite("bun")
	st.InitF = func() {
		st.AddBenchmark("Insert", 200*ORM_MULTI, PgInsert)
		st.AddBenchmark("MultiInsert 100 row", 200*ORM_MULTI, PgInsertMulti)
		st.AddBenchmark("Update", 200*ORM_MULTI, PgUpdate)
		st.AddBenchmark("Read", 200*ORM_MULTI, PgRead)
		st.AddBenchmark("MultiRead limit 100", 200*ORM_MULTI, PgReadSlice)

		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(ORM_SOURCE)))
		bundb = bun.NewDB(sqldb, pgdialect.New())
	}
}

func PgInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})

	for i := 0; i < b.N; i++ {
		m.Id = 0
		if _, err := bundb.NewInsert().Model(m).Exec(context.Background()); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PgInsertMulti(b *B) {
	var ms []*Model
	wrapExecute(b, func() {
		initDB()
		ms = make([]*Model, 0, 100)
		for i := 0; i < 100; i++ {
			ms = append(ms, NewModel())
		}
	})

	for i := 0; i < b.N; i++ {
		for _, m := range ms {
			m.Id = 0
		}
		if _, err := bundb.NewInsert().Model(&ms).Exec(context.Background()); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PgUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := bundb.NewInsert().Model(m).Exec(context.Background()); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		if _, err := bundb.NewUpdate().Model(m).Where("id = ?", m.Id).Exec(context.Background()); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PgRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := bundb.NewInsert().Model(m).Exec(context.Background()); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		if err := bundb.NewSelect().Model(m).Scan(context.Background()); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PgReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < 100; i++ {
			m.Id = 0
			if _, err := bundb.NewInsert().Model(m).Exec(context.Background()); err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
	})

	for i := 0; i < b.N; i++ {
		var models []*Model
		if err := bundb.NewSelect().Model(&models).Where("id > ?", 0).Limit(100).Scan(context.Background()); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}
