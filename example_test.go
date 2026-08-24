package spot_test

import (
	"encoding/json"
	"fmt"

	"github.com/heltonmarx/spot"
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

func Example_geoShapeQueryDSL() {
	shape := spot.NewShape(spot.WithPoint([]float64{13.4, 52.5}))

	shapeBytes, err := json.Marshal(shape)
	if err != nil {
		panic(err)
	}

	// shapeBytes is a json.RawMessage-compatible value, ready for
	// types.GeoShapeFieldQuery{Shape: shapeBytes} from go-elasticsearch/v8.
	fmt.Println(string(shapeBytes))
	// Output: {"type":"point","coordinates":[13.4,52.5]}
}
