// khan
// https://github.com/topfreegames/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package models

import (
	"time"

	"gopkg.in/gorp.v1"
)

//Game identifies uniquely one game
type Game struct {
	ID                            int    `db:"id"`
	PublicID                      string `db:"public_id"`
	Name                          string `db:"name"`
	MinMembershipLevel            int    `db:"min_membership_level"`
	MaxMembershipLevel            int    `db:"max_membership_level"`
	MinLevelToAcceptApplication   int    `db:"min_level_to_accept_application"`
	MinLevelToCreateInvitation    int    `db:"min_level_to_create_invitation"`
	MinLevelOffsetToPromoteMember int    `db:"min_level_offset_to_promote_member"`
	MinLevelOffsetToDemoteMember  int    `db:"min_level_offset_to_demote_member"`
	AllowApplication              bool   `db:"allow_application"`
	Metadata                      string `db:"metadata"`
	CreatedAt                     int64  `db:"created_at"`
	UpdatedAt                     int64  `db:"updated_at"`
}

//PreInsert populates fields before inserting a new game
func (p *Game) PreInsert(s gorp.SqlExecutor) error {
	p.CreatedAt = time.Now().UnixNano()
	p.UpdatedAt = p.CreatedAt
	return nil
}

//PreUpdate populates fields before updating a game
func (p *Game) PreUpdate(s gorp.SqlExecutor) error {
	p.UpdatedAt = time.Now().UnixNano()
	return nil
}

//GetGameByID returns a game by id
func GetGameByID(db DB, id int) (*Game, error) {
	obj, err := db.Get(Game{}, id)
	if err != nil || obj == nil {
		return nil, &ModelNotFoundError{"Game", id}
	}

	game := obj.(*Game)
	return game, nil
}

//GetGameByPublicID returns a game by their public id
func GetGameByPublicID(db DB, publicID string) (*Game, error) {
	var game Game
	err := db.SelectOne(&game, "SELECT * FROM games WHERE public_id=$1", publicID)
	if err != nil || &game == nil {
		return nil, &ModelNotFoundError{"Game", publicID}
	}
	return &game, nil
}

//CreateGame creates a new game
func CreateGame(db DB, publicID, name, metadata string,
	minLevel, maxLevel, minLevelAccept, minLevelCreate, minOffsetPromote, minOffsetDemote int,
	allowApplication bool,
) (*Game, error) {
	game := &Game{
		PublicID:                      publicID,
		Name:                          name,
		MinMembershipLevel:            minLevel,
		MaxMembershipLevel:            maxLevel,
		MinLevelToAcceptApplication:   minLevelAccept,
		MinLevelToCreateInvitation:    minLevelCreate,
		MinLevelOffsetToPromoteMember: minOffsetPromote,
		MinLevelOffsetToDemoteMember:  minOffsetDemote,
		AllowApplication:              allowApplication,
		Metadata:                      metadata,
	}
	err := db.Insert(game)
	if err != nil {
		return nil, err
	}
	return game, nil
}

//UpdateGame updates an existing game
func UpdateGame(db DB, publicID, name, metadata string,
	minLevel, maxLevel, minLevelAccept, minLevelCreate, minOffsetPromote, minOffsetDemote int,
	allowApplication bool,
) (*Game, error) {
	game, err := GetGameByPublicID(db, publicID)

	if err != nil {
		return nil, err
	}

	game.Name = name
	game.MinMembershipLevel = minLevel
	game.MaxMembershipLevel = maxLevel
	game.MinLevelToAcceptApplication = minLevelAccept
	game.MinLevelToCreateInvitation = minLevelCreate
	game.MinLevelOffsetToPromoteMember = minOffsetPromote
	game.MinLevelOffsetToDemoteMember = minOffsetDemote
	game.AllowApplication = allowApplication
	game.Metadata = metadata

	_, err = db.Update(game)

	if err != nil {
		return nil, err
	}

	return game, nil
}