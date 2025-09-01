## Overview

A tool to quickly convert JSON into a Go struct or TypeScript interface.

## Flags

| Flag               | Description                             |
| ------------------ | --------------------------------------- |
| `-h`, `--help`     | help for jsontypify                     |
| `-i`, `--indent`   | output indentation (default: 4)         |
| `-l`, `--language` | output to especific language (go, ts)   |
| `-n`, `--name`     | struct/interface name (default: "Main") |

## Installation

To install jsontypify, run the following command:

```bash
go install github.com/fdanctl/jsontypify
```

## Usage

```bash
jsontypify [flags] <file_path>
```

- `<file_path>` can be a file path or '-' to use stdin.

Examples:

```bash
jsontypify --language go input.json
```

```sh
curl https://db.ygoprodeck.com/api/v7/cardinfo.php?name=Quillbolt%20Hedgehog | jsontypify -
```

## Neovim intergration

Generates the Go struct or TypeScript interface code based on the current file and clipboard contents.
The generated code is inserted directly into the current file being edited.

```lua
function jsontypify()
    local str = vim.fn.getreg('"')
    str = str:gsub("'", "'\\''")

    local fileType = vim.bo.filetype
    if fileType == "typescriptreact" or fileType == "typescript" then
        fileType = "ts"
    end

    -- echo '{"param": "val"}' | jsontypify -i 4 -l go -
    local cmd = "echo '" .. str .. "' | " .. "jsontypify -i " .. vim.bo.tabstop .. " -l " .. fileType .. " -"

    local output = vim.fn.systemlist(cmd) -- runs the command and splits into lines
    vim.api.nvim_put(output, "l", true, true)
end

-- keymap to trigger the function
vim.keymap.set("n", "<leader>jt", jsontypify, { noremap = true, silent = true})
```
