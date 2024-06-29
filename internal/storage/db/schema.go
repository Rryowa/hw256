package db

import (
	"fmt"
	"regexp"
)

// list of regexp pattern for adding schema to the query
var schemaPrefixRegexps = [...]*regexp.Regexp{
	regexp.MustCompile(`(?i)(CREATE TABLE\s+)(\w+)(\s.*)`),
	regexp.MustCompile(`(?i)(CREATE INDEX \w+ ON\s+)(\w+)(\s.*;)`),
	regexp.MustCompile(`(?i)(UPDATE\s+)(\w+)(\s.*)`),
	regexp.MustCompile(`(?i)(INSERT INTO\s+)(\w+)(\s.*)`),
	regexp.MustCompile(`(?i)(DELETE FROM\s+)(\w+)(\s.*)`),
	regexp.MustCompile(`(?i)(SELECT\s+.*?\s+FROM\s+)(\w+)(\s.*)`),
	regexp.MustCompile(`(?i)(DROP INDEX\s+)(\w+)(\s.*)`),
	regexp.MustCompile(`(?i)(DROP TABLE\s+)(\w+)(\s.*)`),
}

func AddSchemaPrefix(schemaName, query string) string {
	prefixedQuery := query
	for _, re := range schemaPrefixRegexps {
		prefixedQuery = re.ReplaceAllString(prefixedQuery, fmt.Sprintf("${1}%s.${2}${3}", schemaName))
	}
	return prefixedQuery
}
