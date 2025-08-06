# Usage

`go build -o bin/`

`.envrc`:
```bash
if [ -f .env.secret ]; then
  dotenv .env.secret

  eval "$(LOG_LEVEL=error /Users/dvd/code/personal/dotenv_gsm/dotenv_gsm bash .env.secret)"
fi
```
