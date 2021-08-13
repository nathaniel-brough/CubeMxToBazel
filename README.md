[![Build Status](https://dev.azure.com/17759661/CubeMxToBazel/_apis/build/status/silvergasp.CubeMxToBazel?branchName=master)](https://dev.azure.com/17759661/CubeMxToBazel/_build/latest?definitionId=1?branchName=master)

# CubeMxToBazel

Converts STM32CUBEMX Projects to bazel projects. This is acheived by making use of the generated .gpdsc files from stm32cubemx and outputing these in a bazel build file.

This project is in the early developement stage with fairly minimal but well tested functionality.

## Usage

### Installation

#### From Source

Requires a valid [golang](https://github.com/golang/go/wiki/Ubuntu) installation

```sh
go get -u github.com/gobuffalo/packr/v2/packr2
go get -u github.com/silvergasp/CubeMxToBazel
packr2 install github.com/silvergasp/CubeMxToBazel
```

#### Binary Installation

```sh
wget https://github.com/silvergasp/CubeMxToBazel/releases/download/v1.0-beta.1/CubeMxToBazel-linux.tar.gz
tar xvzf CubeMxToBazel-linux.tar.gz
mv CubeMxToBazel ~/bin/
```

### Setting up stm32cubemx

Select the "Other Toolchains (GPDSC)" configuration for your project in 'project>settings'.

![cubemxSettings](imgs/project_settings_configuration.png "stm32cubemx settings")

### Running generator

Running the generator is as simple as changing directories into the project and running the converter

```sh
cd YOUR_PROJECT_PATH_HERE
$GOPATH/bin/CubeMxToBazel
```

The output of this is not neccesarily nice to look at, it is recommended to use bazel auto formatter for this.

```sh
# Install Autoformatter
go install github.com/bazelbuild/buildtools/buildifier
# Run formatter on generated build file
$GOPATH/bin/buildifier BUILD
```

### Building using bazel

Build all targets

```sh
cd YOUR_PROJECT_PATH_HERE
bazel build ... --crosstool_top=@bazel_arm_none//tools/arm_compiler:toolchain --cpu=armeabi-v7a
```

The resulting binary executable and binary libraries can be found under `YOUR_PROJECT_PATH_HERE/bazel-bin`. The resulting executable will be named `YOUR_PROJECT_PATH_HERE/bazel-bin/main`.

**NOTE:** On first run this project will download the official arm compiler, as such it is expected that you have a strong internet connection. This downloaded compiler will be cached and does not require re-downloading provided that the arm_none_bazel rule in `WORKSPACE` is not modified.

## Current Functionality

- [x] Converts generated `*.gpdsc` files from stm32cubemx into bazel BUILD files
- [x] Generates bazel WORKSPACE files
- [ ] Generate appropriate bazel compiler flags for; fpu, cpu, hosting specs, optimisations (Currently supports armv7e devices with hardware fpu)
- [x] Implement conditional file inclusion (e.g. conditional inclusion of assembly files based on the compiler)
