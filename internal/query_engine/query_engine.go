package query_engine

import (
	"database/sql"
	"fmt"

	"github.com/cosys-io/cosys/internal/common"
	"github.com/cosys-io/cosys/internal/models"
)

type QueryEngine struct {
	cosys   common.Cosys
	db      *sql.DB
	dialect string
	uid     string
	model   models.Model
}

func NewQueryEngine(cosys common.Cosys, db *sql.DB, dialect string, uid string) (*QueryEngine, error) {
	model, err := cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	return &QueryEngine{
		cosys,
		db,
		dialect,
		uid,
		model,
	}, nil
}

func (q *QueryEngine) FindOne(params *common.QEParams) (models.Entity, error) {
	params.Limit(1)

	lifecycles, err := q.cosys.Lifecycle(q.uid)
	if err != nil {
		return nil, err
	}
	event := common.NewEvent(params)
	if err = lifecycles.Act("beforeFindOne", event); err != nil {
		return nil, err
	}

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

	event.Result = entity
	if err = lifecycles.Act("afterFindOne", event); err != nil {
		return nil, err
	}

	return entity, nil
}

func (q *QueryEngine) FindMany(params *common.QEParams) ([]models.Entity, error) {
	lifecycles, err := q.cosys.Lifecycle(q.uid)
	if err != nil {
		return nil, err
	}
	event := common.NewEvent(params)
	if err = lifecycles.Act("beforeFindMany", event); err != nil {
		return nil, err
	}

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

	event.Result = entities
	if err = lifecycles.Act("afterFindMany", event); err != nil {
		return nil, err
	}

	return entities, nil
}

func (q *QueryEngine) Create(data models.Entity, params *common.QEParams) (models.Entity, error) {
	lifecycles, err := q.cosys.Lifecycle(q.uid)
	if err != nil {
		return nil, err
	}
	event := common.NewEvent(params)
	if err = lifecycles.Act("beforeCreate", event); err != nil {
		return nil, err
	}

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

	event.Result = data
	if err = lifecycles.Act("afterCreate", event); err != nil {
		return nil, err
	}

	// TODO: return new entity

	return data, nil
}

func (q *QueryEngine) CreateMany(datas []models.Entity, params *common.QEParams) ([]models.Entity, error) {
	lifecycles, err := q.cosys.Lifecycle(q.uid)
	if err != nil {
		return nil, err
	}
	event := common.NewEvent(params)
	if err = lifecycles.Act("beforeCreateMany", event); err != nil {
		return nil, err
	}

	query, err := InsertQuery(q.dialect, params, q.model)
	if err != nil {
		return nil, err
	}

	for _, data := range datas {
		values, err := Extract(data, params, q.model)
		if err != nil {
			return nil, err
		}

		_, err = q.db.Exec(query, values...)
		if err != nil {
			return nil, err
		}
	}

	event.Result = datas
	if err = lifecycles.Act("afterCreateMany", event); err != nil {
		return nil, err
	}

	// TODO: return new entity

	return datas, nil
}

func (q *QueryEngine) Update(data models.Entity, params *common.QEParams) (models.Entity, error) {
	params.Limit(1)

	lifecycles, err := q.cosys.Lifecycle(q.uid)
	if err != nil {
		return nil, err
	}
	event := common.NewEvent(params)
	if err = lifecycles.Act("beforeUpdate", event); err != nil {
		return nil, err
	}

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

	event.Result = data
	if err = lifecycles.Act("afterUpdate", event); err != nil {
		return nil, err
	}

	// TODO: return new entity

	return data, nil
}

func (q *QueryEngine) UpdateMany(data models.Entity, params *common.QEParams) (models.Entity, error) {
	lifecycles, err := q.cosys.Lifecycle(q.uid)
	if err != nil {
		return nil, err
	}
	event := common.NewEvent(params)
	if err = lifecycles.Act("beforeUpdateMany", event); err != nil {
		return nil, err
	}

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

	event.Result = data
	if err = lifecycles.Act("afterUpdateMany", event); err != nil {
		return nil, err
	}

	// TODO: return new entity

	return data, nil
}

func (q *QueryEngine) Delete(params *common.QEParams) (models.Entity, error) {
	params.Limit(1)

	lifecycles, err := q.cosys.Lifecycle(q.uid)
	if err != nil {
		return nil, err
	}
	event := common.NewEvent(params)
	if err = lifecycles.Act("beforeDelete", event); err != nil {
		return nil, err
	}

	query, err := DeleteQuery(q.dialect, params, q.model)
	if err != nil {
		return nil, err
	}

	_, err = q.db.Exec(query)
	if err != nil {
		return nil, err
	}

	// TODO: return old entity

	event.Result = nil
	if err = lifecycles.Act("afterDelete", event); err != nil {
		return nil, err
	}

	return nil, nil
}

func (q *QueryEngine) DeleteMany(params *common.QEParams) (models.Entity, error) {
	lifecycles, err := q.cosys.Lifecycle(q.uid)
	if err != nil {
		return nil, err
	}
	event := common.NewEvent(params)
	if err = lifecycles.Act("beforeDeleteMany", event); err != nil {
		return nil, err
	}

	query, err := DeleteQuery(q.dialect, params, q.model)
	if err != nil {
		return nil, err
	}

	_, err = q.db.Exec(query)
	if err != nil {
		return nil, err
	}

	// TODO: return old entity

	event.Result = nil
	if err = lifecycles.Act("afterDeleteMany", event); err != nil {
		return nil, err
	}

	return nil, nil
}
