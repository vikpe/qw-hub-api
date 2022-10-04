# QW Hub API

> Web API serving QuakeWorld info

## Usage

```shell
./qw-hub-api
```

## Config

See [config.json](./config.json)

**Sample config**

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
  "streamers": {
    "annihilazor": "anni",
    "badsebitv": "badsebitv",
    "bogojoker": "bogojoker",
    "bps__": "bps"
  }
}
```

## API endpoints

| URL                      | description                                                               |  
|--------------------------|---------------------------------------------------------------------------|
| `/v2/servers`            | All servers                                                               |  
| `/v2/servers/<address>`  | Server details                                                            |  
| `/v2/servers/mvdsv`      | MVDSV servers                                                             |  
| `/v2/servers/qwfwd`      | QWFWD servers (proxies)                                                   |  
| `/v2/servers/qtv`        | QTV servers                                                               |
|                          |                                                                           |
| `/v2/masters/<address>`  | List of servers on master                                                 |
|                          |                                                                           |
| `/v2/streams`            | Twitch streams casting Quake                                              |  
| `/v2/events`             | Events (from [Wiki](https://wiki.quakeworld.nu/))                         |  
| `/v2/news`               | News (from [QuakeWorld.nu](https://www.quakeworld.nu/))                   |  
| `/v2/forum_posts`        | Forum posts (from [QuakeWorld.nu Forum](https://www.quakeworld.nu/forum)) |  

## Endpoint details

### MVDSV servers

> `/v2/servers/mvdsv`

**Query params**

| URL                 | description                                    |
|---------------------|------------------------------------------------|
| `has_player=xantom` | Servers where xantom is connected as player    |
| `has_client=xantom` | Servers where xantom is connected              |

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
