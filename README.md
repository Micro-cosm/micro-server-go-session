# MicroServerGoSession

A quickstart project that bootstraps a Go-based backend service that spreads a session across your sprawling Mifedom.

**NOTE:**  Built to work seamlessly with [BinGo](https://github.com/wejafoo/bin-go) mife build/deploy utility

----

## Developer Install

Clone the git repository and link project root to your path.

$  `git clone git@github.com:micro-cosm/micro-server-go-session.git`

$  `go mod init weja.us/micro/micro-server-go-session`

$  `go mod tidy`


### Deploy Local Docker

$   `bingo --local`

### Deploy Remote Docker/Cloud

$   `bingo --remote --alias stage`
