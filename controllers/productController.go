package controllers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/satriyoaji/sagara-backend-test/database"
	"github.com/satriyoaji/sagara-backend-test/helper"
	"github.com/satriyoaji/sagara-backend-test/models"
	"os"
	"strconv"
)

func ValidateStructProduct(user *models.Product) []*models.ErrorResponse {
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

func FindAllProduct(c *fiber.Ctx) error {
	_, err := helper.JwtVerify(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	var products []models.Product

	database.DB.Find(&products)

	return c.JSON(products)
}

func FindProductById(c *fiber.Ctx) error {
	_, err := helper.JwtVerify(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var product models.Product

	database.DB.Where("id = ?", productId).First(&product)
	if product.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found !",
		})
	}

	return c.JSON(product)
}

func CreateProduct(c *fiber.Ctx) error {
	_, err := helper.JwtVerify(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	dataProduct := new(models.Product)
	if err := c.BodyParser(dataProduct); err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := ValidateStructProduct(dataProduct)
	if errors != nil {
		return c.Status(fiber.StatusMisdirectedRequest).JSON(fiber.Map{
			"message": errors,
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusMisdirectedRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if errorUpload := c.SaveFile(file, fmt.Sprintf("./uploads/%s", file.Filename)); errorUpload != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errorUpload.Error(),
		})
	} else{
		dataProduct.Image = file.Filename
	}

	database.DB.Model(&models.Product{}).Create(&dataProduct)

	return c.JSON(dataProduct)
}

func UpdateProductById(c *fiber.Ctx) error {
	_, err := helper.JwtVerify(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	dataProduct := new(models.Product)
	if err := c.BodyParser(dataProduct); err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := ValidateStructProduct(dataProduct)
	if errors != nil {
		return c.Status(fiber.StatusMisdirectedRequest).JSON(fiber.Map{
			"message": errors,
		})
	}

	var product models.Product
	database.DB.Where("id = ?", productId).First(&product)
	if product.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found !",
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusMisdirectedRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if errorUpload := c.SaveFile(file, fmt.Sprintf("./uploads/%s", file.Filename)); errorUpload != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errorUpload.Error(),
		})
	} else{
		if _, err := os.Stat("./uploads/"+product.Image); err == nil {
			os.Remove("./uploads/"+product.Image)
		}
		product.Image = file.Filename
		dataProduct.Image = file.Filename
	}

	database.DB.Model(&models.Product{}).Where("id = ?", productId).Updates(dataProduct)

	return c.JSON(product)
}

func DeleteProductById(c *fiber.Ctx) error {
	_, err := helper.JwtVerify(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var product models.Product
	database.DB.Where("id = ?", productId).Delete(&product)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "delete product success!",
	})
}