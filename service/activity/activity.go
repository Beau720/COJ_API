package activity

import (
	"gorm.io/gorm"
)

type Service struct {
	Db *gorm.DB
}

type ActivityError string

// error handling messages
const (
	ActivityErrorNil              ActivityError = ""
	ActivityErrorActivityNotFound ActivityError = "Activity not found"
)

type Activity struct {
	ActivityID   int    `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	ActivityType string `gorm:"column:FirstName;not null" json:"first_name"`
	LastName     string `gorm:"column:LastName;not null" json:"last_name"`
}

type ActiveRecord struct {
	ActiveRecord int    `gorm:"column:ActiveRecord;primaryKey;autoIncrement" json:"active_record "`
	ActivityName string `gorm:"column:ActivityName" json:"activity_name "`
}

func (a *Activity) TableName() string {
	return "activities"
}

// find the user by ID
func (s *Service) SelectByActivityId(activity_id int) (activity *Activity, activity_error ActivityError) {
	activity = &Activity{}
	err := s.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Find(activity, "ActivityID = ?", activity_id).Error
	})

	if err != nil {
		return nil, ActivityError(err.Error())
	}

	if activity.ActivityID < 1 {
		return nil, ActivityErrorActivityNotFound
	}

	return activity, ActivityErrorNil
}

// find the user by email
func (s *Service) SelectByactivitytype(activitytype string) (activity *Activity, activity_error ActivityError) {
	activity = &Activity{}
	err := s.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Find(activity, "ActivityType = ?", activitytype).Error
	})

	if err != nil {
		return nil, ActivityError(err.Error())
	}

	if activity.ActivityID < 1 {
		return nil, ActivityErrorActivityNotFound
	}

	return activity, ActivityErrorNil
}

// fetch all users
func (s *Service) List() (activities []*Activity, activity_error ActivityError) {
	activities = make([]*Activity, 0)

	err := s.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Find(&activities).Error
	})

	if err != nil {
		return nil, ActivityError(err.Error())
	}

	return activities, ActivityErrorNil
}
