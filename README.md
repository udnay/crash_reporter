# Simple CLI for setting up a system to collect core dumps and push to GCS

There are 2 main commands
- `crash_reporter initialize`
- `crash_reporter collect`

## initialize 
Sets up the host system to call `crash_reporter` when ever a core is produced.

Initialize will copy the `crash_reporter` binary to the path specified in the arguments. You can also pass `--collectArgs` flag to initialize which is a string of arguments to be passed to the `collect` comand when a core dump occurs.

This uses the [`core`](https://man7.org/linux/man-pages/man5/core.5.html) linux command. So `/proc/sys/kernel/core_pattern` will be set to:
```
|/path/to/crash_reporter collect <collectArgs>
```
**Note:** The `core_pattern` file has a limit of 128 bytes so keep file paths short and consider specifying a `--config` file to pass to the `collect` command.

## collect
This command reads in a core dump from `stdin` and persists it either locally or to a GCS bucket. You can specify a yaml file with options for the `collect` command.

```yaml
bucket: your-gcs-bucket
gcs: /tmp/gcp_credentials.json
sink: gcs
out: /crash_dumps
```

A full description of the flags can be seen by running `crash_reporter collect -h`

## Future work
Set up this CLI as a Kubernetes daemon set, to monitor and watch pods and services within a cluster and capture any core dumps.