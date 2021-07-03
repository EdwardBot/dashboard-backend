package database

import (
	"context"
	"log"
	"time"
)

type User struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	JoinedAt      time.Time `json:"joined_at"`
	UserName      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	AvatarID      string    `json:"avatar"`
	Guilds        []int64   `json:"guilds"`
	PremiumType   int       `json:"premium_type"`
}

func (u *User) Create() {
	r, err := Conn.Query(context.TODO(), "insert into users (userid, joined, username, discriminator, avatar, guilds, premium) VALUES ($1,now(),$2,$3,$4,$5,$6) returning id", u.UserID, u.UserName, u.Discriminator, u.AvatarID, u.Guilds, u.PremiumType)
	if err != nil {
		log.Printf("SQL error: %s\n", err.Error())
	} else {
		_ = r.Scan(&u.ID)
	}
}

func FindUser(id int64) (User, error) {
	r, err := Conn.Query(context.TODO(), "select * from users where userid=$1", id)
	if err != nil {
		log.Printf("SQL error: %s\n", err.Error())
		return User{}, err
	} else {
		u := User{}
		t := ""
		e := r.Scan(&u.ID, &u.UserID, &t, &u.UserName, &u.Discriminator, &u.AvatarID, &u.Guilds, &u.PremiumType)
		if e != nil {
			log.Printf("SQL error: %s\n", e.Error())
			return u, e
		}
		log.Println(t)
		return u, nil
	}
}
