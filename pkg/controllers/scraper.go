package controllers

import (
	"fmt"
	"net/url"

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

		shows = append(shows, entities.Show{
			Title:     &title,
			Thumbnail: &thumbURL,
			Tags:      tags,
		})
	})

	c.Visit(config.Conf.WebsiteURL)
	c.Wait()

	return shows
}
