package cubemxtobazelinternal

import (
	"testing"
)

func TestOperandBStringParameter(t *testing.T) {
	op := operandBString{Operand: "name", Value: "hello_world"}
	expected := "  name=\"hello_world\",\n"
	got := op.Parameter().String()
	if expected != got {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
}
func TestOperandBBoolParameter(t *testing.T) {
	op := operandBBool{Operand: "linkstatic", Value: true}
	expected := "  linkstatic=True,\n"
	got := op.Parameter().String()
	if expected != got {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
}
func TestBBoolString(t *testing.T) {
	var toggle bBool = true
	expected := "True"
	got := toggle.String()
	if expected != got {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
	toggle = false
	expected = "False"
	got = toggle.String()
	if expected != got {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}

}
