package internal

import (
	"database/sql"
	"fmt"
	"github.com/cosys-io/cosys/common"
)

type Database struct {
	cosys *common.Cosys
	db    *sql.DB
}

func NewDatabase(db *sql.DB, cosys *common.Cosys) *Database {
	return &Database{
		db:    db,
		cosys: cosys,
	}
}

func (d Database) LoadSchema() error {
	for _, model := range d.cosys.Models() {
		schema := schemaToSQL(model.Schema_())
		if _, err := d.db.Exec(schema); err != nil {
			return err
		}
	}

	return nil
}

func (d Database) FindOne(uid string, params common.DBParams) (common.Entity, error) {
	model, err := d.cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	params.Limit = 1

	var state any
	if err = model.CallLifecycle_("beforeFindOne", common.EventQuery{
		Params: params,
		Result: nil,
		State:  &state,
	}); err != nil {
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

	if err = model.CallLifecycle_("afterFindOne", common.EventQuery{
		Params: params,
		Result: entity,
		State:  &state,
	}); err != nil {
		return nil, err
	}

	return entity, nil
}

func (d Database) FindMany(uid string, params common.DBParams) ([]common.Entity, error) {
	model, err := d.cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	var state any
	if err = model.CallLifecycle_("beforeFindMany", common.EventQuery{
		Params: params,
		Result: nil,
		State:  &state,
	}); err != nil {
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

	if err = model.CallLifecycle_("afterFindMany", common.EventQuery{
		Params: params,
		Result: entities,
		State:  &state,
	}); err != nil {
		return nil, err
	}

	return entities, nil
}

func (d Database) Create(uid string, data common.Entity, params common.DBParams) (common.Entity, error) {
	model, err := d.cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	var state any
	if err = model.CallLifecycle_("beforeCreate", common.EventQuery{
		Params: params,
		Result: nil,
		State:  &state,
	}); err != nil {
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

	if err = model.CallLifecycle_("afterCreate", common.EventQuery{
		Params: params,
		Result: data,
		State:  &state,
	}); err != nil {
		return nil, err
	}

	// TODO: return new entity

	return data, nil
}

func (d Database) CreateMany(uid string, datas []common.Entity, params common.DBParams) ([]common.Entity, error) {
	model, err := d.cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	var state any
	if err = model.CallLifecycle_("beforeCreateMany", common.EventQuery{
		Params: params,
		Result: nil,
		State:  &state,
	}); err != nil {
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

	if err = model.CallLifecycle_("afterCreateMany", common.EventQuery{
		Params: params,
		Result: datas,
		State:  &state,
	}); err != nil {
		return nil, err
	}

	// TODO: return new entity

	return datas, nil
}

func (d Database) Update(uid string, data common.Entity, params common.DBParams) (common.Entity, error) {
	params.Limit = 1

	model, err := d.cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	var state any
	if err = model.CallLifecycle_("beforeUpdate", common.EventQuery{
		Params: params,
		Result: nil,
		State:  &state,
	}); err != nil {
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

	if err = model.CallLifecycle_("afterUpdate", common.EventQuery{
		Params: params,
		Result: data,
		State:  &state,
	}); err != nil {
		return nil, err
	}

	// TODO: return new entity

	return data, nil
}

func (d Database) UpdateMany(uid string, data common.Entity, params common.DBParams) ([]common.Entity, error) {
	model, err := d.cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	var state any
	if err = model.CallLifecycle_("beforeUpdateMany", common.EventQuery{
		Params: params,
		Result: nil,
		State:  &state,
	}); err != nil {
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

	if err = model.CallLifecycle_("afterUpdateMany", common.EventQuery{
		Params: params,
		Result: data,
		State:  &state,
	}); err != nil {
		return nil, err
	}

	// TODO: return new entity

	return []common.Entity{data}, nil
}

func (d Database) Delete(uid string, params common.DBParams) (common.Entity, error) {
	params.Limit = 1

	model, err := d.cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	var state any
	if err = model.CallLifecycle_("beforeDelete", common.EventQuery{
		Params: params,
		Result: nil,
		State:  &state,
	}); err != nil {
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

	if err = model.CallLifecycle_("afterDelete", common.EventQuery{
		Params: params,
		Result: nil,
		State:  &state,
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func (d Database) DeleteMany(uid string, params common.DBParams) ([]common.Entity, error) {
	model, err := d.cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	var state any
	if err = model.CallLifecycle_("beforeDeleteMany", common.EventQuery{
		Params: params,
		Result: nil,
		State:  &state,
	}); err != nil {
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

	if err = model.CallLifecycle_("afterDeleteMany", common.EventQuery{
		Params: params,
		Result: nil,
		State:  &state,
	}); err != nil {
		return nil, err
	}

	return nil, nil
}
