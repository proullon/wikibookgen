package clusterer

import (
	"testing"
)

func TestClusteringV2Short(t *testing.T) {
	clu := NewV2()
	ClusteringShort(t, clu)
}

func TestClusteringV2Math(t *testing.T) {
	clu := NewV2()
	ClusteringMath(t, clu)
}
