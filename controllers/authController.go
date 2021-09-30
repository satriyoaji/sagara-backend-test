package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/satriyoaji/sagara-backend-test/database"
	"github.com/satriyoaji/sagara-backend-test/helper"
	"github.com/satriyoaji/sagara-backend-test/models"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

func ValidateStructUser(user *models.User) []*models.ErrorResponse {
	var errors []*models.ErrorResponse
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element models.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func Register(c *fiber.Ctx) error {
	dataUser := new(models.User)

	if err := c.BodyParser(dataUser); err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	errors := ValidateStructUser(dataUser)
	if errors != nil {
		return c.Status(fiber.StatusMisdirectedRequest).JSON(fiber.Map{
			"message": errors,
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(dataUser.Password), 14)
	user := models.User{
		Name: dataUser.Name,
		Email: dataUser.Email,
		Password: password,
	}

	database.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil{
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	// handle error
	if user.Id == 0 { //default Id when return nil
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found!",
		})
	}

	// match password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password!",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer: strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour*24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(helper.GoDotEnvVariable("SECRET_KEY")))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "error when logging in !",
		})
	}

	cookie := fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour*24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Login Succeeded",
		"token": token,
	})
}

func User(c *fiber.Ctx) error {

	token, err := helper.JwtVerify(c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name: "jwt",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "successfully logged out !",
	})
}