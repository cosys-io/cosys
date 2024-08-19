package internal

import (
	"database/sql"
	"fmt"
	"github.com/cosys-io/cosys/common"
	"sync"
)

// Database is an implementation of the Database core service using SQLite3.
type Database struct {
	cosys *common.Cosys
	db    *sql.DB
}

// NewDatabase returns a new Database.
func NewDatabase(cosys *common.Cosys) *Database {
	return &Database{
		db:    new(sql.DB),
		cosys: cosys,
	}
}

// Open starts the connection to the SQLite3 database.
func (d Database) Open(dataSourceName string) error {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}

	*d.db = *db
	return nil
}

// LoadSchema loads the schema of all registered models.
func (d Database) LoadSchema() error {
	for _, model := range d.cosys.Models() {
		schema := schemaQuery(model.Schema_())
		if _, err := d.db.Exec(schema); err != nil {
			return err
		}
	}

	return nil
}

// FindOne returns one entity of the model with the given uid.
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

	query, err := selectQuery(&params, model)
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

	entity, err := scan(rows, &params, model)
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

// FindMany returns multiple entities of the model with the given uid.
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

	query, err := selectQuery(&params, model)
	if err != nil {
		return nil, err
	}

	entities := []common.Entity{}

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		entity, err := scan(rows, &params, model)
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

// Create creates one entity of the model with the given uid with the given data
// and returns the entity after creation.
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

	query, err := insertQuery(&params, model)
	if err != nil {
		return nil, err
	}

	values, err := extract(data, &params, model)
	if err != nil {
		return nil, err
	}

	rows, err := d.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("entity could not be created")
	}

	entity, err := scan(rows, &params, model)
	if err != nil {
		return nil, err
	}

	if err = model.CallLifecycle_("afterCreate", common.EventQuery{
		Params: params,
		Result: entity,
		State:  &state,
	}); err != nil {
		return nil, err
	}

	return entity, nil
}

// CreateMany creates multiple entities of the model with the given uid with the given data
// and returns the entities after creation.
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

	query, err := insertQuery(&params, model)
	if err != nil {
		return nil, err
	}

	entities := make([]common.Entity, len(datas))
	wg := sync.WaitGroup{}
	errCh := make(chan error)

	for index, data := range datas {
		wg.Add(1)
		go func(index int, data common.Entity) {
			values, err := extract(data, &params, model)
			if err != nil {
				errCh <- err
			}

			rows, err := d.db.Query(query, values...)
			if err != nil {
				errCh <- err
			}
			defer rows.Close()

			if !rows.Next() {
				errCh <- fmt.Errorf("entity could not be created")
			}

			entity, err := scan(rows, &params, model)
			if err != nil {
				errCh <- err
			}

			entities[index] = entity
			wg.Done()
		}(index, data)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	select {
	case err, ok := <-errCh:
		if ok {
			return nil, err
		}
	}

	if err = model.CallLifecycle_("afterCreateMany", common.EventQuery{
		Params: params,
		Result: entities,
		State:  &state,
	}); err != nil {
		return nil, err
	}

	return entities, nil
}

// Update updates one entity of the model with the given uid with the given data
// and returns the entity after updating.
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

	query, err := updateQuery(&params, model)
	if err != nil {
		return nil, err
	}

	values, err := extract(data, &params, model)
	if err != nil {
		return nil, err
	}

	rows, err := d.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("entity could not be updated")
	}

	entity, err := scan(rows, &params, model)
	if err != nil {
		return nil, err
	}

	if err = model.CallLifecycle_("afterUpdate", common.EventQuery{
		Params: params,
		Result: entity,
		State:  &state,
	}); err != nil {
		return nil, err
	}

	return entity, nil
}

// UpdateMany updates multiple entities of the model with the given uid with the given data
// and returns the entities after updating.
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

	query, err := updateQuery(&params, model)
	if err != nil {
		return nil, err
	}

	values, err := extract(data, &params, model)
	if err != nil {
		return nil, err
	}

	rows, err := d.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []common.Entity

	for rows.Next() {
		entity, err := scan(rows, &params, model)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	if err = model.CallLifecycle_("afterUpdateMany", common.EventQuery{
		Params: params,
		Result: entities,
		State:  &state,
	}); err != nil {
		return nil, err
	}

	return entities, nil
}

// Delete deletes one entity of the model with the given uid and returns the entity before deletion.
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

	query, err := deleteQuery(&params, model)
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

	entity, err := scan(rows, &params, model)
	if err != nil {
		return nil, err
	}

	if err = model.CallLifecycle_("afterDelete", common.EventQuery{
		Params: params,
		Result: entity,
		State:  &state,
	}); err != nil {
		return nil, err
	}

	return entity, nil
}

// DeleteMany deletes multiple entities of the model with the given uid and returns the entities before deletion.
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

	query, err := deleteQuery(&params, model)
	if err != nil {
		return nil, err
	}

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []common.Entity

	for rows.Next() {
		entity, err := scan(rows, &params, model)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	if err = model.CallLifecycle_("afterDeleteMany", common.EventQuery{
		Params: params,
		Result: entities,
		State:  &state,
	}); err != nil {
		return nil, err
	}

	return entities, nil
}
