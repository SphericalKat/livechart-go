package entities

type Show struct {
	Title     *string  `json:"title"`
	Thumbnail *string  `json:"thumbnail"`
	Tags      []string `json:"tags"`
}
