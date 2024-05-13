package entityservice

import "github.com/cosys-io/cosys/internal/cosys"

func TransformParams(params *cosys.ESParams) *cosys.QEParams {
	return cosys.QEParam().
		Select(params.GetFields...).
		Insert(params.SetFields...).
		Where(params.Filters...).
		Limit(params.LimitVal).
		Offset(params.StartVal).
		OrderBy(params.Sorts...).
		Populate(params.Populates...)
}
