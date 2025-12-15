# indeed_scraper

a go application using http and golang's html parsing library
to scrape data from Indeed.com

this program was created for the CD Block CP Data Science class in 2025

## usage

1. install [go](https://go.dev/)
2. go to [indeed.com](indeed.com)
3. use your browser's developer tools to get your
  cookie from your requests [see this example in chrome](docs/chrome_devtools.png)
4. copy that cookie into a `cookie.txt` file in the directory
  of where you run the project
5. run `go run . -- scrape` to run the program in scraping mode (to collect data)
6. run `go run . -- pack` to run the program in packing mode (which puts collected data into a csv)
7. you can run the program without arguments (`go run .`) to see all possible options for the program
