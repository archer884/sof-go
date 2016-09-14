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

type Cookie struct {
	Category string
	Content  string
}

func (c *Cookie) String() string {
	return fmt.Sprintf("%s: %s", c.Category, c.Content)
}

type Service struct {
	Rng     *rand.Rand
	Cookies []Cookie
}

type Error struct {
	Error string
}

func New(dir string, rng *rand.Rand) (Service, *Error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.txt", dir))
	if err != nil {
		log.Fatal(err)
		return Service{}, &Error{fmt.Sprintf("Unable to read directory: %s", dir)}
	}

	cookies := make([]Cookie, 0)
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

	return Service{rng, cookies}, nil
}

func (s *Service) GetCookie() *Cookie {
	return &s.Cookies[s.Rng.Intn(len(s.Cookies))]
}

func (s *Service) ByCategory(category string) (Cookie, *Error) {
	cookies := make([]Cookie, 0)
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
