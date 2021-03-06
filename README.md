# Cron Parser
A cron expression parser

To build the application run:

```bash
go build
```

Usage examples:

```bash
./cron "1 1 1 1 1 /usr/bin/find"
minute        1
hour          1
day of month  1
month         1
day of week   1
command       /usr/bin/find
```

```bash
./cron "1,2 0,1 2,3 4,5 0,1 /usr/bin/find"
minute        1 2
hour          0 1
day of month  2 3
month         4 5
day of week   0 1
command       /usr/bin/find
```

```bash
./cron "0 * * * * /usr/bin/find"
minute        0
hour          0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23
day of month  1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   0 1 2 3 4 5 6
command       /usr/bin/find
```

To run tests run:

```bash
go test ./...
```