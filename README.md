# slp

slp is a MySQL/PostgreSQL SlowLog Profiler.

[日本語](./README.ja.md)

This tool is similar to [mysqldumpslow](https://dev.mysql.com/doc/refman/8.0/en/mysqldumpslow.html), but can check more metrics.

## Installation

Download from https://github.com/tkuchiki/slp/releases

## Usage

```console
$ slp --help
slp is a MySQL/PostgreSQL SlowLog Profiler

Usage:
  slp [flags]
  slp [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  diff        Show the difference between the two profile results
  help        Help about any command
  my          Profile the slowlogs for MySQL
  pg          Profile the slowlogs for PostgreSQL

Flags:
      --config string   The configuration file
  -h, --help            help for slp
  -v, --version         version for slp

Use "slp [command] --help" for more information about a command.

$ cat example/mysql.slow.log | slp my
+-------+---------------------------------+----------------+----------------+----------------+----------------+
| COUNT |              QUERY              | MIN(QUERYTIME) | MAX(QUERYTIME) | SUM(QUERYTIME) | AVG(QUERYTIME) |
+-------+---------------------------------+----------------+----------------+----------------+----------------+
| 1     | DELETE FROM `t2` WHERE 'S'      | 0.369618       | 0.369618       | 0.369618       | 0.369618       |
|       | < `c1_date` OR `c2` NOT IN      |                |                |                |                |
|       | (SELECT `c3` FROM `t3`)         |                |                |                |                |
| 1     | DELETE FROM `t4` WHERE `c4`     | 7.148949       | 7.148949       | 7.148949       | 7.148949       |
|       | NOT IN (SELECT `c1` FROM `t1`)  |                |                |                |                |
| 1     | INSERT INTO `t2`                | 0.010498       | 0.010498       | 0.010498       | 0.010498       |
|       | (`c2_id`,`c2_string`,`c2_date`) |                |                |                |                |
|       | VALUES (N,'S','S')              |                |                |                |                |
| 1     | INSERT INTO `t2`                | 0.010498       | 0.010498       | 0.010498       | 0.010498       |
|       | (`c2_id`,`c2_string`,`c2_date`) |                |                |                |                |
|       | VALUES (N,'S','S'),(N,'S','S')  |                |                |                |                |
| 1     | SELECT * FROM `t5` WHERE        | 0.010753       | 0.010753       | 0.010753       | 0.010753       |
|       | `c5_id` IN ('S','S','S')        |                |                |                |                |
| 1     | SELECT `t1`.`id` FROM `t1`      | 0.020219       | 0.020219       | 0.020219       | 0.020219       |
|       | JOIN `t2` ON `t2`.`t1_id` =     |                |                |                |                |
|       | `t1`.`id` WHERE `t2`.`t1_id` =  |                |                |                |                |
|       | 'S' ORDER BY `t2`.`t1_id`       |                |                |                |                |
| 2     | UPDATE `t1` SET                 | 1.428614       | 3.504247       | 4.932861       | 2.466430       |
|       | `c1_count`=(SELECT COUNT(N) AS  |                |                |                |                |
|       | `cnt` FROM `t2` WHERE `c3_id`   |                |                |                |                |
|       | = `t3`.`id`)                    |                |                |                |                |
+-------+---------------------------------+----------------+----------------+----------------+----------------+

$ cat example/postgresql.slow.log | slp pg
+-------+--------------------------------+----------------+----------------+----------------+----------------+
| COUNT |             QUERY              | MIN(QUERYTIME) | MAX(QUERYTIME) | SUM(QUERYTIME) | AVG(QUERYTIME) |
+-------+--------------------------------+----------------+----------------+----------------+----------------+
| 1     | DELETE FROM t2 WHERE 'S' <     | 0.369618       | 0.369618       | 0.369618       | 0.369618       |
|       | c1_date OR NOT c2 IN (SELECT   |                |                |                |                |
|       | c3 FROM t3)                    |                |                |                |                |
| 1     | DELETE FROM t4 WHERE NOT c4 IN | 7.148949       | 7.148949       | 7.148949       | 7.148949       |
|       | (SELECT c1 FROM t1)            |                |                |                |                |
| 1     | INSERT INTO t2 (c2_id,         | 0.010498       | 0.010498       | 0.010498       | 0.010498       |
|       | c2_string, c2_date) VALUES (N, |                |                |                |                |
|       | 'S', 'S')                      |                |                |                |                |
| 1     | INSERT INTO t2 (c2_id,         | 0.010498       | 0.010498       | 0.010498       | 0.010498       |
|       | c2_string, c2_date) VALUES (N, |                |                |                |                |
|       | 'S', 'S'), (N, 'S', 'S')       |                |                |                |                |
| 1     | SELECT * FROM t5 WHERE c5_id   | 0.010753       | 0.010753       | 0.010753       | 0.010753       |
|       | IN ('S', 'S', 'S')             |                |                |                |                |
| 1     | SELECT t1.id FROM t1 JOIN      | 0.020219       | 0.020219       | 0.020219       | 0.020219       |
|       | t2 ON t2.t1_id = t1.id WHERE   |                |                |                |                |
|       | t2.t1_id = 'S' ORDER BY        |                |                |                |                |
|       | t2.t1_id                       |                |                |                |                |
| 2     | UPDATE t1 SET c1_count =       | 1.428614       | 3.504247       | 4.932861       | 2.466430       |
|       | (SELECT count(*) AS cnt FROM   |                |                |                |                |
|       | t2 WHERE c3_id = t3.id)        |                |                |                |                |
+-------+--------------------------------+----------------+----------------+----------------+----------------+
```

## print-output-options

You can see the `--output` option values.

```console
$ slp my print-output-options
count
query
min-query-time
max-query-time
sum-query-time
avg-query-time
min-lock-time
max-lock-time
sum-lock-time
avg-lock-time
min-rows-sent
max-rows-sent
sum-rows-sent
avg-rows-sent
min-rows-examined
max-rows-examined
sum-rows-examined
avg-rows-examined
min-rows-affected
max-rows-affected
sum-rows-affected
avg-rows-affected
min-bytes-sent
max-bytes-sent
sum-bytes-sent
avg-bytes-sent

$ slp pg print-output-options
count
query
min-query-time
max-query-time
sum-query-time
avg-query-time

$ slp my print-output-options --percentiles 95,99
count
query
min-query-time
max-query-time
sum-query-time
avg-query-time
p95-query-time
p99-query-time
min-lock-time
max-lock-time
sum-lock-time
avg-lock-time
p95-lock-time
p99-lock-time
min-rows-sent
max-rows-sent
sum-rows-sent
avg-rows-sent
p95-rows-sent
p99-rows-sent
min-rows-examined
max-rows-examined
sum-rows-examined
avg-rows-examined
p95-rows-examined
p99-rows-examined
min-rows-affected
max-rows-affected
sum-rows-affected
avg-rows-affected
p95-rows-affected
p99-rows-affected
min-bytes-sent
max-bytes-sent
sum-bytes-sent
avg-bytes-sent
p95-bytes-sent
p99-bytes-sent

$ slp pg print-output-options --percentiles 95,99
count
query
min-query-time
max-query-time
sum-query-time
avg-query-time
p95-query-time
p99-query-time
```

## diff

- Show the difference between the two profile results
- `+` means an increasing `count`, `rows_sent`, `rows_examined`, `rows_affected`, `bytes_sent`, and `query_time`、`lock_time` are slower
- `-` means a decreasing `count`, `rows_sent`, `rows_examined`, `rows_affected`, `bytes_sent`, and `query_time`、`lock_time` are faster

```console
$ cat /path/to/mysql.slow.log | slp my --dump dumpfile1.yaml

$ cat /path/to/mysql.slow.log | slp my --dump dumpfile2.yaml

$ slp diff dumpfile1.yaml dumpfile2.yaml --show-footers
+---------+---------------------------------+----------------+-------------------+-------------------+-------------------+
|  COUNT  |              QUERY              | MIN(QUERYTIME) |  MAX(QUERYTIME)   |  SUM(QUERYTIME)   |  AVG(QUERYTIME)   |
+---------+---------------------------------+----------------+-------------------+-------------------+-------------------+
| 1       | SELECT * FROM `t5` WHERE        | 0.010753       | 0.010753          | 0.010753          | 0.010753          |
|         | `c5_id` IN ('S','S','S')        |                |                   |                   |                   |
| 1       | DELETE FROM `t2` WHERE 'S'      | 0.369618       | 0.369618          | 0.369618          | 0.369618          |
|         | < `c1_date` OR `c2` NOT IN      |                |                   |                   |                   |
|         | (SELECT `c3` FROM `t3`)         |                |                   |                   |                   |
| 1       | DELETE FROM `t4` WHERE `c4`     | 7.148949       | 7.148949          | 7.148949          | 7.148949          |
|         | NOT IN (SELECT `c1` FROM `t1`)  |                |                   |                   |                   |
| 1       | INSERT INTO `t2`                | 0.010498       | 0.010498          | 0.010498          | 0.010498          |
|         | (`c2_id`,`c2_string`,`c2_date`) |                |                   |                   |                   |
|         | VALUES (N,'S','S')              |                |                   |                   |                   |
| 1       | INSERT INTO `t2`                | 0.010498       | 0.010498          | 0.010498          | 0.010498          |
|         | (`c2_id`,`c2_string`,`c2_date`) |                |                   |                   |                   |
|         | VALUES (N,'S','S'),(N,'S','S')  |                |                   |                   |                   |
| 1       | SELECT `t1`.`id` FROM `t1`      | 0.020219       | 0.020219          | 0.020219          | 0.020219          |
|         | JOIN `t2` ON `t2`.`t1_id` =     |                |                   |                   |                   |
|         | `t1`.`id` WHERE `t2`.`t1_id` =  |                |                   |                   |                   |
|         | 'S' ORDER BY `t2`.`t1_id`       |                |                   |                   |                   |
| 2       | UPDATE `t1` SET                 | 1.428614       | 3.504247          | 4.932861          | 2.466430          |
|         | `c1_count`=(SELECT COUNT(N) AS  |                |                   |                   |                   |
|         | `cnt` FROM `t2` WHERE `c3_id`   |                |                   |                   |                   |
|         | = `t3`.`id`)                    |                |                   |                   |                   |
| 2 (+1)  | DELETE FROM `t1` WHERE 'S' <    | 0.035678       | 1.035678 (+1.000) | 1.071356 (+1.036) | 0.535678 (+0.500) |
|         | `c1_date`                       |                |                   |                   |                   |
+---------+---------------------------------+----------------+-------------------+-------------------+-------------------+
| 10 (+1) |                                                                                                               
+---------+---------------------------------+----------------+-------------------+-------------------+-------------------+
```

## Global options

See: [Usage samples](./docs/usage_samples.md)

- `-c, --config`
    - The configuration file
    - YAML
- `--file=FILE` 
    - The access log file
- `-d, --dump=DUMP`
    - File path for creating the profile results to a file
- `-l, --load=LOAD`
    - File path to read the results of the profile created with the `-d, --dump` option
    - Can expect it to work fast if you change the `--sort` and `--reverse` options for the same profile results
- `--sort=count`
    - Output the results in sorted order
    - Sort in ascending order
    - `count`, `query`
    - `min-query-time`, `max-query-time`, `sum-query-time`, `avg-query-time`
    - `min-lock-time`, `max-lock-time`, `sum-lock-time`, `avg-lock-time`
    - `min-rows-sent`, `max-rows-sent`, `sum-rows-sent`, `avg-rows-sent`
    - `min-rows-examined`, `max-rows-examined`, `sum-rows-examined`, `avg-rows-examined`
    - `min-rows-affected`, `max-rows-affected`, `sum-rows-affected`, `avg-rows-affected`
    - `min-bytes-sent`, `max-bytes-sent`, `sum-bytes-sent`, `avg-bytes-sent`
    - The default is `count`
    - `pN(1~100)-<sort-key>` is modified by the values specified in `--percentiles`
        - The `p` means percentile
        - e.g. `p90-query-time`
        - `count` and `query` does not support
- `-r, --reverse`
    - Sort in desecending order
- `--format=table`
    - Print the profile results in a table, Markdown, TSV, CSV and HTML format
    - The default is table format
- `--noheaders`
    - Print no header when TSV and CSV format
- `--show-footers`
    - Print the total number of each 1xx ~ 5xx in the footer of the table or Markdown format
- `--limit=5000`
    - Maximum number of profile results to be printed
    - This setting is to avoid using too much memory
    - The default is 5000 lines
- `-o, --output="simple"`
    - Specify the profile results to be print, separated by commas
    - `count`, `query`, `min-query-time`, `max-query-time`, `sum-query-time`, `avg-query-time`, `min-lock-time`, `max-lock-time`, `sum-lock-time`, `avg-lock-time`, `min-rows-sent`, `max-rows-sent`, `sum-rows-sent`, `avg-rows-sent`, `min-rows-examined`, `max-rows-examined`, `sum-rows-examined`, `avg-rows-examined`, `min-rows-affected`, `max-rows-affected`, `sum-rows-affected`, `avg-rows-affected`, `min-bytes-sent`, `max-bytes-sent`, `sum-bytes-sent`, `avg-bytes-sent`
        - These outputs are the same for `all`
        - `pN(1~100)-<sort-key>` is modified by the values specified in `--percentiles`
    - The default is `simple`
    - `standard` outputs `all` without `*-rows-affected` and `*-bytes-sent`
- `-m, --matching-groups=PATTERN,...`
    - Treat Queries that match regular expressions as the same Query
    - Evaluate in the specified order. If matched, no further evaluation is performed.
- `-f, --filters=FILTERS`
    - Filters the targets for profile
    - See [Filter](#filter)
- `--pos=POSITION_FILE`
    - Stores the number of bytes to which the file has been read.
    - If the number of bytes is stored in the POSITION_FILE, the data after that number of bytes will be profiled
    - You can profile without truncating the file
        - Also, it is expected to work fast because it seeks and skips files
- `--nosave-pos`
    - Data after the number of bytes specified by `--pos` is profiled, but the number of bytes reads is not stored
- `--percentiles`
    - Specifies the percentile values to output, separated by commas
    - e.g. `90,95,99`
- `-a`, `--noabstract`
    - Do not abstract all numbers to N and strings to 'S'
- `--bundle-values`
    - Bundle VALUES of INSERT statement
    - See: [Usage samples](https://github.com/tkuchiki/slp/blob/main/docs/usage_samples.md#--bundle-values---bundle-where-in)
- `--bundle-where-in`
    - Bundle WHERE IN conditions
    - See: [Usage samples](https://github.com/tkuchiki/slp/blob/main/docs/usage_samples.md#--bundle-values---bundle-where-in)

### `pg` options

- `--log-line-prefix="%m [%p]"`
    - The `log_line_prefix` of postgresql.conf

## Filter

It is a function to include or exclude targets according to the conditions.

### Variables

Filter on the following variables:.

- `Query`
    - SQL
- `QueryTime`
    - The time to acquire queries in seconds
- `LockTime`
    - The time to acquire locks in seconds
- `RowsSent`
    - The number of rows sent to the client
- `RowsExamined`
    - The number of rows examined by the server layer
- `RowsAffected`
    - The number of rows changed
- `BytesSent`
    - The number of bytes sent to all clients

### Operators

The following operators are available:.

- `+`, `-`, `*`, `/`, `%`, `**(pow)` 
- `==`, `!=`, `<`, `>`, `<=`, `>=`
- `not`, `!`
- `and`, `&&`
- `or`, `||`
- `matches`
    - e.g.
       - `Query matches "PATTERN"`
       - `not(Query matches "PATTERN")`
- `contains`
    - e.g.
        - `Query contains "STRING"`
        - `not(Query contains "STRING")`
- `startsWith`
    - e.g.
        - `Query startsWith "PREFIX"`
        - `not(Query startsWith "PREFIX")`
- `endsWith`
    - e.g.
        - `Query endsWith "SUFFIX"`
        - `not(Query endsWith "SUFFIX")`
- `in`
    - e.g.
        - `QueryTime in [0.1, 0.2]`
        - `QueryTime not in [0.1, 0.2]`

See: https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md  

## Usage samples

See: [Usage samples](./docs/usage_samples.md)

## Donation

Donations are welcome as always!  
[:heart: Sponsor](https://github.com/sponsors/tkuchiki)
