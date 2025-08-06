# Usage

`go build`

`~/.config/direnv/direnvrc`:
```bash
#!/bin/bash

dotenv_gsm() {
  if ! has /Users/dvd/code/personal/direnv_gsm/direnv_gsm; then
    echo "please ensure you have direnv_gsm on your PATH"
  fi

  # load the current env file to cache the values.
  dotenv $1

  # fetch gsm secrets and update the env values
  eval "$(/Users/dvd/code/personal/direnv_gsm/direnv_gsm bash $1)"

  # dump and load the current environment into direnv
  direnv_load direnv dump
}
```
