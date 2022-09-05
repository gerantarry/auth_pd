package entity

type Trick struct {
	TrickId         int    `json:"trick_id"`
	DifficultyLevel int    `json:"difficulty_level"`
	VideoId         int    `json:"video_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
}
