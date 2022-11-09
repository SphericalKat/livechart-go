package entities

import "time"

type Show struct {
	Title     *string    `json:"title"`
	Thumbnail *string    `json:"thumbnail"`
	Tags      []string   `json:"tags"`
	AirTime   *time.Time `json:"air_time"`
}
