# streamerslive

Command-line tool and Go library to fetch streaming channels from Youtube or Twitch.

## Commands

### add

Adds channels with the provided URLs to the list displayed by [list](#list). Multiple URLs are divided by one or more spaces.

```shell
streamerslive add https://www.twitch.tv/harukakaribu https://www.twitch.tv/gamesdonequick https://www.youtube.com/channel/UCPZgBtMYoFKypEG2SCvBN9A
```

### list

Displays saved streaming channels.

```shell
$ streamerslive list
+----+----------------+--------------------------------+
| ID |    CHANNEL     |          STREAM TITLE          |
+----+----------------+--------------------------------+
|  1 | arikacaribu    |                                |
|  2 | gamesdonequick | Awesome Games Done Quick 2021  |
|    |                | Online - Benefiting Prevent    |
|    |                | Cancer Foundation - Pokémon    |
|    |                | Platinum                       |
|  3 | rosedoodle     | 🌸「 😳 welcome to the COMFY   |
|    |                | zone 🌹💖|| beepu beepu!       |
|    |                | ✨」🌸【VTuber】               |
+----+----------------+--------------------------------+
```

### remove

Removes a channel from list.

```shell
$ streamerslive remove 3
rosedoodle removed
```

### watch

Opens stream in default browser.

```shell
streamerslive watch 2
```
