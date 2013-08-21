# Retrosheet downloader

... is a bulk retrosheet event file downloader.
currently only supports regular season archives.

This app will skip missing regular season years.

## Usage

```shell
./rs --help
Usage of ./rs:
  -start=1921: Start year. Default: 1940
  -end=2013: End year. Default: this year - 1.
  -v=false: Enable verbose output
  -w=3: Number of workers. Default: 3. Max: 10.
```

or, download all years using 3 (default) workers:

```shell
$ ./rs
```

or, download 1999 - 2003 using 5 workers:

```shell
$ ./rs -start=1999 -end=2003 -w 5
```

Please go easy on the retrosheet site.
