package schema

import (
	"strconv"
	"strings"
)

type Index struct {
	Name    string
	Class   string // UNIQUE | FULLTEXT | SPATIAL
	Type    string // btree, hash, gist, spgist, gin, and brin
	Where   string
	Comment string
	Fields  []IndexOption
}

type IndexOption struct {
	*Field
	Expression string
	Sort       string // DESC, ASC
	Collate    string
	Length     int
}

// ParseIndexes parse schema indexes
func (schema *Schema) ParseIndexes() map[string]Index {
	var indexes = map[string]Index{}

	for _, field := range schema.Fields {
		if field.TagSettings["INDEX"] != "" || field.TagSettings["UNIQUE_INDEX"] != "" {
			for _, index := range parseFieldIndexes(field) {
				idx := indexes[index.Name]
				idx.Name = index.Name
				if idx.Class == "" {
					idx.Class = index.Class
				}
				if idx.Type == "" {
					idx.Type = index.Type
				}
				if idx.Where == "" {
					idx.Where = index.Where
				}
				if idx.Comment == "" {
					idx.Comment = index.Comment
				}
				idx.Fields = append(idx.Fields, index.Fields...)
				indexes[index.Name] = idx
			}
		}
	}

	return indexes
}

func (schema *Schema) LookIndex(name string) *Index {
	if schema != nil {
		indexes := schema.ParseIndexes()
		for _, index := range indexes {
			if index.Name == name {
				return &index
			}

			for _, field := range index.Fields {
				if field.Name == name {
					return &index
				}
			}
		}
	}

	return nil
}

func parseFieldIndexes(field *Field) (indexes []Index) {
	for _, value := range strings.Split(field.Tag.Get("gorm"), ";") {
		if value != "" {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if k == "INDEX" || k == "UNIQUE_INDEX" {
				var (
					name      string
					tag       = strings.Join(v[1:], ":")
					idx       = strings.Index(tag, ",")
					settings  = ParseTagSetting(tag, ",")
					length, _ = strconv.Atoi(settings["LENGTH"])
				)

				if idx == -1 {
					idx = len(tag)
				}

				if idx != -1 {
					name = tag[0:idx]
				}

				if name == "" {
					name = field.Schema.namer.IndexName(field.Schema.Table, field.Name)
				}

				if (k == "UNIQUE_INDEX") || settings["UNIQUE"] != "" {
					settings["CLASS"] = "UNIQUE"
				}

				indexes = append(indexes, Index{
					Name:    name,
					Class:   settings["CLASS"],
					Type:    settings["TYPE"],
					Where:   settings["WHERE"],
					Comment: settings["COMMENT"],
					Fields: []IndexOption{{
						Field:      field,
						Expression: settings["EXPRESSION"],
						Sort:       settings["SORT"],
						Collate:    settings["COLLATE"],
						Length:     length,
					}},
				})
			}
		}
	}

	return
}
