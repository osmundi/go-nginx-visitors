package main

import (
	"bufio"
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

type currentVisitors struct {
	curDay   map[string]struct{}
	curMonth map[string]struct{}
	curYear  map[string]struct{}
}

func (cv *currentVisitors) AddVisitor(ip string) {
	if cv.curDay == nil {
		cv.curDay = make(map[string]struct{})
	}
	if cv.curMonth == nil {
		cv.curMonth = make(map[string]struct{})
	}
	if cv.curYear == nil {
		cv.curYear = make(map[string]struct{})
	}

	if _, exists := cv.curDay[ip]; !exists {
		cv.curDay[ip] = struct{}{}
	}
	if _, exists := cv.curMonth[ip]; !exists {
		cv.curMonth[ip] = struct{}{}
	}
	if _, exists := cv.curYear[ip]; !exists {
		cv.curYear[ip] = struct{}{}
	}
}

func (v visitors) String() string {
	return fmt.Sprintf("New: %d / Old: %d", v.new, v.old)
}

type dailyVisitors struct {
	visitors *visitors
	unique   []string
}

type monthlyVisitors struct {
	visitors *visitors
	unique   []string
	daily    map[int]*dailyVisitors
}

type yearlyVisitors struct {
	visitors *visitors
	unique   []string
	monthly  map[int]*monthlyVisitors
}

type allVisitors struct {
	visitors *visitors
	unique   []string
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

func (all *allVisitors) AddNew(t time.Time) {
	all.yearly[t.Year()].visitors.new++
	all.yearly[t.Year()].monthly[int(t.Month())].visitors.new++
	all.yearly[t.Year()].monthly[int(t.Month())].daily[t.Day()].visitors.new++
	all.visitors.new++
}

func (all allVisitors) GetVisitorsFrom(date string) (error, visitors) {
	t, err := time.Parse("02/01/2006", date)
	if err != nil {
		t, err := time.Parse("01/2006", date)
		if err != nil {
			t, err := time.Parse("2006", date)
			if err != nil {
				return err, visitors{}
			}
			if _, exists := all.yearly[t.Year()]; !exists {
				return nil, visitors{}
			} else {
				return nil, visitors{new: all.yearly[t.Year()].visitors.new, old: all.yearly[t.Year()].visitors.old}
			}
		}
		if _, exists := all.yearly[t.Year()].monthly[int(t.Month())]; !exists {
			return nil, visitors{}
		} else {
			return nil, visitors{new: all.yearly[t.Year()].monthly[int(t.Month())].visitors.new, old: all.yearly[t.Year()].monthly[int(t.Month())].visitors.old}
		}
	}
	if _, exists := all.yearly[t.Year()].monthly[int(t.Month())].daily[t.Day()]; !exists {
		return nil, visitors{}
	} else {
		return nil, visitors{new: all.yearly[t.Year()].monthly[int(t.Month())].daily[t.Day()].visitors.new, old: all.yearly[t.Year()].monthly[int(t.Month())].daily[t.Day()].visitors.old}
	}
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

func (all allVisitors) Visualize() {
	fmt.Printf("draw")
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

	current := currentVisitors{curDay: make(map[string]struct{}), curMonth: make(map[string]struct{}), curYear: make(map[string]struct{})}

	fmt.Println(current)

	for _, arg := range os.Args[1:] {
		readLogFile(arg, &allVisitors, uniqueVisitors, &current)
	}

	err, v := allVisitors.GetVisitorsFrom("20/04/2022")
	if err != nil {
		fmt.Printf("Error parsing date: %v", err)
	}
	fmt.Printf("struct(new): 20/4/2022:%v\n", v.new)
	fmt.Printf("struct(new): 20/4/2022:%v\n", v.old)

	fmt.Println(len(uniqueVisitors))

	allVisitors.ShowMonthlyVisitors()
}

func readLogFile(logfile string, all *allVisitors, unique uniqueVisitors, current *currentVisitors) {
	file, err := os.Open(logfile)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var line, ipAddr string
	var splitLine []string

	// save crawler ips if needed later
	var bannedIps []string

	// keep the day/month/year which are being parsed in memory
	var lastDay, lastMonth, lastYear int

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

			//fmt.Printf("Current:\n - year: %d\n - month: %d\n - day: %d\n", len(current.curYear), len(current.curMonth), len(current.curDay))

			// add recurring users on daily/monthly/yearly basis
			if _, exists := all.yearly[t.Year()]; !exists {
				all.InitYear(t)

				// save unique users from last year
				if lastYear != 0 {
					if _, exists := all.yearly[lastYear]; exists {
						all.yearly[lastYear].visitors.old = len(current.curYear)
					}
				}
				current.curYear = nil
			}
			if _, exists := all.yearly[t.Year()].monthly[int(t.Month())]; !exists {
				all.InitMonth(t)

				// save unique users for current month
				if lastMonth != 0 && lastYear != 0 {
					if _, exists := all.yearly[lastYear].monthly[lastMonth]; exists {
						all.yearly[lastYear].monthly[lastMonth].visitors.old = len(current.curMonth)
					}
				}
				current.curMonth = nil
			}
			if _, exists := all.yearly[t.Year()].monthly[int(t.Month())].daily[t.Day()]; !exists {
				all.InitDate(t)

				// save unique users for current day
				if lastMonth != 0 && lastYear != 0 && lastDay != 0 {
					if _, exists := all.yearly[lastYear].monthly[lastMonth].daily[lastDay]; exists {
						all.yearly[lastYear].monthly[lastMonth].daily[lastDay].visitors.old = len(current.curDay)
					}
				}
				current.curDay = nil
			}

			current.AddVisitor(ipAddr)

			if _, exists := unique[ipAddr]; !exists {
				// increment unique visitors
				unique[ipAddr] = struct{}{}
				all.AddNew(t)
			}

			lastDay = t.Day()
			lastMonth = int(t.Month())
			lastYear = t.Year()
		} else {
			fmt.Printf("Log entry missing data: %v", splitLine)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
