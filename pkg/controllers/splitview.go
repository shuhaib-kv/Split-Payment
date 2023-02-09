package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shuhaib-kv/Split-Gpay-Golang.git/pkg/db"
	"github.com/shuhaib-kv/Split-Gpay-Golang.git/pkg/models"
)

func ViewSplit(c *gin.Context) {
	var body struct {
		expenceid uint
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"status": false,
			"error":  "Invalid JSON",
			"data":   "null",
		})
		return
	}
	var expense models.Expense
	var split []models.Split
	if err := db.DBS.First(&expense, "id=1"); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Group Doesn't exist",
			"error":   "error please enter valid information",
		})
		return
	}
	db.DBS.Find(&split, "expenseid=?", expense.ID).Scan(&split)

	splitData := make([]map[string]interface{}, len(split))
	for i, s := range split {
		splitData[i] = map[string]interface{}{
			"split id":       s.ID,
			"user name":      s.Username,
			"amount":         s.Amount,
			"payment status": s.Paymentstatus,
			"split status":   s.Splitstatus,
		}
	}

	c.JSON(200, gin.H{
		"status":  true,
		"message": "your split",
		"data":    splitData,
	})
	return
}
func ViewMysplit(c *gin.Context) {
	id := c.GetUint("id")
	var gid struct {
		groupid uint
	}
	if err := c.BindJSON(&gid); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"status": false,
			"error":  "Invalid JSON",
			"data":   "null",
		})
		return
	}
	var expense []models.Expense
	db.DBS.Find(&expense, "groupid=? and spliowner=? ", gid.groupid, id).Scan(&expense)
	for _, i := range expense {
		c.JSON(200, gin.H{
			"status":  true,
			"message": "Your Groups",
			"data": gin.H{
				"expenceid": i.ID,
				"tittle":    i.Title,
				"amount ":   i.Amount,
			},
		})
	}
}
func ViewWhoNotPaid(c *gin.Context) {
	var body struct {
		expenceid uint
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"status": false,
			"error":  "Invalid JSON",
			"data":   "null",
		})
		return
	}
	var expense models.Expense
	var split []models.Split
	if err := db.DBS.First(&expense, "id=1"); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Group Doesn't exist",
			"error":   "error please enter valid information",
		})
		return
	}
	db.DBS.Find(&split, "expenseid=? and paymentstatus=?", expense.ID, false).Scan(&split)
	splitData := make([]map[string]interface{}, len(split))
	for i, s := range split {
		splitData[i] = map[string]interface{}{
			"split id":       s.ID,
			"user name":      s.Username,
			"amount":         s.Amount,
			"payment status": s.Paymentstatus,
			"split status":   s.Splitstatus,
		}
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": "persons who not paid",
		"data":    splitData,
	})
}

func ViewWhoPaid(c *gin.Context) {
	var body struct {
		expenceid uint
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"status": false,
			"error":  "Invalid JSON",
			"data":   "null",
		})
		return
	}
	var expense models.Expense
	var split []models.Split
	if err := db.DBS.First(&expense, "id=1"); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Group Doesn't exist",
			"error":   "error please enter valid information",
		})
		return
	}
	db.DBS.Find(&split, "expenseid=? and paymentstatus=?", expense.ID, true).Scan(&split)

	splitData := make([]map[string]interface{}, len(split))
	for i, s := range split {
		splitData[i] = map[string]interface{}{
			"split id":       s.ID,
			"user name":      s.Username,
			"amount":         s.Amount,
			"payment status": s.Paymentstatus,
			"split status":   s.Splitstatus,
		}
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": "persons who  paid",
		"data":    splitData,
	})
}