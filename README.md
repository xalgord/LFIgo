> A faster LFI Fuzzer tool.

## Installation

```bash
go install github.com/xalgord/LFIgo@latest
```

## Usage

```
cat urls.txt | LFIgo
```

or

```
LFIgo --file urls.txt --threads 15
```
