# GTFOcli

`GTFOcli` it's a Command Line Interface for easy binaries search commands that can be used to bypass local security restrictions in misconfigured systems.

![](https://github.com/cmd-tools/gtfocli/blob/main/docs/gtfocli.gif)
---
## Table of Contents

- [Table of Contents](#table-of-contents)
- [Installation](#installation)
- [Usage](#usage)
  - [Search for Unix binaries](#search-for-unix-binaries)
  - [Search for Windows binaries](#search-for-windows-binaries)
  - [Search using dockerized solution](#search-using-dockerized-solution)
  - [CTF](#CTF)
- [Credits](#credits)
- [Contributing](#contributing)

## Installation

Using `go`:

```shell
go install github.com/cmd-tools/gtfocli@latest
```

Using `homebrew`:

```shell
brew tap cmd-tools/homebrew-tap
brew install gtfocli
```

Using `docker`:

```shell
docker pull cmdtoolsowner/gtfocli
```

## Usage
### Search for unix binaries

Search for binary `tar`:

```shell
gtfocli search tar
```

Search for binary `tar` from `stdin`:

```shell
echo "tar" | gtfocli search
```

Search for binaries located into file;

```shell
cat myBinaryList.txt
/bin/bash
/bin/sh
tar
arp
/bin/tail

gtfocli search -f myBinaryList.txt
```

### Search for windows binaries

Search for binary `Winget.exe`:
```shell
gtfocli search Winget --os windows
```

Search for binary `Winget` from `stdin`:

```shell
echo "Winget" | gtfocli search --os windows
```

Search for binaries located into file:

```shell
cat windowsExecutableList.txt
Winget
c:\\Users\\Desktop\\Ssh
Stordiag
Bash
c:\\Users\\Runonce.exe
Cmdkey
c:\dir\subDir\Users\Certreq.exe

gtfocli search -f windowsExecutableList.txt --os windows
```

Search for binary `Winget` and print output in `yaml` format (see `-h` for available formats):

```shell
gtfocli search Winget -o yaml --os windows
```

### Search using dockerized solution

Examples:

Search for binary `Winget` and print output in `yaml` format:

```shell
docker run -i cmdtoolsowner/gtfocli search Winget -o yaml --os windows
```

Search for binary `tar` and print output in `json` format:

```shell
echo 'tar' | docker run -i cmdtoolsowner/gtfocli search -o json
```

Search for binaries located into file mounted as volume in the container:

```shell
cat myBinaryList.txt
/bin/bash
/bin/sh
tar
arp
/bin/tail

docker run -i -v $(pwd):/tmp cmdtoolsowner/gtfocli search -f /tmp/myBinaryList.txt
```

## CTF

An example of common use case for `gtfocli` is together with `find`:

```shell
find / -type f \( -perm 04000 -o -perm -u=s \) -exec gtfocli search {} \; 2>/dev/null
```

or

```shell
find / -type f \( -perm 04000 -o -perm -u=s \) 2>/dev/null | gtfocli search
```

## Credits
Thanks to [GTFOBins](https://gtfobins.github.io/) and [LOLBAS](https://lolbas-project.github.io/), without these projects `gtfocli` would never have come to light.

## Contributing
You want to contribute to this project? Wow, thanks! So please just fork it and send a pull request.
