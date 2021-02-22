package utility

const (
	tonPrice = 350
	tax = 0.01
	oilPrice = 8
	depreciation = 800
	driver = 400
	fixedCost = 0.05
)

// GetCostFromDistance function for get dcost from distance.
func GetCostFromDistance(distance float64) float64 {
	cost := (distance*oilPrice) + (distance*(oilPrice/2)) + ((driver+depreciation)* 1)
	return cost + (cost * fixedCost) + (cost * tax)
}

// GetOfferFromWeight function for get offer from weight.
func GetOfferFromWeight(weight float64) float64 {
	autoPrice := weight * 350
	return (autoPrice * 0.1) + autoPrice
}