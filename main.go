package main

import (
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
)

func ToPostgresURL() string {
	viper.AutomaticEnv()
	postgresConnectionStringFormat := "host=%s port=%d user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Seoul"
	return fmt.Sprintf(postgresConnectionStringFormat, viper.GetString("db_host"), 5432, "postgres", viper.GetString("db_name"), viper.GetString("db_password"))
}
func main() {
	db, err := gorm.Open(postgres.Open(ToPostgresURL()), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		handleError(err)
	}
	err = db.AutoMigrate(&Animal{})
	if err != nil {
		handleError(err)
	}
	err = CreateAnimals(db)
	if err != nil {
		handleError(err)
	}
}

func handleError(err error) {
	fmt.Println(err)
	//log.Fatal(err)
}

type Animal struct {
	Name string `json:"name" gorm:"unique"`
}

func CreateAnimals(db *gorm.DB) error {
	// Note the use of tx as the database handle once you are within a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	fmt.Println("성공합시다")
	if err := tx.Create(&Animal{Name: "만들어지면 안돼"}).Error; err != nil {
		tx.Rollback()
		return err
	}
	//testData := &Animal{}
	//err := tx.First(testData, "name = ?", "test").Error
	//if err != nil {
	//	tx.Rollback()
	//	return err
	//}
	//fmt.Println("find result", testData)
	fmt.Println("실패합시다")
	if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
		tx.Rollback()
		fmt.Println(reflect.TypeOf(err))
		e := err.(*pgconn.PgError)
		if e.Code == "23505" {
			return errors.New("duplicate!!!!!!!!!!!")
		}

		return err
	}
	fmt.Println("성공")

	return tx.Commit().Error
}
