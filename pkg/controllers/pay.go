package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shuhaib-kv/Split-Gpay-Golang.git/pkg/db"
	"github.com/shuhaib-kv/Split-Gpay-Golang.git/pkg/models"
)

func PaySplit(c *gin.Context) {
	id := c.GetUint("id")

	var body struct {
		Amount    uint `json:"amount"`
		Expenceid uint `json:"expenceid"`
		Splitid   uint `json:"splitid"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"status": false,
			"error":  "Invalid JSON",
			"data":   "null",
		})
		return
	}
	var expence models.Expense
	var Split models.Split

	db.DBS.Find(&expence, "id=?", body.Expenceid).Scan(&expence)
	db.DBS.Find(&Split, "expenseid=? and id=? and userid=?", body.Expenceid, body.Splitid, id).Scan(&Split)
	if expence.Status == true {
		c.JSON(http.StatusConflict, gin.H{
			"status": false,
			"error":  "split closed",
			"data":   expence.Status,
		})
		return
	}
	if Split.Splitstatus == true && expence.Status == true {
		c.JSON(http.StatusConflict, gin.H{
			"status": false,
			"error":  "Your split closed by admin",
			"data":   expence.Status,
		})
		return
	}
	if Split.Splitstatus == true {
		c.JSON(http.StatusConflict, gin.H{
			"status": false,
			"error":  "Your paid the split",
			"data":   "",
		})
		return
	}
	if Split.Amount != float64(body.Amount) {
		c.JSON(http.StatusConflict, gin.H{
			"status": false,
			"error":  "Amount doesnt match",
			"data":   Split.Amount,
		})
		return
	}
	var pay = models.Payment{
		Expenseid: body.Expenceid,
		Splitid:   body.Splitid,
		Amount:    body.Amount,
	}
	db.DBS.Create(&pay)
	c.JSON(http.StatusFound, gin.H{
		"status": true,
		"error":  "Your paid the split",
		"data":   pay})
	// var done = models.Split{

	// 	Splitstatus: true,
	// }
	var done = models.Split{

		Paymentid:   pay.ID,
		Splitstatus: true,
	}
	var doneex = models.Expense{

		Status: true,
	}
	db.DBS.Model(&Split).Where("splits.id=?", pay.Splitid).Updates(&done)
	var sp []models.Split
	db.DBS.Raw("select * from splits where splits.expenseid=?", body.Expenceid).Scan(&sp)
	var flag int
	flag = 0
	for _, i := range sp {
		if i.Splitstatus == false {
			flag = 1
		}
	}
	if flag != 1 {
		db.DBS.Model(&expence).Where("id=?", body.Expenceid).Updates(&doneex)

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
	db.DBS.Find(&split, "expenseid=? and splitstatus=?", expense.ID, false).Scan(&split)
	splitData := make([]map[string]interface{}, len(split))
	for i, s := range split {
		splitData[i] = map[string]interface{}{
			"split id":     s.ID,
			"userid":       s.Userid,
			"split owner":  s.Username,
			"amount":       s.Amount,
			"split status": s.Splitstatus,
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
	db.DBS.Find(&split, "expenseid=? and splitstatus=?", expense.ID, true).Scan(&split)

	splitData := make([]map[string]interface{}, len(split))
	for i, s := range split {
		splitData[i] = map[string]interface{}{
			"split id":     s.ID,
			"userid":       s.Userid,
			"split owner":  s.Username,
			"amount":       s.Amount,
			"split status": s.Splitstatus,
		}
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": "persons who  paid",
		"data":    splitData,
	})
}
