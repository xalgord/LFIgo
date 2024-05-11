> **LFIgo**: LFIGo is a lightweight and efficient tool designed for detecting and exploiting Local File Inclusion (LFI) vulnerabilities in web applications. With its streamlined approach, LFIScanGo offers rapid scanning and targeted exploitation capabilities, making it an essential tool for security researchers and Bug Bounty Hunters.

## Installation

```bash
go install github.com/xalgord/LFIgo@latest
```

## Usage

```
 -file string
        File containing URLs
  -h    Show help message (shorthand)
  -help
        Show help message
  -threads int
        Number of threads to use (default 10)
```

## Example

```
cat urls.txt | LFIgo
```

or

```
LFIgo --file urls.txt --threads 15
```


> ⚠️ This tool may pose security risks. It's recommended to use with caution and additional security measures.
