package test

import (
	"fmt"
	"scheduleme/models"
	sq "scheduleme/sqlite"
	"scheduleme/values"
	"time"
)

/*
	Quick and dirty seed'er for testing, not writing tests for this so not bothering with dep injection
*/

func EventsFactory() []*models.Event {
	e1 := &models.Event{
		Name:       "test1",
		Duration:   15 * time.Minute,
		AvailMasks: &values.AvailMasks{(&values.DateSlots{values.DateSlot{Start: time.Now(), End: time.Now().AddDate(0, 1, 0)}}).ToMask(values.AvailMaskINC)},
		UserID:     values.ID(1),
		Visible:    true,
	}
	e2 := &models.Event{
		Name:       "test2",
		Duration:   15 * time.Minute,
		AvailMasks: &values.AvailMasks{(&values.DateSlots{values.DateSlot{Start: time.Now(), End: time.Now().AddDate(0, 1, 0)}}).ToMask(values.AvailMaskINC)},
		UserID:     values.ID(1),
		Visible:    true,
	}
	e3 := &models.Event{
		Name:       "test3",
		Duration:   15 * time.Minute,
		AvailMasks: &values.AvailMasks{(&values.DateSlots{values.DateSlot{Start: time.Now(), End: time.Now().AddDate(0, 1, 0)}}).ToMask(values.AvailMaskINC)},
		UserID:     values.ID(1),
		Visible:    true,
	}
	return []*models.Event{e1, e2, e3}
}

func AuthFactory() []*models.Auth {

	a1 := &models.Auth{
		SourceID:     "test_auth_sourceId",
		Source:       "google",
		AccessToken:  values.Token("test1_refresh_token"),
		Avatar:       "http://example.com/test1_avatar.jpg",
		Name:         "test1 auth username",
		RefreshToken: values.Token("test1_refresh_token"),
		Expiry:       time.Now().Add(time.Hour),
		Email:        "test1@example.com",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		UserID:       values.ID(1),
	}
	return []*models.Auth{a1}
}

func UserFactory() []*models.User {
	u1 := &models.User{
		Name:  "test1",
		Email: "test1@example.com",
	}
	u2 := &models.User{
		Name:  "test2",
		Email: "test2@example.com",
	}
	u3 := &models.User{
		Name:  "test3",
		Email: "test3@example.com",
	}
	return []*models.User{u1, u2, u3}
}

func Seed(db *sq.Db) (err error) {

	//For easing seeding temp turn off foreign key checks
	defer func() {
		if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
			err = fmt.Errorf("foreign keys pragma: %w", err)
			panic(err)
		}
	}()
	if _, err := db.Exec(`PRAGMA foreign_keys = OFF;`); err != nil {
		return fmt.Errorf("foreign keys pragma: %w", err)
	}

	us := models.NewUserService(db)
	for _, u := range UserFactory() {
		_, err := us.CreateUser(u)
		if err != nil {
			return err
		}
	}

	es := models.NewEventService(db)
	for _, e := range EventsFactory() {
		_, err := es.CreateEvent(e)

		if err != nil {
			return err
		}
	}
	as := models.NewAuthService(db)
	for _, a := range AuthFactory() {
		_, err := as.CreateAuth(a)
		if err != nil {
			return err
		}
	}

	return nil
}
