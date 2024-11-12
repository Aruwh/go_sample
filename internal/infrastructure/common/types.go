package common

import (
	"go.mongodb.org/mongo-driver/bson"
)

type (
	// // // // // // // // // // // // // // // // // // // // // //
	// REQUESTS TYPES
	// // // // // // // // // // // // // // // // // // // // // //

	Sort struct {
		Order OrderType  `json:"order" validate:"oneof=asc desc"`
		Field SortByType `json:"field" validate:"oneof=id name"`
	}

	GetManyRequest[F any] struct {
		Filter F     `json:"filter"`
		Sort   Sort  `json:"sort"`
		Skip   int64 `json:"skip" validate:"min=0"`
		Limit  int64 `json:"limit" validate:"min=1,max=50"`
	}

	Projection struct {
		Fields []string
	}
	// // // // // // // // // // // // // // // // // // // // // //
	// RESPONSE TYPES
	// // // // // // // // // // // // // // // // // // // // // //

)

func (p *Projection) ToBson() *bson.M {
	transformed := bson.M{}

	for _, field := range p.Fields {
		transformed[field] = 1
	}

	return &transformed
}

func (s *Sort) ToBson() *bson.M {
	orderMap := make(map[OrderType]int)
	orderMap[OrderASC] = 1
	orderMap[OrderDESC] = -1

	usedField := string(s.Field)
	if usedField == "id" {
		usedField = "_id"
	}

	transformed := bson.M{}
	transformed[usedField] = orderMap[s.Order]

	return &transformed
}
