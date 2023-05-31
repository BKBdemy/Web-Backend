package VideoService

import (
	"EntitlementServer/DatabaseAbstraction"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// GetAllVideosHandler godoc
// @Summary Get all videos
// @Description Get all videos
// @Tags Videos
// @Accept  json
// @Produce  json
// @Success 200 {array} VSVideo
// @Router /api/video [get]
func (V VSService) GetAllVideosHandler(c *gin.Context) {
	videos, err := V.DB.GetAllVideos()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, videos)
}

// GetVideoInfoHandler godoc
// @Summary Get video info
// @Description Get video info
// @Tags Videos
// @Accept  json
// @Produce  json
// @Param number path int true "VSVideo ID"
// @Success 200 {object} VSVideo
// @Router /api/video/{number} [get]
func (V VSService) GetVideoInfoHandler(c *gin.Context) {
	videoID := c.Param("number")

	// Convert videoID to int
	videoIDInt, err := strconv.Atoi(videoID)

	// Retrieve video filename from database
	video, err := V.DB.GetVideoByIndexID(videoIDInt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, video)
}

// MarkFinishedEndpoint godoc
// @Summary Mark video as finished
// @Description Mark video as finished & get points
// @Tags Videos
// @Accept  json
// @Produce  json
// @Param number path int true "VSVideo ID"
// @Success 200 {string} string	"success"
// @Router /api/video/{number}/progress [post]
// @Security		ApiKeyAuth
func (V VSService) MarkFinishedEndpoint(c *gin.Context) {
	videoID := c.Param("number")

	// Convert videoID to int
	videoIDInt, err := strconv.Atoi(videoID)

	// Retrieve video filename from database
	video, err := V.DB.GetVideoByIndexID(videoIDInt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Get user from context
	user := c.MustGet("user").(DatabaseAbstraction.User)

	// Mark video as finished
	err = V.DB.MarkVideoAsWatched(video.IndexID, user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "marked as finished"})
}

// StartVideoStream godoc
// @Summary Start video stream
// @Description Start video stream
// @Tags Videos
// @Param number path int true "VSVideo ID"
// @Success 200
// @Security ApiKeyAuth
// @Router /api/video/{number}/stream [get]
func (V VSService) StartVideoStream(c *gin.Context) {
	videoID := c.Param("number")

	// Convert videoID to int
	videoIDInt, err := strconv.Atoi(videoID)

	// Retrieve video filename from database
	video, err := V.DB.GetVideoByIndexID(videoIDInt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Start stream of file with basepath + filename
	//basepath := os.Getenv("VIDEO_BASE_PATH")
	serveFile(c, video.Filename)
	//c.File(video.Filename)
}

// GetWatchedVideos godoc
// @Summary Get watched videos
// @Description Get watched videos
// @Tags Videos
// @Accept  json
// @Produce  json
// @Success 200 {array} VSVideo
// @Router /api/video/watched [get]
// @Security ApiKeyAuth
func (V VSService) GetWatchedVideos(c *gin.Context) {
	// Get user from context
	user := c.MustGet("user").(DatabaseAbstraction.User)

	// Get watched videos
	videos, err := V.DB.GetWatchedVideosByUser(user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, videos)
}

// Custom ranged file implementation
func serveFile(c *gin.Context, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			println(fmt.Sprintf("Error closing file: %s", err.Error()))
		}
	}(file)

	fileInfo, err := file.Stat()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	fileSize := fileInfo.Size()
	mimeType := mime.TypeByExtension(filepath.Ext(filePath))

	rangeHeader := c.GetHeader("Range")
	if rangeHeader == "" {
		http.ServeContent(c.Writer, c.Request, filePath, fileInfo.ModTime(), file)
		return
	}

	start, end := parseRange(rangeHeader, fileSize)
	if end == -1 {
		end = fileSize - 1
	}

	if _, err := file.Seek(start, 0); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	bytes := make([]byte, end-start+1)
	if _, err := file.Read(bytes); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Header("Accept-Ranges", "bytes")
	c.Data(http.StatusPartialContent, mimeType, bytes)
}

func parseRange(rangeHeader string, fileSize int64) (start, end int64) {
	var err error
	if strings.HasPrefix(rangeHeader, "bytes=") {
		rangeStr := rangeHeader[6:]
		parts := strings.Split(rangeStr, "-")
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return
		}

		if len(parts) > 1 && parts[1] != "" {
			end, err = strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return
			}
		} else {
			end = -1
		}
	}
	return
}
