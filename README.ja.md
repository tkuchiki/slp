# slp

slp は MySQL の slowlog 解析ツールです。

このツールは [mysqldumpslow](https://dev.mysql.com/doc/refman/8.0/en/mysqldumpslow.html) に似ていますが, より多くのメトリクスを確認することができます。

## インストール

https://github.com/tkuchiki/slp/releases から binary をダウンロードして配置してください。

# 使い方

```console
$ slp --help
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
      --bundle-values            Bundle VALUES of INSERT statement
      --bundle-where-in          Bundle WHERE IN conditions
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

$ cat example/slow.log | slp
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

- `cat /path/to/slowlog | slp` のようにパイプでデータを送るか、後述する `-f, --file` オプションでファイルを指定して解析します

## print-output-options

`--output` オプションに指定できる値を確認できます。

```console
$ slp print-output-options
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

$ slp print-output-options --percentiles 95,99
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
```

## diff

- 2つの解析結果のダンプファイルを比較します
- `+` は `count`、`rows_sent`、`rows_examined`、`rows_affected`、`bytes_sent` の増加、`query_time`、`lock_time` が遅くなったことを意味します
- `-` は `count`、`rows_sent`、`rows_examined`、`rows_affected`、`bytes_sent`の減少、`query_time`、`lock_time` が速くなったことを意味します

```console
$ cat /path/to/slow.log | slp --dump dumpfile1.yaml

$ cat /path/to/slow.log | slp --dump dumpfile2.yaml

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

## グローバルオプション

sample は [Usage samples](./docs/usage_samples.md) を参照してください。

- `-c, --config`
    - 各種オプションの設定ファイル
    - YAML
- `--file=FILE` 
    - 解析するファイルのパス
- `-d, --dump=DUMP`
    - 解析結果をファイルに書き出す際のファイルパス
- `-l, --load=LOAD`
    - `-d, --dump` オプションで書き出した解析結果を読み込む際のファイルパス
    - 同じ解析結果に対して、`--sort` や `--reverse` のオプションを変更したい場合に高速に動作することが期待できます
- `--sort=count`
    - 解析結果を表示する際にソートする条件
    - 昇順でソートする 
    - `count`, `query`
    - `min-query-time`, `max-query-time`, `sum-query-time`, `avg-query-time`
    - `min-lock-time`, `max-lock-time`, `sum-lock-time`, `avg-lock-time`
    - `min-rows-sent`, `max-rows-sent`, `sum-rows-sent`, `avg-rows-sent`
    - `min-rows-examined`, `max-rows-examined`, `sum-rows-examined`, `avg-rows-examined`
    - `min-rows-affected`, `max-rows-affected`, `sum-rows-affected`, `avg-rows-affected`
    - `min-bytes-sent`, `max-bytes-sent`, `sum-bytes-sent`, `avg-bytes-sent`
    - デフォルトは `count`
    - `pN(1~100)-<sort-key>` は `--percentiles` で指定したパーセンタイル値によって変更されます
        - e.g. `p90-query-time`
        - `count` と `query` はサポートしていません
- `-r, --reverse`
    - `--sort` オプションのソート結果を降順にします
- `--format=table`
    - 解析結果を テーブル、Markdown, TSV, CSV, HTML 形式で出力する
    - デフォルトはテーブル形式
- `--noheaders`
    - 解析結果を TSV, CSV で出力する際、header を表示しない
- `--show-footers`
    - 解析結果を テーブル, Markdown で出力する際、footer として 1xx ~ 5xx の合計数を表示する
- `--limit=5000`
    - 解析結果の表示上限数
    - 解析結果の表示数が想定より多かった場合でも、リソースを使いすぎないための設定です
    - デフォルトは 5000 行
- `-o, --output="simple"`
    - 出力する解析結果をカンマ区切りで指定する
    - `count`, `query`, `min-query-time`, `max-query-time`, `sum-query-time`, `avg-query-time`, `min-lock-time`, `max-lock-time`, `sum-lock-time`, `avg-lock-time`, `min-rows-sent`, `max-rows-sent`, `sum-rows-sent`, `avg-rows-sent`, `min-rows-examined`, `max-rows-examined`, `sum-rows-examined`, `avg-rows-examined`, `min-rows-affected`, `max-rows-affected`, `sum-rows-affected`, `avg-rows-affected`, `min-bytes-sent`, `max-bytes-sent`, `sum-bytes-sent`, `avg-bytes-sent`
        - `all` でも同様の出力を得られます
        - `pN(1~100)-<sort-key>` は `--percentiles` で指定したパーセンタイル値によって変更されます
    - デフォルトは `simple`
        - `count`, `query`, `*-query-time`
    - `standard` は `all` から `*-rows-affected` and `*-bytes-sent` を除いた出力を得られます
- `-m, --matching-groups=PATTERN,...`
    - 正規表現にマッチした Query を同じ集計対象として扱います
    - 指定した順序で正規表現を評価します。マッチした場合、それ以降の正規表現を評価しません。
- `-f, --filters=FILTERS`
    - 集計対象をフィルタします
    - 後述の[フィルタ](#フィルタ)参照
- `--pos=POSITION_FILE`
    - ファイルをどこまで読み込んだかバイト数を記録します
    - POSITION_FILE にバイト数が書かれていた場合、そのバイト数以降のデータが解析対象になります
    - ファイルを truncate することなく前回解析後からの増分だけを解析することができます
        - また、ファイルを Seek して読み飛ばすので、高速に動作することが見込めます
- `--nosave-pos`
    - `--pos` で指定したバイト数以降のデータを解析対象としますが、読み込んだバイト数の記録はしないようにします
- `--percentiles`
    - 出力するパーセンタイル値をカンマ区切りで指定します
    - e.g. `90,95,99`
- `-a`, `--noabstract`
    - すべての数値と文字列を `N` と `'S'` に置き換えないようにします
- `--bundle-values`
    - INSERT 文の VALUES の値の個数が違うクエリを一つのクエリとして集計します
    - [Usage samples](https://github.com/tkuchiki/slp/blob/main/docs/usage_samples.md#--bundle-values---bundle-where-in) を参照してください
- `--bundle-where-in`
    - WHERE IN の値の個数が違うクエリを一つのクエリとして集計します
    - [Usage samples](https://github.com/tkuchiki/slp/blob/main/docs/usage_samples.md#--bundle-values---bundle-where-in) を参照してください

## フィルタ

集計対象を条件に応じて包含、除外する機能です。

### 変数

以下の変数に対してフィルタをかけることができます。

- `Query`
    - SQL
- `QueryTime`
    - クエリの実行時間(秒)
- `LockTime`
    - ロックを取得した時間(秒)
- `RowsSent`
    - クライアントに送信された行数
- `RowsExamined`
    - サーバーレイヤーで走査された行数
- `RowsAffected`
    - 変更があった行数
- `BytesSent`
    - すべてのクライアントに送信されたバイト数

### 演算子

以下の演算子を使用できます。

- `+`, `-`, `*`, `/`, `%`, `**(べき乗)` 
- `==`, `!=`, `<`, `>`, `<=`, `>=`
- `not`, `!`
- `and`, `&&`
- `or`, `||`
- `matches`
    - 正規表現(`PATTERN`)にマッチするか否か
    - e.g.
       - `Query matches "PATTERN"`
       - `not(Query matches "PATTERN")`
- `contains`
    - 文字列(`STRING`)を含むか否か
    - e.g.
        - `Query contains "STRING"`
        - `not(Query contains "STRING")`
- `startsWith`
    - 文字列に前方一致するか否か
    - e.g.
        - `Query startsWith "PREFIX"`
        - `not(Query startsWith "PREFIX")`
- `endsWith`
    - 文字列に後方一致するか否か
    - e.g.
        - `Query endsWith "SUFFIX"`
        - `not(Query endsWith "SUFFIX")`
- `in`
    - 配列の値を含むか否か
    - e.g.
        - `QueryTime in [0.1, 0.2]`
        - `Method not in [0.1, 0.2]`

詳細は https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md を参照してください。
        
## 利用例

[Usage samples](./docs/usage_samples.md) を参照してください。

## 寄付

寄付はいつでも歓迎します！    
[:heart: Sponsor](https://github.com/sponsors/tkuchiki)
