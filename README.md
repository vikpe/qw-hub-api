# QWS

> Web API serving QuakeWorld server info

## Usage

```sh
qws [-master INTERVAL] [-server INTERVAL] [-active INTERVAL] [-port PORT]
```

| arg      | type  | description                   | default | 
|----------|-------|-------------------------------|---------|
| `port`   | `int` | HTTP listen port              | `3000`  |
| `master` | `int` | Master server update interval | `600`   |
| `server` | `int` | Server update interval        | `30`    |
| `active` | `int` | Active server update interval | `3`     |

## API endpoints

| URL                    | description                                  |  
|------------------------|----------------------------------------------|
| `/v2/server/<address>` | Server details                               |  
| `/v2/mvdsv`            | Mvdsv servers                                |  
| `/v2/qwfwd`            | Qwfwd servers (proxies)                      |  
| `/v2/qtv`              | QTV servers                                  |  
| `/v2/qtv_to_mvdsv`     | Map of QTV streams to mvdsv server addresses |  
| `/v2/mvdsv_to_qtv`     | Map of mvdsv server addresses to QTV streams |
| `/v2/clients`          | Clients                                      |

### Query params

| URL                 | description                                    |
|---------------------|------------------------------------------------|
| `status=started`    | Servers where `Status` is `Started`            |
| `mode=ffa`          | Servers where `Mode` is `ffa`                  |
| `mode=2on2,4on4`    | Servers where `Mode` is `2on2` or `4on4`       |
|                     |                                                |
| `has_player=xantom` | Servers where xantom is connected as player    |
| `has_human_players` | Servers with at least 1 human player           |
| `has_human_players` | Servers with at least 1 human player           |
|                     |                                                |
| `geo.cc=dk`         | Servers where `Country Code` is `DK` (Denmark) |
| `geo.region=asia`   | Servers where `Region` is `Asia`               |
|                     |                                                |
| `sort_by=address`   | Sort by `server address`                       |
| `sort_order=desc`   | Sort in `descending` order                     |
| `limit=5`           | Limit result to `5` servers                    |

## Config

### Master servers

The QuakeWorld master servers to query for servers.

**Example**
`master_servers.json`

```json
[
  "master.quakeworld.nu:27000",
  "master.quakeservers.net:27000",
  "qwmaster.ocrana.de:27000",
  "qwmaster.fodquake.net:27000"
]
```

## Build

```sh
$ go build
```

## Credits

* eb
* Tuna
* XantoM

## See also

* [masterstat](https://github.com/vikpe/masterstat)
* [masterstat-cli](https://github.com/vikpe/masterstat-cli)
* [serverstat](https://github.com/vikpe/serverstat)
* [serverstat-cli](https://github.com/vikpe/serverstat-cli)
