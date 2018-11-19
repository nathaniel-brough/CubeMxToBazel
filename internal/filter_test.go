package cubemxtobazelinternal

import (
	"reflect"
	"testing"
)

func TestFilterAssemblyFiles(t *testing.T) {
	var startupTestFiles MxFiles = []MxFile{
		MxFile{Category: "header", Condition: "", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.h`},
		MxFile{Category: "sourceAsm", Condition: "IAR Toolchain", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Source\Templates\iar\startup_stm32l432xx.s`},
		MxFile{Category: "sourceAsm", Condition: "GCC Toolchain", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Source\Templates\gcc\startup_stm32l432xx.s`},
	}
	var expected MxFiles = []MxFile{
		MxFile{Category: "sourceAsm", Condition: "IAR Toolchain", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Source\Templates\iar\startup_stm32l432xx.s`},
		MxFile{Category: "sourceAsm", Condition: "GCC Toolchain", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Source\Templates\gcc\startup_stm32l432xx.s`},
	}
	got := startupTestFiles.AssemblyFiles().Files()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
}

func TestFilterSourceFiles(t *testing.T) {
	var startupTestFiles MxFiles = []MxFile{
		MxFile{Category: "header", Condition: "", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.h`},
		MxFile{Category: "source", Condition: "", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.c`},
		MxFile{Category: "sourceC", Condition: "", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.c`},
	}
	var expected MxFiles = []MxFile{
		MxFile{Category: "source", Condition: "", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.c`},
		MxFile{Category: "sourceC", Condition: "", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.c`},
	}
	got := startupTestFiles.SourceFiles().Files()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
}

func TestFilterHeaderFiles(t *testing.T) {
	var startupTestFiles MxFiles = []MxFile{
		MxFile{Category: "header", Condition: "", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.h`},
		MxFile{Category: "source", Condition: "", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.c`},
		MxFile{Category: "sourceC", Condition: "", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.c`},
	}
	var expected MxFiles = []MxFile{
		MxFile{Category: "header", Condition: "", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.h`},
	}
	got := startupTestFiles.HeaderFiles().Files()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
}

func TestFilterConditionalFiles(t *testing.T) {
	condition := "GCC Toolchain"
	var startupTestFiles MxFiles = []MxFile{
		MxFile{Category: "header", Condition: "GCC Toolchain", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.h`},
		MxFile{Category: "source", Condition: "", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.c`},
	}
	var expected MxFiles = []MxFile{
		MxFile{Category: "header", Condition: "GCC Toolchain", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.h`},
	}
	got := startupTestFiles.Condition(condition).Files()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
}

func TestFilterConditionalChaining(t *testing.T) {
	condition := "GCC Toolchain"
	var startupTestFiles MxFiles = []MxFile{
		MxFile{Category: "header", Condition: "GCC Toolchain", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.h`},
		MxFile{Category: "header", Name: "Inc/stm32l4xx_hal_conf.h"},
		MxFile{Category: "source", Condition: "GCC Toolchain", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.c`},
	}
	var expected MxFiles = []MxFile{
		MxFile{Category: "header", Condition: "GCC Toolchain", Name: `Drivers\CMSIS\Device\ST\STM32L4xx\Include\stm32l4xx.h`},
	}
	got := startupTestFiles.Condition(condition).HeaderFiles().Files()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected:\n%#v \nGot:\n%#v \n", expected, got)
	}
}
