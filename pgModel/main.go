package pgModel

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/subosito/gotenv"
	"os"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

 
 func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v>", u.Id, u.Username, u.Email)
 }

 func Connect() *pg.DB {
	gotenv.Load()

	db := pg.Connect(&pg.Options{
		User: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DATABASE"),
	 })
  
	 err := createSchema(db)
	 if err != nil {
		fmt.Println(err)
		panic(err)
	 }

	return db
 }


 func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{&User{}} {
	   if err := db.CreateTable(model, nil); err != nil {
		  return nil
	   }
	}
	return nil
 }
