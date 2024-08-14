package routes

import (
	"fmt"
	"github.com/cosys-io/cosys/common"
	"net/http"
	"strconv"
	"strings"
)

// getParams returns the DBParams from the query string.
func getParams(r *http.Request, attrs []common.Attribute) (common.DBParams, error) {
	pageSize, err := getPageSize(r)
	if err != nil {
		return common.DBParams{}, err
	}

	page, err := getPage(r)
	if err != nil {
		return common.DBParams{}, err
	}

	sort, err := getSort(r, attrs)
	if err != nil {
		return common.DBParams{}, err
	}

	var filter []common.Condition

	fields, err := getFields(r, attrs)
	if err != nil {
		return common.DBParams{}, err
	}

	populate, err := getPopulate(r, attrs)
	if err != nil {
		return common.DBParams{}, err
	}

	return common.NewDBParamsBuilder().
		Offset(pageSize * (int64(page) - 1)).
		Limit(pageSize).
		OrderBy(sort...).
		Where(filter...).
		Select(fields...).
		Populate(populate...).
		Build(), nil
}

// getPageSize returns the page size value from the query string.
func getPageSize(r *http.Request) (int64, error) {
	pageSizeString := r.PathValue("pageSize")
	if pageSizeString != "" {
		pageSize, err := strconv.ParseInt(pageSizeString, 10, 64)
		if err != nil {
			return 0, err
		}

		return pageSize, nil
	} else {
		return 20, nil
	}
}

// getPage returns the page number value from the query string.
func getPage(r *http.Request) (int, error) {
	pageString := r.PathValue("page")
	if pageString != "" {
		page, err := strconv.Atoi(pageString)
		if err != nil {
			return 0, err
		}

		return page, nil
	} else {
		return 1, nil
	}
}

// getSort returns the order conditions from the query string.
func getSort(r *http.Request, attrs []common.Attribute) ([]*common.Order, error) {
	sortSliceString := r.PathValue("sort")
	if sortSliceString != "" {
		sortStrings := strings.Split(sortSliceString, ",")
		sort := make([]*common.Order, len(sortStrings))

		for index, sortString := range sortStrings {
			if len(sortString) == 0 {
				return nil, fmt.Errorf("invalid sort format")
			}

			isAsc := true
			if sortString[0] == '-' {
				isAsc = false
				sortString = sortString[1:]
			}

			var sortAttr common.Attribute

			for _, attr := range attrs {
				if attr.CamelName() == sortString {
					sortAttr = attr
				}
			}

			if sortAttr == nil {
				return nil, fmt.Errorf("attribute not found: %s", sortString)
			}

			if isAsc {
				sort[index] = sortAttr.Asc()
			} else {
				sort[index] = sortAttr.Desc()
			}
		}

		return sort, nil
	} else {
		return []*common.Order{}, nil
	}
}

// getFields returns the return fields from the query string.
func getFields(r *http.Request, attrs []common.Attribute) ([]common.Attribute, error) {
	fieldSliceString := r.PathValue("fields")
	if fieldSliceString != "" {
		fieldStrings := strings.Split(fieldSliceString, ",")
		fields := make([]common.Attribute, len(fieldStrings))
		for index, fieldString := range fieldStrings {
			var fieldAttr common.Attribute

			for _, attr := range attrs {
				if attr.CamelName() == fieldString {
					fieldAttr = attr
				}
			}

			if fieldAttr == nil {
				return nil, fmt.Errorf("attribute not found: %s", fieldString)
			}

			fields[index] = fieldAttr
		}

		return fields, nil
	} else {
		return attrs, nil
	}
}

// getPopulate returns the fields to populate from the query string.
func getPopulate(r *http.Request, attrs []common.Attribute) ([]common.Attribute, error) {
	populateSliceString := r.PathValue("populate")
	if populateSliceString != "" {
		populateStrings := strings.Split(populateSliceString, ",")
		populate := make([]common.Attribute, len(populateStrings))

		for _, populateString := range populateStrings {
			var populateAttr common.Attribute

			for _, attr := range attrs {
				if attr.CamelName() == populateString {
					populateAttr = attr
				}
			}

			if populateAttr == nil {
				return nil, fmt.Errorf("attribute not found: %s", populateString)
			}

			populate = append(populate, populateAttr)
		}

		return populate, nil
	} else {
		return []common.Attribute{}, nil
	}
}

// getId returns the entity id from the query params.
func getId(r *http.Request) (int, error) {
	idString := r.PathValue("id")
	if idString == "" {
		return 0, fmt.Errorf("id not found")
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %s", idString)
	}

	return id, nil
}
