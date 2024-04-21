package controllers

import (
	"context"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"os"
)

func UploadFile(c *fiber.Ctx) error {
	ctx := context.Background()
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	cld, _ := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"))
	formFile, err := file.Open()
	resp, err := cld.Upload.Upload(ctx, formFile, uploader.UploadParams{PublicID: uuid.New().String()})
	return c.JSON(fiber.Map{"success": true, "data": resp})
}

func UploadFileForProfile(c *fiber.Ctx) error {
	ctx := context.Background()
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	cld, _ := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"))
	formFile, err := file.Open()
	resp, err := cld.Upload.Upload(ctx, formFile, uploader.UploadParams{PublicID: "profile"})
	return c.JSON(fiber.Map{"success": true, "data": resp})
}
