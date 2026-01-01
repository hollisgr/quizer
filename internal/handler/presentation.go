package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 1. Хендлер загрузки файла
func (h *handler) UploadPresentation(c *gin.Context) {
	file, _ := c.FormFile("file")
	gameIdStr := c.PostForm("game_id")

	// Генерируем уникальное имя
	filename := fmt.Sprintf("%s_%d.pdf", gameIdStr, time.Now().Unix())
	filepath := "./uploads/" + filename

	// Сохраняем файл на диск
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(500, gin.H{"error": "failed to save file"})
		return
	}

	// Сохраняем ПУТЬ в базу данных
	// UPDATE lobbies SET pdf_path = $1 WHERE uuid = $2
	gameId, _ := strconv.Atoi(gameIdStr)
	h.gameSvc.UpdateFilePath(c.Request.Context(), gameId, filepath)

	c.JSON(200, gin.H{"status": "uploaded", "path": filename})
}

// 2. Хендлер отдачи файла (для react-pdf)
func (h *handler) GetPDF(c *gin.Context) {
	gameIdStr := c.Query("game_id")
	gameId, _ := strconv.Atoi(gameIdStr)
	game, _ := h.gameSvc.GameLoad(c.Request.Context(), gameId)

	c.File(game.Link) // Gin сам отдаст файл корректно
}
