package VideoService

import (
	"EntitlementServer/DatabaseAbstraction"
	"github.com/gin-gonic/gin"
)

type VideoService interface {
	GetAllVideos() ([]VSVideo, error)
	GetVideoByIndexID(indexID int) (VSVideo, error)
	GetVideosOfProduct(productID int) ([]VSVideo, error)
	StreamVideo(indexID int) (string, error)
}

type VSVideo struct {
	IndexID     int
	Name        string
	Description string
	Points      int
	Thumbnail   string
	Filename    string
}

type VSService struct {
	DB DatabaseAbstraction.DBOrm
}

func (V VSService) RegisterHandlers(r *gin.Engine, middleware ...gin.HandlerFunc) {
	r.GET("/api/video/:number/stream", V.StartVideoStream)
	r.GET("/api/video", middleware[0], V.GetAllVideosHandler)
	r.GET("/api/video/:number", middleware[0], V.GetVideoInfoHandler)
	r.POST("/api/video/:number/progress", middleware[0], V.MarkFinishedEndpoint)
	r.GET("/api/video/watched", middleware[0], V.GetWatchedVideos)
}

func (V VSService) GetLabel() string {
	return "Video Service"
}

func (V VSService) DBVideoToUpstreamType(video DatabaseAbstraction.Video) VSVideo {
	return VSVideo{
		IndexID:     video.IndexID,
		Name:        video.Name,
		Description: video.Description,
		Points:      video.Points,
		Thumbnail:   video.Thumbnail,
	}
}

func (V VSService) GetAllVideos() ([]VSVideo, error) {
	videos, err := V.DB.GetAllVideos()
	if err != nil {
		return nil, err
	}

	convertedVids := make([]VSVideo, len(videos))
	for i, v := range videos {
		convertedVids[i] = V.DBVideoToUpstreamType(v)
	}

	return convertedVids, nil
}

func (V VSService) GetVideoByIndexID(indexID int) (VSVideo, error) {
	video, err := V.DB.GetVideoByIndexID(indexID)
	if err != nil {
		return VSVideo{}, err
	}

	return V.DBVideoToUpstreamType(video), nil
}

func (V VSService) GetVideosOfProduct(productID int) ([]VSVideo, error) {
	videos, err := V.DB.GetVideosByProductIndexID(productID)
	if err != nil {
		return nil, err
	}

	convertedVids := make([]VSVideo, len(videos))
	for i, v := range videos {
		convertedVids[i] = V.DBVideoToUpstreamType(v)
	}

	return convertedVids, nil
}

func (V VSService) StreamVideo(indexID int, userID int) (string, error) {
	// Get the video filename from the database,
	// then find the file in the filesystem and return it after checking the user's product ownership

	// Get the video from the database
	video, err := V.DB.GetVideoByIndexID(indexID)
	if err != nil {
		return "", err
	}

	// Get the product that the video belongs to
	product, err := V.DB.GetProductByVideoIndexID(indexID)
	if err != nil {
		return "", err
	}

	// Check if the user owns the product
	ownsProduct, err := V.DB.GetOwnedProducts(userID)
	if err != nil {
		return "", err
	}

	// Check if the user owns the product
	for _, p := range ownsProduct {
		if p.IndexID == product.IndexID {
			return video.Filename, nil
		}
	}

	return video.Name, nil
}
