// !!! DO NOT EDIT THIS FILE
package models

import (
	"database/sql"
	"github.com/mylxsw/eloquent/query"
	"gopkg.in/guregu/null.v3"
	"time"
)

func init() {

	// AddProjectGlobalScope assign a global scope to a model for soft delete
	AddProjectGlobalScope("soft_delete", func(builder query.Condition) {
		builder.WhereNull("deleted_at")
	})

}

// Project is a Project object
type Project struct {
	original *projectOriginal

	Id          int64
	Name        string
	Description string
	Visibility  int
	UserId      int64
	SortLevel   int
	CatalogId   null.Int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt null.Time
}

// projectOriginal is an object which stores original Project from database
type projectOriginal struct {
	Id          int64
	Name        string
	Description string
	Visibility  int
	UserId      int64
	SortLevel   int
	CatalogId   null.Int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt null.Time
}

// Staled identify whether the object has been modified
func (Project *Project) Staled() bool {
	if Project.original == nil {
		Project.original = &projectOriginal{}
	}

	if Project.Id != Project.original.Id {
		return true
	}

	if Project.Name != Project.original.Name {
		return true
	}
	if Project.Description != Project.original.Description {
		return true
	}
	if Project.Visibility != Project.original.Visibility {
		return true
	}
	if Project.UserId != Project.original.UserId {
		return true
	}
	if Project.SortLevel != Project.original.SortLevel {
		return true
	}
	if Project.CatalogId != Project.original.CatalogId {
		return true
	}

	if Project.CreatedAt != Project.original.CreatedAt {
		return true
	}
	if Project.UpdatedAt != Project.original.UpdatedAt {
		return true
	}
	if Project.DeletedAt != Project.original.DeletedAt {
		return true
	}

	return false
}

// StaledKV return all fields has been modified
func (Project *Project) StaledKV() query.KV {
	kv := make(query.KV, 0)

	if Project.original == nil {
		Project.original = &projectOriginal{}
	}

	if Project.Id != Project.original.Id {
		kv["id"] = Project.Id
	}

	if Project.Name != Project.original.Name {
		kv["name"] = Project.Name
	}
	if Project.Description != Project.original.Description {
		kv["description"] = Project.Description
	}
	if Project.Visibility != Project.original.Visibility {
		kv["visibility"] = Project.Visibility
	}
	if Project.UserId != Project.original.UserId {
		kv["user_id"] = Project.UserId
	}
	if Project.SortLevel != Project.original.SortLevel {
		kv["sort_level"] = Project.SortLevel
	}
	if Project.CatalogId != Project.original.CatalogId {
		kv["catalog_id"] = Project.CatalogId
	}

	if Project.CreatedAt != Project.original.CreatedAt {
		kv["created_at"] = Project.CreatedAt
	}
	if Project.UpdatedAt != Project.original.UpdatedAt {
		kv["updated_at"] = Project.UpdatedAt
	}
	if Project.DeletedAt != Project.original.DeletedAt {
		kv["deleted_at"] = Project.DeletedAt
	}

	return kv
}

// ProjectDelegate is an delegate which add some model powers to object
type ProjectDelegate struct {
	delegate *ProjectModel
	project  *Project
}

// Delegate create a Project for Project
func (Project *Project) Delegate(m *ProjectModel) *ProjectDelegate {
	return &ProjectDelegate{
		delegate: m,
		project:  Project,
	}
}

// Save create a new model or update it
func (d *ProjectDelegate) Save() error {
	id, _, err := d.delegate.SaveOrUpdate(*d.project)
	if err != nil {
		return err
	}

	d.project.Id = id
	return nil
}

// Delete remove a Project
func (d *ProjectDelegate) Delete() error {
	_, err := d.delegate.DeleteById(d.project.Id)
	if err != nil {
		return err
	}

	return nil
}

type projectWrap struct {
	Id          null.Int
	Name        null.String
	Description null.String
	Visibility  null.Int
	UserId      null.Int
	SortLevel   null.Int
	CatalogId   null.Int

	CreatedAt null.Time
	UpdatedAt null.Time
	DeletedAt null.Time
}

func (w projectWrap) ToProject() Project {
	return Project{
		original: &projectOriginal{
			Id:          w.Id.Int64,
			Name:        w.Name.String,
			Description: w.Description.String,
			Visibility:  int(w.Visibility.Int64),
			UserId:      w.UserId.Int64,
			SortLevel:   int(w.SortLevel.Int64),
			CatalogId:   w.CatalogId,

			CreatedAt: w.CreatedAt.Time,
			UpdatedAt: w.UpdatedAt.Time,
			DeletedAt: w.DeletedAt,
		},
		Id:          w.Id.Int64,
		Name:        w.Name.String,
		Description: w.Description.String,
		Visibility:  int(w.Visibility.Int64),
		UserId:      w.UserId.Int64,
		SortLevel:   int(w.SortLevel.Int64),
		CatalogId:   w.CatalogId,

		CreatedAt: w.CreatedAt.Time,
		UpdatedAt: w.UpdatedAt.Time,
		DeletedAt: w.DeletedAt,
	}
}

// ProjectModel is a model which encapsulates the operations of the object
type ProjectModel struct {
	db        *sql.DB
	tableName string

	excludeGlobalScopes []string
	includeLocalScopes  []string
}

type projectScope struct {
	name  string
	apply func(builder query.Condition)
}

var projectGlobalScopes = make([]projectScope, 0)
var projectLocalScopes = make([]projectScope, 0)

// NewProjectModel create a ProjectModel
func NewProjectModel(db *sql.DB) *ProjectModel {
	return &ProjectModel{
		db:                  db,
		tableName:           "wz_projects",
		excludeGlobalScopes: make([]string, 0),
		includeLocalScopes:  make([]string, 0),
	}
}

// AddProjectGlobalScope assign a global scope to a model
func AddProjectGlobalScope(name string, apply func(builder query.Condition)) {
	projectGlobalScopes = append(projectGlobalScopes, projectScope{name: name, apply: apply})
}

// AddProjectLocalScope assign a local scope to a model
func AddProjectLocalScope(name string, apply func(builder query.Condition)) {
	projectLocalScopes = append(projectLocalScopes, projectScope{name: name, apply: apply})
}

func (m *ProjectModel) applyScope() query.Condition {
	scopeCond := query.ConditionBuilder()
	for _, g := range projectGlobalScopes {
		if m.globalScopeEnabled(g.name) {
			g.apply(scopeCond)
		}
	}

	for _, s := range projectLocalScopes {
		if m.localScopeEnabled(s.name) {
			s.apply(scopeCond)
		}
	}

	return scopeCond
}

func (m *ProjectModel) localScopeEnabled(name string) bool {
	for _, n := range m.includeLocalScopes {
		if name == n {
			return true
		}
	}

	return false
}

func (m *ProjectModel) globalScopeEnabled(name string) bool {
	for _, n := range m.excludeGlobalScopes {
		if name == n {
			return false
		}
	}

	return true
}

// WithTrashed force soft deleted models to appear in a result set
func (m *ProjectModel) WithTrashed() *ProjectModel {
	return m.WithoutGlobalScopes("soft_delete")
}

func (m *ProjectModel) clone() *ProjectModel {
	return &ProjectModel{
		db:                  m.db,
		tableName:           m.tableName,
		excludeGlobalScopes: append([]string{}, m.excludeGlobalScopes...),
		includeLocalScopes:  append([]string{}, m.includeLocalScopes...),
	}
}

// WithoutGlobalScopes remove a global scope for given query
func (m *ProjectModel) WithoutGlobalScopes(names ...string) *ProjectModel {
	mc := m.clone()
	mc.excludeGlobalScopes = append(mc.excludeGlobalScopes, names...)

	return mc
}

// WithLocalScopes add a local scope for given query
func (m *ProjectModel) WithLocalScopes(names ...string) *ProjectModel {
	mc := m.clone()
	mc.includeLocalScopes = append(mc.includeLocalScopes, names...)

	return mc
}

// Find retrieve a model by its primary key
func (m *ProjectModel) Find(id int64) (Project, error) {
	return m.First(query.Builder().Where("id", "=", id))
}

// Count return model count for a given query
func (m *ProjectModel) Count(builder query.SQLBuilder) (int64, error) {
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
func (m *ProjectModel) Get(builder query.SQLBuilder) ([]Project, error) {
	builder = builder.Table(m.tableName).Select("id", "created_at", "updated_at", "name", "description", "visibility", "user_id", "sort_level", "catalog_id", "deleted_at")
	sqlStr, params := builder.AppendCondition(m.applyScope()).ResolveQuery()

	rows, err := m.db.Query(sqlStr, params...)
	if err != nil {
		return nil, err
	}

	projects := make([]Project, 0)
	for rows.Next() {
		var projectVar projectWrap
		if err := rows.Scan(&projectVar.Id, &projectVar.CreatedAt, &projectVar.UpdatedAt, &projectVar.Name, &projectVar.Description, &projectVar.Visibility, &projectVar.UserId, &projectVar.SortLevel, &projectVar.CatalogId, &projectVar.DeletedAt); err != nil {
			return nil, err
		}

		projects = append(projects, projectVar.ToProject())
	}

	return projects, nil
}

// First return first result for given query
func (m *ProjectModel) First(builder query.SQLBuilder) (Project, error) {
	res, err := m.Get(builder.Limit(1))
	if err != nil {
		return Project{}, err
	}

	if len(res) == 0 {
		return Project{}, sql.ErrNoRows
	}

	return res[0], nil
}

// Create save a new Project to database
func (m *ProjectModel) Create(kv query.KV) (int64, error) {
	kv["created_at"] = time.Now()
	kv["updated_at"] = time.Now()

	sqlStr, params := query.Builder().Table(m.tableName).ResolveInsert(kv)

	res, err := m.db.Exec(sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// SaveAll save all Projects to database
func (m *ProjectModel) SaveAll(projects []Project) ([]int64, error) {
	ids := make([]int64, 0)
	for _, project := range projects {
		id, err := m.Save(project)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// Save save a Project to database
func (m *ProjectModel) Save(project Project) (int64, error) {
	return m.Create(query.KV{
		"name":        project.Name,
		"description": project.Description,
		"visibility":  project.Visibility,
		"user_id":     project.UserId,
		"sort_level":  project.SortLevel,
		"catalog_id":  project.CatalogId,
	})
}

// SaveOrUpdate save a new Project or update it when it has a id > 0
func (m *ProjectModel) SaveOrUpdate(project Project) (id int64, updated bool, err error) {
	if project.Id > 0 {
		_, _err := m.UpdateById(project.Id, project)
		return project.Id, true, _err
	}

	_id, _err := m.Save(project)
	return _id, false, _err
}

// UpdateFields update kv for a given query
func (m *ProjectModel) UpdateFields(builder query.SQLBuilder, kv query.KV) (int64, error) {
	if len(kv) == 0 {
		return 0, nil
	}

	kv["updated_at"] = time.Now()

	builder = builder.AppendCondition(m.applyScope())
	sqlStr, params := builder.Table(m.tableName).ResolveUpdate(kv)

	res, err := m.db.Exec(sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// Update update a model for given query
func (m *ProjectModel) Update(builder query.SQLBuilder, project Project) (int64, error) {
	return m.UpdateFields(builder, project.StaledKV())
}

// UpdateById update a model by id
func (m *ProjectModel) UpdateById(id int64, project Project) (int64, error) {
	return m.Update(query.Builder().Where("id", "=", id), project)
}

// ForceDelete permanently remove a soft deleted model from the database
func (m *ProjectModel) ForceDelete(builder query.SQLBuilder) (int64, error) {
	m2 := m.WithTrashed()

	sqlStr, params := builder.AppendCondition(m2.applyScope()).Table(m2.tableName).ResolveDelete()

	res, err := m2.db.Exec(sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// ForceDeleteById permanently remove a soft deleted model from the database by id
func (m *ProjectModel) ForceDeleteById(id int64) (int64, error) {
	return m.ForceDelete(query.Builder().Where("id", "=", id))
}

// Restore restore a soft deleted model into an active state
func (m *ProjectModel) Restore(builder query.SQLBuilder) (int64, error) {
	m2 := m.WithTrashed()
	return m2.UpdateFields(builder, query.KV{
		"deleted_at": nil,
	})
}

// RestoreById restore a soft deleted model into an active state by id
func (m *ProjectModel) RestoreById(id int64) (int64, error) {
	return m.Restore(query.Builder().Where("id", "=", id))
}

// Delete remove a model
func (m *ProjectModel) Delete(builder query.SQLBuilder) (int64, error) {

	return m.UpdateFields(builder, query.KV{
		"deleted_at": time.Now(),
	})

}

// DeleteById remove a model by id
func (m *ProjectModel) DeleteById(id int64) (int64, error) {
	return m.Delete(query.Builder().Where("id", "=", id))
}
