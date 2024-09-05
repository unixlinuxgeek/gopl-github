// The github package provides a Go API for Github-hosting .
// Пакет github предоставляет Go API для Github-хостинга .
// См.https://developer.github.com/v3/search/$searchissues.

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}
type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

// SearchIssues запрашивает Github.
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	// Необходимо закрыть resp.Body на всех путях выполнения.
	// (В главе 5 вы познакомитесь с более простым решением: `defer`.)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("Сбой запроса: %s", resp.Status)
	}
	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

// SearchIssues2 запрашивает Github. (добавлена пагинация по страницам)
func SearchIssues2(terms []string, page int, perPage int) (*IssuesSearchResult, error) {
	if perPage == 0 {
		perPage = 20
	}
	if page == 0 {
		page = 1
	}
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q + "&page=" + strconv.Itoa(page) + "&per_page=" + strconv.Itoa(perPage))
	if err != nil {
		return nil, err
	}
	// Необходимо закрыть resp.Body на всех путях выполнения.
	// (В главе 5 вы познакомитесь с более простым решением: `defer`.)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("Сбой запроса: %s", resp.Status)
	}
	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}
