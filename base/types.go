package base

type table struct {
	name    string
	columns []column
	labels  map[string]any
}

type column struct {
	name       string
	dataType   string
	primaryKey bool
	notNull    bool
}

func newColumn(name string, dataType string, primaryKey bool, notNull bool) column {

	return column{
		name:       name,
		dataType:   dataType,
		primaryKey: primaryKey,
		notNull:    notNull,
	}
}

func (r column) toString() string {

	ret := r.name + " " + r.dataType
	if r.primaryKey {
		ret += " PRIMARY KEY"
	}
	if r.notNull {
		ret += " NOT NULL"
	}
	return ret
}
