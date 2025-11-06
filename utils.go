package geo

func position_string(position Position) string {
	return position.String()
}

func position_slice(position Position) []float64 {
	return position.MarshalSlice()
}
