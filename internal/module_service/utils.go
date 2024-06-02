package module_service

import "github.com/cosys-io/cosys/internal/common"

func TransformParams(params *common.ESParams) *common.QEParams {
	return common.QEParam().
		Select(params.GetFields...).
		Insert(params.SetFields...).
		Where(params.Filters...).
		Limit(params.LimitVal).
		Offset(params.StartVal).
		OrderBy(params.Sorts...).
		Populate(params.Populates...)
}
