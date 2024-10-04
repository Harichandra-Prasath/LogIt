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
  
Quick-Notes on `DefaultLogger`  
1. `DefaultLogger` will write to os.Stderr for ERROR logs and os.Stdout for rest  
2. Standard Flags are used for creating a log record which is Date and Time.  
3. `DefaultLogger` will be created with level `LogIT.LEVEL_INFO`. Debug Level logs will be ignored.  
   
You can create a custom logger with `LogIt.NewLogger()`  

```go
opts := LogIt.LoggerOptions{
    Level: LogIt.LEVEL_DEBUG,
    RecordOptions: LogIt.RecordOptions{
        Spacing:   1,
        Flags:     LogIt.DATE_FLAG | LogIt.TIME_FLAG,
    },
}

logger := LogIt.NewLogger(opts, LogIt.NewTextHandler(os.Stdout, os.Stderr))
defer logger.Flush()

logger.Debug("Hello", "World")
```

Quick-Notes on Custom Loggers  
1. You can create a logger with prefered level.  
2. You can direct the output of logs to preferred `io.Writer`.    
3. Custom Options for log records includes Flags, Pretty Output, Spacing between Flags.  

## Future-Work  
  
1. Support for JSON  
2. Additional Flags  

## Contributions 

This is a small project. Any contributions and improvements will be a massive help to make it better.  