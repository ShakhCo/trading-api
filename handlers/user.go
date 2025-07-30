package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
	"trading_api/db"
	"trading_api/models"
)

func RegisterUserHandler(ctx *fasthttp.RequestCtx) {
	// Ensure it's a POST request
	if !ctx.IsPost() {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		fmt.Fprintf(ctx, "Only POST allowed")
		return
	}

	var payload struct {
		TelegramID int64  `json:"telegram_id"`
		FirstName  string `json:"first"`
		LastName   string `json:"last"` // Optional
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
		// Update existing user
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
	if err := db.DB.Find(&users).Error; err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "❌ Failed to fetch users: %v", err)
		return
	}

	// Convert users to JSON
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
