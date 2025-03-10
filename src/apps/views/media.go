package views

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"socious/src/apps/auth"
	"socious/src/apps/models"
	"socious/src/apps/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/socious-io/goaccount"
)

func mediaGroup(router *gin.Engine) {
	g := router.Group("media")
	g.Use(auth.LoginRequired())

	g.POST("", auth.LoginRequired(), sociousIdSession(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		sessionToken := c.MustGet("socious_id_session").(goaccount.SessionToken)

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file received"})
			return
		}

		// Open the file
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open file"})
			return
		}
		defer src.Close()

		// Upload file
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		checksum, err := utils.GenerateChecksum(src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open file"})
			return
		}

		src.Seek(0, io.SeekStart)

		fileName := fmt.Sprintf("%s%s", checksum, filepath.Ext(file.Filename))
		fileURL, err := c.MustGet("uploader").(*utils.GCSUploader).UploadFile(ctx, fileName, file.Header.Get("Content-Type"), src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		sessionMedia := new(models.Media)
		if err := sessionToken.UploadMedia(src, &sessionMedia); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(sessionMedia)

		media := &models.Media{
			Filename:   file.Filename,
			URL:        fileURL,
			IdentityID: identity.ID,
		}

		if err := media.Create(c.MustGet("ctx").(context.Context)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, media)

	})

	// g.GET("/:id", func(c *gin.Context) {
	// 	identity := c.MustGet("identity").(*models.Identity)

	// })
}
