package cubemxtobazelinternal

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

func TestMxProjectToCCLibrary(t *testing.T) {
	proj := projectImpl{components: MxComponents{
		Components: []MxComponent{
			MxComponent{
				Class:       "Device",
				Group:       "Startup",
				Version:     "2.1.0",
				Description: "System Startup for STMicroelectronics",
				Files: []MxFile{
					MxFile{Category: "header", Name: `example.h`},
					MxFile{Category: "sourceAsm", Name: `example.s`},
					MxFile{Category: "source", Name: `example.cc`},
				}}},
	},
	}
	expected := []CcLibraryRule{CcLibraryRule{rule{
		Operands: attributeList{
			attributeBString{Operand: "name", Value: bString(ccLibraryTargetName(proj.Components()[0]))},
			attributeBStringList{Operand: attSrcs, Value: bStringList{"example.cc", "example.s", "example.h"}},
			attributeBStringList{Operand: attHdrs, Value: bStringList{"example.h"}},
			attributeBBool{Operand: attLinkStatic, Value: true},
		},
		comment: comment{Comment: "# System Startup for STMicroelectronics  Device:Startup:, version:2.1.0"},
	}}}
	got := MxProjectToCcLibraryRules(proj)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n Diff:\n", expected, got)
		if diff := deep.Equal(got, expected); diff != nil {
			t.Error(diff)
		}
	}
}

// TODO: Fix this test to use project files instead of components
func TestMxProjectToCCBinary(t *testing.T) {
	project := projectImpl{generator: MxGenerator{
		ProjectFiles: []MxFile{
			MxFile{Category: "header", Name: `example.h`},
			MxFile{Category: "sourceAsm", Name: `example.s`},
			MxFile{Category: "source", Name: `example.cc`},
		}},
		components: MxComponents{Components: []MxComponent{
			MxComponent{
				Class:       "Device",
				Group:       "Startup",
				Version:     "2.1.0",
				Description: "System Startup for STMicroelectronics",
			},
		}}}
	expected := ccBinaryRule{rule{
		Operands: attributeList{
			attributeBString{Operand: "name", Value: "main"},
			attributeBStringList{Operand: attSrcs, Value: bStringList{"example.cc", "example.s", "example.h"}},
			attributeBStringList{Operand: attHdrs, Value: bStringList{"example.h"}},
			attributeBStringList{Operand: "deps", Value: bStringList{":" + ccLibraryTargetName(project.components.Components[0])}},
		},
		comment: comment{Comment: "# Main target"},
	}}
	got := MxProjectToCcBinaryRule(project)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
}
