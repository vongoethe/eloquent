package generator

var temp = `
// !!! DO NOT EDIT THIS FILE
package {{ .PackageName }} 

import (
	{{ range $i, $imp := packages }}
	"{{ $imp }}"{{ end }}
)

func init() {
{{ range $i, $m := .Models }}{{ if $m.Definition.SoftDelete }}
	// Add{{camel $m.ColumnName }}GlobalScope assign a global scope to a model for soft delete
	Add{{ camel $m.ColumnName }}GlobalScope("soft_delete", func(builder query.Condition) {
		builder.WhereNull("deleted_at")
	})
{{ end }}{{ end }}
}

{{ range $i, $m := .Models }}
// {{ camel $m.ColumnName }} is a {{ camel $m.ColumnName }} object
type {{ camel $m.ColumnName }} struct {
	original *{{ lowercase $m.ColumnName }}Original

	Id int64{{ range $j, $f := $m.Definition.Fields }}
	{{ camel $f.ColumnName }} {{ $f.Type }}{{ end }}
	{{ if not $m.Definition.WithoutCreateTime }}
	CreatedAt time.Time{{ end }}{{ if not $m.Definition.WithoutUpdateTime }}
	UpdatedAt time.Time{{ end }}{{ if $m.Definition.SoftDelete }}
	DeletedAt null.Time{{ end }}
}

// {{ lowercase $m.ColumnName }}Original is an object which stores original {{ camel $m.ColumnName }} from database
type {{ lowercase $m.ColumnName }}Original struct {
	Id int64{{ range $j, $f := $m.Definition.Fields }}
	{{ camel $f.ColumnName }} {{ $f.Type }}{{ end }}
	{{ if not $m.Definition.WithoutCreateTime }}
	CreatedAt time.Time{{ end }}{{ if not $m.Definition.WithoutUpdateTime }}
	UpdatedAt time.Time{{ end }}{{ if $m.Definition.SoftDelete }}
	DeletedAt null.Time{{ end }}
}

// Staled identify whether the object has been modified
func ({{ $m.ColumnName }} *{{ camel $m.ColumnName }}) Staled() bool {
	if {{ $m.ColumnName }}.original == nil {
		{{ $m.ColumnName }}.original = &{{ lowercase $m.ColumnName }}Original {}
	}

	if {{ $m.ColumnName }}.Id != {{ $m.ColumnName }}.original.Id {
		return true
	}

	{{ range $j, $f := $m.Definition.Fields }}
	if {{ $m.ColumnName }}.{{ camel $f.ColumnName }} != {{ $m.ColumnName }}.original.{{ camel $f.ColumnName }} {
		return true
	}{{ end }}

	{{ if not $m.Definition.WithoutCreateTime }}
	if {{ $m.ColumnName }}.CreatedAt != {{ $m.ColumnName }}.original.CreatedAt {
		return true
	}{{ end }}{{ if not $m.Definition.WithoutUpdateTime }}
	if {{ $m.ColumnName }}.UpdatedAt != {{ $m.ColumnName }}.original.UpdatedAt {
		return true
	}{{ end }}{{ if $m.Definition.SoftDelete }}
	if {{ $m.ColumnName }}.DeletedAt != {{ $m.ColumnName }}.original.DeletedAt {
		return true
	}{{ end }}

	return false
}

// StaledKV return all fields has been modified
func ({{ $m.ColumnName }} *{{ camel $m.ColumnName }}) StaledKV() query.KV {
	kv := make(query.KV, 0)

	if {{ $m.ColumnName }}.original == nil {
		{{ $m.ColumnName }}.original = &{{ lowercase $m.ColumnName }}Original {}
	}

	if {{ $m.ColumnName }}.Id != {{ $m.ColumnName }}.original.Id {
		kv["id"] = {{ $m.ColumnName }}.Id
	}

	{{ range $j, $f := $m.Definition.Fields }}
	if {{ $m.ColumnName }}.{{ camel $f.ColumnName }} != {{ $m.ColumnName }}.original.{{ camel $f.ColumnName }} {
		kv["{{ snake $f.ColumnName }}"] = {{ $m.ColumnName }}.{{ camel $f.ColumnName }}
	}{{ end }}

	{{ if not $m.Definition.WithoutCreateTime }}
	if {{ $m.ColumnName }}.CreatedAt != {{ $m.ColumnName }}.original.CreatedAt {
		kv["created_at"] = {{ $m.ColumnName }}.CreatedAt
	}{{ end }}{{ if not $m.Definition.WithoutUpdateTime }}
	if {{ $m.ColumnName }}.UpdatedAt != {{ $m.ColumnName }}.original.UpdatedAt {
		kv["updated_at"] = {{ $m.ColumnName }}.UpdatedAt
	}{{ end }}{{ if $m.Definition.SoftDelete }}
	if {{ $m.ColumnName }}.DeletedAt != {{ $m.ColumnName }}.original.DeletedAt {
		kv["deleted_at"] = {{ $m.ColumnName }}.DeletedAt
	}{{ end }}

	return kv
}

// {{ camel $m.ColumnName }}Delegate is an delegate which add some model powers to object
type {{ camel $m.ColumnName }}Delegate struct {
	delegate *{{ camel $m.ColumnName }}Model
	{{ lowercase $m.ColumnName }} *{{ camel $m.ColumnName }}
}

// Delegate create a {{ camel $m.ColumnName }} for {{ $m.ColumnName }}
func ({{ $m.ColumnName }} *{{ camel $m.ColumnName }}) Delegate(m *{{ camel $m.ColumnName }}Model) *{{ camel $m.ColumnName }}Delegate {
	return &{{ camel $m.ColumnName }}Delegate {
		delegate: m,
		{{ lowercase $m.ColumnName }}: {{ $m.ColumnName }},
	}
}

// Save create a new model or update it 
func (d *{{ camel $m.ColumnName }}Delegate) Save() error {
	id, _, err := d.delegate.SaveOrUpdate(*d.{{ lowercase $m.ColumnName }})
	if err != nil {
		return err 
	}

	d.{{ lowercase $m.ColumnName }}.Id = id
	return nil
}

// Delete remove a {{ $m.ColumnName }}
func (d *{{ camel $m.ColumnName }}Delegate) Delete() error {
	_, err := d.delegate.DeleteById(d.{{ lowercase $m.ColumnName }}.Id)
	if err != nil {
		return err 
	}

	return nil
}

type {{ lower_camel $m.ColumnName }}Wrap struct { 
	Id null.Int{{ range $j, $f := $m.Definition.Fields }}	
	{{ camel $f.ColumnName }} {{ wrap_type $f.Type }}{{ end }}
	{{ if not $m.Definition.WithoutCreateTime }}
	CreatedAt null.Time{{ end }}{{ if not $m.Definition.WithoutUpdateTime }}
	UpdatedAt null.Time{{ end }}{{ if $m.Definition.SoftDelete }}
	DeletedAt null.Time{{ end }}
}

func (w {{ lower_camel $m.ColumnName }}Wrap) To{{ camel $m.ColumnName }} () {{ camel $m.ColumnName }} {
	return {{ camel $m.ColumnName }} {
		original: &{{ lowercase $m.ColumnName }}Original {
			Id: w.Id.Int64,{{ range $j, $f := $m.Definition.Fields }}
			{{ camel $f.ColumnName }}: {{ unwrap_type $f.ColumnName $f.Type }},{{ end }}
			{{ if not $m.Definition.WithoutCreateTime }}
			CreatedAt: w.CreatedAt.Time,{{ end }}{{ if not $m.Definition.WithoutUpdateTime }}
			UpdatedAt: w.UpdatedAt.Time,{{ end }}{{ if $m.Definition.SoftDelete }}
			DeletedAt: w.DeletedAt,{{ end }}
		},
		Id: w.Id.Int64,{{ range $j, $f := $m.Definition.Fields }}
		{{ camel $f.ColumnName }}: {{ unwrap_type $f.ColumnName $f.Type }},{{ end }}
		{{ if not $m.Definition.WithoutCreateTime }}
		CreatedAt: w.CreatedAt.Time,{{ end }}{{ if not $m.Definition.WithoutUpdateTime }}
		UpdatedAt: w.UpdatedAt.Time,{{ end }}{{ if $m.Definition.SoftDelete }}
		DeletedAt: w.DeletedAt,{{ end }}
	}
}

// {{ camel $m.ColumnName }}Model is a model which encapsulates the operations of the object
type {{ camel $m.ColumnName }}Model struct {
	db *sql.DB
	tableName string

	excludeGlobalScopes []string
	includeLocalScopes []string
}

type {{ lowercase $m.ColumnName }}Scope struct {
	name  string
	apply func(builder query.Condition)
}

var {{ lowercase $m.ColumnName }}GlobalScopes = make([]{{ lowercase $m.ColumnName }}Scope, 0)
var {{ lowercase $m.ColumnName }}LocalScopes = make([]{{ lowercase $m.ColumnName }}Scope, 0)

// New{{ camel $m.ColumnName }}Model create a {{ camel $m.ColumnName }}Model
func New{{ camel $m.ColumnName }}Model (db *sql.DB) *{{ camel $m.ColumnName }}Model {
	return &{{ camel $m.ColumnName }}Model {
		db: db, 
		tableName: "{{ table $i }}",
		excludeGlobalScopes: make([]string, 0),
		includeLocalScopes: make([]string, 0),
	}
}

// Add{{ camel $m.ColumnName }}GlobalScope assign a global scope to a model
func Add{{ camel $m.ColumnName }}GlobalScope(name string, apply func(builder query.Condition)) {
	{{ lowercase $m.ColumnName }}GlobalScopes = append({{ lowercase $m.ColumnName }}GlobalScopes, {{ lowercase $m.ColumnName }}Scope{name: name, apply: apply})
}

// Add{{ camel $m.ColumnName }}LocalScope assign a local scope to a model
func Add{{ camel $m.ColumnName }}LocalScope(name string, apply func(builder query.Condition)) {
	{{ lowercase $m.ColumnName }}LocalScopes = append({{ lowercase $m.ColumnName }}LocalScopes, {{ lowercase $m.ColumnName }}Scope{name: name, apply: apply})
}

func (m *{{ camel $m.ColumnName }}Model) applyScope() query.Condition {
	scopeCond := query.ConditionBuilder()
	for _, g := range {{ lowercase $m.ColumnName }}GlobalScopes {
		if m.globalScopeEnabled(g.name) {
			g.apply(scopeCond)
		}
	}

	for _, s := range {{ lowercase $m.ColumnName }}LocalScopes {
		if m.localScopeEnabled(s.name) {
			s.apply(scopeCond)
		}
	}

	return scopeCond
}

func (m *{{ camel $m.ColumnName }}Model) localScopeEnabled(name string) bool {
	for _, n := range m.includeLocalScopes {
		if name == n {
			return true
		}
	}

	return false
}

func (m *{{ camel $m.ColumnName }}Model) globalScopeEnabled(name string) bool {
	for _, n := range m.excludeGlobalScopes {
		if name == n {
			return false
		}
	}
	
	return true
}

{{ if $m.Definition.SoftDelete }}
// WithTrashed force soft deleted models to appear in a result set
func (m *{{ camel $m.ColumnName }}Model) WithTrashed() *{{ camel $m.ColumnName }}Model {
	return m.WithoutGlobalScopes("soft_delete")
}
{{ end }}

func (m *{{ camel $m.ColumnName }}Model) clone() *{{ camel $m.ColumnName }}Model {
	return &{{ camel $m.ColumnName }}Model{
		db: m.db, 
		tableName: m.tableName,
		excludeGlobalScopes: append([]string{}, m.excludeGlobalScopes...),
		includeLocalScopes: append([]string{}, m.includeLocalScopes...),
	}
}

// WithoutGlobalScopes remove a global scope for given query
func (m *{{ camel $m.ColumnName }}Model) WithoutGlobalScopes(names ...string) *{{ camel $m.ColumnName }}Model {
	mc := m.clone()
	mc.excludeGlobalScopes = append(mc.excludeGlobalScopes, names...)

	return mc
}

// WithLocalScopes add a local scope for given query
func (m *{{ camel $m.ColumnName }}Model) WithLocalScopes(names ...string) *{{ camel $m.ColumnName }}Model {
	mc := m.clone()
	mc.includeLocalScopes = append(mc.includeLocalScopes, names...)

	return mc
}

// Find retrieve a model by its primary key
func (m *{{ camel $m.ColumnName }}Model) Find(id int64) ({{ camel $m.ColumnName }}, error) {
	return m.First(query.Builder().Where("id", "=", id))
}

// Count return model count for a given query
func (m *{{ camel $m.ColumnName }}Model) Count(builder query.SQLBuilder) (int64, error) {
	sqlStr, params := builder.Table(m.tableName).ResolveCount()
	
	rows, err := m.db.Query(sqlStr, params...)
	if err != nil {
		return 0, err
	}

	rows.Next()
	var res int64
	if err := rows.Scan(&res); err != nil {
		return 0, err
	}

	return res, nil
}

// Get retrieve all results for given query
func (m *{{ camel $m.ColumnName }}Model) Get(builder query.SQLBuilder) ([]{{ camel $m.ColumnName }}, error) {
	builder = builder.Table(m.tableName).Select("id"{{ if not $m.Definition.WithoutCreateTime }}, "created_at"{{ end }}{{ if not $m.Definition.WithoutUpdateTime }}, "updated_at"{{ end }}{{ range $j, $f := assignable_fields $m.Definition.Fields }}, "{{ snake $f.ColumnName }}"{{ end }}{{ if $m.Definition.SoftDelete }}, "deleted_at"{{ end }})
	sqlStr, params := builder.AppendCondition(m.applyScope()).ResolveQuery()
	
	rows, err := m.db.Query(sqlStr, params...)
	if err != nil {
		return nil, err
	}

	{{ lowercase $m.ColumnName }}s := make([]{{ camel $m.ColumnName }}, 0)
	for rows.Next() {
		var {{ lowercase $m.ColumnName }}Var {{ lower_camel $m.ColumnName }}Wrap
		if err := rows.Scan(&{{ lowercase $m.ColumnName }}Var.Id{{ if not $m.Definition.WithoutCreateTime }}, &{{ lowercase $m.ColumnName }}Var.CreatedAt{{ end }}{{ if not $m.Definition.WithoutUpdateTime }}, &{{ lowercase $m.ColumnName }}Var.UpdatedAt{{ end }}{{ range $j, $f := assignable_fields $m.Definition.Fields }}, &{{ lowercase $m.ColumnName }}Var.{{ camel $f.ColumnName }}{{ end }}{{ if $m.Definition.SoftDelete }}, &{{ lowercase $m.ColumnName }}Var.DeletedAt{{ end }}); err != nil {
			return nil, err
		}

		{{ lowercase $m.ColumnName }}s = append({{ lowercase $m.ColumnName }}s, {{ lowercase $m.ColumnName }}Var.To{{ camel $m.ColumnName }}())
	}

	return {{ lowercase $m.ColumnName }}s, nil
}

// First return first result for given query
func (m *{{ camel $m.ColumnName }}Model) First(builder query.SQLBuilder) ({{ camel $m.ColumnName }}, error) {
	res, err := m.Get(builder.Limit(1))
	if err != nil {
		return {{ camel $m.ColumnName }}{}, err 
	}

	if len(res) == 0 {
		return {{ camel $m.ColumnName }}{}, sql.ErrNoRows
	}

	return res[0], nil
}

// Create save a new {{ $m.ColumnName }} to database
func (m *{{ camel $m.ColumnName }}Model) Create(kv query.KV) (int64, error) {
	{{ if not $m.Definition.WithoutCreateTime }}kv["created_at"] = time.Now(){{ end }}
	{{ if not $m.Definition.WithoutUpdateTime }}kv["updated_at"] = time.Now(){{ end }}

	sqlStr, params := query.Builder().Table(m.tableName).ResolveInsert(kv)

	res, err := m.db.Exec(sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()	
}

// SaveAll save all {{ $m.ColumnName }}s to database
func (m *{{ camel $m.ColumnName }}Model) SaveAll({{ lowercase $m.ColumnName }}s []{{ camel $m.ColumnName }}) ([]int64, error) {
	ids := make([]int64, 0)
	for _, {{ lowercase $m.ColumnName }} := range {{ lowercase $m.ColumnName }}s {
		id, err := m.Save({{ lowercase $m.ColumnName }})
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// Save save a {{ $m.ColumnName }} to database
func (m *{{ camel $m.ColumnName }}Model) Save({{ lowercase $m.ColumnName }} {{ camel $m.ColumnName }}) (int64, error) {
	return m.Create(query.KV{ {{ range $j, $f := assignable_fields $m.Definition.Fields }}
		"{{ snake $f.ColumnName }}": {{ lowercase $m.ColumnName }}.{{ camel $f.ColumnName }},{{ end }}
	})	
}

// SaveOrUpdate save a new {{ $m.ColumnName }} or update it when it has a id > 0
func (m *{{ camel $m.ColumnName }}Model) SaveOrUpdate({{ lowercase $m.ColumnName }} {{ camel $m.ColumnName }}) (id int64, updated bool, err error) {
	if {{ lowercase $m.ColumnName }}.Id > 0 {
		_, _err := m.UpdateById({{ lowercase $m.ColumnName }}.Id, {{ lowercase $m.ColumnName }})
		return {{ lowercase $m.ColumnName }}.Id, true, _err
	}

	_id, _err := m.Save({{ lowercase $m.ColumnName }})
	return _id, false, _err
}

// UpdateFields update kv for a given query
func (m *{{ camel $m.ColumnName }}Model) UpdateFields(builder query.SQLBuilder, kv query.KV) (int64, error) {
	if len(kv) == 0 {
		return 0, nil
	}

	{{ if not $m.Definition.WithoutUpdateTime }}
	kv["updated_at"] = time.Now()
	{{ end }} 
	
	builder = builder.AppendCondition(m.applyScope())
	sqlStr, params := builder.Table(m.tableName).ResolveUpdate(kv)

	res, err := m.db.Exec(sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// Update update a model for given query
func (m *{{ camel $m.ColumnName }}Model) Update(builder query.SQLBuilder, {{ lowercase $m.ColumnName }} {{ camel $m.ColumnName }}) (int64, error) {
	return m.UpdateFields(builder, {{ lowercase $m.ColumnName }}.StaledKV())
}

// UpdateById update a model by id
func (m *{{ camel $m.ColumnName }}Model) UpdateById(id int64, {{ lowercase $m.ColumnName }} {{ camel $m.ColumnName }}) (int64, error) {
	return m.Update(query.Builder().Where("id", "=", id), {{ lowercase $m.ColumnName }})
}

{{ if $m.Definition.SoftDelete }}
// ForceDelete permanently remove a soft deleted model from the database
func (m *{{ camel $m.ColumnName }}Model) ForceDelete(builder query.SQLBuilder) (int64, error) {
	m2 := m.WithTrashed()

	sqlStr, params := builder.AppendCondition(m2.applyScope()).Table(m2.tableName).ResolveDelete()

	res, err := m2.db.Exec(sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// ForceDeleteById permanently remove a soft deleted model from the database by id
func (m *{{ camel $m.ColumnName }}Model) ForceDeleteById(id int64) (int64, error) {
	return m.ForceDelete(query.Builder().Where("id", "=", id))
}

// Restore restore a soft deleted model into an active state
func (m *{{ camel $m.ColumnName }}Model) Restore(builder query.SQLBuilder) (int64, error) {
	m2 := m.WithTrashed()
	return m2.UpdateFields(builder, query.KV {
		"deleted_at": nil,
	})
}

// RestoreById restore a soft deleted model into an active state by id
func (m *{{ camel $m.ColumnName }}Model) RestoreById(id int64) (int64, error) {
	return m.Restore(query.Builder().Where("id", "=", id))
}
{{ end }}

// Delete remove a model
func (m *{{ camel $m.ColumnName }}Model) Delete(builder query.SQLBuilder) (int64, error) {
	{{ if $m.Definition.SoftDelete }}
	return m.UpdateFields(builder, query.KV {
		"deleted_at": time.Now(),
	})
	{{ else }}
	sqlStr, params := builder.AppendCondition(m.applyScope()).Table(m.tableName).ResolveDelete()

	res, err := m.db.Exec(sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
	{{ end }}
}

// DeleteById remove a model by id
func (m *{{ camel $m.ColumnName }}Model) DeleteById(id int64) (int64, error) {
	return m.Delete(query.Builder().Where("id", "=", id))
}

{{ end }}

`

func GetTemplate() string {
	return temp
}
