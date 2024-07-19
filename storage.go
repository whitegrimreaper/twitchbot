package main

import (
	"fmt"
	"time"
	"errors"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

// Table of all users, contains users and the points they have
type UserPoints struct {
	// ignoring this until i get a cleaner model with a better init func
	//Username string
	UserID int     				`gorm:"primaryKey"`
	Points int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Need table to store relationship between users and bosses they request
// uniquely keyed to each user + boss pair (each user can have up to one entry per boss)
// i.e. Pigeon requests 50 Kril kills and 20 ED2 runs
//       -> two entries under Pigeon's UserID, one for Kril, one for ED2
type UserBossRequest struct {
	UUID   int					`gorm:"primaryKey"`
	UserID int
	BossID int
	BossKillsDone int
	BossKillsLeft int

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Table of all bosses, one entry per boss
// Shows the current 'queue' of boss kills
// Kills left -> kills left in queue to do
// kills done -> tracking progress over the whole event
// Boss Multiple -> # of kills purchased per event
//         if this is set to 5, you can only request 5, 10, 15, etc.
//		   not 100% on this one ngl
// Cost -> number of points per BossMultiple # of kills
type BossEntry struct {
	BossID 	      int 			`gorm:"primaryKey"`
	BossName      string
	BossCost      int
	BossKillsLeft int
	BossKillsDone int

	CreatedAt time.Time
	UpdatedAt time.Time
}

func pointsDBInit()(db *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("points_test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&UserPoints{})
	return db
}

func bossDBInit()(db *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("boss_test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&BossEntry{})
	return db
}

func doesUserExist(targetUser int)(respCode bool, respMessage string, exists bool) {
	var user UserPoints
	err  := PointsDB.First(&user, "user_id = ?", targetUser).Error
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
	err := PointsDB.Where(UserPoints{UserID: targetUser}).FirstOrCreate(&user).Error
	if err != nil {
		fmt.Printf("Error in init %+v\n", err)
	}
	points := user.Points
	err = PointsDB.Model(&user).Update("Points", user.Points+pointAmount).Error
	if err != nil {
		fmt.Printf("Error in Update %+v\n", err)
	}
	fmt.Printf("Added %d points for user %d(%d), was at %d now has %d\n", pointAmount, targetUser, user.UserID, points, user.Points)
	return false, ""
}

func findUserPoints(targetUser int)(respCode bool, respMessage string, points int) {
	var user UserPoints
	err  := PointsDB.First(&user, "user_id = ?", targetUser).Error
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
	err := PointsDB.Where(UserPoints{UserID: targetUser}).First(&user).Error
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}
