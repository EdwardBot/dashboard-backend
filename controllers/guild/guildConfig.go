package guild

import (
	"errors"
	"github.com/edward-backend/database"
	"github.com/edward-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func HandleGetGuildConfig(ctx *gin.Context) {
	if !ctx.MustGet("hasAuth").(bool) {
		return
	}
	var conf database.GuildConfig
	var wallet database.Wallet
	r := database.Conn.Table("guild-configs").First(&conf, ctx.Param("id"))
	database.Conn.Model(&database.Wallet{}).Where("guild = ? and userid = ?", ctx.Param("id"), ctx.MustGet("userId")).First(&wallet)

	wallet.GuildId = strconv.FormatInt(wallet.GID, 10)
	wallet.UserId = strconv.FormatInt(wallet.UID, 10)

	conf.WelcomeChannel = strconv.FormatInt(int64(conf.Wch), 10)
	conf.LeaveChannel = strconv.FormatInt(int64(conf.LCh), 10)
	conf.LogChannel = strconv.FormatInt(int64(conf.LoCh), 10)

	s, _ := utils.GetDiscordInstance().GuildChannels(wallet.GuildId)

	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(404, gin.H{
			"status":   "errorMsg",
			"errorMsg": "No guild found",
		})
		return
	}

	conf.BotAdmins = utils.IToSArray(utils.PSQLToArray(conf.BAs))
	ctx.JSON(200, gin.H{
		"guild":    conf,
		"wallet":   wallet,
		"channels": s,
	})
}
