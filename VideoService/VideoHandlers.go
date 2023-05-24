package VideoService

import (
	"EntitlementServer/DatabaseAbstraction"
	"github.com/gin-gonic/gin"
	"strconv"
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
	c.File(video.Filename)
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
