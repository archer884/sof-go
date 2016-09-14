package cookies

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

// Cookie stores a quote and a category for that quote.
type Cookie struct {
	Category string
	Content  string
}

func (c *Cookie) String() string {
	return fmt.Sprintf("%s: %s", c.Category, c.Content)
}

// Service provides access to cookies on a random or categorized basis.
type Service struct {
	Rng     *rand.Rand
	Cookies []Cookie
}

// Error is returned in the event that this library experiences a failure.
type Error struct {
	Error string
}

// New creates a new cookies service based on a directory and a random number
// generator.
func New(dir string, rng *rand.Rand) (Service, *Error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.txt", dir))
	if err != nil {
		log.Fatal(err)
		return Service{}, &Error{fmt.Sprintf("Unable to read directory: %s", dir)}
	}

	var cookies []Cookie
	for _, file := range files {
		info, err1 := os.Stat(file)
		category := strings.Replace(info.Name(), ".txt", "", 1)
		content, err2 := ioutil.ReadFile(file)
		if err1 == nil && err2 == nil {
			content := string(content)
			for _, quote := range strings.Split(content, "%") {
				quote := strings.TrimSpace(quote)
				if len(quote) > 0 {
					cookie := Cookie{category, quote}
					cookies = append(cookies, cookie)
				}
			}
		}
	}

	fmt.Printf("%v cookies loaded\n", len(cookies))
	return Service{rng, cookies}, nil
}

// GetCookie returns a random cookie.
func (s *Service) GetCookie() Cookie {
	return s.Cookies[s.Rng.Intn(len(s.Cookies))]
}

// ByCategory returns a random cookie from the category provided or an Error
// value instead, in the event that the category requested has no quotes
// available.
func (s *Service) ByCategory(category string) (Cookie, *Error) {
	var cookies []Cookie
	for idx := range s.Cookies {
		if category == s.Cookies[idx].Category {
			cookies = append(cookies, s.Cookies[idx])
		}
	}

	if len(cookies) == 0 {
		return Cookie{}, &Error{fmt.Sprintf("Category not supported: %s", category)}
	}

	return cookies[s.Rng.Intn(len(cookies))], nil
}
