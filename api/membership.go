// khan
// https://github.com/topfreegames/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package api

import (
	"time"

	"github.com/kataras/iris"
	"github.com/topfreegames/khan/models"
	"github.com/uber-go/zap"
)

// Helper methods are located in membership_helpers module.
// This module is only for handlers

// ApplyForMembershipHandler is the handler responsible for applying for new memberships
func ApplyForMembershipHandler(app *App) func(c *iris.Context) {
	return func(c *iris.Context) {
		start := time.Now()
		gameID := c.Param("gameID")
		clanPublicID := c.Param("clanPublicID")

		l := app.Logger.With(
			zap.String("source", "membershipHandler"),
			zap.String("operation", "applyForMembership"),
			zap.String("gameID", gameID),
			zap.String("clanPublicID", clanPublicID),
		)

		var payload applyForMembershipPayload
		if err := LoadJSONPayload(&payload, c, l); err != nil {
			FailWith(400, err.Error(), c)
			return
		}

		optional, err := getMembershipOptionalParameters(app, c)
		if err != nil {
			FailWith(400, err.Error(), c)
			return
		}

		l = l.With(
			zap.String("level", payload.Level),
			zap.String("playerPublicID", payload.PlayerPublicID),
		)

		l.Debug("Getting DB connection...")
		db, err := GetCtxDB(c)
		if err != nil {
			l.Error("Failed to connect to DB.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}
		l.Debug("DB Connection successful.")

		game, err := app.GetGame(gameID)
		if err != nil {
			l.Warn("Could not find game.")
			FailWith(404, err.Error(), c)
			return
		}

		l.Debug("Applying for membership...")
		membership, err := models.CreateMembership(
			db,
			game,
			gameID,
			payload.Level,
			payload.PlayerPublicID,
			clanPublicID,
			payload.PlayerPublicID,
			optional.Message,
		)

		if err != nil {
			l.Error("Membership application failed.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}

		l.Info("Membership application successful.")

		err = dispatchMembershipHookByID(
			app, db, models.MembershipApplicationCreatedHook,
			membership.GameID, membership.ClanID, membership.PlayerID,
			membership.RequestorID, membership.Message,
		)
		if err != nil {
			l.Error("Membership application created dispatch hook failed.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}

		l.Info(
			"Membership application created successfully.",
			zap.Duration("duration", time.Now().Sub(start)),
		)

		SucceedWith(map[string]interface{}{}, c)
	}
}

// InviteForMembershipHandler is the handler responsible for creating new memberships
func InviteForMembershipHandler(app *App) func(c *iris.Context) {
	return func(c *iris.Context) {
		start := time.Now()
		gameID := c.Param("gameID")
		clanPublicID := c.Param("clanPublicID")

		l := app.Logger.With(
			zap.String("source", "membershipHandler"),
			zap.String("operation", "inviteForMembership"),
			zap.String("gameID", gameID),
			zap.String("clanPublicID", clanPublicID),
		)

		var payload inviteForMembershipPayload
		if err := LoadJSONPayload(&payload, c, l); err != nil {
			FailWith(400, err.Error(), c)
			return
		}

		optional, err := getMembershipOptionalParameters(app, c)
		if err != nil {
			FailWith(400, err.Error(), c)
			return
		}

		l = l.With(
			zap.String("level", payload.Level),
			zap.String("playerPublicID", payload.PlayerPublicID),
			zap.String("requestorPublicID", payload.RequestorPublicID),
		)

		l.Debug("Getting DB connection...")
		db, err := GetCtxDB(c)
		if err != nil {
			l.Error("Failed to connect to DB.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}
		l.Debug("DB Connection successful.")

		game, err := app.GetGame(gameID)
		if err != nil {
			l.Warn("Could not find game.")
			FailWith(404, err.Error(), c)
			return
		}

		l.Debug("Inviting for membership...")
		membership, err := models.CreateMembership(
			db,
			game,
			gameID,
			payload.Level,
			payload.PlayerPublicID,
			clanPublicID,
			payload.RequestorPublicID,
			optional.Message,
		)

		if err != nil {
			l.Error("Membership invitation failed.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}

		err = dispatchMembershipHookByID(
			app, db, models.MembershipApplicationCreatedHook,
			membership.GameID, membership.ClanID, membership.PlayerID,
			membership.RequestorID, membership.Message,
		)
		if err != nil {
			l.Error("Membership invitation dispatch hook failed.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}

		l.Info(
			"Membership invitation created successfully.",
			zap.Duration("duration", time.Now().Sub(start)),
		)

		SucceedWith(map[string]interface{}{}, c)
	}
}

// ApproveOrDenyMembershipApplicationHandler is the handler responsible for approving or denying a membership invitation
func ApproveOrDenyMembershipApplicationHandler(app *App) func(c *iris.Context) {
	return func(c *iris.Context) {
		start := time.Now()
		action := c.Param("action")
		gameID := c.Param("gameID")
		clanPublicID := c.Param("clanPublicID")

		l := app.Logger.With(
			zap.String("source", "membershipHandler"),
			zap.String("operation", "approveOrDenyApplication"),
			zap.String("gameID", gameID),
			zap.String("clanPublicID", clanPublicID),
			zap.String("action", action),
		)

		var payload basePayloadWithRequestorAndPlayerPublicIDs
		if err := LoadJSONPayload(&payload, c, l); err != nil {
			FailWith(400, err.Error(), c)
			return
		}

		l = l.With(
			zap.String("playerPublicID", payload.PlayerPublicID),
			zap.String("requestorPublicID", payload.RequestorPublicID),
		)

		l.Debug("Getting DB connection...")
		db, err := GetCtxDB(c)
		if err != nil {
			l.Error("Failed to connect to DB.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}
		l.Debug("DB Connection successful.")

		game, err := app.GetGame(gameID)
		if err != nil {
			l.Warn("Could not find game.")
			FailWith(404, err.Error(), c)
			return
		}

		l.Debug("Approving/Denying membership application.")
		membership, err := models.ApproveOrDenyMembershipApplication(
			db,
			game,
			gameID,
			payload.PlayerPublicID,
			clanPublicID,
			payload.RequestorPublicID,
			action,
		)

		if err != nil {
			l.Error("Approving/Denying membership application failed.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}

		l.Debug("Retrieving requestor details.")
		requestor, err := models.GetPlayerByPublicID(db, gameID, payload.RequestorPublicID)
		if err != nil {
			l.Error("Requestor details retrieval failed.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}
		l.Debug("Requestor details retrieved successfully.")

		hookType := models.MembershipApprovedHook
		if action == "deny" {
			hookType = models.MembershipDeniedHook
		}

		err = dispatchMembershipHookByID(
			app, db, hookType,
			membership.GameID, membership.ClanID, membership.PlayerID,
			requestor.ID, membership.Message,
		)
		if err != nil {
			l.Error("Membership approved/denied application dispatch hook failed.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}

		l.Info(
			"Membership application approved/denied successfully.",
			zap.Duration("duration", time.Now().Sub(start)),
		)

		SucceedWith(map[string]interface{}{}, c)
	}
}

// ApproveOrDenyMembershipInvitationHandler is the handler responsible for approving or denying a membership invitation
func ApproveOrDenyMembershipInvitationHandler(app *App) func(c *iris.Context) {
	return func(c *iris.Context) {
		start := time.Now()
		action := c.Param("action")
		gameID := c.Param("gameID")
		clanPublicID := c.Param("clanPublicID")

		l := app.Logger.With(
			zap.String("source", "membershipHandler"),
			zap.String("operation", "approveOrDenyInvitation"),
			zap.String("gameID", gameID),
			zap.String("clanPublicID", clanPublicID),
			zap.String("action", action),
		)

		var payload approveOrDenyMembershipInvitationPayload
		if err := LoadJSONPayload(&payload, c, l); err != nil {
			FailWith(400, err.Error(), c)
			return
		}

		l = l.With(
			zap.String("playerPublicID", payload.PlayerPublicID),
		)

		l.Debug("Getting DB connection...")
		db, err := GetCtxDB(c)
		if err != nil {
			l.Error("Failed to connect to DB.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}
		l.Debug("DB Connection successful.")

		game, err := app.GetGame(gameID)
		if err != nil {
			l.Warn("Could not find game.")
			FailWith(404, err.Error(), c)
			return
		}

		l.Debug("Approving/Denying membership invitation...")
		membership, err := models.ApproveOrDenyMembershipInvitation(
			db,
			game,
			gameID,
			payload.PlayerPublicID,
			clanPublicID,
			action,
		)

		if err != nil {
			l.Error("Membership invitation approval/deny failed.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}

		hookType := models.MembershipApprovedHook
		if action == "deny" {
			hookType = models.MembershipDeniedHook
		}

		err = dispatchMembershipHookByID(
			app, db, hookType,
			membership.GameID, membership.ClanID, membership.PlayerID,
			membership.PlayerID, membership.Message,
		)
		if err != nil {
			FailWith(500, err.Error(), c)
			return
		}

		l.Info(
			"Membership invitation approved/denied successfully.",
			zap.Duration("duration", time.Now().Sub(start)),
		)

		SucceedWith(map[string]interface{}{}, c)
	}
}

// DeleteMembershipHandler is the handler responsible for deleting a member
func DeleteMembershipHandler(app *App) func(c *iris.Context) {
	return func(c *iris.Context) {
		start := time.Now()
		clanPublicID := c.Param("clanPublicID")

		l := app.Logger.With(
			zap.String("source", "membershipHandler"),
			zap.String("operation", "deleteMembership"),
			zap.String("clanPublicID", clanPublicID),
		)

		payload, game, status, err := getPayloadAndGame(app, c, l)
		if err != nil {
			FailWith(status, err.Error(), c)
			return
		}

		l = l.With(
			zap.String("gameID", game.PublicID),
			zap.String("playerPublicID", payload.PlayerPublicID),
			zap.String("requestorPublicID", payload.RequestorPublicID),
		)

		l.Debug("Getting DB connection...")
		db, err := GetCtxDB(c)
		if err != nil {
			l.Error("Failed to connect to DB.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}
		l.Debug("DB Connection successful.")

		l.Debug("Deleting membership...")
		err = models.DeleteMembership(
			db,
			game,
			game.PublicID,
			payload.PlayerPublicID,
			clanPublicID,
			payload.RequestorPublicID,
		)

		if err != nil {
			l.Error("Membership delete failed.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}

		err = dispatchMembershipHookByPublicID(
			app, db, models.MembershipLeftHook,
			game.PublicID, clanPublicID, payload.PlayerPublicID,
			payload.RequestorPublicID,
		)
		if err != nil {
			l.Error("Membership deleted hook dispatch failed.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}

		l.Info(
			"Membership deleted successfully.",
			zap.Duration("duration", time.Now().Sub(start)),
		)

		SucceedWith(map[string]interface{}{}, c)
	}
}

// PromoteOrDemoteMembershipHandler is the handler responsible for promoting or demoting a member
func PromoteOrDemoteMembershipHandler(app *App, action string) func(c *iris.Context) {
	return func(c *iris.Context) {
		start := time.Now()
		clanPublicID := c.Param("clanPublicID")

		l := app.Logger.With(
			zap.String("source", "membershipHandler"),
			zap.String("operation", "promoteOrDemoteMembership"),
			zap.String("clanPublicID", clanPublicID),
			zap.String("action", action),
		)

		payload, game, status, err := getPayloadAndGame(app, c, l)
		if err != nil {
			FailWith(status, err.Error(), c)
			return
		}

		l = l.With(
			zap.String("gameID", game.PublicID),
			zap.String("playerPublicID", payload.PlayerPublicID),
			zap.String("requestorPublicID", payload.RequestorPublicID),
		)

		l.Debug("Getting DB connection...")
		db, err := GetCtxDB(c)
		if err != nil {
			l.Error("Failed to connect to DB.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}
		l.Debug("DB Connection successful.")

		l.Debug("Promoting/Demoting member...")
		membership, err := models.PromoteOrDemoteMember(
			db,
			game,
			game.PublicID,
			payload.PlayerPublicID,
			clanPublicID,
			payload.RequestorPublicID,
			action,
		)

		if err != nil {
			l.Error("Member promotion/demotion failed.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}
		l.Info("Member promoted/demoted successful.")

		l.Debug("Retrieving promoter/demoter member...")
		requestor, err := models.GetPlayerByPublicID(db, membership.GameID, payload.RequestorPublicID)
		if err != nil {
			l.Error("Promoter/Demoter member retrieval failed.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}
		l.Debug("Promoter/Demoter member retrieved successfully.")

		hookType := models.MembershipPromotedHook
		if action == "demote" {
			hookType = models.MembershipDemotedHook
		}

		err = dispatchMembershipHookByID(
			app, db, hookType,
			membership.GameID, membership.ClanID, membership.PlayerID,
			requestor.ID, membership.Message,
		)
		if err != nil {
			l.Error("Promote/Demote member hook dispatch failed.", zap.Error(err))
			FailWith(500, err.Error(), c)
			return
		}

		l.Info(
			"Member promoted/demoted successfully.",
			zap.Duration("duration", time.Now().Sub(start)),
		)

		SucceedWith(map[string]interface{}{
			"level": membership.Level,
		}, c)
	}
}
