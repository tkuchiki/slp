## Usage Samples

### Basic

```console
$ cat example/mysql.slow.log | slp my
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

### `--sort sum-query-time`

```console
$ cat example/mysql.slow.log | slp my --sort sum-query-time
+-------+---------------------------------+----------------+----------------+----------------+----------------+
| COUNT |              QUERY              | MIN(QUERYTIME) | MAX(QUERYTIME) | SUM(QUERYTIME) | AVG(QUERYTIME) |
+-------+---------------------------------+----------------+----------------+----------------+----------------+
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
| 1     | DELETE FROM `t1` WHERE 'S' <    | 0.035678       | 0.035678       | 0.035678       | 0.035678       |
|       | `c1_date`                       |                |                |                |                |
| 1     | DELETE FROM `t2` WHERE 'S'      | 0.369618       | 0.369618       | 0.369618       | 0.369618       |
|       | < `c1_date` OR `c2` NOT IN      |                |                |                |                |
|       | (SELECT `c3` FROM `t3`)         |                |                |                |                |
| 2     | UPDATE `t1` SET                 | 1.428614       | 3.504247       | 4.932861       | 2.466430       |
|       | `c1_count`=(SELECT COUNT(N) AS  |                |                |                |                |
|       | `cnt` FROM `t2` WHERE `c3_id`   |                |                |                |                |
|       | = `t3`.`id`)                    |                |                |                |                |
| 1     | DELETE FROM `t4` WHERE `c4`     | 7.148949       | 7.148949       | 7.148949       | 7.148949       |
|       | NOT IN (SELECT `c1` FROM `t1`)  |                |                |                |                |
+-------+---------------------------------+----------------+----------------+----------------+----------------+
```

### `--reverse`

```console
$ cat example/mysql.slow.log | slp my --sort sum-query-time -r
+-------+---------------------------------+----------------+----------------+----------------+----------------+
| COUNT |              QUERY              | MIN(QUERYTIME) | MAX(QUERYTIME) | SUM(QUERYTIME) | AVG(QUERYTIME) |
+-------+---------------------------------+----------------+----------------+----------------+----------------+
| 1     | DELETE FROM `t4` WHERE `c4`     | 7.148949       | 7.148949       | 7.148949       | 7.148949       |
|       | NOT IN (SELECT `c1` FROM `t1`)  |                |                |                |                |
| 2     | UPDATE `t1` SET                 | 1.428614       | 3.504247       | 4.932861       | 2.466430       |
|       | `c1_count`=(SELECT COUNT(N) AS  |                |                |                |                |
|       | `cnt` FROM `t2` WHERE `c3_id`   |                |                |                |                |
|       | = `t3`.`id`)                    |                |                |                |                |
| 1     | DELETE FROM `t2` WHERE 'S'      | 0.369618       | 0.369618       | 0.369618       | 0.369618       |
|       | < `c1_date` OR `c2` NOT IN      |                |                |                |                |
|       | (SELECT `c3` FROM `t3`)         |                |                |                |                |
| 1     | DELETE FROM `t1` WHERE 'S' <    | 0.035678       | 0.035678       | 0.035678       | 0.035678       |
|       | `c1_date`                       |                |                |                |                |
| 1     | SELECT `t1`.`id` FROM `t1`      | 0.020219       | 0.020219       | 0.020219       | 0.020219       |
|       | JOIN `t2` ON `t2`.`t1_id` =     |                |                |                |                |
|       | `t1`.`id` WHERE `t2`.`t1_id` =  |                |                |                |                |
|       | 'S' ORDER BY `t2`.`t1_id`       |                |                |                |                |
| 1     | SELECT * FROM `t5` WHERE        | 0.010753       | 0.010753       | 0.010753       | 0.010753       |
|       | `c5_id` IN ('S','S','S')        |                |                |                |                |
| 1     | INSERT INTO `t2`                | 0.010498       | 0.010498       | 0.010498       | 0.010498       |
|       | (`c2_id`,`c2_string`,`c2_date`) |                |                |                |                |
|       | VALUES (N,'S','S')              |                |                |                |                |
| 1     | INSERT INTO `t2`                | 0.010498       | 0.010498       | 0.010498       | 0.010498       |
|       | (`c2_id`,`c2_string`,`c2_date`) |                |                |                |                |
|       | VALUES (N,'S','S'),(N,'S','S')  |                |                |                |                |
+-------+---------------------------------+----------------+----------------+----------------+----------------+
```

### `--format md or markdown`

```console
$ cat example/mysql.slow.log | slp my --format md
| COUNT |              QUERY              | MIN(QUERYTIME) | MAX(QUERYTIME) | SUM(QUERYTIME) | AVG(QUERYTIME) |
|-------|---------------------------------|----------------|----------------|----------------|----------------|
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
```

### `--format tsv`

```console
$ cat example/mysql.slow.log | slp my --format tsv
Count	Query	Min(QueryTime)	Max(QueryTime)	Sum(QueryTime)	Avg(QueryTime)
1	DELETE FROM `t1` WHERE 'S' < `c1_date`	0.035678	0.035678	0.035678	0.035678
1	DELETE FROM `t2` WHERE 'S' < `c1_date` OR `c2` NOT IN (SELECT `c3` FROM `t3`)	0.369618	0.369618	0.369618	0.369618
1	DELETE FROM `t4` WHERE `c4` NOT IN (SELECT `c1` FROM `t1`)	7.148949	7.148949	7.148949	7.148949
1	INSERT INTO `t2` (`c2_id`,`c2_string`,`c2_date`) VALUES (N,'S','S')	0.010498	0.010498	0.010498	0.010498
1	INSERT INTO `t2` (`c2_id`,`c2_string`,`c2_date`) VALUES (N,'S','S'),(N,'S','S')	0.010498	0.010498	0.010498	0.010498
1	SELECT * FROM `t5` WHERE `c5_id` IN ('S','S','S')	0.010753	0.010753	0.010753	0.010753
1	SELECT `t1`.`id` FROM `t1` JOIN `t2` ON `t2`.`t1_id` = `t1`.`id` WHERE `t2`.`t1_id` = 'S' ORDER BY `t2`.`t1_id`	0.020219	0.020219	0.020219	0.020219
2	UPDATE `t1` SET `c1_count`=(SELECT COUNT(N) AS `cnt` FROM `t2` WHERE `c3_id` = `t3`.`id`)	1.428614	3.504247	4.932861	2.466430
```

### `--format csv`

```console
$ cat example/mysql.slow.log | slp my --format csv
Count,Query,Min(QueryTime),Max(QueryTime),Sum(QueryTime),Avg(QueryTime)
1,DELETE FROM `t1` WHERE 'S' < `c1_date`,0.035678,0.035678,0.035678,0.035678
1,DELETE FROM `t2` WHERE 'S' < `c1_date` OR `c2` NOT IN (SELECT `c3` FROM `t3`),0.369618,0.369618,0.369618,0.369618
1,DELETE FROM `t4` WHERE `c4` NOT IN (SELECT `c1` FROM `t1`),7.148949,7.148949,7.148949,7.148949
1,INSERT INTO `t2` (`c2_id`,`c2_string`,`c2_date`) VALUES (N,'S','S'),0.010498,0.010498,0.010498,0.010498
1,INSERT INTO `t2` (`c2_id`,`c2_string`,`c2_date`) VALUES (N,'S','S'),(N,'S','S'),0.010498,0.010498,0.010498,0.010498
1,SELECT * FROM `t5` WHERE `c5_id` IN ('S','S','S'),0.010753,0.010753,0.010753,0.010753
1,SELECT `t1`.`id` FROM `t1` JOIN `t2` ON `t2`.`t1_id` = `t1`.`id` WHERE `t2`.`t1_id` = 'S' ORDER BY `t2`.`t1_id`,0.020219,0.020219,0.020219,0.020219
2,UPDATE `t1` SET `c1_count`=(SELECT COUNT(N) AS `cnt` FROM `t2` WHERE `c3_id` = `t3`.`id`),1.428614,3.504247,4.932861,2.466430
```

### `--noheaders`

Only TSV, CSV

```console
$ cat example/mysql.slow.log | slp my --format tsv --noheaders
1,DELETE FROM `t1` WHERE 'S' < `c1_date`,0.035678,0.035678,0.035678,0.035678
1,DELETE FROM `t2` WHERE 'S' < `c1_date` OR `c2` NOT IN (SELECT `c3` FROM `t3`),0.369618,0.369618,0.369618,0.369618
1,DELETE FROM `t4` WHERE `c4` NOT IN (SELECT `c1` FROM `t1`),7.148949,7.148949,7.148949,7.148949
1,INSERT INTO `t2` (`c2_id`,`c2_string`,`c2_date`) VALUES (N,'S','S'),0.010498,0.010498,0.010498,0.010498
1,INSERT INTO `t2` (`c2_id`,`c2_string`,`c2_date`) VALUES (N,'S','S'),(N,'S','S'),0.010498,0.010498,0.010498,0.010498
1,SELECT * FROM `t5` WHERE `c5_id` IN ('S','S','S'),0.010753,0.010753,0.010753,0.010753
1,SELECT `t1`.`id` FROM `t1` JOIN `t2` ON `t2`.`t1_id` = `t1`.`id` WHERE `t2`.`t1_id` = 'S' ORDER BY `t2`.`t1_id`,0.020219,0.020219,0.020219,0.020219
2,UPDATE `t1` SET `c1_count`=(SELECT COUNT(N) AS `cnt` FROM `t2` WHERE `c3_id` = `t3`.`id`),1.428614,3.504247,4.932861,2.466430
```

### `--limit N`

```console
$ cat example/mysql.slow.log | slp my --limit 8
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

$ cat example/mysql.slow.log | slp my --limit 7
2022/07/26 09:46:34 Too many Queries (7 or less)
```

### `-o count,query,avg-query-time,p99-query-time`

```console
$ cat example/mysql.slow.log | slp my -o count,query,avg-query-time,p99-query-time --percentiles 99
+-------+---------------------------------+----------------+----------------+
| COUNT |              QUERY              | AVG(QUERYTIME) | P99(QUERYTIME) |
+-------+---------------------------------+----------------+----------------+
| 1     | DELETE FROM `t1` WHERE 'S' <    | 0.035678       | 0.035678       |
|       | `c1_date`                       |                |                |
| 1     | DELETE FROM `t2` WHERE 'S'      | 0.369618       | 0.369618       |
|       | < `c1_date` OR `c2` NOT IN      |                |                |
|       | (SELECT `c3` FROM `t3`)         |                |                |
| 1     | DELETE FROM `t4` WHERE `c4`     | 7.148949       | 7.148949       |
|       | NOT IN (SELECT `c1` FROM `t1`)  |                |                |
| 1     | INSERT INTO `t2`                | 0.010498       | 0.010498       |
|       | (`c2_id`,`c2_string`,`c2_date`) |                |                |
|       | VALUES (N,'S','S')              |                |                |
| 1     | INSERT INTO `t2`                | 0.010498       | 0.010498       |
|       | (`c2_id`,`c2_string`,`c2_date`) |                |                |
|       | VALUES (N,'S','S'),(N,'S','S')  |                |                |
| 1     | SELECT * FROM `t5` WHERE        | 0.010753       | 0.010753       |
|       | `c5_id` IN ('S','S','S')        |                |                |
| 1     | SELECT `t1`.`id` FROM `t1`      | 0.020219       | 0.020219       |
|       | JOIN `t2` ON `t2`.`t1_id` =     |                |                |
|       | `t1`.`id` WHERE `t2`.`t1_id` =  |                |                |
|       | 'S' ORDER BY `t2`.`t1_id`       |                |                |
| 2     | UPDATE `t1` SET                 | 2.466430       | 3.504247       |
|       | `c1_count`=(SELECT COUNT(N) AS  |                |                |
|       | `cnt` FROM `t2` WHERE `c3_id`   |                |                |
|       | = `t3`.`id`)                    |                |                |
+-------+---------------------------------+----------------+----------------+
```

### `--show-footers`

```console
$ cat example/mysql.slow.log | slp my --show-footers
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
| 9     |
+-------+---------------------------------+----------------+----------------+----------------+----------------+
```
 
### `--pos /tmp/slp.pos`
 
```console
$ stat -c %s example/mysql.slow.log
2395

$ slp my --file example/mysql.slow.log --pos /tmp/slp.pos
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

$ cat /tmp/slp.pos
2395

$ cat << EOS >> example/mysql.slow.log
# Time: 2022-07-20T00:31:55.988806Z
# User@Host: root[root] @ localhost [127.0.0.1]  Id:     8
# Query_time: 0.035678  Lock_time: 0.000002 Rows_sent: 0  Rows_examined: 30004
SET timestamp=1658277115;
DELETE FROM `t1` WHERE '2022-05-13 09:00:00.000' < `c1_date`
EOS

$ slp my --file example/mysql.slow.log --pos /tmp/slp.pos
+-------+--------------------------------+----------------+----------------+----------------+----------------+
| COUNT |             QUERY              | MIN(QUERYTIME) | MAX(QUERYTIME) | SUM(QUERYTIME) | AVG(QUERYTIME) |
+-------+--------------------------------+----------------+----------------+----------------+----------------+
| 1     | DELETE FROM `t1` WHERE 'S' <   | 0.035678       | 0.035678       | 0.035678       | 0.035678       |
|       | `c1_date`                      |                |                |                |                |
+-------+--------------------------------+----------------+----------------+----------------+----------------+

$ cat /tmp/slp.pos
2657

$ stat -c %s example/mysql.slow.log
2657
```

### `--nosave-pos`

```console
$ stat -c %s example/mysql.slow.log
2395

$ cat example/mysql.slow.log | slp my --pos /tmp/slp.pos
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

$ cat /tmp/slp.pos
2395

$ cat << EOS >> example/mysql.slow.log
# Time: 2022-07-20T00:31:55.988806Z
# User@Host: root[root] @ localhost [127.0.0.1]  Id:     8
# Query_time: 0.035678  Lock_time: 0.000002 Rows_sent: 0  Rows_examined: 30004
SET timestamp=1658277115;
DELETE FROM `t1` WHERE '2022-05-13 09:00:00.000' < `c1_date`
EOS

$ slp my --file example/mysql.slow.log --pos /tmp/slp.pos --nosave-pos
+-------+--------------------------------+----------------+----------------+----------------+----------------+
| COUNT |             QUERY              | MIN(QUERYTIME) | MAX(QUERYTIME) | SUM(QUERYTIME) | AVG(QUERYTIME) |
+-------+--------------------------------+----------------+----------------+----------------+----------------+
| 1     | DELETE FROM `t1` WHERE 'S' <   | 0.035678       | 0.035678       | 0.035678       | 0.035678       |
|       | `c1_date`                      |                |                |                |                |
+-------+--------------------------------+----------------+----------------+----------------+----------------+

$ cat /tmp/slp.pos
2395
```

### `--dump /tmp/slp.dump / --load /tmp/slp.dump`

```console
$ cat example/mysql.slow.log | slp my --dump /tmp/slp.dump
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

$ slp my --load /tmp/slp.dump
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
```

### `-a`, `--noabstract`

```console
$ cat example/mysql.slow.log | slp my -a
+-------+--------------------------------+----------------+----------------+----------------+----------------+
| COUNT |             QUERY              | MIN(QUERYTIME) | MAX(QUERYTIME) | SUM(QUERYTIME) | AVG(QUERYTIME) |
+-------+--------------------------------+----------------+----------------+----------------+----------------+
| 1     | DELETE FROM `t2` WHERE         | 0.369618       | 0.369618       | 0.369618       | 0.369618       |
|       | '2022-05-13 09:00:00.000'      |                |                |                |                |
|       | < `c1_date` OR `c2` NOT IN     |                |                |                |                |
|       | (SELECT `c3` FROM `t3`)        |                |                |                |                |
| 1     | DELETE FROM `t4` WHERE `c4`    | 7.148949       | 7.148949       | 7.148949       | 7.148949       |
|       | NOT IN (SELECT `c1` FROM `t1`) |                |                |                |                |
| 1     | INSERT INTO t2 (`c2_id`,       | 0.010498       | 0.010498       | 0.010498       | 0.010498       |
|       | `c2_string`, `c2_date`) VALUES |                |                |                |                |
|       | (123, 'abc', '2022-07-20       |                |                |                |                |
|       | 00:32:19.086200468')           |                |                |                |                |
| 1     | INSERT INTO t2 (`c2_id`,       | 0.010498       | 0.010498       | 0.010498       | 0.010498       |
|       | `c2_string`, `c2_date`) VALUES |                |                |                |                |
|       | (123, 'abc', '2022-07-20       |                |                |                |                |
|       | 00:32:19.086200468'),(456,     |                |                |                |                |
|       | 'def', '2022-07-21             |                |                |                |                |
|       | 00:32:19.086200468')           |                |                |                |                |
| 1     | SELECT * FROM t5 WHERE `c5_id` | 0.010753       | 0.010753       | 0.010753       | 0.010753       |
|       | IN ('id-123', 'id-456',        |                |                |                |                |
|       | 'id-789')                      |                |                |                |                |
| 1     | SELECT `t1`.`id` FROM `t1`     | 0.020219       | 0.020219       | 0.020219       | 0.020219       |
|       | JOIN `t2` ON `t2`.`t1_id` =    |                |                |                |                |
|       | `t1`.`id` WHERE `t2`.`t1_id` = |                |                |                |                |
|       | 'id-123' ORDER BY t2.t1_id     |                |                |                |                |
| 2     | UPDATE `t1` SET `c1_count`     | 1.428614       | 3.504247       | 4.932861       | 2.466430       |
|       | = (SELECT COUNT(*) AS cnt      |                |                |                |                |
|       | FROM `t2` where `c3_id` =      |                |                |                |                |
|       | `t3`.`id`)                     |                |                |                |                |
+-------+--------------------------------+----------------+----------------+----------------+----------------+
```

### `--bundle-values`, `--bundle-where-in`

```console
$ cat example/mysql.slow.bundle.log | slp my
+-------+---------------------------------+----------------+----------------+----------------+----------------+
| COUNT |              QUERY              | MIN(QUERYTIME) | MAX(QUERYTIME) | SUM(QUERYTIME) | AVG(QUERYTIME) |
+-------+---------------------------------+----------------+----------------+----------------+----------------+
| 1     | SELECT * FROM `t1` WHERE `c1`   | 1.148949       | 1.148949       | 1.148949       | 1.148949       |
|       | IN ('S','S','S')                |                |                |                |                |
| 1     | SELECT * FROM `t1` WHERE `c1`   | 2.148949       | 2.148949       | 2.148949       | 2.148949       |
|       | IN ('S','S','S','S','S','S')    |                |                |                |                |
| 1     | INSERT INTO `t2`                | 0.010498       | 0.010498       | 0.010498       | 0.010498       |
|       | (`c2_id`,`c2_string`,`c2_date`) |                |                |                |                |
|       | VALUES (N,'S','S')              |                |                |                |                |
| 1     | INSERT INTO `t2`                | 0.010498       | 0.010498       | 0.010498       | 0.010498       |
|       | (`c2_id`,`c2_string`,`c2_date`) |                |                |                |                |
|       | VALUES (N,'S','S'),(N,'S','S')  |                |                |                |                |
+-------+---------------------------------+----------------+----------------+----------------+----------------+

$ cat example/mysql.slow.bundle.log | slp my --bundle-values --bundle-where-in
+-------+---------------------------------+----------------+----------------+----------------+----------------+
| COUNT |              QUERY              | MIN(QUERYTIME) | MAX(QUERYTIME) | SUM(QUERYTIME) | AVG(QUERYTIME) |
+-------+---------------------------------+----------------+----------------+----------------+----------------+
| 2     | SELECT * FROM `t1` WHERE `c1`   | 1.148949       | 2.148949       | 3.297898       | 1.648949       |
|       | IN ('S')                        |                |                |                |                |
| 2     | INSERT INTO `t2`                | 0.010498       | 0.010498       | 0.020996       | 0.010498       |
|       | (`c2_id`,`c2_string`,`c2_date`) |                |                |                |                |
|       | VALUES ('S','S','S')            |                |                |                |                |
+-------+---------------------------------+----------------+----------------+----------------+----------------+
```
