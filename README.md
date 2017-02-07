# boom

Boom is a tool for manipulating bosh manifests. This tool targets not only to help with the quick manifest modification, but also to with the maintenance of Bosh configuration.

## Installation

* Mac OS

```

$ wget https://github.com/jszroberto/boom/releases/download/0.1/boom ; chmod +x boom; mv boom /usr/local/bin
```

* Linux:

```

$ wget https://github.com/jszroberto/boom/releases/download/0.1/boom-linux ; chmod +x boom-linux; mv boom-linux /usr/local/bin/boom
```


## Usage

```
NAME:
   boom - a simple and quick tool for bosh manifest maintenance

USAGE:
   boom [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     set-instances, si    Set the number of instances
     scale-instances, sc  Scale number of instances
     help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
