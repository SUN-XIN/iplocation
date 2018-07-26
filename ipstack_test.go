package locip

import (
	"fmt"
	"testing"
)

func TestIPLocationFromIPStack(t *testing.T) {
	testIP := "193.253.222.92:8800"

	r, err := IPLocationFromIPStack(testIP)
	if err != nil {
		t.Errorf("Failed IPLocationFromIPStack: %+v", err)
	}

	fmt.Printf("get (%+v)", *r)
}
