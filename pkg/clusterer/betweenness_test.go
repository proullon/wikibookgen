package clusterer

import (
	"testing"
)

func TestBetweennessShort(t *testing.T) {
	clu := NewBetweenness()
	ClusteringShort(t, clu)
}
