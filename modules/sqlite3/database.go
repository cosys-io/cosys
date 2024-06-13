package sqlite3

import (
	"database/sql"
	"fmt"

	"github.com/cosys-io/cosys/common"
)

type Database struct {
	cosys *common.Cosys
	db    *sql.DB
}

func (d Database) FindOne(uid string, params common.DBParams) (common.Entity, error) {
	model, ok := d.cosys.Models[uid]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", uid)
	}

	params = params.Limit(1)

	lifecycle := model.Lifecycle_()
	event := common.NewEvent(params)
	before, ok := lifecycle["beforeFindOne"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: beforeFindOne")
	}
	if err := before(event); err != nil {
		return nil, err
	}

	query, err := SelectQuery(&params, model)
	if err != nil {
		return nil, err
	}

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("entity not found")
	}

	entity, err := Scan(rows, &params, model)
	if err != nil {
		return nil, err
	}

	event.Result = entity
	after, ok := lifecycle["afterFindOne"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: afterFindOne")
	}
	if err := after(event); err != nil {
		return nil, err
	}

	return entity, nil
}

func (d Database) FindMany(uid string, params common.DBParams) ([]common.Entity, error) {
	model, ok := d.cosys.Models[uid]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", uid)
	}

	lifecycle := model.Lifecycle_()
	event := common.NewEvent(params)
	before, ok := lifecycle["beforeFindMany"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: beforeFindMany")
	}
	if err := before(event); err != nil {
		return nil, err
	}

	query, err := SelectQuery(&params, model)
	if err != nil {
		return nil, err
	}

	var entities []common.Entity

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		entity, err := Scan(rows, &params, model)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	event.Result = entities
	after, ok := lifecycle["afterFindMany"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: afterFindMany")
	}
	if err := after(event); err != nil {
		return nil, err
	}

	return entities, nil
}

func (d Database) Create(uid string, data common.Entity, params common.DBParams) (common.Entity, error) {
	model, ok := d.cosys.Models[uid]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", uid)
	}

	lifecycle := model.Lifecycle_()
	event := common.NewEvent(params)
	before, ok := lifecycle["beforeCreate"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: beforeCreate")
	}
	if err := before(event); err != nil {
		return nil, err
	}

	query, err := InsertQuery(&params, model)
	if err != nil {
		return nil, err
	}

	values, err := Extract(data, &params, model)
	if err != nil {
		return nil, err
	}

	_, err = d.db.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	event.Result = data
	after, ok := lifecycle["afterCreate"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: afterCreate")
	}
	if err := after(event); err != nil {
		return nil, err
	}

	// TODO: return new entity

	return data, nil
}

func (d Database) CreateMany(uid string, datas []common.Entity, params common.DBParams) ([]common.Entity, error) {
	model, ok := d.cosys.Models[uid]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", uid)
	}

	lifecycle := model.Lifecycle_()
	event := common.NewEvent(params)
	before, ok := lifecycle["beforeCreateMany"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: beforeCreateMany")
	}
	if err := before(event); err != nil {
		return nil, err
	}

	query, err := InsertQuery(&params, model)
	if err != nil {
		return nil, err
	}

	for _, data := range datas {
		values, err := Extract(data, &params, model)
		if err != nil {
			return nil, err
		}

		_, err = d.db.Exec(query, values...)
		if err != nil {
			return nil, err
		}
	}

	event.Result = datas
	after, ok := lifecycle["afterCreateMany"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: afterCreateMany")
	}
	if err := after(event); err != nil {
		return nil, err
	}

	// TODO: return new entity

	return datas, nil
}

func (d Database) Update(uid string, data common.Entity, params common.DBParams) (common.Entity, error) {
	params = params.Limit(1)

	model, ok := d.cosys.Models[uid]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", uid)
	}

	lifecycle := model.Lifecycle_()
	event := common.NewEvent(params)
	before, ok := lifecycle["beforeUpdate"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: beforeUpdate")
	}
	if err := before(event); err != nil {
		return nil, err
	}

	query, err := UpdateQuery(&params, model)
	if err != nil {
		return nil, err
	}

	values, err := Extract(data, &params, model)
	if err != nil {
		return nil, err
	}

	_, err = d.db.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	event.Result = data
	after, ok := lifecycle["afterUpdate"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: afterUpdate")
	}
	if err := after(event); err != nil {
		return nil, err
	}

	// TODO: return new entity

	return data, nil
}

func (d Database) UpdateMany(uid string, data common.Entity, params common.DBParams) ([]common.Entity, error) {
	model, ok := d.cosys.Models[uid]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", uid)
	}

	lifecycle := model.Lifecycle_()
	event := common.NewEvent(params)
	before, ok := lifecycle["beforeUpdateMany"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: beforeUpdateMany")
	}
	if err := before(event); err != nil {
		return nil, err
	}

	query, err := UpdateQuery(&params, model)
	if err != nil {
		return nil, err
	}

	values, err := Extract(data, &params, model)
	if err != nil {
		return nil, err
	}

	_, err = d.db.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	event.Result = data
	after, ok := lifecycle["afterUpdateMany"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: afterUpdateMany")
	}
	if err := after(event); err != nil {
		return nil, err
	}

	// TODO: return new entity

	return []common.Entity{data}, nil
}

func (d Database) Delete(uid string, params common.DBParams) (common.Entity, error) {
	params = params.Limit(1)

	model, ok := d.cosys.Models[uid]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", uid)
	}

	lifecycle := model.Lifecycle_()
	event := common.NewEvent(params)
	before, ok := lifecycle["beforeDelete"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: beforeDelete")
	}
	if err := before(event); err != nil {
		return nil, err
	}

	query, err := DeleteQuery(&params, model)
	if err != nil {
		return nil, err
	}

	_, err = d.db.Exec(query)
	if err != nil {
		return nil, err
	}

	// TODO: return old entity

	event.Result = nil
	after, ok := lifecycle["afterDelete"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: afterDelete")
	}
	if err := after(event); err != nil {
		return nil, err
	}

	return nil, nil
}

func (d Database) DeleteMany(uid string, params common.DBParams) ([]common.Entity, error) {
	model, ok := d.cosys.Models[uid]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", uid)
	}

	lifecycle := model.Lifecycle_()
	event := common.NewEvent(params)
	before, ok := lifecycle["beforeDeleteMany"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: beforeDeleteMany")
	}
	if err := before(event); err != nil {
		return nil, err
	}

	query, err := DeleteQuery(&params, model)
	if err != nil {
		return nil, err
	}

	_, err = d.db.Exec(query)
	if err != nil {
		return nil, err
	}

	// TODO: return old entity

	event.Result = nil
	after, ok := lifecycle["afterDeleteMany"]
	if !ok {
		return nil, fmt.Errorf("lifecycle not found: afterDeleteMany")
	}
	if err := after(event); err != nil {
		return nil, err
	}

	return nil, nil
}
