ghligh
======

ghligh can be used to manipulate pdf files in various ways.

### Usage:
-  ghligh [flags]
-  ghligh [command]

### Available Commands:
- `cat`         cat shows highlights pdf files
- `completion`  Generate the autocompletion script for the specified shell
- `export`      export pdf highlights into json
- `hash`        display the ghligh hash used to identify a documet [json]
- `help`        Help about any command
- `import`      import highlights from json file
- `info`        display info about pdf documents [json]
- `ls`          show files with highlights or tagged with 'ls' [unix]
- `tag`         manage pdf tags


### Flags:
  -h, --help   help for ghligh

Use `ghligh [command] --help` for more information about a command.

### todo
#### `ghligh ls`
- implement `-t` option [ ]
#### `ghligh browse`
- show file path if retrieved pdf info doens't have a name [ ]
- show highlights in text [ ]
- list documents opened [ ]
- command execution [ ]
- if no document is specified blank screen [ ]

#### `ghligh bookmark`
- like ghligh tag but use different magic string to store bookmarks [ ]

#### `ghligh action`
- show poppler actions of a pdf file [ ]

#### `ghligh serve / sync`
