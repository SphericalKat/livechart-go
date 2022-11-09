package entities

type Link struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

func GetType(class string) string {
	switch class {
	case "website-icon":
		return "Official Website"
	case "preview-icon":
		return "Video preview"
	case "watch-icon":
		return "Legal watch sources"
	case "twitter-icon":
		return "Twitter"
	case "anilist-icon":
		return "Anilist"
	case "mal-icon":
		return "MyAnimeList"
	case "anidb-icon":
		return "AniDB"
	case "anime-planet-icon":
		return "Anime Planet"
	case "anisearch-icon":
		return "AniSearch"
	case "kitsu-icon":
		return "Kitsu"
	case "crunchyroll-icon":
		return "Crunchyroll"
	default:
		return "Unknown"
	}
}
