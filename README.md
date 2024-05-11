> A faster LFI Fuzzer tool.

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
