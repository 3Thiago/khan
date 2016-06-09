package models

import (
	"time"

	"gopkg.in/gorp.v1"
)

//Membership relates a player to a clan
type Membership struct {
	ID        int    `db:"id"`
	GameID    string `db:"game_id"`
	Level     int    `db:"membership_level"`
	Approved  bool   `db:"approved"`
	Denied    bool   `db:"denied"`
	PlayerID  int    `db:"player_id"`
	ClanID    int    `db:"clan_id"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
}

//PreInsert populates fields before inserting a new clan
func (m *Membership) PreInsert(s gorp.SqlExecutor) error {
	m.CreatedAt = time.Now().UnixNano()
	m.UpdatedAt = m.CreatedAt
	return nil
}

//PreUpdate populates fields before updating a clan
func (m *Membership) PreUpdate(s gorp.SqlExecutor) error {
	m.UpdatedAt = time.Now().UnixNano()
	return nil
}

//GetMembershipByID returns a membership by id
func GetMembershipByID(id int) (*Membership, error) {
	obj, err := db.Get(Membership{}, id)
	if err != nil {
		return nil, err
	}
	return obj.(*Membership), nil
}
