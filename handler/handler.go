package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sap200/TinyFundConnect/pangeaauth"
)

func TestPathHandler(c *gin.Context) {
	fmt.Println(c.GetQuery("code"))
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}

func CheckTokenHandler(c *gin.Context) {
	token := c.Query("token")
	email := c.Query("email")

	err := pangeaauth.ValidateEmailPangeaToken(email, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "FAILURE",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
	})

}
