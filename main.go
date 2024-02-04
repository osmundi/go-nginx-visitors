package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

// types
type visitors struct {
	new int
	old int
}

func (v visitors) String() string {
	return fmt.Sprintf("New: %d / Old: %d", v.new, v.old)
}

type dailyVisitors struct {
	visitors *visitors
}

type monthlyVisitors struct {
	visitors *visitors
	daily    map[int]*dailyVisitors
}

type yearlyVisitors struct {
	visitors *visitors
	monthly  map[int]*monthlyVisitors
}

type allVisitors struct {
	visitors *visitors
	yearly   map[int]*yearlyVisitors
}

func (all *allVisitors) InitDate(t time.Time) {
	all.yearly[t.Year()].monthly[int(t.Month())].daily[t.Day()] = new(dailyVisitors)
	all.yearly[t.Year()].monthly[int(t.Month())].daily[t.Day()].visitors = new(visitors)
}

func (all *allVisitors) InitMonth(t time.Time) {
	all.yearly[t.Year()].monthly[int(t.Month())] = new(monthlyVisitors)
	all.yearly[t.Year()].monthly[int(t.Month())].daily = make(map[int]*dailyVisitors)
	all.yearly[t.Year()].monthly[int(t.Month())].visitors = new(visitors)
}

func (all *allVisitors) InitYear(t time.Time) {
	all.yearly[t.Year()] = new(yearlyVisitors)
	all.yearly[t.Year()].monthly = make(map[int]*monthlyVisitors)
	all.yearly[t.Year()].visitors = new(visitors)
}

func (all *allVisitors) AddNew(t time.Time, newVisitor bool) {
	if newVisitor {
		all.yearly[t.Year()].visitors.new++
		all.yearly[t.Year()].monthly[int(t.Month())].visitors.new++
		all.yearly[t.Year()].monthly[int(t.Month())].daily[t.Day()].visitors.new++
		all.visitors.new++
	} else {
		all.yearly[t.Year()].visitors.old++
		all.yearly[t.Year()].monthly[int(t.Month())].visitors.old++
		all.yearly[t.Year()].monthly[int(t.Month())].daily[t.Day()].visitors.old++
		all.visitors.old++
	}
}

func (all allVisitors) GetNew(date string) (error, int) {
	t, err := time.Parse("02/01/2006", date)
	if err != nil {
		t, err := time.Parse("01/2006", date)
		if err != nil {
			t, err := time.Parse("2006", date)
			if err != nil {
				return err, 0
			}
			return nil, all.yearly[t.Year()].visitors.new
		}
		return nil, all.yearly[t.Year()].monthly[int(t.Month())].visitors.new
	}
	return nil, all.yearly[t.Year()].monthly[int(t.Month())].daily[t.Day()].visitors.new
}

func (all allVisitors) ShowMonthlyVisitors() {
	months := make([]int, 0)
	for year, yearly := range all.yearly {
		if year < 1970 {
			continue
		}
		fmt.Printf("Year: %d\n", year)

		for month, _ := range yearly.monthly {
			months = append(months, month)
		}
		sort.Ints(months)
		for _, k := range months {
			fmt.Printf("Month: %d\n", k)
			fmt.Printf("new: %d / old: %d\n", yearly.monthly[k].visitors.new, yearly.monthly[k].visitors.old)
		}
		months = nil
	}
}

// util
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func appendUnique(uniqueStrings map[string]struct{}, value string) {
	// Check if the value is already in the list
	if _, exists := uniqueStrings[value]; !exists {
		// Add the value to the list
		uniqueStrings[value] = struct{}{}
	}
}

func prettyPrint(data map[int]map[int]*visitors) (string, error) {
	prettyJSON, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}

func isCrawler(line string, bannedIps *[]string) bool {
	crawlers := []string{"YandexBot", "SoftDev", "UptimeRobot", "Nessus", "GoogleBot"}
	for _, s := range crawlers {
		if strings.Contains(line, s) {
			*bannedIps = append(*bannedIps, line)
			return true
		}
	}
	return false
}

type uniqueVisitors map[string]struct{}

func main() {
	// Check if there are command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <log.file.1> <log.file.2> <log.file.n>")
		return
	}

	// datastructure to hold visitor data
	allVisitors := allVisitors{visitors: new(visitors), yearly: make(map[int]*yearlyVisitors)}

	// ip addresses from logs (with map its possible to distinct unique visitors)
	uniqueVisitors := make(map[string]struct{})

	for _, arg := range os.Args[1:] {
		readLogFile(arg, &allVisitors, uniqueVisitors)
	}

	err, new := allVisitors.GetNew("20/04/2022")
	if err != nil {
		fmt.Printf("Error parsing date: %v", err)
	}
	fmt.Printf("struct(new): 20/4/2022:%v\n", new)

	fmt.Println(len(uniqueVisitors))

	allVisitors.ShowMonthlyVisitors()
}

func readLogFile(logfile string, all *allVisitors, unique uniqueVisitors) {
	file, err := os.Open(logfile)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var line, ipAddr string
	var splitLine []string

	var bannedIps []string
	for scanner.Scan() {
		line = scanner.Text()

		if isCrawler(line, &bannedIps) {
			continue
		}

		splitLine = strings.Split(line, " ")

		if len(splitLine) > 2 {
			ipAddr = splitLine[0]

			t, err := time.Parse(
				"02/Jan/2006:15:04:05",
				strings.TrimLeft(strings.Split(line, " ")[3], "["),
			)
			if err != nil {
				fmt.Println(err)
				fmt.Printf("%v\n", splitLine)
			}

			// init datastructures
			if _, exists := all.yearly[t.Year()]; !exists {
				all.InitYear(t)
			}
			if _, exists := all.yearly[t.Year()].monthly[int(t.Month())]; !exists {
				all.InitMonth(t)
			}
			if _, exists := all.yearly[t.Year()].monthly[int(t.Month())].daily[t.Day()]; !exists {
				all.InitDate(t)
			}

			if _, exists := unique[ipAddr]; !exists {
				// increment unique visitors
				unique[ipAddr] = struct{}{}
				all.AddNew(t, true)
			} else {
				// increment recurring visitors
				all.AddNew(t, false)
			}

		} else {
			fmt.Printf("Log entry missing data: %v", splitLine)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
