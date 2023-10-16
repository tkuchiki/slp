package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/google/go-cmp/cmp"
	"github.com/tkuchiki/slp/internal/testutil"
	"github.com/tkuchiki/slp/options"
)

func Test_createOptionsFromConfig(t *testing.T) {
	viper.Reset()
	command := NewCommand("test")

	tempDir := t.TempDir()
	sort := "query"
	dummyOpts := testutil.DummyOptions(sort)

	var err error
	command.flags.config, err = testutil.CreateTempDirAndFile(tempDir, "test_create_options_from_config_config", testutil.DummyConfigFile(sort, dummyOpts))
	if err != nil {
		t.Fatal(err)
	}

	var opts *options.Options
	opts, err = command.flags.createOptionsFromConfig(command.rootCmd)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(dummyOpts, opts); diff != "" {
		t.Errorf("%s", diff)
	}
}

func Test_createOptionsFromConfig_overwrite(t *testing.T) {
	command := NewCommand("test")

	tempDir := t.TempDir()
	sort := "query"

	overwrittenSort := "max-query-time"
	overwrittenOpts := testutil.DummyOverwrittenOptions(overwrittenSort)

	var err error
	command.flags.config, err = testutil.CreateTempDirAndFile(tempDir, "test_create_options_from_config_overwrite_config", testutil.DummyConfigFile(sort, overwrittenOpts))
	if err != nil {
		t.Fatal(err)
	}

	viper.Set("file", overwrittenOpts.File)
	viper.Set("dump", overwrittenOpts.Dump)
	viper.Set("load", overwrittenOpts.Load)
	viper.Set("sort", overwrittenSort)
	viper.Set("reverse", overwrittenOpts.Reverse)
	viper.Set("format", overwrittenOpts.Format)
	viper.Set("noheaders", overwrittenOpts.NoHeaders)
	viper.Set("show_footers", overwrittenOpts.ShowFooters)
	viper.Set("limit", overwrittenOpts.Limit)
	viper.Set("matching_groups", strings.Join(overwrittenOpts.MatchingGroups, ","))
	viper.Set("filters", overwrittenOpts.Filters)
	viper.Set("pos_file", overwrittenOpts.PosFile)
	viper.Set("nosave_pos", overwrittenOpts.NoSavePos)
	viper.Set("output", overwrittenOpts.Output)
	viper.Set("percentiles", testutil.IntSliceToString(overwrittenOpts.Percentiles))
	viper.Set("pagination_limit", overwrittenOpts.PaginationLimit)
	viper.Set("bundle_where_in", overwrittenOpts.BundleWhereIn)
	viper.Set("bundle_values", overwrittenOpts.BundleValues)
	viper.Set("noabstract", overwrittenOpts.NoAbstract)

	var opts *options.Options
	opts, err = command.flags.createOptionsFromConfig(command.rootCmd)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(overwrittenOpts, opts); diff != "" {
		t.Errorf("%s", diff)
	}
}
