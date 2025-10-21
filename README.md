# pyupdate

A project for updating the dependencies in `pyproject.toml` to the latest version.

- [x] Support `uv`
- [x] Update main dependencies
- [x] Update dev dependencies
- [ ] Support `poetry`

```bash
pyupdate is a command-line tool that helps you update Python packages in your environment.

Usage:
  pyupdate [flags]

Flags:
  -d, --directory string   Path to directory containing pyproject.toml (default ".")
  -h, --help               help for pyupdate
  -s, --save-exact         Save exact versions of updated packages (default true)
```
