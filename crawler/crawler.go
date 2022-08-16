package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)

type Engine struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Price    string `json:"price,optional"`
	Shipping string `json:"shipping,optional"`
	Img      string `json:"img,optional"`
	Grade    string `json:"grade,optional"`
}

func GetEngines(w http.ResponseWriter, r *http.Request) {

	//var engines []Engine
	vars := mux.Vars(r)
	vin := vars["vin"]

	// INITIATE MEGA CRAWLERS

	// Crawler that gets link to actual engine page
	link, err := getEngineLink(vin)
	if err != nil {
		log.Println("Issue getting Engine Page Link")
	}

	// Crawler for actual Engine Page
	engines, err := getEngines(*link)
	if err != nil {
		log.Println("Issue getting engines.")
	}
	// Crawler to get individual engine Data
	engineData, err := getIndividualEngineData(engines)
	if err != nil {
		log.Println("Issue getting engines.")
	}

	w.Header().Add("Content-Type", "application/json")

	j, err := json.Marshal(engineData)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(j)

}

func getEngineLink(vin string) (*string, error) {
	var links []Engine
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.OnHTML("div[class=searchColOne]", func(h *colly.HTMLElement) {
		h.ForEach("div", func(i int, h *colly.HTMLElement) {
			if h.ChildText("a") == "Engine" && h.ChildAttr("a", "href") != "" {
				e := Engine{
					Name: h.ChildText("a"),
					URL:  h.ChildAttr("a", "href"),
				}
				links = append(links, e)
				c.Visit(h.Request.AbsoluteURL(h.ChildAttr("a", "href")))
			}
		})

	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.PostMultipart("https://www.hollanderparts.com/Home", map[string][]byte{
		"hdnVIN": []byte(vin),
	})
	c.Wait()

	engineLink := links[1]

	return &engineLink.URL, nil

}

func getEngines(url_suffix string) ([]Engine, error) {
	var engines []Engine
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnHTML("div[class=searchColOne]", func(h *colly.HTMLElement) {
		h.ForEach("div", func(i int, h *colly.HTMLElement) {
			e := Engine{
				Name: h.ChildText("a"),
				URL:  h.ChildAttr("a", "href"),
			}
			engines = append(engines, e)
			c.Visit(h.Request.AbsoluteURL(h.ChildAttr("a", "href")))
		})

	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})
	c.Visit("https://www.hollanderparts.com/" + url_suffix)
	c.Wait()

	return engines, nil

}

func getIndividualEngineData(links []Engine) ([]Engine, error) {

	var engines []Engine

	for _, engine := range links {
		c := colly.NewCollector(
			colly.AllowURLRevisit(),
		)

		c.OnHTML("div[class=individualPartHolder]", func(h *colly.HTMLElement) {

			name := strings.Split(h.Response.Request.URL.String(), "/")
			price := h.ChildText("div[class=partPrice]")
			shipping := h.ChildText("div[class=partShipping]")
			img := h.ChildAttr("img", "src")
			grade := h.ChildText("div[class=gradeText]")

			e := Engine{
				Name:     name[len(name)-1],
				URL:      h.Response.Request.URL.String(),
				Grade:    grade,
				Img:      img,
				Price:    price,
				Shipping: shipping,
			}

			engines = append(engines, e)

		})
		c.OnRequest(func(r *colly.Request) {
			log.Println("visiting", r.URL.String())
		})

		c.Visit("https://www.hollanderparts.com/" + engine.URL)

		c.Wait()

	}

	fmt.Println(engines)

	return engines, nil

}
