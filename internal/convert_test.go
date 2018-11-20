package cubemxtobazelinternal

import (
	"reflect"
	"testing"
)

func TestMxComponentToCCLibrary(t *testing.T) {
	component := MxComponent{
		Class:       "Device",
		Group:       "Startup",
		Version:     "2.1.0",
		Description: "System Startup for STMicroelectronics",
		Files: []MxFile{
			MxFile{Category: "header", Name: `example.h`},
			MxFile{Category: "sourceAsm", Name: `example.s`},
			MxFile{Category: "source", Name: `example.cc`},
		},
	}
	expected := CcLibraryRule{rule{
		Operands: attributeList{
			attributeBStringList{Operand: attSrcs, Value: bStringList{"example.cc", "example.s", "example.h"}},
			attributeBStringList{Operand: attHdrs, Value: bStringList{"example.h"}},
			attributeBBool{Operand: attLinkStatic, Value: true},
		},
		comment: comment{Comment: "# System Startup for STMicroelectronics, Device:Startup, version:2.1.0"},
	}}
	got := MxComponentToCcLibraryRule(component)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
}

func TestMxComponentToCCBinary(t *testing.T) {
	component := MxComponent{
		Class:       "Device",
		Group:       "Startup",
		Version:     "2.1.0",
		Description: "System Startup for STMicroelectronics",
		Files: []MxFile{
			MxFile{Category: "header", Name: `example.h`},
			MxFile{Category: "sourceAsm", Name: `example.s`},
			MxFile{Category: "source", Name: `example.cc`},
		},
	}
	expected := ccBinaryRule{rule{
		Operands: attributeList{
			attributeBStringList{Operand: attSrcs, Value: bStringList{"example.cc", "example.s", "example.h"}},
			attributeBStringList{Operand: attHdrs, Value: bStringList{"example.h"}},
		},
		comment: comment{Comment: "# System Startup for STMicroelectronics, Device:Startup, version:2.1.0"},
	}}
	got := mxComponentToCcBinaryRule(component)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
}
