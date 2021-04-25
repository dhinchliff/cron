# Cron Parser
A cron expression parser.

You will need Go 1.14 or higher on your system to build and test.

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
./cron "*/15 0-8/2 1-7 * 1-5"
minute        0 15 30 45
hour          0 2 4 6 8
day of month  1 2 3 4 5 6 7
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
```

To run tests run:

```bash
go test ./...
```