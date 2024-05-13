package queryengine

import (
	"database/sql"
	"fmt"

	"github.com/cosys-io/cosys/internal/cosys"
	"github.com/cosys-io/cosys/internal/models"
)

type QueryEngine struct {
	db      *sql.DB
	dialect string
	model   models.Model
}

func NewQueryEngine(cosys cosys.Cosys, db *sql.DB, dialect string, uid string) (*QueryEngine, error) {
	model, err := cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	return &QueryEngine{
		db,
		dialect,
		model,
	}, nil
}

func (q *QueryEngine) FindOne(params *cosys.QEParams) (models.Entity, error) {
	params.Limit(1)

	query, err := SelectQuery(q.dialect, params, q.model)
	if err != nil {
		return nil, err
	}

	rows, err := q.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("entity not found")
	}

	entity, err := Scan(rows, params, q.model)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (q *QueryEngine) FindMany(params *cosys.QEParams) ([]models.Entity, error) {
	query, err := SelectQuery(q.dialect, params, q.model)
	if err != nil {
		return nil, err
	}

	var entities []models.Entity

	rows, err := q.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		entity, err := Scan(rows, params, q.model)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

func (q *QueryEngine) Create(data models.Entity, params *cosys.QEParams) (models.Entity, error) {
	query, err := InsertQuery(q.dialect, params, q.model)
	if err != nil {
		return nil, err
	}

	values, err := Extract(data, params, q.model)
	if err != nil {
		return nil, err
	}

	_, err = q.db.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	// TODO: return new entity

	return data, nil
}

func (q *QueryEngine) Update(data models.Entity, params *cosys.QEParams) (models.Entity, error) {
	query, err := UpdateQuery(q.dialect, params, q.model)
	if err != nil {
		return nil, err
	}

	values, err := Extract(data, params, q.model)
	if err != nil {
		return nil, err
	}

	_, err = q.db.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	// TODO: return new entity

	return data, nil
}

func (q *QueryEngine) Delete(params *cosys.QEParams) (models.Entity, error) {
	query, err := DeleteQuery(q.dialect, params, q.model)
	if err != nil {
		return nil, err
	}

	_, err = q.db.Exec(query)
	if err != nil {
		return nil, err
	}

	// TODO: return old entity

	return nil, nil
}
