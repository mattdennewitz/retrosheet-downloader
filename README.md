# Retrosheet downloader

... is a bulk retrosheet event file downloader.
currently only supports regular season archives.

This app will skip missing regular season years.

## Building

Building from source:

```
$ git clone https://github.com/mattdennewitz/retrosheet-downloader/
$ make all
```

This will create a binary called `rs`.

## Usage

```shell
./rs --help
Usage of ./rs:
  -start=1921: Start year. Default: 1940
  -end=2013: End year. Default: this year - 1.
  -w=3: Number of workers. Default: 3. Max: 10.
  -out=".": Download output path. Default: "."
```

`rs` will, by default, spawn 3 download threads.
Files will be downloaded to ".". This can be changed using the `-out` flag.

### Examples:

Download all years using 3 (default) workers to "~/work/retrosheet/src/"

```shell
$ ./rs -start 1921 -out ~/work/retrosheet/src
```

Download 1999-2003 using 5 workers:

```shell
$ ./rs -start=1999 -end=2003 -w 5
```

Please go easy on the Retrosheet site.
