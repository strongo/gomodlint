# gomodlint

A linter for `go.mod` file - validates NO `replace` directive is committed

## Usage in pre-commit hook

```shell
if git status | grep -q ":\s*go\.mod"; then
  gomodlint
  exit_code=$?
  if [ $exit_code -ne 0 ]; then
    exit $exit_code
  fi
fi
```

## Usage in CI

```yaml
    - name: go.mod validation
      run: go run github.com/strongo/gomodlint@latest
```

## [License](LICENSE)

Free to use & modify - [MIT License](https://opensource.org/license/mit/)

## Used by:

- [github.com/sneat-co](https://github.com/sneat-co/) - an extendable platform that allows to build apps with extension modules.
- [DALgo](https://github.com/dal-go) - Database Abstraction Layer in GO language
