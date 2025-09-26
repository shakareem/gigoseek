package concerts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/shakareem/gigoseek/pkg/config"
)

type Event struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description_short"`
	StartsAt    string `json:"starts_at"`
	Location    struct {
		City    string `json:"city"`
		Address string `json:"address"`
	} `json:"location"`
	URL string `json:"url"`
}

type EventsResponse struct {
	Values []Event `json:"values"`
	Total  int     `json:"total"`
}

func GetTimepadConcerts(artists []string, city string) []Event {
	concerts := []Event{}
	for _, artist := range artists {
		artistConcerts, err := getArtistConcert(artist, city)
		if err != nil {
			log.Printf("Ошибка при получении событий для артиста %s: %v", artist, err)
			continue
		}
		concerts = append(concerts, artistConcerts...)
	}

	return concerts
}

func getArtistConcert(artist, city string) ([]Event, error) {
	params := url.Values{}
	params.Add("category_ids", config.Get().Timepad.ConcertsCategoryID)
	params.Add("cities", city)
	params.Add("keywords", artist)

	fullURL := config.Get().Timepad.ApiURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.Get().TimepadApiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response: %s", resp.Status)
	}

	var result EventsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Values, nil
}
