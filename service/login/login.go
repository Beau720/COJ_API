package login

import (
	"COJ_API/service/database"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Service struct {
	Db *gorm.DB
}

type Input struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
}

type LoginError string

const (
	LoginErrorNil                 LoginError = ""
	LoginErrorCredintialsNotFound LoginError = "User doesn't exist/ Not found!!"
)

func Checker(db_config *database.Config) *Service {
	var err error
	s := Service{}

	// Open a new GORM database connection
	s.Db, err = gorm.Open(mysql.Open(db_config.ConnectString()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &s
}

func (s *Service) GetUserByUsername(username string) (input *Input, login_error LoginError) {

	input = &Input{}
	/*var input Input

	if err := s.Db.Where("Email = ? AND Password = ?", username, password).First(&input).Error; err != nil {
		return Input{}, err
	}*/

	err := s.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Find(input, "Email = ?", username).Error
	})

	if err != nil {
		return nil, LoginError(err.Error())
	}

	if input.ID < 1 {
		return nil, LoginErrorCredintialsNotFound
	}

	return input, LoginErrorNil
}

func (s *Service) ValidatePassword(hashedPassword, Password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(Password))
	return err == nil
}
