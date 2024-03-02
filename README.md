# GLNotify

A CLI tool to get updates from your gitlab projects.
At the moment it will only notify you about new comments in merge requests.

## Install

Download the proper file in the [release section](https://github.com/pugkong/glnotify/releases) or build it yourself.

And create `glnotify/config.json` file in your config directory (`$XDG_CONFIG_HOME`, `$HOME/.config`,
`$HOME/Library/Application Support` or `%AppData%`)

```json
{
  "base_url": "https://gitlab.com",
  "token": "<your-token>",
  "user_id": 42,
  "project_ids": [43, 44, 45]
}
```

## Usage

Simply run `glnotify`, e.g.

```bash
$ glnotify
MR: Some MR
Commented by: John, Jane
https://gitlab.com/some-mr

$ glnotify

$
```
