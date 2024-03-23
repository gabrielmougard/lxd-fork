package wasmer

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Config tests

func testGetBytes(moduleFileName string) []byte {
	_, filename, _, _ := runtime.Caller(0)
	modulePath := path.Join(path.Dir(filename), "testdata", moduleFileName)
	bytes, _ := os.ReadFile(modulePath)

	return bytes
}

func TestCompilerKind(t *testing.T) {
	assert.Equal(t, CRANELIFT.String(), "cranelift")
	assert.Equal(t, LLVM.String(), "llvm")
	assert.Equal(t, SINGLEPASS.String(), "singlepass")
}

func TestEngineKind(t *testing.T) {
	assert.Equal(t, UNIVERSAL.String(), "universal")
}

func TestConfig(t *testing.T) {
	config := NewConfig()

	engine := NewEngineWithConfig(config)
	store := NewStore(engine)
	module, err := NewModule(store, testGetBytes("tests.wasm"))
	assert.NoError(t, err)

	instance, err := NewInstance(module, NewImportObject())
	assert.NoError(t, err)

	sum, err := instance.Exports.GetFunction("sum")
	assert.NoError(t, err)

	result, err := sum(37, 5)
	assert.NoError(t, err)
	assert.Equal(t, result, int32(42))
}

func TestConfigForMetering(t *testing.T) {
	opmap := map[Opcode]uint32{
		End:      1,
		LocalGet: 1,
		I32Add:   4,
	}

	// We allocate 800000000 'gas' points for the metering middleware
	config := NewConfig().PushMeteringMiddleware(800000000, opmap)
	engine := NewEngineWithConfig(config)
	store := NewStore(engine)
	module, err := NewModule(store, testGetBytes("tests.wasm"))
	assert.NoError(t, err)

	instance, err := NewInstance(module, NewImportObject())
	assert.NoError(t, err)

	sum, err := instance.Exports.GetFunction("sum")
	assert.NoError(t, err)

	result, err := sum(37, 5)
	assert.NoError(t, err)
	assert.Equal(t, result, int32(42))
	rp := instance.GetRemainingPoints()
	assert.Equal(t, int(rp), 800000000-7)
	// total instruction count should be 7:
	// 1 (`local.get`, for the first param of the function) +
	// 1 (`local.get`, for the second param of the function) +
	// 4 (`i32.add`, the addition within the function) +
	// 1 (`end`, the end of the function)
}

// func TestConfigForMeteringFn(t *testing.T) {
// 	config := NewConfig().PushMeteringMiddlewarePtr(800000000, getInternalCPointer())
// 	engine := NewEngineWithConfig(config)
// 	store := NewStore(engine)
// 	module, err := NewModule(store, testGetBytes("tests.wasm"))
// 	assert.NoError(t, err)

// 	instance, err := NewInstance(module, NewImportObject())
// 	assert.NoError(t, err)

// 	sum, err := instance.Exports.GetFunction("sum")
// 	assert.NoError(t, err)

// 	result, err := sum(37, 5)
// 	assert.NoError(t, err)
// 	assert.Equal(t, result, int32(42))
// 	rp := instance.GetRemainingPoints()
// 	assert.Equal(t, int(rp), 800000000-7)
// 	// total instruction count should be 7
// }

func TestConfig_AllCombinations(t *testing.T) {
	type Test struct {
		compilerName string
		engineName   string
		config       *Config
	}

	var configs = []Test{}

	is_amd64 := runtime.GOARCH == "amd64"
	is_arm64 := runtime.GOARCH == "arm64"
	is_linux := runtime.GOOS == "linux"
	has_universal := IsEngineAvailable(UNIVERSAL)

	if IsCompilerAvailable(CRANELIFT) {
		// Cranelift with the Universal engine works everywhere
		if has_universal {
			configs = append(configs, Test{"Cranelift", "Universal", NewConfig().UseCraneliftCompiler().UseUniversalEngine()})
		}
	}

	// The LLVM backend for Wasmer is disabled in the C API for the moment because it causes the linker to fail since LLVM is not statically linked.
	// TODO: Open issue in https://github.com/wasmerio/wasmer/issues
	// TODO: I noticed that Wasmer uses the 0.1.1 version of Inkwell to interface with LLVM in a safe way. Starting after 0.3.0, Inkwell introduced new compilation features
	//       enabling 'LLVM linking preferences' which could solve the issue (https://github.com/TheDan64/inkwell/compare/0.2.0...0.3.0).
	//       Basically, we'd need the `compiler-llvm` crate to be compiled with the 0.3.0 (or higher) version of Inkwell with either the `llvm15-0-force-static` or `llvm15-0-force-dynamic` feature.
	//       Then we could add back the `llvm` feature in the wasmer Makefile for the `build-capi` target (`c-api` wasmer crate) as LLVM will be linked with Inkwell.
	//
	//       Should we try this approach and file a PR to be merged with Wasmer upstream? Or should we create a local PATCH file for our use case if this got delayed by the Wasmer maintainers?
	if IsCompilerAvailable(LLVM) {
		// LLVM with the Universal engine works on Linux  (TODO: Darwin support should be ok too, but let's not test it for now)
		if has_universal && (is_amd64 || is_arm64) && is_linux {
			configs = append(configs, Test{"LLVM", "Universal", NewConfig().UseLLVMCompiler().UseUniversalEngine()})
		}
	} else {
		t.Log("LLVM compiler is disabled. See above comment for more information.")
	}

	if IsCompilerAvailable(SINGLEPASS) {
		// Singlepass with the Universal engine works on Linux (TODO: Darwin support should be ok too, but let's not test it for now)
		if has_universal && (is_amd64 || is_arm64) && is_linux {
			configs = append(configs, Test{"Singlepass", "Universal", NewConfig().UseSinglepassCompiler().UseUniversalEngine()})
		}
	}

	for _, test := range configs {
		t.Run(
			fmt.Sprintf("compiler=%s, engine=%s", test.compilerName, test.engineName),
			func(t *testing.T) {
				engine := NewEngineWithConfig(test.config)
				store := NewStore(engine)
				module, err := NewModule(store, testGetBytes("tests.wasm"))
				assert.NoError(t, err)

				instance, err := NewInstance(module, NewImportObject())
				assert.NoError(t, err)

				sum, err := instance.Exports.GetFunction("sum")
				assert.NoError(t, err)

				result, err := sum(37, 5)
				assert.NoError(t, err)
				assert.Equal(t, result, int32(42))
			},
		)
	}
}
