package controllers

import (
	"assignmentday23/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var usersInMemory = []models.User{}

type UserController struct {
	db *gorm.DB
}

func GetUserWithoutDB(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": usersInMemory})
}

func CreateUserWithoutDB(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = uint(len(usersInMemory) + 1)
	usersInMemory = append(usersInMemory, user)
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		db: db,
	}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := uc.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (uc *UserController) GetUser(c *gin.Context) {
	var users []models.User
	if err := uc.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}
