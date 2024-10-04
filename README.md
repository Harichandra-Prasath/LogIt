# LogIt

Simple Open Source Asynchronous Thread-Safe Logger built for Go  

## Table of Contents

1. [Installation](#installation)
2. [Usage](#usage)
3. [Future Work](#future-work)
4. [Contributions](#contributions)

## Installation  

1. Get the package  

```bash
go get -u github.com/Harichandra-Prasath/LogIt  
```
  
## Usage

The easiest way to get started is using `DefaultLogger`   

```go
logger := LogIt.DefaultLogger()  
defer logger.Flush()  
logger.Info("Hello","World")  
```

As logging calls are non-blocking, the main go-routine should wait to prevent logs being lost. Flush will log the remaining logs before the main go-routine exits.    

## Future-Work  
  
1. Support for JSON  
2. Additional Flags  

## Contributions 

This is a small project. Any contributions and improvements will be a massive help to make it better.  