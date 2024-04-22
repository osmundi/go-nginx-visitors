# go-nginx-visitors
Calculate daily/monthly/yearly unique and recurring IP addresses from nginx access logs

### Disclaimer
Identifying a user solely based on their IP address can be challenging, especially if the IP address changes frequently. This program helps provide insight into the number of users visiting a website in cases where there is no session management or sophisticated IP tracking in use.


### Usage
Pass access logs as command line arguments. The program logic assumes that the log files are sorted by date in ascending order.
```go run main.go access.log.oldest access.log acces.log.latest```

Depending of the API endpoints in your application it might be good idea to filter only the requests containing specific string by passing a filter parameter(s) to the program. In this case only the log rows containing "api/login", "api/search" or "home" are counted as real users and all other rows are skipped:
```go run main.go -filter=api/login -filter=api/search -filter=home access.log```


### Outputting and visualizing data

By default the program outputs data as tab seperated values so its up to user how it will be visualized, e.g. with gnuplot:

```go run main.go access.log > test.dat```


```gnuplot -e "set title 'Nginx access.log unique and recurring visitors'; set auto x; set terminal png; set output 'output.png'; set yrange [0:*]; set style data histogram; set style histogram cluster gap 1; set style fill solid border -1; set boxwidth 0.9; set xtic rotate by -45 scale 0; plot 'test.dat' using 2:xtic(1) ti col, '' u 3 ti col"```

Its possible to read visitor data from specifig data, month or year also.

Usage examples:
~~~
fmt.Println(len(uniqueVisitors))

err, v := allVisitors.GetVisitorsFrom("20/04/2022")
if err != nil {
	fmt.Printf("Error parsing date: %v", err)
}
fmt.Printf("New visitors in 20/4/2022:%v\n", v.new)
fmt.Printf("Old visitors in 20/4/2022:%v\n", v.old)

err, v2 := allVisitors.GetVisitorsFrom("2023")
if err != nil {
	fmt.Printf("Error parsing date: %v", err)
}
fmt.Printf("New visitors in 2023:%v\n", v2.new)
fmt.Printf("Old visitors in 2023:%v\n", v2.old)
~~~

### Testing
go test
