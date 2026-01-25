package note

import (
	"fmt"
	"note/internal/utils"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *NoteHandler) UploadImage(c *gin.Context) {
	// 1. 获取上传的文件
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		utils.Error(c, 400, "请上传图片")
		return
	}
	defer file.Close()

	// 2. 生成唯一文件名 (防止重名覆盖)
	// 例如: uuid.jpg 或 timestamp_filename.jpg
	ext := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// 3. 上传到 MinIO
	url, err := h.storageService.UploadImage(c, newFileName, header.Size, file, header.Header.Get("Content-Type"))
	if err != nil {
		utils.Error(c, 500, "图片上传失败")
		return
	}

	// 4. 返回 URL 给前端
	// 前端拿到这个 URL 后，把它塞到笔记的 content 里，或者作为 cover_image 字段传给 CreateNote
	utils.Success(c, gin.H{"url": url})
}
