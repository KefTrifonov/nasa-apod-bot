package nasa

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type APOD struct {
	Date        string
	Title       string
	Explanation string
	MediaType   string // "image" или "video"
	URL         string // для видео — ссылка на ролик, для картинки — ссылка на картинку
	HDURL       string // ссылка на картинку (у видео-дней бывает превью)
}

const feedURL = "https://science.nasa.gov/feed/apod-basic/"

// feed — структура RSS-ленты: <rss><channel><item>...</item></channel></rss>
type feed struct {
	Items []item `xml:"channel>item"`
}

// item — одна запись ленты (один день APOD)
type item struct {
	Title       string `xml:"title"`
	PubDate     string `xml:"pubDate"`
	HDURL       string `xml:"hdurl"`
	Explanation string `xml:"explanation"`
	Content     string `xml:"encoded"` // <content:encoded> — полная HTML-страница дня
}

// videoSrc ищет в HTML встроенный ролик: <iframe src="...youtube..."> или <source src="...mp4">
var videoSrc = regexp.MustCompile(`(?i)<(?:iframe|source)[^>]+src="([^"]*(?:youtube\.com|youtu\.be|vimeo\.com|\.mp4|\.webm)[^"]*)"`)

// Fetch возвращает запись ленты по номеру: 0 — сегодня, 1 — вчера и т.д.
func Fetch(index int) (*APOD, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(feedURL)
	if err != nil {
		return nil, fmt.Errorf("request apod feed: %w", errors.Unwrap(err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("apod feed: unexpected status %d", resp.StatusCode)
	}

	var f feed
	if err := xml.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, fmt.Errorf("decode apod feed: %w", err)
	}

	if index < 0 || index >= len(f.Items) {
		return nil, fmt.Errorf("apod feed: no item #%d (feed has %d)", index, len(f.Items))
	}
	it := f.Items[index]

	date := it.PubDate
	if t, err := time.Parse(time.RFC1123Z, it.PubDate); err == nil {
		date = t.Format("2006-01-02")
	}

	apod := &APOD{
		Date:        date,
		Title:       it.Title,
		Explanation: it.Explanation,
		MediaType:   "image",
		URL:         it.HDURL,
		HDURL:       it.HDURL,
	}

	if m := videoSrc.FindStringSubmatch(it.Content); m != nil {
		apod.MediaType = "video"
		apod.URL = m[1]
	}

	if apod.URL == "" {
		return nil, fmt.Errorf("apod %q: neither image nor video found", it.Title)
	}

	return apod, nil
}
