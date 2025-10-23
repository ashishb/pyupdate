# pyupdate

[![Lint GitHub Actions](https://github.com/ashishb/pyupdate/actions/workflows/lint-github-actions.yaml/badge.svg)](https://github.com/ashishb/pyupdate/actions/workflows/lint-github-actions.yaml)
[![Lint YAML](https://github.com/ashishb/pyupdate/actions/workflows/lint-yaml.yaml/badge.svg)](https://github.com/ashishb/pyupdate/actions/workflows/lint-yaml.yaml)
[![Lint Markdown](https://github.com/ashishb/pyupdate/actions/workflows/lint-markdown.yaml/badge.svg)](https://github.com/ashishb/pyupdate/actions/workflows/lint-markdown.yaml)

[![Lint Go](https://github.com/ashishb/pyupdate/actions/workflows/lint-go.yaml/badge.svg)](https://github.com/ashishb/pyupdate/actions/workflows/lint-go.yaml)
[![Validate Go code formatting](https://github.com/ashishb/pyupdate/actions/workflows/format-go.yaml/badge.svg)](https://github.com/ashishb/pyupdate/actions/workflows/format-go.yaml)

A project for updating the dependencies in `pyproject.toml` to their latest versions.

- [x] Support `uv`
- [x] Update main dependencies
- [x] Update dev dependencies
- [ ] Support `poetry`

```bash
pyupdate is a command-line tool that helps you update packages in your Python project.

Usage:
  pyupdate [flags]

Flags:
  -d, --directory string   Path to directory containing pyproject.toml (default ".")
  -h, --help               help for pyupdate
  -s, --save-exact         Save exact versions of updated packages (default true)
```
