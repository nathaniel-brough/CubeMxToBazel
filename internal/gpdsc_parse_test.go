package cubemxtobazelinternal

import (
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/silvergasp/CubeMxToBazel/data"
)

func TestParsePackageRequirements(t *testing.T) {
	gpdsc := data.SampleStm32Gpdsc()
	got := MxRequirements{}
	err := xml.Unmarshal(gpdsc, &got)
	if err != nil {
		t.Errorf("Unmarshal Failed: %#v", err)
	}
	langC := MxLanguage{Name: "C", Version: "99"}
	langCC := MxLanguage{Name: "C++", Version: "11"}
	languages := []MxLanguage{langC, langCC}
	expected := MxRequirements{Languages: languages}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
}

func TestParsePackageInfo(t *testing.T) {
	gpdsc := data.SampleStm32Gpdsc()
	got := MxInfo{}
	err := xml.Unmarshal(gpdsc, &got)
	if err != nil {
		t.Errorf("Unmarshal Failed: %#v", err)
	}
	expected := MxInfo{Vendor: "STMicroelectronics", Name: "stm32cubeTest", Description: "STM32CubeMX generated pack description"}
	if expected != got {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
}

func TestParsePackageOptions(t *testing.T) {
	gpdsc := data.SampleStm32Gpdsc()
	got := MxOptions{}
	err := xml.Unmarshal(gpdsc, &got)
	if err != nil {
		t.Errorf("Unmarshal Failed: %#v", err)
	}
	expected := MxOptions{
		StackSize: MxStack{"0x400"},
		HeapSize:  MxHeap{"0x200"},
		DebugProbe: MxDebugProbe{
			Name:     "ST-Link",
			Protocol: "swd",
		},
	}
	if expected != got {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
}

func TestParsePackageGenerator(t *testing.T) {
	gpdsc := data.SampleStm32Gpdsc()
	got := MxGenerator{}
	err := xml.Unmarshal(gpdsc, &got)
	if err != nil {
		t.Errorf("Unmarshal Failed: %#v", err)
	}

	files := []MxFile{
		MxFile{Category: "header", Name: "Inc/stm32l4xx_it.h"},
		MxFile{Category: "header", Name: "Inc/stm32l4xx_hal_conf.h"},
		MxFile{Category: "header", Name: "Inc/main.h"},
		MxFile{Category: "source", Name: "Src/stm32l4xx_it.c"},
		MxFile{Category: "source", Name: "Src/stm32l4xx_hal_msp.c"},
		MxFile{Category: "source", Name: "Src/main.c"},
	}
	expected := MxGenerator{Select: MxSelect{DeviceName: "STM32L432KCUx"}, ProjectFiles: files}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}

}

func TestParsePackageComponents(t *testing.T) {
	gpdsc := data.SampleStm32Gpdsc()
	got := MxComponents{}
	err := xml.Unmarshal(gpdsc, &got)
	if err != nil {
		t.Errorf("Unmarshal Failed: %#v", err)
	}

	startupFiles := []MxFile{
		MxFile{Category: "header", Condition: "", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.h`},
		MxFile{Category: "sourceAsm", Condition: "IAR Toolchain", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Source\Templates\iar\startup_stm32l432xx.s`},
		MxFile{Category: "sourceAsm", Condition: "GCC Toolchain", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Source\Templates\gcc\startup_stm32l432xx.s`},
	}

	expected := MxComponents{Components: []MxComponent{
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
	}}
	if diff := deep.Equal(got, expected); diff != nil {
		t.Error(diff)
	}

}

func TestParsePackageConditions(t *testing.T) {
	gpdsc := data.SampleStm32Gpdsc()
	got := MxConditions{}
	err := xml.Unmarshal(gpdsc, &got)
	if err != nil {
		t.Errorf("Unmarshal Failed: %#v", err)
	}

	expected := MxConditions{Conditions: []MxCondition{
		MxCondition{
			ID:          "ARM Toolchain",
			Description: "ARM compiler for C and C++ Filter",
			Require:     MxRequire{Compiler: "ARMCC"},
		},
		MxCondition{
			ID:          "GCC Toolchain",
			Description: "GNU Tools for ARM Embedded Processors Filter",
			Require:     MxRequire{Compiler: "GCC"},
		},
		MxCondition{
			ID:          "IAR Toolchain",
			Description: "IAR compiler for C and C++ Filter",
			Require:     MxRequire{Compiler: "IAR"},
		},
	},
	}
	if diff := deep.Equal(got, expected); diff != nil {
		t.Error(diff)
	}
}
