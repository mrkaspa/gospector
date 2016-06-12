# Installation

Run the command

```bash
    go get github.com/mrkaspa/gospector
```

Use the command 

```bash
    gospector --help
    NAME:
       gospector - Check the README.md here httpds://github.com/mrkaspa/gospector
    
    USAGE:
       gospector [global options] command [command options] [arguments...]
    
    VERSION:
       1.0.0
    
    COMMANDS:
    GLOBAL OPTIONS:
       --dir value     Directory to gospect
       --config value  Config file for gospector
       --help, -h      show help
       --version, -v   print the version
```

Provide a valid dir and if it contains a file gospector.json valid you don't have to provide it explicitly.

A gospector.json looks like this:
 
```json
{
  "rules": [
    {
      "extensions": [
        ".rb"
      ],
      "words": [
        "puts"
      ]
    },
    {
      "extensions": [
        ".js"
      ],
      "words": [
        "console.log",
        "console.info"
      ]
    }
  ],
  "subdirs": ["spec", "webpack"]
}

```

With this configuration I want to check when I have some puts inside a *rb files and console.log, console.info inside *.js files inside the spec and webpack folders. 