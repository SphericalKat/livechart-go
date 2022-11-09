package entities

import "time"

type Show struct {
	Title         *string    `json:"title"`
	Thumbnail     *string    `json:"thumbnail"`
	Source        *string    `json:"source"`
	EpisodeFormat *string    `json:"episode_format"`
	Summary       *string    `json:"summary"`
	Tags          []string   `json:"tags"`
	Studios       []string   `json:"studios"`
	AirTime       *time.Time `json:"air_time"`
	RelatedLinks  []Link     `json:"related_links"`
}
