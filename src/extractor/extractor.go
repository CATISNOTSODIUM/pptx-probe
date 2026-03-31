package extractor

import (
	"ppt-probe/src/models"
	"slices"
)

func Filter(nodes []models.Node, filter func(models.Node) bool) []models.Node {

	result := make([]models.Node, 0)
	for _, n := range nodes {
		if filter(n) {
			result = append(result, n)
		}
		result = slices.Concat(result, Filter(n.Nodes, filter))
	}

	return result
}

func FilterWithParent(nodes []models.Node, filter func(models.Node) bool, parent *models.Node) []models.Node {

	result := make([]models.Node, 0)
	for _, n := range nodes {
		if filter(n) {
			n.Parent = parent
			result = append(result, n)
		}
		result = slices.Concat(result, FilterWithParent(n.Nodes, filter, &n))
	}

	return result
}
