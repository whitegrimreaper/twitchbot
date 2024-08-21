package main

import (
	"fmt"
	"time"
	"errors"
	"encoding/json"
	"gorm.io/gorm"
	"database/sql/driver"
	"gorm.io/driver/sqlite"

	"golang.org/x/exp/slices"
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
	reqId  int					`gorm:"primaryKey;autoIncrement:true"`
	UserID int
	BossID int
	BossKillsDone int
	BossKillsLeft int

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Table of all bosses, one entry per boss
// BossName is the most common name
// Other nicknames can be saved in BossNicknames table
// Shows the current 'queue' of boss kills
// Kills left -> kills left in queue to do
// kills done -> tracking progress over the whole event
// requests frozen - > stop allowing new requests for this boss
// Boss Multiple -> # of kills purchased per event
//         if this is set to 5, you can only request 5, 10, 15, etc.
//		   not 100% on this one ngl
// Cost -> number of points per BossMultiple # of kills
type BossEntry struct {
	BossID 	      	   int 			`gorm:"primaryKey"`
	BossName      	   string
	BossCost      	   int
	BossKillsLeft 	   int
	BossKillsDone 	   int
	BossRequestsFrozen bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	*a, _ = value.([]string)
	return nil
}

func (a StringArray) Value() (driver.Value, error) {
	val, err := json.Marshal([]string(a))
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Bosses may have many nicknames
// They don't change a ton so we can store them in a table and look later
// Primary key is ID, same as BossEntry
type BossNicknames struct {
	BossID 	      int 			 `gorm:"primaryKey"`
	BossName      string
	BossNicks	  StringArray	 //`gorm:"type:text[]"`

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

func reqQueueDBInit()(db *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("request_queue_test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&UserBossRequest{})
	return db
}

func bossNicksDBInit()(db *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("boss_nicks_test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&BossNicknames{})
	return db
}

func doesUserExist(targetUser int)(respCode int, respMessage string, exists bool) {
	var user UserPoints
	err  := PointsDB.First(&user, "user_id = ?", targetUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Printf("User doesn't exist!\n")
			return 1, "", false
		}
		fmt.Printf("User exist call returned a different error? %v\n", err)
		return -1, err.Error(), false
	}
	fmt.Printf("User exists!\n")
	return 0, "", true
}

func writePointGainEvent(targetUser int, pointAmount int)(respCode int, respMessage string) {
	var user UserPoints
	err := PointsDB.Where(UserPoints{UserID: targetUser}).FirstOrCreate(&user).Error
	if err != nil {
		fmt.Printf("Error in init %+v\n", err)
	}
	points := user.Points
	err = PointsDB.Model(&user).Update("Points", user.Points + pointAmount).Error
	if err != nil {
		fmt.Printf("Error in Update %+v\n", err)
	}
	fmt.Printf("Added %d points for user %d(%d), was at %d now has %d\n", pointAmount, targetUser, user.UserID, points, user.Points)
	return 0, ""
}

func findUserPoints(targetUser int)(respCode int, respMessage string, points int) {
	var user UserPoints
	err  := PointsDB.First(&user, "user_id = ?", targetUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Printf("User doesn't exist!\n")
			return 1, "", 0
		}
		fmt.Printf("User exist call returned a different error? %v\n", err)
		return -1, err.Error(), 0
	}
	fmt.Printf("User exists!\n")
	return 0, "", user.Points
}

func findBossInfo(bossId int)(respCode int, respMessage string, boss BossEntry) {
	err := BossDB.First(&boss, "boss_id = ?", bossId).Error
	if err != nil {
		return 1, "Boss not found!", boss
	}
	return 0, "", boss
}

func writePointSpendEvent(targetUser int, pointAmount int)(respCode int, respMessage string) {
	var user UserPoints
	err := PointsDB.Where(UserPoints{UserID: targetUser}).First(&user).Error
	if err != nil {
		return -1, err.Error()
	}
	user.Points -= pointAmount
	err = PointsDB.Model(&user).Update("points", user.Points).Error
	if err != nil {
		return -1, err.Error()
	}
	return 0, ""
}

func addBossKillsQueue(targetUser int, bossId int, bossKills int)(respCode int, respMessage string) {
	var boss BossEntry
	var queue UserBossRequest

	err := BossDB.Where("boss_id = ?", bossId).First(&boss).Error
	if err != nil {
		return -1, err.Error()
	}
	err = ReqQueueDB.Where(UserBossRequest{UserID: targetUser, BossID: bossId}).FirstOrCreate(&queue).Error
	if err != nil {
		return -1, err.Error()
	}
	err = ReqQueueDB.Model(&queue).Where("user_id = ? AND boss_id = ?", targetUser, bossId).
		Update("boss_kills_left", queue.BossKillsLeft + bossKills).Error
	if err != nil {
		return -1, err.Error()
	}
	return 0, ""
}

func addBossKillsMain(targetUser int, bossId int, bossKills int)(respCode int, respMessage string) {
	var boss BossEntry

	err := BossDB.Where("boss_id = ?", bossId).First(&boss).Error
	if err != nil {
		return -1, err.Error()
	}
	err = ReqQueueDB.Model(&boss).Where("boss_id = ?", bossId).
		Update("boss_kills_left", boss.BossKillsLeft + bossKills).Error
	if err != nil {
		return -1, err.Error()
	}
	return 0, ""
}

func isBossKillLocked(bossId int)(isLocked bool, respCode int, respMessage string) {
	var boss BossEntry

	err := BossDB.Where("boss_id = ?", bossId).First(&boss).Error
	if err != nil {
		return true, -1, err.Error()
	}
	return boss.BossRequestsFrozen, 0, ""
}

func getBossNicks(bossId int)(respCode int, respMessage string, nicks BossNicknames) {
	err := BossNickDB.Where("boss_id = ?", bossId).Find(&nicks).Error
	if err != nil {
		return -1, "Uhhhh", nicks
	}
	return 0, "", nicks
}

func getBossTrueName(inputString string)(respCode int, respMessage string, name string) {
	_, _, bosses := getBossList()
	for _, boss := range bosses {
		_, _, nicks := getBossNicks(boss.BossID)
		if(slices.Contains(nicks.BossNicks, inputString) || inputString == boss.BossName) {
			return 0, "", boss.BossName
		}
	}
	return -1, "Boss not found (nickname might be wrong)", ""
}

func getBossNameList()(respCode int, respMessage string, names []string) {
	var bosses []BossEntry
	err := BossDB.Find(&bosses).Error
	if err != nil {
		return -1, "Error finding literally anything", nil
	}
	for _, boss := range bosses {
		names = append(names, boss.BossName)
	}
	return 0, "", names
}

func getBossList()(respCode int, respMessage string, bosses []BossEntry) {
	err := BossDB.Find(&bosses).Error
	if err != nil {
		return -1, "Error finding literally anything", nil
	}
	return 0, "", bosses
}

func getBossWithName(name string)(respCode int, respMessage string, boss BossEntry) {
	err := BossDB.First(&boss, "boss_name = ?", name).Error
	if err != nil {
		return 1, "Probably no boss with that name bud", boss
	}
	return 0, "", boss
}