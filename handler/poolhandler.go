package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sap200/TinyFundConnect/db"
	e "github.com/sap200/TinyFundConnect/errors"
	"github.com/sap200/TinyFundConnect/pangeaauth"
	"github.com/sap200/TinyFundConnect/pangealogger"
	"github.com/sap200/TinyFundConnect/secret"
	"github.com/sap200/TinyFundConnect/types"
	"github.com/sap200/TinyFundConnect/utils"
)

func CreatePoolHandler(c *gin.Context) {

	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	var pool types.Pool
	err = c.ShouldBindJSON(&pool)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	// generate a unique pool Id
	uuid := uuid.New().String()
	pool.PoolId = uuid
	pool.MembersList = append(pool.MembersList, types.Member{
		Email:     pool.CreatorEmail,
		HasJoined: true,
	})

	// save it in db
	m, err := pool.EncodeToMap()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}
	repo := db.New()
	repo.Save(secret.POOL_COLLECTION, uuid, m)

	event := pangealogger.New(
		"User with email "+pool.CreatorEmail+"created a pool with id "+uuid,
		"POOL_CREATION_EVENT",
		"",
		pool.CreatorEmail,
		uuid,
	)
	pangealogger.Log(event)

	c.JSON(http.StatusOK, gin.H{
		"code":    e.POOL_SUCCESS_CODE,
		"status":  types.SUCCESS_MESSAGE_STATUS_VALUE,
		"message": e.POOL_CREATION_SUCCESS_MESSAGE,
		"poolId":  pool.PoolId,
	})
}

func GetPoolDetailsHandler(c *gin.Context) {

	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	poolId := c.Query("poolId")
	fmt.Println("poolId", poolId)

	repo := db.New()
	ok, err := repo.DoesRecordExists(secret.POOL_COLLECTION, poolId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": "Pool with given Id not found",
		})

		return
	}

	x, err := repo.GetPoolDetails(secret.POOL_COLLECTION, poolId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	var p types.Pool
	fmt.Println(x)
	json.Unmarshal(utils.GetBytesFromInterface(x), &p)
	fmt.Println(p)

	r, _ := p.EncodeToMap()
	c.JSON(http.StatusOK, r)

}

func UpdatePoolMemberStatus(c *gin.Context) {

	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	var pool types.Pool
	err = c.ShouldBindJSON(&pool)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	poolId := pool.PoolId

	// fetch the pool record from db
	repo := db.New()
	ok, err := repo.DoesRecordExists(secret.POOL_COLLECTION, poolId)
	fmt.Println(ok, err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": "Pool with given Id not found",
		})

		return
	}

	x, err := repo.GetPoolDetails(secret.POOL_COLLECTION, poolId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	var p types.Pool
	fmt.Println(x)
	json.Unmarshal(utils.GetBytesFromInterface(x), &p)
	fmt.Println(p)

	for i := 0; i < len(pool.MembersList); i++ {
		m := pool.MembersList[i]
		emailToBeUpdated := m.Email
		for j := 0; j < len(p.MembersList); j++ {
			mem := p.MembersList[j]
			memberEmail := mem.Email
			if emailToBeUpdated == memberEmail {
				p.MembersList[j].HasJoined = pool.MembersList[i].HasJoined
			}
		}
	}

	fmt.Println("After")
	fmt.Println(p)

	nm, _ := p.EncodeToMap()
	// update db
	repo.Save(secret.POOL_COLLECTION, p.PoolId, nm)

	c.JSON(http.StatusOK, nm)
}

// Add pool member
func PoolAddMemberHandler(c *gin.Context) {

	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	var pool types.Pool
	err = c.ShouldBindJSON(&pool)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	poolId := pool.PoolId

	// fetch the pool record from db
	repo := db.New()
	ok, err := repo.DoesRecordExists(secret.POOL_COLLECTION, poolId)
	fmt.Println(ok, err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": "Pool with given Id not found",
		})

		return
	}

	x, err := repo.GetPoolDetails(secret.POOL_COLLECTION, poolId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	var p types.Pool
	fmt.Println(x)
	json.Unmarshal(utils.GetBytesFromInterface(x), &p)
	fmt.Println(p)

	// check if userEmail is in Pool
	doesHeExists := false
	for i := 0; i < len(p.MembersList); i++ {
		if p.MembersList[i].Email == pool.OneMember.Email {
			doesHeExists = true
			break
		}
	}

	if doesHeExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": "user with email already exists in db",
		})

		return
	}

	p.MembersList = append(p.MembersList, types.Member{Email: pool.OneMember.Email})
	p.OneMember = types.Member{}

	pnm, _ := p.EncodeToMap()

	repo.Save(secret.POOL_COLLECTION, p.PoolId, pnm)

	event := pangealogger.New(
		"User with email "+pool.OneMember.Email+" was added to pool with id "+poolId,
		"POOL_UPDATION_EVENT",
		"",
		"",
		poolId,
	)
	pangealogger.Log(event)

	c.JSON(http.StatusOK, pnm)
}

func ExitFromPoolHandler(c *gin.Context) {
	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	var pool types.Pool
	err = c.ShouldBindJSON(&pool)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	poolId := pool.PoolId

	// fetch the pool record from db
	repo := db.New()
	ok, err := repo.DoesRecordExists(secret.POOL_COLLECTION, poolId)
	fmt.Println(ok, err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": "Pool with given Id not found",
		})

		return
	}

	x, err := repo.GetPoolDetails(secret.POOL_COLLECTION, poolId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	var p types.Pool
	fmt.Println(x)
	json.Unmarshal(utils.GetBytesFromInterface(x), &p)
	fmt.Println(p)

	// check if userEmail is in Pool
	doesHeExists := false
	hisIndex := -1
	for i := 0; i < len(p.MembersList); i++ {
		if p.MembersList[i].Email == pool.OneMember.Email {
			doesHeExists = true
			hisIndex = i
			break
		}
	}

	if !doesHeExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": "User not in the pool",
		})

		return
	}

	member := p.MembersList[hisIndex]
	newslice := []types.Member{}
	newslice = append(newslice, p.MembersList[:hisIndex]...)
	newslice = append(newslice, p.MembersList[hisIndex+1:]...)
	p.MembersList = newslice
	p.OneMember = types.Member{}
	p.PoolBalance -= member.Balance

	pnm, _ := p.EncodeToMap()

	repo.Save(secret.POOL_COLLECTION, p.PoolId, pnm)

	// TODO: add secure auditing logs
	event := pangealogger.New(
		"User with email "+pool.OneMember.Email+"has exited from pool with id"+poolId,
		"POOL_EXIT_EVENT",
		"",
		pool.OneMember.Email,
		poolId,
	)
	pangealogger.Log(event)

	c.JSON(http.StatusOK, pnm)

}

func GetAllPoolsWithListedEmail(c *gin.Context) {
	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	email := c.Query("email")

	fmt.Println("query email: ", email)

	repo := db.New()
	allPools, err := repo.GetAllPools(secret.POOL_COLLECTION)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	//fmt.Println(allPools)

	ap := *allPools
	uswps := []types.UserWithPoolAndStatus{}

	for i := 0; i < len(ap); i++ {
		for j := 0; j < len(ap[i].MembersList); j++ {
			if ap[i].MembersList[j].Email == email {
				activeStatus := ap[i].MembersList[j].HasJoined
				uwps := types.UserWithPoolAndStatus{
					CurrentPool: ap[i],
					IsActive:    activeStatus,
				}
				uswps = append(uswps, uwps)
				break
			}
		}
	}

	fmt.Println("new pool::", uswps)
	result := types.PoolListWithUserMail{
		UserEmail: email,
		PoolList:  uswps,
	}

	newMap, err := result.EncodeToMap()
	fmt.Println(newMap)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, newMap)
}
