package benchs

import (
	"fmt"

	gormv1 "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var gormV1DB *gormv1.DB

func init() {
	st := NewSuite("gorm_v1")
	st.InitF = func() {
		st.AddBenchmark("Insert", 200*ORM_MULTI, GormV1Insert)
		st.AddBenchmark("MultiInsert 100 row", 200*ORM_MULTI, GormV1InsertMulti)
		st.AddBenchmark("Update", 200*ORM_MULTI, GormV1Update)
		st.AddBenchmark("Read", 200*ORM_MULTI, GormV1Read)
		st.AddBenchmark("MultiRead limit 100", 200*ORM_MULTI, GormV1ReadSlice)
		dsn := ORM_SOURCE
		conn, err := gormv1.Open("postgres", dsn)
		if err != nil {
			panic(err)
		}
		gormV1DB = conn

		gormV1DB.Callback().Create().Remove("gorm:begin_transaction")
		gormV1DB.Callback().Update().Remove("gorm:begin_transaction")
		gormV1DB.Callback().Delete().Remove("gorm:begin_transaction")

		gormV1DB.Callback().Create().Remove("gorm:commit_or_rollback_transaction")
		gormV1DB.Callback().Update().Remove("gorm:commit_or_rollback_transaction")
		gormV1DB.Callback().Delete().Remove("gorm:commit_or_rollback_transaction")
	}
}

func GormV1Insert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})

	for i := 0; i < b.N; i++ {
		m.Id = 0
		if err := gormV1DB.Create(m).Error; err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func GormV1InsertMulti(b *B) {
	var ms []*Model
	wrapExecute(b, func() {
		initDB()
		ms = make([]*Model, 0, 100)
		for i := 0; i < 100; i++ {
			ms = append(ms, NewModel())
		}
	})

	for i := 0; i < b.N; i++ {
		panic(fmt.Errorf("doesn't work"))
		for _, m := range ms {
			m.Id = 0
		}
		if err := gormV1DB.Create(&ms).Error; err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func GormV1Update(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if err := gormV1DB.Create(m).Error; err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		if err := gormV1DB.Model(m).Updates(m).Error; err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func GormV1Read(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if err := gormV1DB.Create(m).Error; err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})
	for i := 0; i < b.N; i++ {
		if err := gormV1DB.Take(m).Error; err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func GormV1ReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < 100; i++ {
			m.Id = 0
			if err := gormV1DB.Create(m).Error; err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
	})

	for i := 0; i < b.N; i++ {
		var models []*Model
		if err := gormV1DB.Where("id > ?", 0).Limit(100).Find(&models).Error; err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}
