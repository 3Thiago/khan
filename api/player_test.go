// khan
// https://github.com/topfreegames/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	. "github.com/franela/goblin"
	"github.com/satori/go.uuid"
	"github.com/topfreegames/khan/models"
)

func TestPlayerHandler(t *testing.T) {
	g := Goblin(t)

	testDb, err := models.GetTestDB()
	g.Assert(err == nil).IsTrue()

	g.Describe("Create Player Handler", func() {
		g.It("Should create player", func() {
			a := GetDefaultTestApp()
			game := models.GameFactory.MustCreate().(*models.Game)
			err := a.Db.Insert(game)
			AssertNotError(g, err)

			payload := map[string]interface{}{
				"publicID": randomdata.FullName(randomdata.RandomGender),
				"name":     randomdata.FullName(randomdata.RandomGender),
				"metadata": map[string]interface{}{"x": 1},
			}
			res := PostJSON(a, GetGameRoute(game.PublicID, "/players"), t, payload)

			g.Assert(res.Raw().StatusCode).Equal(http.StatusOK)
			var result map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsTrue()
			g.Assert(result["publicID"]).Equal(payload["publicID"].(string))

			dbPlayer, err := models.GetPlayerByPublicID(
				a.Db, game.PublicID, payload["publicID"].(string),
			)
			AssertNotError(g, err)
			g.Assert(dbPlayer.GameID).Equal(game.PublicID)
			g.Assert(dbPlayer.PublicID).Equal(payload["publicID"])
			g.Assert(dbPlayer.Name).Equal(payload["name"])
			g.Assert(dbPlayer.Metadata).Equal(payload["metadata"])
		})

		g.It("Should not create player if missing parameters", func() {
			a := GetDefaultTestApp()
			route := GetGameRoute("game-id", "/players")
			res := PostJSON(a, route, t, map[string]interface{}{})

			g.Assert(res.Raw().StatusCode).Equal(http.StatusBadRequest)
			var result map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsFalse()
			g.Assert(result["reason"]).Equal("publicID is required, name is required, metadata is required")
		})

		g.It("Should not create player if invalid payload", func() {
			a := GetDefaultTestApp()
			route := GetGameRoute("game-id", "/players")
			res := PostBody(a, route, t, "invalid")

			g.Assert(res.Raw().StatusCode).Equal(http.StatusBadRequest)
			var result map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsFalse()
			g.Assert(strings.Contains(result["reason"].(string), "While trying to read JSON")).IsTrue()
		})

		g.It("Should not create player if invalid data", func() {
			a := GetDefaultTestApp()
			game := models.GameFactory.MustCreate().(*models.Game)
			err := a.Db.Insert(game)
			AssertNotError(g, err)

			payload := map[string]interface{}{
				"publicID": strings.Repeat("s", 256),
				"name":     randomdata.FullName(randomdata.RandomGender),
				"metadata": map[string]interface{}{"x": 1},
			}
			res := PostJSON(a, GetGameRoute(game.PublicID, "/players"), t, payload)

			g.Assert(res.Raw().StatusCode).Equal(http.StatusInternalServerError)
			var result map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsFalse()
			g.Assert(result["reason"]).Equal("pq: value too long for type character varying(255)")
		})
	})

	g.Describe("Update Player Handler", func() {
		g.It("Should update player", func() {
			a := GetDefaultTestApp()
			_, player, err := models.CreatePlayerFactory(a.Db, "")
			AssertNotError(g, err)

			metadata := map[string]interface{}{"y": 10}
			payload := map[string]interface{}{
				"name":     player.Name,
				"metadata": metadata,
			}

			route := GetGameRoute(player.GameID, fmt.Sprintf("/players/%s", player.PublicID))
			res := PutJSON(a, route, t, payload)
			g.Assert(res.Raw().StatusCode).Equal(http.StatusOK)
			var result map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsTrue()

			dbPlayer, err := models.GetPlayerByPublicID(a.Db, player.GameID, player.PublicID)
			AssertNotError(g, err)
			g.Assert(dbPlayer.GameID).Equal(player.GameID)
			g.Assert(dbPlayer.PublicID).Equal(player.PublicID)
			g.Assert(dbPlayer.Name).Equal(player.Name)
			g.Assert(dbPlayer.Metadata).Equal(metadata)
		})

		g.It("Should not update player if missing parameters", func() {
			a := GetDefaultTestApp()
			route := GetGameRoute("game-id", "/players/player-id")
			res := PutJSON(a, route, t, map[string]interface{}{})

			g.Assert(res.Raw().StatusCode).Equal(http.StatusBadRequest)
			var result map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsFalse()
			g.Assert(result["reason"]).Equal("name is required, metadata is required")
		})

		g.It("Should not update player if invalid payload", func() {
			a := GetDefaultTestApp()
			route := GetGameRoute("game-id", "/players/fake")
			res := PutBody(a, route, t, "invalid")

			g.Assert(res.Raw().StatusCode).Equal(http.StatusBadRequest)
			var result map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsFalse()
			g.Assert(strings.Contains(result["reason"].(string), "While trying to read JSON")).IsTrue()
		})

		g.It("Should not update player if invalid data", func() {
			a := GetDefaultTestApp()
			_, player, err := models.CreatePlayerFactory(a.Db, "")
			AssertNotError(g, err)

			payload := map[string]interface{}{
				"publicID": player.PublicID,
				"name":     strings.Repeat("s", 256),
				"metadata": map[string]interface{}{},
			}
			route := GetGameRoute(player.GameID, fmt.Sprintf("/players/%s", player.PublicID))
			res := PutJSON(a, route, t, payload)

			g.Assert(res.Raw().StatusCode).Equal(http.StatusInternalServerError)
			var result map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsFalse()
			g.Assert(result["reason"]).Equal("pq: value too long for type character varying(255)")
		})
	})

	g.Describe("Retrieve Player", func() {
		g.It("Should retrieve player", func() {
			a := GetDefaultTestApp()
			gameID := uuid.NewV4().String()
			player, err := models.GetTestPlayerWithMemberships(testDb, gameID, 5, 2, 3, 8)
			g.Assert(err == nil).IsTrue()

			route := GetGameRoute(player.GameID, fmt.Sprintf("/players/%s", player.PublicID))
			res := Get(a, route, t)

			g.Assert(res.Raw().StatusCode).Equal(http.StatusOK)
			var playerDetails map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &playerDetails)
			g.Assert(playerDetails["success"]).IsTrue()

			// Player Details
			g.Assert(playerDetails["publicID"]).Equal(player.PublicID)
			g.Assert(playerDetails["name"]).Equal(player.Name)
			g.Assert(playerDetails["metadata"] != nil).IsTrue()

			//Memberships
			g.Assert(len(playerDetails["memberships"].([]interface{}))).Equal(18)

			clans := playerDetails["clans"].(map[string]interface{}) // can't be map[string]interface{}
			approved := clans["approved"].([]interface{})
			denied := clans["denied"].([]interface{})
			banned := clans["banned"].([]interface{})
			pendingApplications := clans["pendingApplications"].([]interface{})
			pendingInvites := clans["pendingInvites"].([]interface{})

			g.Assert(len(approved)).Equal(5)
			g.Assert(len(denied)).Equal(2)
			g.Assert(len(banned)).Equal(3)
			g.Assert(len(pendingApplications)).Equal(0)
			g.Assert(len(pendingInvites)).Equal(8)
		})
		g.It("Should return 404 for invalid player", func() {
			a := GetDefaultTestApp()
			route := GetGameRoute("some-game", "/players/invalid-player")
			res := Get(a, route, t)

			g.Assert(res.Raw().StatusCode).Equal(http.StatusNotFound)

			var playerDetails map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &playerDetails)
			g.Assert(playerDetails["success"]).IsFalse()
			g.Assert(playerDetails["reason"]).Equal("Player was not found with id: invalid-player")
		})
	})

	g.Describe("Player Hooks", func() {
		g.It("Should call create player hook", func() {
			hooks, err := models.GetHooksForRoutes(testDb, []string{
				"http://localhost:52525/playercreated",
			}, models.PlayerCreatedHook)
			g.Assert(err == nil).IsTrue()
			responses := startRouteHandler([]string{"/playercreated"}, 52525)

			app := GetDefaultTestApp()
			time.Sleep(time.Second)

			gameID := hooks[0].GameID
			payload := map[string]interface{}{
				"publicID": randomdata.FullName(randomdata.RandomGender),
				"name":     randomdata.FullName(randomdata.RandomGender),
				"metadata": map[string]interface{}{"x": "a"},
			}
			res := PostJSON(app, GetGameRoute(gameID, "/players"), t, payload)

			g.Assert(res.Raw().StatusCode).Equal(http.StatusOK)
			var result map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsTrue()
			g.Assert(result["publicID"]).Equal(payload["publicID"].(string))

			app.Dispatcher.Wait()

			g.Assert(len(*responses)).Equal(1)

			player := (*responses)[0]["payload"].(map[string]interface{})
			g.Assert(player["gameID"]).Equal(gameID)
			g.Assert(player["publicID"]).Equal(payload["publicID"])
			g.Assert(player["name"]).Equal(payload["name"])
			g.Assert(str(player["membershipCount"])).Equal("0")
			g.Assert(str(player["ownershipCount"])).Equal("0")
			playerMetadata := player["metadata"].(map[string]interface{})
			metadata := payload["metadata"].(map[string]interface{})
			for k, v := range playerMetadata {
				g.Assert(v).Equal(metadata[k])
			}
		})

		g.It("Should call update player hook", func() {
			hooks, err := models.GetHooksForRoutes(testDb, []string{
				"http://localhost:52525/updated",
			}, models.PlayerUpdatedHook)
			g.Assert(err == nil).IsTrue()
			responses := startRouteHandler([]string{"/updated"}, 52525)

			player := models.PlayerFactory.MustCreateWithOption(map[string]interface{}{"GameID": hooks[0].GameID}).(*models.Player)
			err = testDb.Insert(player)
			AssertNotError(g, err)

			app := GetDefaultTestApp()
			time.Sleep(time.Second)

			gameID := hooks[0].GameID
			payload := map[string]interface{}{
				"publicID": player.PublicID,
				"name":     player.Name,
				"metadata": player.Metadata,
			}
			res := PutJSON(app, GetGameRoute(gameID, fmt.Sprintf("/players/%s", player.PublicID)), t, payload)

			g.Assert(res.Raw().StatusCode).Equal(http.StatusOK)
			var result map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsTrue()

			app.Dispatcher.Wait()
			g.Assert(len(*responses)).Equal(1)

			playerPayload := (*responses)[0]["payload"].(map[string]interface{})
			g.Assert(playerPayload["gameID"]).Equal(gameID)
			g.Assert(playerPayload["publicID"]).Equal(payload["publicID"])
			g.Assert(playerPayload["name"]).Equal(payload["name"])
			g.Assert(str(playerPayload["membershipCount"])).Equal("0")
			g.Assert(str(playerPayload["ownershipCount"])).Equal("0")
			playerMetadata := playerPayload["metadata"].(map[string]interface{})
			metadata := payload["metadata"].(map[string]interface{})
			for k, v := range playerMetadata {
				g.Assert(v).Equal(metadata[k])
			}
		})
	})
}
