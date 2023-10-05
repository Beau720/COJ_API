package form

import (
	"log"
	"reflect"
	"strings"
	"time"

	"COJ_API/service/database"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Service struct {
	Db *gorm.DB
}

type FormError string

/*type Config struct {
	Username       string
	Password       string
	ConnectionType string
	Host           string
	Port           int
	Name           string
}

func (c *Config) ConnectString() string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s", c.Username, c.Password, c.ConnectionType, c.Host, c.Port, c.Name)
}*/

const (
	FormErrorNil          FormError = ""
	FormErrorFormNotFound FormError = "Form not found"
)

type Form struct {
	ReferenceNumber    int    `gorm:" column:ReferenceNumber;primaryKey" json:"referenceNumber"`
	DateOfArrival      string `gorm:" column:DateOfArrival;not null" json:"dateOfArrival"`
	InspectorOfficer   string `gorm:" column:InspectorOfficer;not null" json:"inspector_officer"`
	InspectionNumber   int    `gorm:" column:InspectionNumber;not null" json:"inspector_No"`
	Basement           int    `gorm:" column:Basement;not null json:basement"`
	Floor              int    `gorm:" column:Floor;not null json:floor"`
	Meeting            int    `gorm:" column:Meeting;not null json:meeting"`
	Contact            string `gorm:" column:Contact;not null json:contact"`
	Telephone          string `gorm:" column:Telephone;not null" json:"telephone"`
	BuildingName       string `gorm:" column:BuildingName;not null" json:" building_name"`
	Address            string `gorm:" column:Address;not null" json:"address"`
	Premise            string `gorm:" column:Premise;not null" json:"premise"`
	Activity           string `gorm:" column:Activity;not null" json:"activity"`
	Finding            string `gorm:" column:Finding;not null" json:"finding"`
	Grading            string `gorm:" column:Grading;not null" json:"grading"`
	Comment            string `gorm:" column:Comment;not null" json:"comment"`
	FollowUpDate       string `gorm:" column:FollowUpDate;not null" json:"follow_up_date"`
	InspectorSignature string `gorm:" column:InspectorSignature;not null" json:"inspector_signature"`
	CreatedDateD       []byte `gorm:"column:CreatedDate;type:DATETIME DEFAULT CURRENT_TIMESTAMP" json:"-"`
	LastUpdatedDateD   []byte `gorm:"column:LastUpdatedDate;type:DATETIME DEFAULT CURRENT_TIMESTAMP" json:"-"`
	CreatedDateJ       string `gorm:"-" json:"created_date"`
	LastUpdatedDateJ   string `gorm:"-" json:"last_updated_date"`
}

func (u *Form) TableName() string {
	return "form"
}

func NewForm(db_config *database.Config) *Service {
	var err error

	s := Service{}

	// Open a new GORM database connection
	s.Db, err = gorm.Open(mysql.Open(db_config.ConnectString()), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Perform auto-migration to create or update the table based on the User struct
	err = s.Db.AutoMigrate(&Form{})
	if err != nil {
		panic("Auto-migration failed: " + err.Error())
	}

	log.Println("Auto-migration completed successfully.")

	// Remember to close the database connection when done.
	//s.Db.Close()

	return &s
}

func (s *Service) Create(form *Form) (form_resp *Form, form_error FormError) {
	// Reset the ID at it will be omitted and created on insert
	form.ReferenceNumber = 0

	err := s.Db.Transaction(func(tx *gorm.DB) error {
		form.CreatedDateD = []byte(time.Now().Format("2006-01-02 15:04:05"))
		form.LastUpdatedDateD = []byte(time.Now().Format("2006-01-02 15:04:05"))

		err := tx.Omit("ReferenceNumber").Create(form).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, FormError(err.Error())
	}

	form, form_error = s.SelectByRefNo(form.ReferenceNumber)
	if form_error != FormErrorNil {
		return nil, form_error
	}

	return form, FormErrorNil
}

func (s *Service) Update(form *Form) (form_resp *Form, form_error FormError) {
	_, form_error = s.SelectByRefNo(form.ReferenceNumber)
	if form_error != FormErrorNil {
		return nil, form_error
	}

	update_data := make(map[string]interface{})
	form_struct := reflect.ValueOf(form).Elem()

	form.LastUpdatedDateD = []byte(time.Now().Format("2006-01-02 15:04:05"))

	for i := 0; i < form_struct.NumField(); i++ {
		name := form_struct.Type().Field(i).Name
		value := form_struct.Field(i).Interface()

		if name == "ReferenceNumber" || name == "InspectorSignature" || strings.Contains(name, "CreatedDate") || strings.Contains(name, "LastUpdatedDate") {
			continue
		}

		for j := 0; j < form_struct.NumField(); j++ {
			update_data[name] = value
		}

	}

	err := s.Db.Transaction(func(tx *gorm.DB) error {
		err := tx.Table(form.TableName()).Where("ReferenceNumber = ?", form.ReferenceNumber).Updates(&update_data).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, FormError(err.Error())
	}

	form, form_error = s.SelectByRefNo(form.ReferenceNumber)
	if form_error != FormErrorNil {
		return nil, form_error
	}

	return form, FormErrorNil
}

func (s *Service) SelectByRefNo(refNo int) (form *Form, form_error FormError) {
	form = &Form{}
	err := s.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Find(form, "ReferenceNumber = ?", refNo).Error
	})

	if err != nil {
		return nil, FormError(err.Error())
	}

	if form.ReferenceNumber < 1 {
		return nil, FormErrorFormNotFound
	}

	form.CreatedDateJ = byteTimeStampToString(form.CreatedDateD)
	form.LastUpdatedDateJ = byteTimeStampToString(form.CreatedDateD)

	return form, FormErrorNil
}

func (s *Service) FindByInspectorName(inspector_name string) (form *Form, form_error FormError) {
	form = &Form{}
	err := s.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Find(form, "InspectorOfficer = ?", inspector_name).Error
	})

	if err != nil {
		return nil, FormError(err.Error())
	}

	form.CreatedDateJ = byteTimeStampToString(form.CreatedDateD)

	if form.ReferenceNumber < 1 {
		return nil, FormErrorFormNotFound
	}

	return form, FormErrorNil
}

func (s *Service) FindbydateofArrival(dateOfArrival string) (form *Form, form_error FormError) {
	form = &Form{}
	err := s.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Find(&form, "DateOfArrival = ?", dateOfArrival).Error
	})

	if err != nil {
		return form, FormError(err.Error())
	}

	if form.ReferenceNumber < 1 {
		return form, FormErrorFormNotFound
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", form.CreatedDateJ)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	form.DateOfArrival = parsedTime.Format("2006-01-02")

	return form, FormErrorNil
}

func (s *Service) List() (forms []*Form, form_error FormError) {
	forms = make([]*Form, 0)

	err := s.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Find(&forms).Error
	})

	if err != nil {
		return nil, FormError(err.Error())
	}

	return forms, FormErrorNil
}

func byteTimeStampToString(str []byte) string {
	time_str, _ := time.Parse("2006-01-02 15:04:05", string(str))
	t, err := time.Parse(time.RFC3339, time_str.Format(time.RFC3339))
	if err != nil {
		panic(err)
	}

	return t.Format("2006-01-02 15:04:05")
}
