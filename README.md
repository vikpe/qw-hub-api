# QW Hub API [![Test](https://github.com/vikpe/qw-hub-api/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/vikpe/qw-hub-api/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/vikpe/qw-hub-api/branch/main/graph/badge.svg)](https://codecov.io/gh/vikpe/qw-hub-api) [![Go Report Card](https://goreportcard.com/badge/github.com/vikpe/qw-hub-api)](https://goreportcard.com/report/github.com/vikpe/qw-hub-api)

> Web API serving QuakeWorld info

## Usage

1) Rename `config.sample.json` to `config.json`.
2) Build
3) `./qw-hub-api`

## Config

See [config.sample.json](./config.sample.json)

**Example**

```json
{
  "port": 3000,
  "servers": {
    "active_server_interval": 4,
    "server_interval": 30,
    "master_interval": 14400,
    "master_servers": [
      "master.quakeworld.nu:27000",
      "master.quakeservers.net:27000",
      "qwmaster.ocrana.de:27000",
      "qwmaster.fodquake.net:27000"
    ]
  },
   "qtv_demo_sources": [
    {"address": "qw.irc.ax:28000", "demo_date_format": "ymd"},
    {"address": "troopers.fi:28000", "demo_date_format": "ymd"},
    {"address": "qw.foppa.dk:28000", "demo_date_format": "dmy"}
  ],
  "streamers": {
    "annihilazor": "anni",
    "quakeworld": "[streambot]",
    "suddendeathTV": "suddendeathTV",
    "vikpe": "XantoM"
  }
}
```

## API endpoints

| URL                     | description                                                               |  
|-------------------------|---------------------------------------------------------------------------|
| `/v2/servers`           | All servers                                                               |  
| `/v2/servers/<address>` | Server details                                                            |  
| `/v2/servers/mvdsv`     | MVDSV servers                                                             |  
| `/v2/servers/qwfwd`     | QWFWD servers (proxies)                                                   |  
| `/v2/servers/qtv`       | QTV servers                                                               |
|                         |                                                                           |
| `/v2/masters/<address>` | List of servers on master                                                 |
|                         |                                                                           |
| `/v2/demos`             | Demos from popular servers                                                |  
| `/v2/streams`           | Twitch streams casting Quake                                              |  
| `/v2/events`            | Events (from [Wiki](https://wiki.quakeworld.nu/))                         |  
| `/v2/news`              | News (from [QuakeWorld.nu](https://www.quakeworld.nu/))                   |  
| `/v2/forum_posts`       | Forum posts (from [QuakeWorld.nu Forum](https://www.quakeworld.nu/forum)) |  

## Endpoint details

### MVDSV servers

> `/v2/servers/mvdsv`

**Query params**

| URL                 | description                                    |
|---------------------|------------------------------------------------|
| `has_player=xantom` | Servers where xantom is connected as player    |
| `has_client=xantom` | Servers where xantom is connected              |

### Demos

> `/v2/demos`

**Query params**

| URL                       | description                                             |
|---------------------------|---------------------------------------------------------|
| `query=2on2 xantom dm3`   | Demos where filename matches `2on2`, `xantom` and `dm3` |
| `mode=2on2`               | Demos with mode `2on2`                                  |
| `qtv_address=qw.foppa.dk` | Demos from `qw.foppa.dk` servers                        |
| `limit=10`                | Limit to `10` demos                                     |

## Build

```she
go build
```

## Development

Run locally.

```shell
./qw-hub-api
```

Now you try an endpoint, e.g. http://localhost:3000/v2/servers

## Credits

* eb
* Tuna
* XantoM

## Related projects

* [hub.quakeworld.nu](https://github.com/quakeworldnu/hub.quakeworld.nu)
* [streambot](https://github.com/vikpe/qw-streambot)
* [masterstat](https://github.com/vikpe/masterstat)
* [masterstat-cli](https://github.com/vikpe/masterstat-cli)
* [serverstat](https://github.com/vikpe/serverstat)
* [serverstat-cli](https://github.com/vikpe/serverstat-cli)
