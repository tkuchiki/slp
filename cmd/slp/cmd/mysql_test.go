package cmd

import "testing"

func Test_newMySQLCmd(t *testing.T) {
	args := []string{
		"my",
		"--file", "../../../example/mysql.slow.log",
		"--format", "csv",
		"--output", "count,query,min-query-time,max-query-time,avg-lock-time,avg-rows-examined",
	}
	want := "Count,Query,Min(QueryTime),Max(QueryTime),Avg(LockTime),Avg(RowsExamined)\n" +
		"1,DELETE FROM `t2` WHERE 'S' < `c1_date` OR `c2` NOT IN (SELECT `c3` FROM `t3`),0.369618,0.369618,0.000002,258959.000000\n" +
		"1,DELETE FROM `t4` WHERE `c4` NOT IN (SELECT `c1` FROM `t1`),7.148949,7.148949,0.000004,20720537.000000\n" +
		"1,INSERT INTO `t2` (`c2_id`,`c2_string`,`c2_date`) VALUES (N,'S','S'),0.010498,0.010498,0.000000,0.000000\n" +
		"1,INSERT INTO `t2` (`c2_id`,`c2_string`,`c2_date`) VALUES (N,'S','S'),(N,'S','S'),0.010498,0.010498,0.000000,0.000000\n" +
		"1,SELECT * FROM `t5` WHERE `c5_id` IN ('S','S','S'),0.010753,0.010753,0.000001,67.000000\n" +
		"1,SELECT `t1`.`id` FROM `t1` JOIN `t2` ON `t2`.`t1_id` = `t1`.`id` WHERE `t2`.`t1_id` = 'S' ORDER BY `t2`.`t1_id`,0.020219,0.020219,0.000001,174.000000\n" +
		"2,UPDATE `t1` SET `c1_count`=(SELECT COUNT(N) AS `cnt` FROM `t2` WHERE `c3_id` = `t3`.`id`),1.428614,3.504247,0.000002,6811337.500000\n"

	if got := executeForOutput(t, args); got != want {
		t.Fatalf("output:\n%s\nwant:\n%s", got, want)
	}
}
