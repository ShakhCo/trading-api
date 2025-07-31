package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/valyala/fasthttp"
	"trading_api/db"
	"trading_api/models"
)

func RegisterUserHandler(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		fmt.Fprintf(ctx, "Only POST allowed")
		return
	}

	var payload struct {
		TelegramID int64  `json:"telegram_id"`
		FirstName  string `json:"first"`
		LastName   string `json:"last"`
	}

	if err := json.Unmarshal(ctx.PostBody(), &payload); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprintf(ctx, "Invalid JSON: %v", err)
		return
	}

	if payload.FirstName == "" || payload.TelegramID == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprintf(ctx, "Missing first name or telegram_id")
		return
	}

	var existing models.User
	tx := db.DB.First(&existing, "telegram_id = ?", payload.TelegramID)

	if tx.Error == nil {
		existing.FirstName = payload.FirstName
		if payload.LastName != "" {
			existing.LastName = &payload.LastName
		}
		if err := db.DB.Save(&existing).Error; err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			fmt.Fprintf(ctx, "Failed to update user: %v", err)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusOK)
		fmt.Fprintf(ctx, "✅ User updated!")
		return
	}

	user := models.User{
		TelegramID: payload.TelegramID,
		FirstName:  payload.FirstName,
		CreatedAt:  time.Now(),
	}
	if payload.LastName != "" {
		user.LastName = &payload.LastName
	}

	if err := db.DB.Create(&user).Error; err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "DB error: %v", err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	fmt.Fprintf(ctx, "✅ User registered!")
}

func GetAllUsersHandler(ctx *fasthttp.RequestCtx) {
	var users []models.User
	if err := db.DB.Preload("Photos").Find(&users).Error; err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "❌ Failed to fetch users: %v", err)
		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "❌ JSON marshal error: %v", err)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(response)
}

func UploadPhotoHandler(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		fmt.Fprintf(ctx, "Only POST allowed")
		return
	}

	telegramID := ctx.UserValue("telegram_id")
	if telegramID == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprintf(ctx, "Missing telegram_id in URL")
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprintf(ctx, "Failed to parse form: %v", err)
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprintf(ctx, "No file provided")
		return
	}

	fileHeader := files[0]
	src, err := fileHeader.Open()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Can't open file: %v", err)
		return
	}
	defer src.Close()

	// Save the file locally
	os.MkdirAll("uploads", os.ModePerm)
	filename := fmt.Sprintf("%d_%d_%s", time.Now().Unix(), ctx.Time().Nanosecond(), fileHeader.Filename)
	fullPath := fmt.Sprintf("uploads/%s", filename)
	dst, err := os.Create(fullPath)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Can't create file: %v", err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Can't save file: %v", err)
		return
	}

	// Save metadata in database
	var user models.User
	if err := db.DB.First(&user, "telegram_id = ?", telegramID).Error; err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		fmt.Fprintf(ctx, "User not found")
		return
	}

	photo := models.UserPhoto{
		UserID:   user.TelegramID,
		FilePath: fullPath,
	}
	if err := db.DB.Create(&photo).Error; err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "DB error: %v", err)
		return
	}

	// Return the download link
	downloadURL := fmt.Sprintf("http://localhost:9000/uploads/%s", filename) // Adjust base URL if needed

	response := map[string]string{
		"message":      "✅ Photo uploaded successfully!",
		"download_url": downloadURL,
	}

	respBytes, _ := json.Marshal(response)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody(respBytes)
}
