package models

// !!! DO NOT EDIT THIS FILE

import (
	"context"
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
	original     *projectOriginal
	projectModel *ProjectModel

	Id          int64 `json:"id"`
	Name        string
	Description string `json:"description"`
	Visibility  int    `json:"visibility" yaml:"visibility"`
	UserId      int64
	SortLevel   int `yaml:"sort_level"`
	CatalogId   null.Int
	Status      string `json:"status" yaml:"status"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   null.Time
}

// SetModel set model for Project
func (projectSelf *Project) SetModel(projectModel *ProjectModel) {
	projectSelf.projectModel = projectModel
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
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   null.Time
}

// Staled identify whether the object has been modified
func (projectSelf *Project) Staled() bool {
	if projectSelf.original == nil {
		projectSelf.original = &projectOriginal{}
	}

	if projectSelf.Id != projectSelf.original.Id {
		return true
	}
	if projectSelf.Name != projectSelf.original.Name {
		return true
	}
	if projectSelf.Description != projectSelf.original.Description {
		return true
	}
	if projectSelf.Visibility != projectSelf.original.Visibility {
		return true
	}
	if projectSelf.UserId != projectSelf.original.UserId {
		return true
	}
	if projectSelf.SortLevel != projectSelf.original.SortLevel {
		return true
	}
	if projectSelf.CatalogId != projectSelf.original.CatalogId {
		return true
	}
	if projectSelf.Status != projectSelf.original.Status {
		return true
	}
	if projectSelf.CreatedAt != projectSelf.original.CreatedAt {
		return true
	}
	if projectSelf.UpdatedAt != projectSelf.original.UpdatedAt {
		return true
	}
	if projectSelf.DeletedAt != projectSelf.original.DeletedAt {
		return true
	}

	return false
}

// StaledKV return all fields has been modified
func (projectSelf *Project) StaledKV() query.KV {
	kv := make(query.KV, 0)

	if projectSelf.original == nil {
		projectSelf.original = &projectOriginal{}
	}

	if projectSelf.Id != projectSelf.original.Id {
		kv["id"] = projectSelf.Id
	}
	if projectSelf.Name != projectSelf.original.Name {
		kv["name"] = projectSelf.Name
	}
	if projectSelf.Description != projectSelf.original.Description {
		kv["description"] = projectSelf.Description
	}
	if projectSelf.Visibility != projectSelf.original.Visibility {
		kv["visibility"] = projectSelf.Visibility
	}
	if projectSelf.UserId != projectSelf.original.UserId {
		kv["user_id"] = projectSelf.UserId
	}
	if projectSelf.SortLevel != projectSelf.original.SortLevel {
		kv["sort_level"] = projectSelf.SortLevel
	}
	if projectSelf.CatalogId != projectSelf.original.CatalogId {
		kv["catalog_id"] = projectSelf.CatalogId
	}
	if projectSelf.Status != projectSelf.original.Status {
		kv["status"] = projectSelf.Status
	}
	if projectSelf.CreatedAt != projectSelf.original.CreatedAt {
		kv["created_at"] = projectSelf.CreatedAt
	}
	if projectSelf.UpdatedAt != projectSelf.original.UpdatedAt {
		kv["updated_at"] = projectSelf.UpdatedAt
	}
	if projectSelf.DeletedAt != projectSelf.original.DeletedAt {
		kv["deleted_at"] = projectSelf.DeletedAt
	}

	return kv
}

func (projectSelf *Project) Pages() *PageModel {

	q := query.Builder().Where("page_id", projectSelf.Id)

	return NewPageModel(projectSelf.projectModel.GetDB()).Query(q)
}

// Save create a new model or update it
func (projectSelf *Project) Save() error {
	if projectSelf.projectModel == nil {
		return query.ErrModelNotSet
	}

	id, _, err := projectSelf.projectModel.SaveOrUpdate(*projectSelf)
	if err != nil {
		return err
	}

	projectSelf.Id = id
	return nil
}

// Delete remove a Project
func (projectSelf *Project) Delete() error {
	if projectSelf.projectModel == nil {
		return query.ErrModelNotSet
	}

	_, err := projectSelf.projectModel.DeleteById(projectSelf.Id)
	if err != nil {
		return err
	}

	return nil
}

type projectScope struct {
	name  string
	apply func(builder query.Condition)
}

var projectGlobalScopes = make([]projectScope, 0)
var projectLocalScopes = make([]projectScope, 0)

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

type projectWrap struct {
	Id          null.Int
	Name        null.String
	Description null.String
	Visibility  null.Int
	UserId      null.Int
	SortLevel   null.Int
	CatalogId   null.Int
	Status      null.String
	CreatedAt   null.Time
	UpdatedAt   null.Time
	DeletedAt   null.Time
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
			Status:      w.Status.String,
			CreatedAt:   w.CreatedAt.Time,
			UpdatedAt:   w.UpdatedAt.Time,
			DeletedAt:   w.DeletedAt,
		},

		Id:          w.Id.Int64,
		Name:        w.Name.String,
		Description: w.Description.String,
		Visibility:  int(w.Visibility.Int64),
		UserId:      w.UserId.Int64,
		SortLevel:   int(w.SortLevel.Int64),
		CatalogId:   w.CatalogId,
		Status:      w.Status.String,
		CreatedAt:   w.CreatedAt.Time,
		UpdatedAt:   w.UpdatedAt.Time,
		DeletedAt:   w.DeletedAt,
	}
}

// ProjectModel is a model which encapsulates the operations of the object
type ProjectModel struct {
	db        *query.DatabaseWrap
	tableName string

	excludeGlobalScopes []string
	includeLocalScopes  []string

	query query.SQLBuilder
}

var projectTableName = "wz_projects"

func SetProjectTable(tableName string) {
	projectTableName = tableName
}

// NewProjectModel create a ProjectModel
func NewProjectModel(db query.Database) *ProjectModel {
	return &ProjectModel{
		db:                  query.NewDatabaseWrap(db),
		tableName:           projectTableName,
		excludeGlobalScopes: make([]string, 0),
		includeLocalScopes:  make([]string, 0),
		query:               query.Builder(),
	}
}

// GetDB return database instance
func (m *ProjectModel) GetDB() query.Database {
	return m.db.GetDB()
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
		query:               m.query,
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

// Query add query builder to model
func (m *ProjectModel) Query(builder query.SQLBuilder) *ProjectModel {
	mm := m.clone()
	mm.query = mm.query.Merge(builder)

	return mm
}

// Find retrieve a model by its primary key
func (m *ProjectModel) Find(id int64) (Project, error) {
	return m.First(m.query.Where("id", "=", id))
}

// Exists return whether the records exists for a given query
func (m *ProjectModel) Exists(builders ...query.SQLBuilder) (bool, error) {
	count, err := m.Count(builders...)
	return count > 0, err
}

// Count return model count for a given query
func (m *ProjectModel) Count(builders ...query.SQLBuilder) (int64, error) {
	sqlStr, params := m.query.Merge(builders...).Table(m.tableName).ResolveCount()

	rows, err := m.db.QueryContext(context.Background(), sqlStr, params...)
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

func (m *ProjectModel) Paginate(page int64, perPage int64, builders ...query.SQLBuilder) ([]Project, query.PaginateMeta, error) {
	if page <= 0 {
		page = 1
	}

	if perPage <= 0 {
		perPage = 15
	}

	meta := query.PaginateMeta{
		PerPage: perPage,
		Page:    page,
	}

	count, err := m.Count(builders...)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = count
	meta.LastPage = count / perPage
	if count%perPage != 0 {
		meta.LastPage += 1
	}

	res, err := m.Get(append([]query.SQLBuilder{query.Builder().Limit(perPage).Offset((page - 1) * perPage)}, builders...)...)
	if err != nil {
		return res, meta, err
	}

	return res, meta, nil
}

// Get retrieve all results for given query
func (m *ProjectModel) Get(builders ...query.SQLBuilder) ([]Project, error) {
	sqlStr, params := m.query.Merge(builders...).
		Table(m.tableName).
		Select(
			"id",
			"name",
			"description",
			"visibility",
			"user_id",
			"sort_level",
			"catalog_id",
			"status",
			"created_at",
			"updated_at",
			"deleted_at",
		).AppendCondition(m.applyScope()).
		ResolveQuery()

	rows, err := m.db.QueryContext(context.Background(), sqlStr, params...)
	if err != nil {
		return nil, err
	}

	projects := make([]Project, 0)
	for rows.Next() {
		var projectVar projectWrap
		if err := rows.Scan(
			&projectVar.Id,
			&projectVar.Name,
			&projectVar.Description,
			&projectVar.Visibility,
			&projectVar.UserId,
			&projectVar.SortLevel,
			&projectVar.CatalogId,
			&projectVar.Status,
			&projectVar.CreatedAt,
			&projectVar.UpdatedAt,
			&projectVar.DeletedAt); err != nil {
			return nil, err
		}

		projectReal := projectVar.ToProject()
		projectReal.SetModel(m)
		projects = append(projects, projectReal)
	}

	return projects, nil
}

// First return first result for given query
func (m *ProjectModel) First(builders ...query.SQLBuilder) (Project, error) {
	res, err := m.Get(append(builders, query.Builder().Limit(1))...)
	if err != nil {
		return Project{}, err
	}

	if len(res) == 0 {
		return Project{}, query.ErrNoResult
	}

	return res[0], nil
}

// Create save a new Project to database
func (m *ProjectModel) Create(kv query.KV) (int64, error) {
	kv["created_at"] = time.Now()
	kv["updated_at"] = time.Now()

	sqlStr, params := m.query.Table(m.tableName).ResolveInsert(kv)

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
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
		"status":      project.Status,
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
func (m *ProjectModel) UpdateFields(kv query.KV, builders ...query.SQLBuilder) (int64, error) {
	if len(kv) == 0 {
		return 0, nil
	}

	kv["updated_at"] = time.Now()

	sqlStr, params := m.query.Merge(builders...).AppendCondition(m.applyScope()).
		Table(m.tableName).
		ResolveUpdate(kv)

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// Update update a model for given query
func (m *ProjectModel) Update(project Project) (int64, error) {
	return m.UpdateFields(project.StaledKV())
}

// UpdateById update a model by id
func (m *ProjectModel) UpdateById(id int64, project Project) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Update(project)
}

// ForceDelete permanently remove a soft deleted model from the database
func (m *ProjectModel) ForceDelete(builders ...query.SQLBuilder) (int64, error) {
	m2 := m.WithTrashed()

	sqlStr, params := m2.query.Merge(builders...).AppendCondition(m2.applyScope()).Table(m2.tableName).ResolveDelete()

	res, err := m2.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// ForceDeleteById permanently remove a soft deleted model from the database by id
func (m *ProjectModel) ForceDeleteById(id int64) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).ForceDelete()
}

// Restore restore a soft deleted model into an active state
func (m *ProjectModel) Restore(builders ...query.SQLBuilder) (int64, error) {
	m2 := m.WithTrashed()
	return m2.UpdateFields(query.KV{
		"deleted_at": nil,
	}, builders...)
}

// RestoreById restore a soft deleted model into an active state by id
func (m *ProjectModel) RestoreById(id int64) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Restore()
}

// Delete remove a model
func (m *ProjectModel) Delete(builders ...query.SQLBuilder) (int64, error) {

	return m.UpdateFields(query.KV{
		"deleted_at": time.Now(),
	}, builders...)

}

// DeleteById remove a model by id
func (m *ProjectModel) DeleteById(id int64) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Delete()
}
