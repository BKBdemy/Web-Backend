package VideoService_test

import (
	"EntitlementServer/DatabaseAbstraction"
	"EntitlementServer/DatabaseAbstraction/mocks"
	"EntitlementServer/VideoService"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetAllVideosHandler(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	// Mock DB
	mockDB := new(mocks.DBOrm)

	// Set expectations on DB
	videos := []DatabaseAbstraction.Video{
		{
			IndexID:     1,
			Name:        "Video 1",
			Description: "This is Video 1",
			Points:      100,
			Thumbnail:   "thumbnail1.png",
			Filename:    "filename1.mp4",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	mockDB.On("GetAllVideos").Return(videos, nil)

	videoSvc := VideoService.VSService{DB: mockDB}

	// Register handler
	r.GET("/api/video", videoSvc.GetAllVideosHandler)

	req, _ := http.NewRequest(http.MethodGet, "/api/video", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDB.AssertExpectations(t)
}

func TestGetVideoByIndexIDHandler(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	// Mock DB
	mockDB := new(mocks.DBOrm)

	// Set expectations on DB
	video := DatabaseAbstraction.Video{
		IndexID:     1,
		Name:        "Video 1",
		Description: "This is Video 1",
		Points:      100,
		Thumbnail:   "thumbnail1.png",
		Filename:    "filename1.mp4",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockDB.On("GetVideoByIndexID", 1).Return(video, nil)

	videoSvc := VideoService.VSService{DB: mockDB}

	// Register handler
	r.GET("/api/video/:number/info", videoSvc.GetVideoInfoHandler)

	req, _ := http.NewRequest(http.MethodGet, "/api/video/1/info", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDB.AssertExpectations(t)
}
