package user

import (
	"COJ_API/service/database"
	"crypto/sha256"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Service struct {
	Db *gorm.DB
}

type UserError string

// error handling messages
const (
	UserErrorNil             UserError = ""
	UserErrorUserEmailExists UserError = "User with the same email address already exist"
	UserErrorUserRSAIDExists UserError = "User with the same RSA ID number already exist"
	UserErrorUserNotFound    UserError = "User not found"
)

// User Table attributes and constraites
type User struct {
	ID               int    `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	FirstName        string `gorm:"column:FirstName;not null" json:"first_name"`
	LastName         string `gorm:"column:LastName;not null" json:"last_name"`
	Email            string `gorm:"column:Email;not null" json:"email"`
	Gender           string `gorm:"column:Gender;not null" json:"gender"`
	Province         string `gorm:"column:Province;not null" json:"province"`
	SAID             string `gorm:"column:SAID" json:"said,omitempty"`
	CreatedDateD     []byte `gorm:"column:CreatedDate;type:DATETIME DEFAULT CURRENT_TIMESTAMP" json:"-"`
	LastUpdatedDateD []byte `gorm:"column:LastUpdatedDate;type:DATETIME DEFAULT CURRENT_TIMESTAMP" json:"-"`
	CreatedDateJ     string `gorm:"-" json:"created_date"`
	LastUpdatedDateJ string `gorm:"-" json:"last_updated_date"`
	Password         string `gorm:"column:Password;not null" json:"-"`
	RoleID           int    `gorm:"column:RoleID" json:"role_id,omitempty"`
	Active           int    `gorm:"column:Active" json:"active,omitempty"`
}

// naming the user table
func (u *User) TableName() string {
	return "user"
}

func NewUser(db_config *database.Config) *Service {
	var err error

	s := Service{}

	// Open a new GORM database connection
	s.Db, err = gorm.Open(mysql.Open(db_config.ConnectString()), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Perform auto-migration to create or update the table based on the User struct
	err = s.Db.AutoMigrate(&User{})
	if err != nil {
		panic("Auto-migration failed: " + err.Error())
	}

	log.Println("Auto-migration completed successfully.")

	// Remember to close the database connection when done.
	//s.Db.Close()

	return &s
}

// create a new user
func (s *Service) Create(user *User) (user_resp *User, user_error UserError) {
	_, user_error = s.SelectByEmail(user.Email)
	if user_error != UserErrorUserNotFound {
		return nil, UserErrorUserEmailExists
	}

	_, user_error = s.SelectByRSAID(user.SAID)
	if user_error != UserErrorUserNotFound {
		return nil, UserErrorUserRSAIDExists
	}

	// Reset the ID at it will be omitted and created on insert
	user.ID = 0

	err := s.Db.Transaction(func(tx *gorm.DB) error {
		user.CreatedDateD = []byte(time.Now().Format("2006-01-02 15:04:05"))
		user.LastUpdatedDateD = []byte(time.Now().Format("2006-01-02 15:04:05"))

		//hashing the password
		user.Password = hash256(user.Password)

		//restricting the id to be updated
		err := tx.Omit("ID").Create(user).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, UserError(err.Error())
	}

	user, user_error = s.SelectById(user.ID)
	if user_error != UserErrorNil {
		return nil, user_error
	}

	return user, UserErrorNil
}

// update the user information
func (s *Service) Update(user *User) (user_resp *User, user_error UserError) {
	_, user_error = s.SelectById(user.ID)
	if user_error != UserErrorNil {
		return nil, user_error
	}

	update_data := make(map[string]interface{})
	user_struct := reflect.ValueOf(user).Elem()

	user.LastUpdatedDateD = []byte(time.Now().Format("2006-01-02 15:04:05"))
	//assigning name to the user structre for it to be updated
	for i := 0; i < user_struct.NumField(); i++ {
		name := user_struct.Type().Field(i).Name
		value := user_struct.Field(i).Interface()

		//blocking te user to updating the spefic fields
		if name == "ID" || name == "Password" || strings.Contains(name, "CreatedDate") || strings.Contains(name, "LastUpdatedDate") {
			continue
		}
		//updatable fields
		for j := 0; j < user_struct.NumField(); j++ {
			update_data[name] = value
		}

	}

	err := s.Db.Transaction(func(tx *gorm.DB) error {
		err := tx.Table(user.TableName()).Where("ID = ?", user.ID).Updates(&update_data).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, UserError(err.Error())
	}

	user, user_error = s.SelectById(user.ID)
	if user_error != UserErrorNil {
		return nil, user_error
	}

	return user, UserErrorNil
}

// find the user by ID
func (s *Service) SelectById(user_id int) (user *User, user_error UserError) {
	user = &User{}
	err := s.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Find(user, "ID = ?", user_id).Error
	})

	if err != nil {
		return nil, UserError(err.Error())
	}

	if user.ID < 1 {
		return nil, UserErrorUserNotFound
	}

	user.CreatedDateJ = byteTimeStampToString(user.CreatedDateD)
	user.LastUpdatedDateJ = byteTimeStampToString(user.CreatedDateD)

	return user, UserErrorNil
}

// find the user by email
func (s *Service) SelectByEmail(email string) (user *User, user_error UserError) {
	user = &User{}
	err := s.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Find(user, "Email = ?", email).Error
	})

	if err != nil {
		return nil, UserError(err.Error())
	}

	if user.ID < 1 {
		return nil, UserErrorUserNotFound
	}

	return user, UserErrorNil
}

// find user by the SA ID number
func (s *Service) SelectByRSAID(said string) (user *User, user_error UserError) {
	user = &User{}
	err := s.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Find(user, "SAID = ?", said).Error
	})

	if err != nil {
		return nil, UserError(err.Error())
	}

	if user.ID < 1 {
		return nil, UserErrorUserNotFound
	}

	return user, UserErrorNil
}

// fetch all users
func (s *Service) List() (users []*User, user_error UserError) {
	users = make([]*User, 0)

	err := s.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Find(&users).Error
	})

	if err != nil {
		return nil, UserError(err.Error())
	}

	return users, UserErrorNil
}

// converting the time stamp to a string so it can be used by the front end
func byteTimeStampToString(str []byte) string {
	time_str, _ := time.Parse("2006-01-02 15:04:05", string(str))
	t, err := time.Parse(time.RFC3339, time_str.Format(time.RFC3339))
	if err != nil {
		panic(err)
	}

	return t.Format("2006-01-02 15:04:05")
}

// hashing the password so it store securely
func hash256(in string) string {
	// Hash the password with sha256 encryption method
	s := sha256.New()
	s.Write([]byte(in))
	// Convert print the bytes to a string
	return fmt.Sprintf("%x", s.Sum(nil))
}
