package cubemxtobazelinternal

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/silvergasp/CubeMxToBazel/data"
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
		Keys: attributeList{
			attributeBString{Key: "name", Value: bString(ccLibraryTargetName(proj.Components()[0]))},
			attributeBStringList{Key: attSrcs, Value: bStringList{"example.cc", "example.s", "example.h"}},
			attributeBStringList{Key: attHdrs, Value: bStringList{"example.h"}},
			attributeBBool{Key: attLinkStatic, Value: true},
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
		Keys: attributeList{
			attributeBString{Key: "name", Value: "main"},
			attributeBStringList{Key: attSrcs, Value: bStringList{"example.cc", "example.s", "example.h"}},
			attributeBStringList{Key: "deps", Value: bStringList{":" + ccLibraryTargetName(project.components.Components[0])}},
		},
		comment: comment{Comment: "# Main target"},
	}}
	got := MxProjectToCcBinaryRule(project)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
}

func TestParsePackageProjectInitComponents(t *testing.T) {
	gpdsc := data.SampleStm32Gpdsc()

	got := ProjectInit(gpdsc)
	startupFiles := []MxFile{
		MxFile{Category: "header", Condition: "", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.h`},
		MxFile{Category: "sourceAsm", Condition: "IAR Toolchain", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Source\Templates\iar\startup_stm32l432xx.s`},
		MxFile{Category: "sourceAsm", Condition: "GCC Toolchain", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Source\Templates\gcc\startup_stm32l432xx.s`},
	}
	expectedComponents := []MxComponent{
		MxComponent{
			Class:       "CMSIS",
			Group:       "CORE",
			Version:     "4.0.0",
			Description: "CMSIS-CORE for ARM",
			Files:       []MxFile{MxFile{Category: "header", Name: `Drivers\CMSIS\Include\core_cm4.h`}},
		},
		MxComponent{
			Class:       "Device",
			Group:       "Startup",
			Version:     "2.1.0",
			Description: "System Startup for STMicroelectronics",
			Files:       startupFiles,
		},
	}
	if diff := deep.Equal(got.Components(), expectedComponents); diff != nil {
		t.Error(diff)
	}
}
