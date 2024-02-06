# go-nginx-visitors
Calculate daily/monthly/yearly unique and recurring IP addresses from nginx access logs

### Testing

go test

### Outputting and visualizing data

```go run main.go access.log > test.dat```

go-nginx-visitors output data as tab seperated values so its up to user how it will be visualized, e.g. with gnuplot:

```gnuplot -e "set title 'Nginx access.log unique and recurring visitors'; set auto x; set terminal png; set output 'output.png'; set yrange [0:*]; set style data histogram; set style histogram cluster gap 1; set style fill solid border -1; set boxwidth 0.9; set xtic rotate by -45 scale 0; plot 'test.dat' using 2:xtic(1) ti col, '' u 3 ti col"```



### Disclaimer
Identifying a user solely based on their IP address can be challenging, especially if the IP address changes frequently. This program helps provide insight into the number of users visiting a website in cases where there is no session management or sophisticated IP tracking in use.

gnuplot -e "set title 'Nginx access.log unique and recurring visitors';set auto x; set terminal png; set output 'output.png';set yrange [0:*]; set style data histogram; set style histogram cluster gap 1; set style fill solid border -1; set boxwidth 0.9; set xtic rotate by -45 scale 0; plot 'test.dat' using 2:xtic(1) ti col, '' u 3 ti col"
