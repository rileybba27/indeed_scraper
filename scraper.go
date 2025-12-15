package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func Scraper(additionalArguments string) {
	err := os.MkdirAll("data/listings", os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	foundFile, err := os.OpenFile("data/found.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	defer foundFile.Close()

	foundBytes, err := io.ReadAll(foundFile)
	if err != nil {
		log.Fatalln(err)
	}

	skipList := strings.Split(string(foundBytes), "\n")

	cookieBytes, err := os.ReadFile("cookie.txt")
	if err != nil {
		log.Fatalln("Missing `cookie.txt`!")
	}

	cookie := string(cookieBytes)
	cookie = strings.TrimSpace(cookie)

	log.Println("Found cookie, starting request...")

	httpClient := &http.Client{}

	extraParams := additionalArguments
	initialPage, err := GetURL(BASE_URL+"/jobs?q=programmer"+extraParams, cookie, httpClient)
	if err != nil {
		log.Fatalln("Failed to get initial job listings page, cannot proceed without listings to scrape! Error:", err)
	}

	initialPageReader := strings.NewReader(initialPage)

	document, err := html.Parse(initialPageReader)
	if err != nil {
		log.Fatalln("Failed to parse initial job listings page HTML, cannot proceed. Error:", err)
	}

	listingNodes := FindNodesOfClass(document, "resultContent")

	log.Println("Parsing", len(listingNodes), "nodes...")
	if err != nil {
		log.Fatalf("Error with main context: %s\n", err)
	}

	for _, listingNode := range listingNodes {
		time.Sleep(10 * time.Second)

		var job ProgrammingJob
		var jobURL string
		// var jobFound bool
		var jobPriceString string

		jobTitleNode, err := FindNodeBasedOnPath(listingNode, []HTMLPath{
			HTMLPath{
				Tag:   "div",
				Index: 0,
			},
			HTMLPath{
				Tag:   "h2",
				Index: 0,
			},
			HTMLPath{
				Tag:   "a",
				Index: 0,
			},
		})
		if err != nil {
			log.Println("Couldn't find job title node, error:", err)
			continue
		}
		job.Title = CollectNodeText(jobTitleNode)

		jobCompanyNode, err := FindNodeBasedOnPath(listingNode, []HTMLPath{
			HTMLPath{
				Tag:   "div",
				Index: 1,
			},
			HTMLPath{
				Tag:   "div",
				Index: 0,
			},
			HTMLPath{
				Tag:   "div",
				Index: -1,
			},
			HTMLPath{
				Tag:   "div",
				Index: 0,
			},
		})
		if err != nil {
			log.Println("Couldn't find job company node, error:", err)
			continue
		}
		job.Company = CollectNodeText(jobCompanyNode)

		jobLocationNode, err := FindNodeBasedOnPath(listingNode, []HTMLPath{
			HTMLPath{
				Tag:   "div",
				Index: 1,
			},
			HTMLPath{
				Tag:   "div",
				Index: 0,
			},
			HTMLPath{
				Tag:   "div",
				Index: -1,
			},
			HTMLPath{
				Tag:   "div",
				Index: -1,
			},
		})
		if err != nil {
			log.Println("Couldn't find job company node, error:", err)
			continue
		}
		job.Location = CollectNodeText(jobLocationNode)

		jobPriceNode, err := FindNodeBasedOnPath(listingNode, []HTMLPath{
			HTMLPath{
				Tag:   "div",
				Index: 1,
			},
			HTMLPath{
				Tag:   "div",
				Index: 1,
			},
		})
		if err != nil {
			log.Println("Couldn't find job price node, error:", err)
			continue
		}
		jobPriceString = CollectNodeText(jobPriceNode)

		jobURL, err = GetNodeAttr(jobTitleNode, "href")
		if err != nil {
			log.Println("Couldn't find job url href, error:", err)
			continue
		}

		jobURL = BASE_URL + jobURL
		if slices.Contains(skipList, jobURL) {
			log.Println("Found previously parsed listing, skipping...")
			continue
		}

		log.Println("Parsing", jobURL)

		jobListingHTML, err := GetURL(jobURL, cookie, httpClient)
		if err != nil {
			if strings.Contains(err.Error(), "200 OK") && strings.Contains(err.Error(), "429") {
				log.Println("Got rate limited, waiting extra 20 seconds to hopefully help.")
				time.Sleep(20 * time.Second)
				continue
			}

			log.Println("Failed to get job listing HTML, Error:", err)
			continue
		}

		skipList = append(skipList, jobURL)
		n, err := foundFile.WriteString(jobURL + "\n")
		if err != nil {
			log.Println("Failed to write jobURL to file, wrote", n, "characters")
			continue
		}

		reader := strings.NewReader(jobListingHTML)
		jobListingDocument, err := html.Parse(reader)
		if err != nil {
			log.Println("Error parsing response html:", err)
			continue
		}

		jobDescriptionNode, err := FindNodeByAttr(jobListingDocument, "id", "jobDescriptionText")
		if err != nil {
			log.Println("Error finding jobDescriptionText Node:", err)
			continue
		}

		jobDescription := CollectNodeText(jobDescriptionNode)
		job.ParseDescription(jobDescription)

		if strings.HasPrefix(jobPriceString, "From") {
			jobPriceString = strings.TrimSpace(strings.ReplaceAll(jobPriceString, "From", ""))
		}

		if strings.Contains(jobPriceString, "an hour") {
			jobPriceString = strings.TrimSpace(strings.ReplaceAll(jobPriceString, "an hour", ""))
			job.PayMethod = Wage
			job.ParsePayment(jobPriceString)
		} else if strings.Contains(jobPriceString, "a month") {
			jobPriceString = strings.TrimSpace(strings.ReplaceAll(jobPriceString, "a month", ""))
			job.PayMethod = Salary
			job.ParsePayment(jobPriceString)
			job.AveragePay *= 12
		} else if strings.Contains(jobPriceString, "a year") {
			jobPriceString = strings.TrimSpace(strings.ReplaceAll(jobPriceString, "a year", ""))
			job.PayMethod = Salary
			job.ParsePayment(jobPriceString)
		} else {
			job.PayMethod = Unknown
			job.AveragePay = 0
		}

		job.Title = strings.ReplaceAll(job.Title, "\n", "")
		job.Location = strings.ReplaceAll(job.Location, "\n", "")
		job.Location = strings.ReplaceAll(job.Location, "\t", "")
		job.Location = strings.ReplaceAll(job.Location, " ", " ")
		job.Company = strings.ReplaceAll(job.Company, "\n", "")

		hash := sha256.New()
		_, err = fmt.Fprintf(hash, "%v", job)
		if err != nil {
			log.Println("Failed to write job struct data into hash, error:", err)
			continue
		}

		sum := hash.Sum(nil)
		fileName := hex.EncodeToString(sum)

		_, err = os.Stat("data/listings/" + fileName + ".json")
		if err == nil {
			log.Println("Listing already found based on hash, skipping...")
			continue
		}

		jsonBytes, err := json.Marshal(job)
		if err != nil {
			log.Println("Failed to marshal job struct into JSON, error:", err)
			continue
		}

		err = os.WriteFile("data/listings/"+fileName+".json", jsonBytes, os.ModePerm)
		if err != nil {
			log.Println("Failed to write JSON to file, error:", err)
			continue
		}

		log.Println("Parsed and wrote", fileName, "to file!")
	}
}
