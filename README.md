# QWS

> Web API serving QuakeWorld info

## Usage

```shell
qws [-master INTERVAL] [-server INTERVAL] [-active INTERVAL] [-port PORT]
```

| arg      | type  | description                   | default | 
|----------|-------|-------------------------------|---------|
| `port`   | `int` | HTTP listen port              | `3000`  |
| `master` | `int` | Master server update interval | `600`   |
| `server` | `int` | Server update interval        | `30`    |
| `active` | `int` | Active server update interval | `3`     |

## API endpoints

| URL                     | description                                                               |  
|-------------------------|---------------------------------------------------------------------------|
| `/v2/servers`           | All servers                                                               |  
| `/v2/servers/<address>` | Server details                                                            |  
| `/v2/servers/mvdsv`     | MVDSV servers                                                             |  
| `/v2/servers/qwfwd`     | QWFWD servers (proxies)                                                   |  
| `/v2/servers/qtv`       | QTV servers                                                               |
|                         |                                                                           |
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

```she
go build
```

## Development

Run locally on port `4000`.

```shell
./qws -port=4000
```

Now you try an endpoint, e.g. http://localhost:4000/v2/servers

## Credits

* eb
* Tuna
* XantoM

## See also

* [streambot](https://github.com/vikpe/qw-streambot)
* [masterstat](https://github.com/vikpe/masterstat)
* [masterstat-cli](https://github.com/vikpe/masterstat-cli)
* [serverstat](https://github.com/vikpe/serverstat)
* [serverstat-cli](https://github.com/vikpe/serverstat-cli)
