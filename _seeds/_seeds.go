package seeds

import (
	"fmt"
	"scheduleme/config"
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

func Seed(cfg *config.ConfigStruct) {
	// config.SetConfigForTestingOnly(config.ConfigFromEnv())

	db, err := sq.NewOpenDB(cfg.Dsn)
	fmt.Printf("seeding dsn: %s\n", cfg.Dsn)
	if err != nil {
		panic(err)
	}
	us := models.NewUserService(db)
	for _, u := range UserFactory() {
		_, err := us.CreateUser(u)
		if err != nil {
			panic(err)
		}
	}

	es := models.NewEventService(db)
	for _, e := range EventsFactory() {
		_, err := es.CreateEvent(e)
		if err != nil {
			panic(err)
		}
	}
	println("done seeding")
}
