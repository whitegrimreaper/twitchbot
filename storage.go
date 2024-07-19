package main

import (
	"fmt"
	"time"
	"errors"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

type UserPoints struct {
	// ignoring this until i get a cleaner model with a better init func
	//Username string
	UserID int     				`gorm:"primaryKey"`
	Points int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func dbinit()(db *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("example.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&UserPoints{})
	return db
}

func doesUserExist(targetUser int)(respCode bool, respMessage string, exists bool) {
	var user UserPoints
	err  := Db.First(&user, "user_id = ?", targetUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Printf("User doesn't exist!\n")
			return false, "", false
		}
		fmt.Printf("User exist call returned a different error? %v\n", err)
		return false, err.Error(), false
	}
	fmt.Printf("User exists!\n")
	return true, "", true
}

func writePointGainEvent(targetUser int, pointAmount int)(respCode bool, respMessage string) {
	var user UserPoints
	err := Db.Where(UserPoints{UserID: targetUser}).FirstOrCreate(&user).Error
	if err != nil {
		fmt.Printf("Error in init %+v\n", err)
	}
	points := user.Points
	err = Db.Model(&user).Update("Points", user.Points+pointAmount).Error
	if err != nil {
		fmt.Printf("Error in Update %+v\n", err)
	}
	fmt.Printf("Added %d points for user %d(%d), was at %d now has %d\n", pointAmount, targetUser, user.UserID, points, user.Points)
	return false, ""
}

func findUserPoints(targetUser int)(respCode bool, respMessage string, points int) {
	var user UserPoints
	err  := Db.First(&user, "user_id = ?", targetUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Printf("User doesn't exist!\n")
			return false, "", 0
		}
		fmt.Printf("User exist call returned a different error? %v\n", err)
		return false, err.Error(), 0
	}
	fmt.Printf("User exists!\n")
	return true, "", user.Points
}

func writePointSpendEvent(targetUser int, pointAmount int)(respCode bool, respMessage string) {
	var user UserPoints
	err := Db.Where(UserPoints{UserID: targetUser}).First(&user).Error
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}
