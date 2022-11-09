package controllers

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/SphericalKat/livechart-go/internal/config"
	"github.com/SphericalKat/livechart-go/pkg/entities"
	"github.com/gocolly/colly/v2"
)

func GetLatest() []entities.Show {
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(),
	)

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	shows := make([]entities.Show, 0)

	c.OnHTML(".anime-card", func(h *colly.HTMLElement) {
		// get title
		title := h.ChildText(".main-title")

		// get tags
		tags := make([]string, 0)
		h.ForEach(".anime-tags > li", func(i int, h *colly.HTMLElement) {
			tags = append(tags, h.Text)
		})

		// get thumbnail URL
		var thumbURL string
		rawThumb := h.ChildAttr(".poster-container > img", "data-src")
		parsed, err := url.Parse(rawThumb)
		if err != nil {
			thumbURL = ""
		}
		thumbURL = fmt.Sprintf(`%s://%s%s?style=large&format=jpg`, parsed.Scheme, parsed.Host, parsed.Path)

		// get timestamp
		var airTime *time.Time
		timestamp := h.ChildAttr(".poster-container > .episode-countdown", "data-timestamp")
		t, err := strconv.ParseInt(timestamp, 10, 64)
		if err == nil {
			utcTime := time.Unix(t, 0).UTC()
			airTime = &utcTime
		}

		// get anime studios
		studios := make([]string, 0)
		h.ForEach(".anime-studios > li", func(i int, h *colly.HTMLElement) {
			studios = append(studios, h.Text)
		})

		source := h.ChildText(".anime-source")
		eps := h.ChildText(".anime-episodes")

		shows = append(shows, entities.Show{
			Title:         &title,
			Thumbnail:     &thumbURL,
			Tags:          tags,
			Studios:       studios,
			AirTime:       airTime,
			Source:        &source,
			EpisodeFormat: &eps,
		})
	})

	c.Visit(config.Conf.WebsiteURL)
	c.Wait()

	return shows
}
