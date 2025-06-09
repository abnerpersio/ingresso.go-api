package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"ingresso.go/internal/infra/config"
)

type MovieService struct {
	cachedGenres []Genre
}

func makeRequest(method string, path string, body io.Reader) (*http.Request, error) {
	baseUrl := config.GetEnv("MOVIES_API_URL")
	token := config.GetEnv("MOVIES_API_KEY")

	req, err := http.NewRequest(method, baseUrl+path, body)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	query := req.URL.Query()
	query.Set("language", "pt-BR")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Authorization", "Bearer "+token)

	return req, nil
}

type movieListResponse struct {
	Page    int `json:"page"`
	Results []struct {
		Id            int32  `json:"id"`
		OriginalTitle string `json:"original_title"`
		Overview      string `json:"overview"`
		PosterPath    string `json:"poster_path"`
		ReleaseDate   string `json:"release_date"`
		Title         string `json:"title"`
		GenreIds      []int  `json:"genre_ids"`
	} `json:"results"`
}

type genreListResponse struct {
	Genres []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
}

type Genre struct {
	Id   int
	Name string
}

type Movie struct {
	Id            string   `json:"id"`
	Title         string   `json:"title"`
	Overview      string   `json:"overview"`
	PosterPath    string   `json:"poster_path"`
	ReleaseDate   string   `json:"release_date"`
	Genres        []string `json:"genres"`
	OriginalTitle string   `json:"original_title"`
}

func NewMovieService() *MovieService {
	return &MovieService{}
}

func formatImageUrl(path string) string {
	if path == "" {
		return ""
	}

	return config.GetEnv("MOVIES_IMAGE_BASE_URL") + "/w500" + path
}

func (service *MovieService) listGenres() []Genre {
	if len(service.cachedGenres) > 0 {
		return service.cachedGenres
	}

	req, err := makeRequest(http.MethodGet, "/genre/movie/list", nil)

	if err != nil {
		return []Genre{}
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return []Genre{}
	}

	defer resp.Body.Close()
	var response genreListResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	genres := make([]Genre, 0)

	for _, genre := range response.Genres {
		genres = append(genres, Genre{Id: genre.Id, Name: genre.Name})
	}

	service.cachedGenres = genres
	return genres
}

func (service *MovieService) formatGenres(genreIds []int) []string {
	genres := service.listGenres()
	result := make([]string, 0)

	for _, genre := range genres {
		for _, id := range genreIds {
			if genre.Id == id {
				result = append(result, genre.Name)
			}
		}
	}

	return result
}

func (service *MovieService) List() ([]Movie, error) {
	req, err := makeRequest(http.MethodGet, "/movie/now_playing", nil)

	if err != nil {
		return []Movie{}, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return []Movie{}, err
	}

	defer resp.Body.Close()
	var response movieListResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)
	list := make([]Movie, 0)

	for _, movie := range response.Results {
		list = append(list, Movie{
			Id:            fmt.Sprint(movie.Id),
			Title:         movie.Title,
			OriginalTitle: movie.OriginalTitle,
			Overview:      movie.Overview,
			PosterPath:    formatImageUrl(movie.PosterPath),
			ReleaseDate:   movie.ReleaseDate,
			Genres:        service.formatGenres(movie.GenreIds),
		})
	}

	return list, nil
}
