# Alfred MetaCPAN Workflow

A workflow for [Alfred](http://www.alfredapp.com/) that searches Perl modules from [MetaCPAN](https://metacpan.org/) incrementally.

## Example

![Searching modules](https://cloud.githubusercontent.com/assets/115636/5282423/7918f9f6-7b49-11e4-8ce5-b75eb2b77c39.png)

## Installation

1. Go to the [latest release](https://github.com/oalders/alfred-metacpan-workflow/releases/latest)
2. Download the `.alfredworkflow` file for your Mac:
   - `metacpan-*-darwin-arm64.alfredworkflow` — Apple Silicon (M1/M2/M3/M4)
   - `metacpan-*-darwin-amd64.alfredworkflow` — Intel
3. Double-click the downloaded file to install it in Alfred

## Building locally

```sh
mkdir -p dist && make && open dist/metacpan-*.alfredworkflow
```

## Licence

[MIT](LICENSE)

## Author

[handlename](https://github.com/handlename) (original) / [oalders](https://github.com/oalders) (fork)
