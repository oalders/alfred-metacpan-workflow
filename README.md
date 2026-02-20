# Alfred MetaCPAN Workflow

A workflow for [Alfred](http://www.alfredapp.com/) that searches Perl modules from [MetaCPAN](https://metacpan.org/) incrementally.

## Example

![Searching modules](https://cloud.githubusercontent.com/assets/115636/5282423/7918f9f6-7b49-11e4-8ce5-b75eb2b77c39.png)

## Installation

### Build locally (recommended)

If you have Go installed, this is the simplest option — no Gatekeeper issues:

```sh
mkdir -p dist && make && open dist/metacpan-*.alfredworkflow
```

### Download a release

1. Go to the [latest release](https://github.com/oalders/alfred-metacpan-workflow/releases/latest)
2. Download the `.alfredworkflow` file for your Mac:
   - `metacpan-*-darwin-arm64.alfredworkflow` — Apple Silicon (M1/M2/M3/M4)
   - `metacpan-*-darwin-amd64.alfredworkflow` — Intel
3. Double-click the downloaded file to install it in Alfred

Because the binary is not signed with an Apple Developer certificate, macOS may
block it when Alfred tries to run it. If that happens, run this after installing:

```sh
dir=$(find ~/Library/Application\ Support/Alfred/Alfred.alfredpreferences/workflows \
  -name info.plist -exec grep -l "oalders.alfredmetacpan" {} \; \
  | sed 's|/info.plist||') && xattr -dr com.apple.quarantine "$dir"
```

## Licence

[MIT](LICENSE)

## Author

[handlename](https://github.com/handlename) (original) / [oalders](https://github.com/oalders) (fork)
