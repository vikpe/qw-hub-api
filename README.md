# QW Hub API [![Test](https://github.com/vikpe/qw-hub-api/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/vikpe/qw-hub-api/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/vikpe/qw-hub-api/branch/main/graph/badge.svg)](https://codecov.io/gh/vikpe/qw-hub-api) [![Go Report Card](https://goreportcard.com/badge/github.com/vikpe/qw-hub-api)](https://goreportcard.com/report/github.com/vikpe/qw-hub-api)

> Web API serving QuakeWorld info

## Usage

1) Rename/edit `config.sample.json` to `config.json`.
2) Rename/edit `.env.example` to `.env`.
3) Build
4) `./qw-hub-api`

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
      "qwmaster.fodquake.net:27000"
    ]
  },
  "qtv_demo_sources": [
    {
      "address": "de.quake.world:28000",
      "demo_date_format": "ymd",
      "timezone": "UTC"
    },
    {
      "address": "troopers.fi:28000",
      "demo_date_format": "ymd",
      "timezone": "Europe/Helsinki"
    },
    {
      "address": "quake.se:28000",
      "demo_date_format": "Ymd",
      "timezone": "Europe/Stockholm"
    }
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

| URL                                | description                                                               |  
|------------------------------------|---------------------------------------------------------------------------|
| `/v2/servers`                      | All servers                                                               |  
| `/v2/servers/mvdsv`                | MVDSV servers                                                             |  
| `/v2/servers/qwfwd`                | QWFWD servers (proxies)                                                   |  
| `/v2/servers/qtv`                  | QTV servers                                                               |
|                                    |                                                                           |
| `/v2/servers/<address>`            | Server details                                                            |
| `/v2/servers/<address>/lastscores` | Server lastscores                                                         |
| `/v2/servers/<address>/laststats`  | Server laststats                                                          |
|                                    |                                                                           |
| `/v2/masters/<address>`            | List of servers on master                                                 |
|                                    |                                                                           |
| `/v2/demos`                        | Demos from popular servers                                                |  
| `/v2/streams`                      | Twitch streams casting Quake                                              |  
| `/v2/events`                       | Events (from [Wiki](https://wiki.quakeworld.nu/))                         |  
| `/v2/news`                         | News (from [QuakeWorld.nu](https://www.quakeworld.nu/))                   |  
| `/v2/forum_posts`                  | Forum posts (from [QuakeWorld.nu Forum](https://www.quakeworld.nu/forum)) |  

## Endpoint details

### MVDSV servers

> `/v2/servers/mvdsv`

| Param          | Type                     | Example             | Description                                   |
|----------------|--------------------------|---------------------|-----------------------------------------------|
| **empty**      | `include\|exclude\|only` | `empty=include`     | Include empty servers (default `exclude`)     |
| **hostname**   | `string`                 | `hostname=quake.se` | Servers matching hostname `quake.se`          |
| **has_player** | `string`                 | `has_player=xantom` | Servers where `xantom` is connected as player |
| **has_client** | `string`                 | `has_client=xantom` | Servers where `xantom` is connected           |
| **limit**      | `int`                    | `limit=10`          | Limit to `10` servers                         |

### Demos

> `/v2/demos`

| Param           | Type          | Default | Example                | Description                                             |
|-----------------|---------------|---------|------------------------|---------------------------------------------------------|
| **q**           | `string`      |         | `q=2on2 xantom dm3`    | Demos where filename matches `2on2`, `xantom` and `dm3` |
| **mode**        | `string`      |         | `mode=2on2`            | Demos with mode `2on2`                                  |
| **qtv_address** | `string`      |         | `qtv_address=quake.se` | Demos from `quake.se` qtv server                        |
| **limit**       | `int [1-500]` | `100`   | `limit=10`             | Limit to `10` demos                                     |

## Build

```shell
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
