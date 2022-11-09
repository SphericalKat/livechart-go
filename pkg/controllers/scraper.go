package controllers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
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

	c.OnHTML(".anime", func(h *colly.HTMLElement) {
		// get title
		title := h.ChildText(".main-title")

		// get tags
		tags := h.ChildTexts(".anime-tags > li")

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
		studios := h.ChildTexts(".anime-studios > li")

		// get more metadata
		source := h.ChildText(".anime-source")
		eps := h.ChildText(".anime-episodes")
		summary := h.ChildText(".anime-synopsis")

		// get related links
		relatedLinks := make([]entities.Link, 0)
		h.ForEach(".related-links > li", func(i int, h *colly.HTMLElement) {
			url := h.ChildAttr("a", "href")
			if strings.HasPrefix(url, "/") {
				url = fmt.Sprintf("%s%s", config.Conf.WebsiteURL, url)
			}

			linkType := entities.GetType(h.ChildAttr("a", "class"))
			relatedLinks = append(relatedLinks, entities.Link{
				Type: linkType,
				URL:  url,
			})
		})

		shows = append(shows, entities.Show{
			Title:         &title,
			Thumbnail:     &thumbURL,
			Tags:          tags,
			Studios:       studios,
			AirTime:       airTime,
			Source:        &source,
			EpisodeFormat: &eps,
			Summary:       &summary,
			RelatedLinks:  relatedLinks,
		})
	})

	c.Visit(config.Conf.WebsiteURL)
	c.Wait()

	return shows
}
