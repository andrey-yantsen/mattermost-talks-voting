package storage

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	migration "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mattermost/mattermost-server/model"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"time"
)

type Storage struct {
	uri string
	db  *sql.DB
}

type Registration struct {
	ChannelId                     string
	Timezone                      *time.Location
	OwnerId                       string
	TalksToShowOnVoting           int
	ScheduleDow                   int8
	ScheduleTimeInMinutes         int16
	MinimalPeopleRequired         int
	FirstReminderHours            int8
	SecondReminderMinutes         int8
	AddRandomTalkOnSecondReminder bool
	FinalCallMinutes              int8
	ShowVoteResult                bool
	CreatedAt                     time.Time
	UpdatedAt                     time.Time
	Active                        bool
}

func LoadRegistrationFromMap(data map[string]interface{}) (*Registration, map[string]error) {
	var err error

	errors := make(map[string]error)
	r := &Registration{}

	r.ChannelId = data["channel_id"].(string)

	if len(r.ChannelId) == 0 || len(r.ChannelId) > 32 {
		errors["channel_id"] = fmt.Errorf("`channel_id` length should be between 1 and 32 chars")
	}

	r.OwnerId = data["channel_id"].(string)

	if len(r.OwnerId) == 0 || len(r.OwnerId) > 32 {
		errors["owner_id"] = fmt.Errorf("`owner_id` length should be between 1 and 32 chars")
	}

	tz := data["timezone"].(string)

	r.Timezone, err = time.LoadLocation(tz)
	if err != nil {
		errors["timezone"] = err
	}

	v, err := strconv.ParseInt(data["talks_to_show_on_voting"].(string), 10, 64)
	if err != nil {
		errors["talks_to_show_on_voting"] = err
	} else if v < 2 || v > 10 {
		errors["talks_to_show_on_voting"] = fmt.Errorf("value should be between 2 and 10")
	} else {
		r.TalksToShowOnVoting = int(v)
	}

	v, err = strconv.ParseInt(data["minimal_people_required"].(string), 10, 64)
	if err != nil {
		errors["minimal_people_required"] = err
	} else if v < 3 || v > 10 {
		errors["minimal_people_required"] = fmt.Errorf("value should be between 3 and 10")
	} else {
		r.MinimalPeopleRequired = int(v)
	}

	v, err = strconv.ParseInt(data["first_reminder_hours"].(string), 10, 64)
	if err != nil {
		errors["first_reminder_hours"] = err
	} else if v < 1 || v > 12 {
		errors["first_reminder_hours"] = fmt.Errorf("value should be between 1 and 12")
	} else {
		r.FirstReminderHours = int8(v)
	}

	v, err = strconv.ParseInt(data["second_reminder_minutes"].(string), 10, 64)
	if err != nil {
		errors["second_reminder_minutes"] = err
	} else if v < 10 || v > 60 {
		errors["second_reminder_minutes"] = fmt.Errorf("value should be between 10 and 60")
	} else if v >= int64(r.FirstReminderHours)*60 {
		errors["second_reminder_minutes"] = fmt.Errorf("second reminder should be after the first")
	} else {
		r.SecondReminderMinutes = int8(v)
	}

	r.AddRandomTalkOnSecondReminder = data["add_random_talk_on_second_reminder"].(string) == "1"

	v, err = strconv.ParseInt(data["final_call_minutes"].(string), 10, 64)
	if err != nil {
		errors["final_call_minutes"] = err
	} else if v < 0 || v > 60 {
		errors["final_call_minutes"] = fmt.Errorf("value should be between 0 and 60")
	} else if v >= int64(r.SecondReminderMinutes) {
		errors["final_call_minutes"] = fmt.Errorf("voting termination should be after second reminder")
	} else {
		r.FinalCallMinutes = int8(v)
	}

	if len(errors) > 0 {
		return nil, errors
	} else {
		return r, nil
	}
}

func (r *Registration) GetDialogElements() []model.DialogElement {
	return []model.DialogElement{
		//{
		//	DisplayName: "Talks in which languages to show?",
		//	Name:        "languages",
		//	Type:        "select",
		//	Default:     "en",
		//	Options: []*model.PostActionOptions{
		//		{"English", "en"},
		//		{"Russian", "ru"},
		//		{"English + Russian", "en,ru"},
		//	},
		//},
		{
			DisplayName: "Talks to display for each voting",
			Name:        "talks_to_show_on_voting",
			Default:     "4",
			SubType:     "int",
			Type:        "text",
		},
		{
			DisplayName: "Minimal peoples required for successful voting",
			Name:        "minimal_people_required",
			Default:     "3",
			SubType:     "int",
			Type:        "text",
		},
		{
			DisplayName: "Starting voting process in X hours before the scheduled watching time",
			Name:        "first_reminder_hours",
			Default:     "6",
			SubType:     "int",
			Type:        "text",
		},
		{
			DisplayName: "Send a reminder in X minutes before the time if not enough votes",
			Name:        "second_reminder_minutes",
			Default:     "30",
			SubType:     "int",
			Type:        "text",
		},
		{
			DisplayName: "Add one more talk with second reminder?",
			Name:        "add_random_talk_on_second_reminder",
			Default:     "0",
			Type:        "select",
			Options: []*model.PostActionOptions{
				{"Yes", "1"},
				{"No", "0"},
			},
		},
		{
			DisplayName: "Close voting in X minutes before the time",
			Name:        "final_call_minutes",
			Default:     "1",
			Type:        "text",
			SubType:     "int",
		},
		//{
		//	DisplayName: "Day of Week for the talks",
		//	Name:        "schedule_dow",
		//	Type:        "select",
		//	Placeholder: "Select the day of week",
		//	Default:     "5",
		//	Options: []*model.PostActionOptions{
		//		{"Monday", "1"},
		//		{"Tuesday", "2"},
		//		{"Wednesday", "3"},
		//		{"Thursday", "4"},
		//		{"Friday", "5"},
		//	},
		//},
		//{
		//	DisplayName: "Time to watch the talk",
		//	Name:        "schedule_time_in_minutes",
		//	Type:        "select",
		//	Default:     "1020",
		//	Placeholder: "Select the time",
		//	Options: []*model.PostActionOptions{
		//		{"09:00", "540"},
		//		{"10:00", "600"},
		//		{"11:00", "660"},
		//		{"12:00", "720"},
		//		{"13:00", "780"},
		//		{"14:00", "840"},
		//		{"15:00", "900"},
		//		{"16:00", "960"},
		//		{"17:00", "1020"},
		//		{"18:00", "1080"},
		//	},
		//},
		{
			DisplayName: "Timezone",
			Name:        "timezone",
			Type:        "select",
			Default:     "Europe/London",
			Options: []*model.PostActionOptions{
				{"London", "Europe/London"},
				{"Moscow", "Europe/Moscow"},
				{"Austin (CST)", "US/Central"},
			},
		},
	}
}

func DbConnect(uri string) (*Storage, error) {
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return &Storage{
		uri: uri,
		db:  db,
	}, nil
}

func (s *Storage) Migrate() error {
	cfg := &migration.Config{}
	db, err := migration.WithInstance(s.db, cfg)
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "", db)
	if err != nil {
		return err
	}
	dbVersion, _, _ := db.Version()
	mVersion, _, _ := m.Version()
	if int(mVersion) > dbVersion {
		return m.Up()
	}
	return nil
}

func (s *Storage) IsRegistered(channelId string) bool {
	ret := 0
	err := s.db.QueryRow("SELECT 1 FROM registrations WHERE channel_id = ?", channelId).Scan(&ret)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err)
		}
	}
	return false
}

func (s *Storage) GetRegistration(channelId string) *Registration {
	var tz string
	r := &Registration{}
	err := s.db.QueryRow("SELECT channel_id, timezone, owner_id, talks_to_show_on_voting, "+
		"schedule_dow, schedule_time_in_minutes, minimal_people_required, "+
		"first_reminder_hours, second_reminder_minutes, add_random_talk_on_second_reminder, "+
		"final_call_minutes, show_vote_result, active, created_at, updated_at "+
		"FROM registrations WHERE channel_id = ?", channelId).Scan(
		&r.ChannelId,
		&tz,
		&r.OwnerId,
		&r.TalksToShowOnVoting,
		&r.ScheduleDow,
		&r.ScheduleTimeInMinutes,
		&r.MinimalPeopleRequired,
		&r.FirstReminderHours,
		&r.SecondReminderMinutes,
		&r.AddRandomTalkOnSecondReminder,
		&r.FinalCallMinutes,
		&r.ShowVoteResult,
		&r.Active,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			panic(err)
		}
	}

	r.Timezone, err = time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}

	return r
}

func (s *Storage) SaveRegistration(r *Registration) error {
	_, err := s.db.Exec("REPLACE INTO registrations (channel_id, timezone, owner_id,"+
		"talks_to_show_on_voting, schedule_dow, schedule_time_in_minutes, minimal_people_required,"+
		"first_reminder_hours, second_reminder_minutes, add_random_talk_on_second_reminder,"+
		"final_call_minutes, show_vote_result, active, updated_at)"+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())",
		r.ChannelId, r.Timezone.String(), r.OwnerId, r.TalksToShowOnVoting, r.ScheduleDow,
		r.ScheduleTimeInMinutes, r.MinimalPeopleRequired, r.FirstReminderHours, r.SecondReminderMinutes,
		r.AddRandomTalkOnSecondReminder, r.FinalCallMinutes, r.ShowVoteResult, r.Active)
	return err
}
