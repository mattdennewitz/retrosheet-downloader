package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func worker(id int, years <-chan int, results chan<- string, v bool) {
	for year := range years {
		/* create file */
		fn := fmt.Sprintf("dl/%deve.zip", year)
		out_f, err := os.Create(fn)
		defer out_f.Close()

		if err != nil {
		 	fmt.Printf("Could not create %d.zip: %s\n", year, err.Error())
		 	continue
		}

		/* download archive */
		url := fmt.Sprintf("http://www.retrosheet.org/events/%deve.zip", year)

		resp, err := http.Get(url)
		defer resp.Body.Close()

		if err != nil {
			fmt.Printf("Could not download %s: %s\n", url, err.Error())
			continue
		}

		/* write output */
		_, err = io.Copy(out_f, resp.Body)

		if err != nil {
			fmt.Printf("Could not write %d.zip: %s\n", year, err.Error())			
		}

		fmt.Printf("+ [%d] Saved %s\n", id, url)

		results <- fn
	}
}

func minmax(v int, min int, max int) (int) {
	if v < min {
		return min
	} else if v > max {
		return max
	}

	return v
}

func main() {
	var s_year, e_year, wrx int
	var v bool

	this_year := time.Now().Year()

	/* read flags */
	flag.IntVar(&s_year, "start", 1921, "Start year. Default: 1940")
	flag.IntVar(&e_year, "end", this_year, "End year. Default: this year - 1.")
	flag.IntVar(&wrx, "w", 3, "Number of workers. Default: 3. Max: 10.")
	flag.BoolVar(&v, "v", false, "Enable verbose output")
	flag.Parse()

	/* display usage info */
	welcome_msg := `
Retrosheet Downloader

Config:

  - Range: %d - %d
  - Workers: %d

`
	fmt.Printf(welcome_msg, s_year, e_year, wrx)


	/*
	 create worker threads
	 */

	s_year = minmax(s_year, 1921, this_year)
	e_year = minmax(e_year, s_year, this_year)
	dx := e_year - s_year
	fmt.Println(dx)

	years := make(chan int, dx)
	results := make(chan string, dx)

	wrx = minmax(wrx, 1, 10)

	for i := 0; i < wrx; i ++ {
		go worker(i, years, results)
	}


	/*
	 feed work into threads
	 */

	/* remember that not all years are available */
	skippable_years := [14]int{
	  	1923, 1924, 1925, 1926, 1928, 1929,
  	  	1930, 1932, 1933, 1934, 1935, 1936, 1937, 1939}

	 Y:
	for year := s_year; year < e_year; year++ {
		for sy := range skippable_years {
		 	if year == sy {
		 		continue Y
		 	}
		}

		years <- year
		fmt.Println("Queued ", year)
	}

	/* no more work to be done across this channel */
	close(years)

	/* drain work pool */
	for year := s_year; year < e_year; year++ {
		<-results
	}

	fmt.Println("done")
}