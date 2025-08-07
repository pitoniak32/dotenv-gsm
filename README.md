# Usage

`go build -o bin/` (or somewhere on your path)

`.envrc`:
```bash
if [ -f .env.secret ]; then
  dotenv .env.secret

  if has dotenv_gsm; then
    eval "$(LOG_LEVEL=debug dotenv_gsm bash .env.secret)"
  fi
fi
```
