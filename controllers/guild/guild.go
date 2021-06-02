package guild

import (
	"github.com/bwmarrin/discordgo"
	"github.com/edward-backend/database"
	"github.com/edward-backend/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

var (
	errorMsg = gin.H{
		"status":   "errorMsg",
		"errorMsg": "Invalid id provided",
	}
	roles = gin.H{}
)

func HandleGuild(c *gin.Context) {
	if !c.MustGet("hasAuth").(bool) {
		return
	}
	i, e := strconv.ParseInt(c.Param("id"), 10, 64)
	if e != nil {
		c.JSON(500, errorMsg)
		return
	}
	g, _ := database.FindGuilds(i)
	client := utils.GetDiscordInstance()
	//Loop through the users guilds
	for i := range g {
		guild := g[i]
		member, _ := client.GuildMember(guild.GID, c.Param("id"))
		//Get all the roles
		gRoles, _ := client.GuildRoles(guild.GID)
		//Check his roles
		var tmpPerm int64 = 0
		for r := range member.Roles {
			roleId := member.Roles[r]
			//Check if it is in the cache
			if roles[roleId] != nil {
				tmpPerm |= (roles[roleId]).(gin.H)["role"].(*discordgo.Role).Permissions
			} else {
				role := utils.Find(gRoles, func(e *discordgo.Role) bool {
					return e.ID == roleId
				})
				roles[roleId] = gin.H{
					"role": role,
				}
				tmpPerm |= role.Permissions
			}
		}
		guild.Permissions = strconv.FormatInt(tmpPerm, 2)
		g[i] = guild
	}
	c.JSON(200, g)
}
