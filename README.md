# dotenv-gsm

Extends the already powerful [direnv](https://github.com/direnv/direnv) to enable sharing of sensitive environment variables through [Google Secret Manager](https://cloud.google.com/secret-manager/docs/overview)!

## Ideas
- add `dotenv-gsm init > ~/.config/direnv/direnvrc` command to generate the helper function automatically
- I would be welcome to making this more extensible to enable other secrets backends also if there is a desire!

# Usage

Add `dotenv-gsm` to your `PATH` (***NOTE: there's a difference between `dotenv-gsm` the binary, and `dotenv_gsm` the bash utility function***)

Add a helper function either to `~/.config/direnv/direnvrc`, or your `.envrc`:
```bash
# Usage: dotenv_gsm [<dotenv>]
#
# Loads a ".env.gsm" file and fetches values from Google Secret Manager
dotenv_gsm() {
  local path=${1:-}
  if [[ -z $path ]]; then
    path=$PWD/.env
  elif [[ -d $path ]]; then
    path=$path/.env
  fi
  watch_file "$path"
  if ! [[ -f $path ]]; then
    log_error ".env at $path not found"
    return 1
  fi
  if has dotenv-gsm; then
    eval "$("$direnv" dotenv bash "$@" | LOG_LEVEL=error dotenv-gsm -)"
  else
    log_error "ERROR: ensure you have 'dotenv-gsm' installed, and its available on your path!"
    return 1
  fi
}
```

***!! IMPORTANT !!*** Add a `.gitignore` to your project:
```bash
# If you want to do a similar pattern to this project you can use this example
.env*
!.env/
!.env/.env.gsm
!.env/.env.non-secret
```

Create a `.env.gsm` file with your secret paths:
```bash
TOP_SECRET_LATEST="project/test-project-id-1234/secrets/top-secret-name"
TOP_SECRET_OLDER="project/test-project-id-1234/secrets/top-secret-name/versions/1"
```

Add the new function to your `.envrc`:
```bash
dotenv_gsm .env.gsm
```

```bash
❯ direnv allow
direnv: loading ~/code/personal/dotenv-gsm/.envrc
direnv: export +TEST_SECRET +TOP_SECRET_LATEST +TOP_SECRET_OLDER

❯ echo $TOP_SECRET_LATEST
value-from-gsm-secret!
```

## Thank you!

- [direnv](https://github.com/direnv/direnv): Thank you `direnv` for such an amazing developer experience!
- [goreleaser](https://github.com/goreleaser/goreleaser): Thank you `goreleaser` for making the release process a pleasure!
