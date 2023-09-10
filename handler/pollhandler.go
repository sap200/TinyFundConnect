package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sap200/TinyFundConnect/db"
	e "github.com/sap200/TinyFundConnect/errors"
	"github.com/sap200/TinyFundConnect/pangeaauth"
	"github.com/sap200/TinyFundConnect/pangeachecks"
	"github.com/sap200/TinyFundConnect/pangealogger"
	"github.com/sap200/TinyFundConnect/secret"
	"github.com/sap200/TinyFundConnect/types"
)

func CreatePollHandler(c *gin.Context) {

	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	var poll types.Poll
	err = c.ShouldBindJSON(&poll)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	redactedMessage, err := pangeachecks.RedactAMessage(poll.Message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	poll.Message = redactedMessage

	poll.PollId = uuid.New().String()
	poll.IsActive = true

	nm, _ := poll.EncodeToMap()

	repo := db.New()
	repo.Save(secret.POLLING_COLLECTION, poll.PollId, nm)

	event := pangealogger.New(
		"User with email "+poll.PollCreatorEmail+"created a poll with id "+poll.PollId+"in pool with Id "+poll.PoolId,
		"POLL_CREATION_EVENT",
		"",
		poll.PollCreatorEmail,
		poll.PoolId,
	)
	pangealogger.Log(event)

	c.JSON(http.StatusOK, nm)
}

func InactivatePollHandler(c *gin.Context) {

	var poll types.Poll
	err := c.ShouldBindJSON(&poll)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	pollId := poll.PollId

	// get the repo by pollId
	repo := db.New()
	p, err := repo.GetPollDetailsByPollId(secret.POLLING_COLLECTION, pollId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	p.IsActive = false
	nm, _ := p.EncodeToMap()

	repo.Save(secret.POLLING_COLLECTION, pollId, nm)

	c.JSON(http.StatusOK, nm)
}

func CreateVoteHandler(c *gin.Context) {
	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	var vote types.Vote
	err = c.ShouldBindJSON(&vote)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	repo := db.New()
	poll, err := repo.GetPollDetailsByPollId(secret.POLLING_COLLECTION, vote.PollId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	vote.PoolId = poll.PoolId

	if !poll.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": "Poll is inactive and it cannot be voted",
		})
		return
	}

	if poll.Votes == nil {
		poll.Votes = []types.Vote{vote}
	} else {
		alreadyVoted := false
		for i := 0; i < len(poll.Votes); i++ {
			if poll.Votes[i].VoterEmail == vote.VoterEmail {
				alreadyVoted = true
			}
		}
		if alreadyVoted {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    e.POLL_ERROR_CODE,
				"message": "User has already voted",
			})
			return
		}

		if vote.Choice < 0 || vote.Choice >= len(poll.Options) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    e.POLL_ERROR_CODE,
				"message": "Invalid option selected",
			})
			return
		}

		poll.Votes = append(poll.Votes, vote)
	}

	nm, _ := poll.EncodeToMap()

	repo.Save(secret.POLLING_COLLECTION, vote.PollId, nm)

	repm, _ := vote.EncodeToMap()

	event := pangealogger.New(
		"User with email "+vote.VoterEmail+"voted for poll with id "+vote.PollId+"in pool "+vote.PoolId,
		"VOTING_EVENT",
		"",
		vote.VoterEmail,
		vote.PoolId,
	)
	pangealogger.Log(event)

	c.JSON(http.StatusOK, repm)

}

func GetAllPollById(c *gin.Context) {

	pollId := c.Query("pollId")

	repo := db.New()
	exists, err := repo.DoesRecordExists(secret.POLLING_COLLECTION, pollId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": "Poll doesn't exists",
		})
		return
	}

	poll, err := repo.GetPollDetailsByPollId(secret.POLLING_COLLECTION, pollId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	nm, _ := poll.EncodeToMap()

	c.JSON(http.StatusOK, nm)
}

func GetAllPollByPoolHandler(c *gin.Context) {
	poolId := c.Query("poolId")

	repo := db.New()

	polls, err := repo.GetAllPollsByPoolId(secret.POLLING_COLLECTION, poolId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	var nm []map[string]interface{}
	binary, _ := json.Marshal(polls)
	json.Unmarshal(binary, &nm)

	c.JSON(http.StatusOK, nm)

}

func CheckIfUserHasAlreadyVotedForPoll(c *gin.Context) {
	pollId := c.Query("pollId")
	email := c.Query("email")

	repo := db.New()
	exists, err := repo.DoesRecordExists(secret.POLLING_COLLECTION, pollId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": "Poll doesn't exists",
		})
		return
	}

	poll, err := repo.GetPollDetailsByPollId(secret.POLLING_COLLECTION, pollId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	if poll.Votes == nil {
		c.JSON(http.StatusOK, gin.H{
			"hasVoted": false,
		})
		return
	}

	hasVoted := false
	for i := 0; i < len(poll.Votes); i++ {
		if poll.Votes[i].VoterEmail == email {
			hasVoted = true
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"hasVoted": hasVoted,
	})

}

func GetPollResultHandler(c *gin.Context) {
	pollId := c.Query("pollId")

	repo := db.New()
	exists, err := repo.DoesRecordExists(secret.POLLING_COLLECTION, pollId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": "Poll doesn't exists",
		})
		return
	}

	poll, err := repo.GetPollDetailsByPollId(secret.POLLING_COLLECTION, pollId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	m, err := poll.GetResults()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POLL_ERROR_CODE,
			"message": err.Error(),
		})
		return
	}

	rm := make(map[string]interface{})
	rm["pollResultByOptionIndex"] = m

	sm := make(map[string]int)
	totalVoteCount := 0
	for k, v := range m {
		sm[poll.Options[k]] = v
		totalVoteCount += v
	}

	rm["pollResultByOptionNames"] = sm
	rm["totalVoteCount"] = totalVoteCount

	c.JSON(http.StatusOK, rm)
}
