package spot_test

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/heltonmarx/spot"
	"github.com/heltonmarx/spot/eshape"
)

func Example_buildPointShape() {
	shape := spot.NewShape(spot.WithPoint([]float64{13.4, 52.5}))

	data, err := json.Marshal(shape)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output: {"type":"point","coordinates":[13.4,52.5]}
}

func Example_buildPolygonShape() {
	shape := spot.NewShape(spot.WithPolygon([][][]float64{
		{
			{13.4, 52.5},
			{14.4, 52.5},
			{14.4, 53.5},
			{13.4, 52.5},
		},
	}))

	data, err := json.Marshal(shape)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output: {"type":"polygon","coordinates":[[[13.4,52.5],[14.4,52.5],[14.4,53.5],[13.4,52.5]]]}
}

func ExampleGeoShapeQuery() {
	shape := spot.NewShape(spot.WithPoint([]float64{13.4, 52.5}))

	query := eshape.NewGeoShapeQuery("location").
		Shape(shape).
		Relation(eshape.RelationWithin)

	body, err := query.Body()
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output: {"geo_shape":{"location":{"relation":"within","shape":{"type":"point","coordinates":[13.4,52.5]}}}}
}
