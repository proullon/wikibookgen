package clusterer

import (
	"testing"
)

func TestClusteringV1Short(t *testing.T) {
	clu := NewV1()
	ClusteringShort(t, clu)
}

func TestClusteringV1Math(t *testing.T) {
	clu := NewV1()
	ClusteringMath(t, clu)
}
