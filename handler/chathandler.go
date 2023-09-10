package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sap200/TinyFundConnect/db"
	e "github.com/sap200/TinyFundConnect/errors"
	"github.com/sap200/TinyFundConnect/pangeachecks"
	"github.com/sap200/TinyFundConnect/secret"
	"github.com/sap200/TinyFundConnect/types"
)

func SendChatHandler(c *gin.Context) {
	var chat types.Chat
	err := c.ShouldBindJSON(&chat)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.CHAT_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	chat.TimeStamp = time.Now().UnixMilli()

	redactedMessage, err := pangeachecks.RedactAMessage(chat.Message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.CHAT_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	chat.Message = redactedMessage

	repo := db.New()

	ok, err := repo.DoesRecordExists(secret.CHAT_COLLECTION, chat.PoolId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.CHAT_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	if ok {
		// get the record
		clogs, err := repo.GetChatDetailsByPool(secret.CHAT_COLLECTION, chat.PoolId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    e.CHAT_ERROR_CODE,
				"message": err.Error(),
			})
			return
		}

		clogs.ChatRecord = append(clogs.ChatRecord, chat)

		nm, _ := clogs.EncodeToMap()

		// append the chat and save the record
		repo.Save(secret.CHAT_COLLECTION, chat.PoolId, nm)

	} else {
		//save first chat
		c := types.Chats{
			ChatRecord: []types.Chat{chat},
		}
		nm, _ := c.EncodeToMap()
		repo.Save(secret.CHAT_COLLECTION, chat.PoolId, nm)

	}

	// no need to log chats, its unnecessary too much of logs

	c.JSON(http.StatusOK, gin.H{
		"code":    e.CHAT_SUCCESS_CODE,
		"message": "Chat sent successfully",
	})

}

func GetChatByPoolHandler(c *gin.Context) {
	poolId := c.Query("poolId")

	repo := db.New()

	clogs, err := repo.GetChatDetailsByPool(secret.CHAT_COLLECTION, poolId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.CHAT_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	nm, _ := clogs.EncodeToMap()
	c.JSON(http.StatusOK, nm)
}
