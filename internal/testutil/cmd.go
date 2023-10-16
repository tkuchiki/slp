package testutil

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/tkuchiki/slp/options"
	"github.com/tkuchiki/slp/stats"
)

func CreateTempDirAndFile(dir, filename, content string) (string, error) {
	fpath := filepath.Join(dir, filename)
	err := os.WriteFile(fpath, []byte(content), 0644)

	return fpath, err
}

func ConfigFile() string {
	return `sort: query
reverse: true
`
}

func DummyOptions(sort string) *options.Options {
	sortOptions := stats.NewSortOptions()
	sortOptions.SetAndValidate(sort)

	return &options.Options{
		File:        "/path/to/file",
		Sort:        sortOptions.SortType(),
		Reverse:     false,
		Format:      "markdown",
		Limit:       100,
		NoHeaders:   false,
		ShowFooters: false,
		MatchingGroups: []string{
			"SELECT .+",
		},
		Filters:         "Query matches 'SELECT'",
		Output:          "count,query,min,max",
		PosFile:         "/path/to/pos",
		NoSavePos:       false,
		Percentiles:     []int{1, 5},
		PaginationLimit: 10,
		BundleWhereIn:   false,
		BundleValues:    false,
		NoAbstract:      false,
	}
}

func DummyOverwrittenOptions(sort string) *options.Options {
	sortOptions := stats.NewSortOptions()
	sortOptions.SetAndValidate(sort)

	return &options.Options{
		File:        "/path/to/overwritten/file",
		Sort:        sortOptions.SortType(),
		Reverse:     true,
		Format:      "tsv",
		Limit:       200,
		NoHeaders:   true,
		ShowFooters: true,
		MatchingGroups: []string{
			"SELECT .+",
			"INSERT .+",
		},
		Filters:         "Query matches 'SELECT'",
		Output:          "query,avg",
		PosFile:         "/path/to/overwritten/pos",
		NoSavePos:       true,
		Percentiles:     []int{5, 9},
		PaginationLimit: 20,
		BundleWhereIn:   true,
		BundleValues:    true,
		NoAbstract:      true,
	}
}

func DummyConfigFile(sort string, dummyOpts *options.Options) string {
	configTmpl := `file: {{ .File }}
sort: ` + sort + `
reverse: {{ .Reverse }}
format: {{ .Format }}
limit: {{ .Limit }}
noheaders: {{ .NoHeaders }}
show_footers: {{ .ShowFooters }}
matching_groups:
{{ range .MatchingGroups }}
  - {{ . }}
{{ end }}
filters: {{ .Filters }}
output: {{ .Output }}
pos_file: {{ .PosFile }}
nosave_pos: {{ .NoSavePos }}
percentiles:
{{ range .Percentiles }}
  - {{ . }}
{{ end }}
pagination_limit: {{ .PaginationLimit }}
bundle_where_in: {{ .BundleWhereIn }}
bundle_values: {{ .BundleValues }}
noabstract: {{ .NoAbstract }}
`
	t, err := template.New("dummy_config").Parse(configTmpl)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, dummyOpts); err != nil {
		panic(err)
	}

	return buf.String()
}
