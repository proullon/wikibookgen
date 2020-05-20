package clusterer

import (
	"testing"
)

func TestClusteringV1Short(t *testing.T) {
	clu := NewV1()
	ClusteringShort(t, clu)
}
