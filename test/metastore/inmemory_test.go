package metastore_test

import (
	"testing"

	"com.ak.gooverlord/metastore"
)

func TestInMemoryGet(t *testing.T) {
	imb := metastore.GetInMemoryBacked()
	_, b := imb.Set("key", true)
	if !b {
		t.Fatalf("Set ops failed with b: %t", b)
	}
	_, b = imb.Set("key2", false)
	if !b {
		t.Fatalf("Set ops failed for key2, b: %t", b)
	}

	val, b := imb.Get("key")
	if val != true && b {
		t.Fatalf("Returned wrong value stored for key %t", val)
	}

	val2, b2 := imb.Get("key2")
	if val2 != false && b2 {
		t.Fatalf("Returned wrong value stored for key2 %t", val2)
	}
}
