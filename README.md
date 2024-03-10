# GLNotify

[![build](https://github.com/PugKong/glnotify/actions/workflows/release.yml/badge.svg)](https://github.com/PugKong/glnotify/actions/workflows/release.yml)
[![build](https://github.com/PugKong/glnotify/actions/workflows/test.yml/badge.svg)](https://github.com/PugKong/glnotify/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/github/PugKong/glnotify/badge.svg?branch=master)](https://coveralls.io/github/PugKong/glnotify?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/PugKong/glnotify)](https://goreportcard.com/report/github.com/PugKong/glnotify)
[![LoC](https://tokei.rs/b1/github/PugKong/glnotify)](https://github.com/PugKong/glnotify)
[![Release](https://img.shields.io/github/release/PugKong/glnotify.svg?style=flat-square)](https://github.com/PugKong/glnotify/releases/latest)

A CLI tool to get updates from your gitlab projects.
At the moment it will only notify you about new comments and labels in merge requests.

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
