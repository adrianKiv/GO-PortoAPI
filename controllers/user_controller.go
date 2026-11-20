package controllers

import (
	"net/http"
	"time"
	"os"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"api/portodata/models"
	"api/portodata/config"
)

type AuthImputRegister struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthImputLogin struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}


func RegisterUser(context *gin.Context) {
	var input AuthImputRegister

	// validate input register
	err := context.ShouldBindJSON(&input); 
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"massage": "invalid input",
			"error": err.Error(),
		})
		return
	}	

	//hash pass
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"massage": "failed to hash password",
			"error": err.Error(),
		})
		return
	}

	// simpan ke db tabel user
	user := models.User{
		Name: input.Name,
		Email: input.Email,
		Password: string(hashedPassword),
	}

	userCreated := config.DB.Create(&user).Error
	if userCreated != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"massage": "failed to create user",
			"error": userCreated.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"massage": "user registered successfully",
		"data": gin.H{
			"id": user.ID,
			"name": user.Name,
			"email": user.Email,

		},
	})
}

func LoginUser(context *gin.Context) {
	var input AuthImputLogin

	//validate input login
	err := context.ShouldBindJSON(&input);
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"massage": "invalid input",
			"error": err.Error(),
		})
		return
	}

	var user models.User
	userData := config.DB.Where("email = ?", input.Email).First(&user).Error
	if userData != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"massage": "invalid email or password",
			"error": userData.Error(),
		})
		return
	}

	errMatchpass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if errMatchpass != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"massage": "invalid email or password",
			"error": errMatchpass.Error(),
		})
		return
	}

	//token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email": user.Email,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"massage": "failed to generate token",
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"massage": "login success",
		"token": tokenString,
		"user": gin.H{
			"id": user.ID,
			"name": user.Name,
			"email": user.Email,
		},
	})
}

func GetCurrentUser(context *gin.Context) {
	userID, exists := context.Get("userID")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "user not authenticated",
		})
		return
	}

	var user models.User
	userData := config.DB.Select("id", "name", "email").First(&user, userID).Error
	if userData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"error": userData.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "success get current user",
		"data": gin.H{
			"id": user.ID,
			"name": user.Name,
			"email": user.Email,	
		},
	})
}