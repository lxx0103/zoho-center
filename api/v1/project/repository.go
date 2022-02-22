package project

import (
	"database/sql"
	"time"
)

type projectRepository struct {
	tx *sql.Tx
}

func NewProjectRepository(transaction *sql.Tx) ProjectRepository {
	return &projectRepository{
		tx: transaction,
	}
}

type ProjectRepository interface {
	//Project Management
	CreateProject(ProjectNew, int64) (int64, error)
	UpdateProject(int64, ProjectNew) (int64, error)
	GetProjectByID(int64, int64) (*Project, error)
	CheckNameExist(string, int64) (int, error)
}

func (r *projectRepository) CreateProject(info ProjectNew, organizationID int64) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO projects
		(
			organization_id,
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, organizationID, info.Name, info.Status, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *projectRepository) UpdateProject(id int64, info ProjectNew) (int64, error) {
	result, err := r.tx.Exec(`
		Update projects SET 
		name = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Status, time.Now(), info.User, id)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (r *projectRepository) GetProjectByID(id int64, organizationID int64) (*Project, error) {
	var res Project
	var row *sql.Row
	if organizationID != 0 {
		row = r.tx.QueryRow(`SELECT id, organization_id, name, status, created, created_by, updated, updated_by FROM projects WHERE id = ? AND organization_id = ? LIMIT 1`, id, organizationID)
	} else {
		row = r.tx.QueryRow(`SELECT id, organization_id, name, status, created, created_by, updated, updated_by FROM projects WHERE id = ? LIMIT 1`, id)
	}
	err := row.Scan(&res.ID, &res.OrganizationID, &res.Name, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *projectRepository) CheckNameExist(name string, organizationID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM projects WHERE name = ? AND organization_id = ? LIMIT 1`, name, organizationID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}
