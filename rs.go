package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func worker(id int, years <-chan int, results chan<- string) {
	for year := range years {
		/* create file */
		fn := fmt.Sprintf("%deve.zip", year)
		out_path := filepath.Join(os.TempDir(), fn)
		out_f, err := os.Create(out_path)
		defer out_f.Close()

		if err != nil {
		 	fmt.Printf("Error: could not create %d.zip: %s\n", year, err.Error())
		 	continue
		}

		/* download archive */
		url := fmt.Sprintf("http://www.retrosheet.org/events/%deve.zip", year)

		resp, err := http.Get(url)
		defer resp.Body.Close()

		if err != nil {
			fmt.Printf("Error: could not download %s: %s\n", url, err.Error())
			continue
		}

		/* write output */
		_, err = io.Copy(out_f, resp.Body)

		if err != nil {
			fmt.Printf("Error: could not write %d.zip: %s\n", year, err.Error())			
		}

		fmt.Printf("+ Saved %s\n", year)

		results <- out_path
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
	var user_path string

	this_year := time.Now().Year()

	/* read flags */
	flag.IntVar(&s_year, "start", 1921, "Start year. Default: 1940")
	flag.IntVar(&e_year, "end", this_year, "End year, inclusive. Default: this year - 1.")
	flag.IntVar(&wrx, "w", 3, "Number of workers. Default: 3. Max: 10.")
	flag.StringVar(&user_path, "out", ".", "Download output path. Default: '.'")
	flag.Parse()


	/*
	 ensure output path is aces
	 */

	/* create full output path */
	if user_path, err := filepath.Abs(user_path); err != nil {
		fmt.Println("Error: could not resolve path:", user_path)
		os.Exit(1)
	}

	/* ... and ensure it exists */
	if s, err := os.Stat(user_path); err != nil || !s.IsDir() {
		fmt.Printf("Error: path %s must exist and be a directory\n", user_path)
		os.Exit(1)
	}


	/* display usage info */
	welcome_msg := `
Retrosheet Downloader
  - Range: %d - %d
  - Workers: %d
  - Output to: %s

`
	fmt.Printf(welcome_msg, s_year, e_year, wrx, user_path)


	/*
	 create worker threads
	 */

	s_year = minmax(s_year, 1921, this_year)
	e_year = minmax(e_year, s_year, this_year)
	dx := e_year - s_year

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
	for year := s_year; year <= e_year; year++ {
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
	for year := s_year; year <= e_year; year++ {
		out_path := <-results

		/* move file into place */
		os.Rename(out_path, filepath.Join(user_path, filepath.Base(out_path)))
	}
}