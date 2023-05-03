package DatabaseAbstraction

import (
	"context"
	"time"
)

type Video struct {
	IndexID     int
	Name        string
	Description string
	Points      int
	Thumbnail   string
	Filename    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Get a video by its indexID
func (dbc DBConnector) GetVideoByIndexID(indexID int) (Video, error) {
	// Get the video from the database
	row := dbc.DB.QueryRow(context.Background(), "SELECT id, name, description, points, thumbnail, filename FROM video WHERE id = $1", indexID)

	var video Video
	err := row.Scan(&video.IndexID, &video.Name, &video.Description, &video.Points, &video.Thumbnail, &video.Filename)
	if err != nil {
		return Video{}, err
	}

	return video, nil
}

// Get video parent product
func (dbc DBConnector) GetProductByVideoIndexID(indexID int) (Product, error) {
	// Get the product from the database
	row := dbc.DB.QueryRow(context.Background(), "SELECT id, name, description, price, image FROM products WHERE id = (SELECT parent_product_id FROM video WHERE id = $1)", indexID)

	var product Product
	err := row.Scan(&product.IndexID, &product.Name, &product.Description, &product.Price, &product.Image)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func (dbc DBConnector) GetAllVideos() ([]Video, error) {
	// Get all the videos from the database
	rows, err := dbc.DB.Query(context.Background(), "SELECT id, name, description, points, thumbnail, filename FROM video")
	if err != nil {
		return []Video{}, err
	}

	// Iterate over the rows and add them to the slice
	var videos []Video

	for rows.Next() {
		var video Video
		err := rows.Scan(&video.IndexID, &video.Name, &video.Description, &video.Points, &video.Thumbnail, &video.Filename)
		if err != nil {
			return []Video{}, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

// Get all videos related to a Product
func (dbc DBConnector) GetVideosByProductIndexID(indexID int) ([]Video, error) {
	// Get all the videos from the database
	rows, err := dbc.DB.Query(context.Background(), "SELECT id, name, description, points, thumbnail, filename FROM video WHERE parent_product_id = $1", indexID)
	if err != nil {
		return []Video{}, err
	}

	// Iterate over the rows and add them to the slice
	var videos []Video

	for rows.Next() {
		var video Video
		err := rows.Scan(&video.IndexID, &video.Name, &video.Description, &video.Points, &video.Thumbnail, &video.Filename)
		if err != nil {
			return []Video{}, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (dbc DBConnector) MarkVideoAsWatched(indexID int, user User) error {
	_, err := dbc.DB.Exec(context.Background(), "INSERT INTO user_watched_videos (user_id, video_id) VALUES ($1, $2)", user.IndexID, indexID)
	if err != nil {
		return err
	}

	// Increase user points by the video's points
	_, err = dbc.DB.Exec(context.Background(), "UPDATE users SET points = points + (SELECT points FROM video WHERE id = $1) WHERE id = $2", indexID, user.IndexID)
	if err != nil {
		return err
	}

	return nil
}

func (dbc DBConnector) GetWatchedVideosByUser(user User) ([]Video, error) {
	// Get all the videos from the database
	rows, err := dbc.DB.Query(context.Background(), "SELECT id, name, description, points, thumbnail, filename FROM video WHERE id IN (SELECT video_id FROM user_watched_videos WHERE user_id = $1)", user.IndexID)
	if err != nil {
		return []Video{}, err
	}

	// Iterate over the rows and add them to the slice
	var videos []Video

	for rows.Next() {
		var video Video
		err := rows.Scan(&video.IndexID, &video.Name, &video.Description, &video.Points, &video.Thumbnail, &video.Filename)
		if err != nil {
			return []Video{}, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}
