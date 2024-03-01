# GTFOcli

`GTFOcli` it's Command Line Interface for easy binaries search commands that can be used to bypass local security restrictions in misconfigured systems.

---
## Table of Contents

- [Table of Contents](#table-of-contents)
- [Installation](#installation)
- [Usage](#usage)
  - [Search for Unix binaries](#search-for-unix-binaries)
  - [Search for Windows binaries](#search-for-windows-binaries)
- [Contributing](#contributing)

## Installation

Using `go`

```shell
go get github.com/cmd-tools/gtfocli
```

## Usage
### Search for unix binaries
Search for binary `tar`
```shell
gtfocli search tar
```

Search for binary `tar` from `stdin`
```shell
echo "tar" | gtfocli search
```

Search for binaries located into file
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
Search for binary `Winget.exe`
```shell
gtfocli search Winget --os windows
```

Search for binary `Winget` from `stdin`
```shell
echo "Winget" | gtfocli search --os windows
```

Search for binaries located into file
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

Search for binary `Winget` and print output in `yaml` format (see `-h` for available formats)
```shell
gtfocli search tar -o yaml --os windows
```

## Contributing
You want to contribute to this project? Wow, thanks! So please just fork it and send me a pull request.

