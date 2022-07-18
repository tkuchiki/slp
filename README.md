# slp

slp is a MySQL SlowLog Profiler.

## Installation

Download from https://github.com/tkuchiki/slp/releases

## Usage

```console
$ slp --help
$ ./slp --help
slp is a (MySQL) SlowLog Profiler

Usage:
  slp [flags]
  slp [command]

Available Commands:
  completion           Generate the autocompletion script for the specified shell
  diff                 Show the difference between the two profile results
  help                 Help about any command
  print-output-options Print --output options

Flags:
      --assemble-values          Assemble VALUES of INSERT statement
      --assemble-where-in        Assemble WHERE IN statement
      --config string            The configuration file
      --dump string              Dump profiled data as YAML
      --file string              The slowlog file
  -f, --filters string           Only the logs are profiled that match the conditions
      --format string            The output format (table, markdown, tsv, csv and html) (default "table")
  -h, --help                     help for slp
      --limit int                The maximum number of results to display (default 5000)
      --load string              Load the profiled YAML data
  -m, --matching-groups string   Specifies Query matching groups separated by commas
  -a, --noabstract               Do not abstract all numbers to N and strings to 'S'
      --noheaders                Output no header line at all (only --format=tsv, csv)
      --nosave-pos               Do not save position file
  -o, --output string            Specifies the results to display, separated by commas (default "simple")
      --page int                 Number of pages of pagination (default 100)
      --percentiles string       Specifies the percentiles separated by commas
      --pos string               The position file
  -r, --reverse                  Sort results in reverse order
      --show-footers             Output footer line at all (only --format=table, markdown)
      --sort string              Output the results in sorted order (default "count")
  -v, --version                  version for slp

Use "slp [command] --help" for more information about a command.
```

## Example

```console
$ cat example/slow.log | slp
+-------+---------------------------------+----------------+----------------+----------------+----------------+
| COUNT |              QUERY              | MIN(QUERYTIME) | MAX(QUERYTIME) | SUM(QUERYTIME) | AVG(QUERYTIME) |
+-------+---------------------------------+----------------+----------------+----------------+----------------+
| 1     | DELETE FROM `t1` WHERE 'S' <    | 0.035678       | 0.035678       | 0.035678       | 0.035678       |
|       | `c1_date`                       |                |                |                |                |
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
```
