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

func testGetInstance(t *testing.T) *Instance {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(store, testGetBytes("tests.wasm"))
	assert.NoError(t, err)

	instance, err := NewInstance(module, NewImportObject())
	assert.NoError(t, err)

	return instance
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

// Engine tests

func testEngine(t *testing.T, engine *Engine) {
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

func TestEngine(t *testing.T) {
	testEngine(t, NewEngine())
}

func TestEngineWithTarget(t *testing.T) {
	triple, err := NewTriple("aarch64-unknown-linux-gnu")
	assert.NoError(t, err)

	cpuFeatures := NewCpuFeatures()
	assert.NoError(t, err)

	target := NewTarget(triple, cpuFeatures)

	config := NewConfig()
	config.UseTarget(target)

	engine := NewEngineWithConfig(config)
	store := NewStore(engine)

	module, err := NewModule(store, testGetBytes("tests.wasm"))
	assert.NoError(t, err)

	_ = module
}

// Function tests

func TestRawFunction(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(
		store,
		[]byte(`
			(module
			  (type $sum_t (func (param i32 i32) (result i32)))
			  (func $sum_f (type $sum_t) (param $x i32) (param $y i32) (result i32)
			    local.get $x
			    local.get $y
			    i32.add)
			  (export "sum" (func $sum_f)))
		`),
	)
	assert.NoError(t, err)

	instance, err := NewInstance(module, NewImportObject())
	assert.NoError(t, err)

	sum, err := instance.Exports.GetRawFunction("sum")
	assert.NoError(t, err)
	assert.Equal(t, sum.ParameterArity(), uint(2))
	assert.Equal(t, sum.ResultArity(), uint(1))

	result, err := sum.Call(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, result, int32(3))
}

func TestFunctionCallReturnZeroValue(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(
		store,
		[]byte(`
			(module
			  (type $test_t (func (param i32 i32)))
			  (func $test_f (type $test_t) (param $x i32) (param $y i32))
			  (export "test" (func $test_f)))
		`),
	)
	assert.NoError(t, err)

	instance, err := NewInstance(module, NewImportObject())
	assert.NoError(t, err)

	test, err := instance.Exports.GetFunction("test")
	assert.NoError(t, err)

	result, err := test(1, 2)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestFunctionCallReturnMultipleValues(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(
		store,
		[]byte(`
			(module
			  (type $swap_t (func (param i32 i64) (result i64 i32)))
			  (func $swap_f (type $swap_t) (param $x i32) (param $y i64) (result i64 i32)
			    local.get $y
			    local.get $x)
			  (export "swap" (func $swap_f)))
		`),
	)
	assert.NoError(t, err)

	instance, err := NewInstance(module, NewImportObject())
	assert.NoError(t, err)

	swap, err := instance.Exports.GetFunction("swap")
	assert.NoError(t, err)

	results, err := swap(41, 42)
	assert.NoError(t, err)
	assert.Equal(t, results, []interface{}{int64(42), int32(41)})
}

func TestFunctionSum(t *testing.T) {
	instance := testGetInstance(t)

	f, _ := instance.Exports.GetFunction("sum")
	result, err := f(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, result, int32(3))
}

func TestFunctionArity0(t *testing.T) {
	instance := testGetInstance(t)

	f, _ := instance.Exports.GetFunction("arity_0")
	result, err := f()
	assert.NoError(t, err)
	assert.Equal(t, result, int32(42))
}

func TestFunctionI32I32(t *testing.T) {
	instance := testGetInstance(t)

	f, _ := instance.Exports.GetFunction("i32_i32")
	result, err := f(7)
	assert.NoError(t, err)
	assert.Equal(t, result, int32(7))

	result, _ = f(int8(7))
	assert.Equal(t, result, int32(7))

	result, _ = f(uint8(7))
	assert.Equal(t, result, int32(7))

	result, _ = f(int16(7))
	assert.Equal(t, result, int32(7))

	result, _ = f(uint16(7))
	assert.Equal(t, result, int32(7))

	result, _ = f(int32(7))
	assert.Equal(t, result, int32(7))

	result, _ = f(int(7))
	assert.Equal(t, result, int32(7))

	result, _ = f(uint(7))
	assert.Equal(t, result, int32(7))
}

func TestFunctionI64I64(t *testing.T) {
	instance := testGetInstance(t)

	f, _ := instance.Exports.GetFunction("i64_i64")
	result, err := f(7)
	assert.NoError(t, err)
	assert.Equal(t, result, int64(7))

	result, _ = f(int8(7))
	assert.Equal(t, result, int64(7))

	result, _ = f(uint8(7))
	assert.Equal(t, result, int64(7))

	result, _ = f(int16(7))
	assert.Equal(t, result, int64(7))

	result, _ = f(uint16(7))
	assert.Equal(t, result, int64(7))

	result, _ = f(int32(7))
	assert.Equal(t, result, int64(7))

	result, _ = f(int64(7))
	assert.Equal(t, result, int64(7))

	result, _ = f(int(7))
	assert.Equal(t, result, int64(7))

	result, _ = f(uint(7))
	assert.Equal(t, result, int64(7))
}

func TestFunctionF32F32(t *testing.T) {
	instance := testGetInstance(t)

	f, _ := instance.Exports.GetFunction("f32_f32")
	result, err := f(float32(7.42))
	assert.NoError(t, err)
	assert.Equal(t, result, float32(7.42))
}

func TestFunctionF64F64(t *testing.T) {
	instance := testGetInstance(t)

	f, _ := instance.Exports.GetFunction("f64_f64")
	result, err := f(7.42)
	assert.NoError(t, err)
	assert.Equal(t, result, float64(7.42))

	result, _ = f(float64(7.42))
	assert.Equal(t, result, float64(7.42))
}

func TestFunctionI32I64F32F64F64(t *testing.T) {
	instance := testGetInstance(t)

	f, _ := instance.Exports.GetFunction("i32_i64_f32_f64_f64")
	result, err := f(1, 2, float32(3.4), 5.6)
	assert.NoError(t, err)
	assert.Equal(t, float64(int(result.(float64)*10000000))/10000000, 1+2+3.4+5.6)
}

func TestFunctionBoolCastedtoI32(t *testing.T) {
	instance := testGetInstance(t)

	f, _ := instance.Exports.GetFunction("bool_casted_to_i32")
	result, err := f()
	assert.NoError(t, err)
	assert.Equal(t, result, int32(1))
}

func TestHostFunction(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(
		store,
		[]byte(`
			(module
			  (import "math" "sum" (func $sum (param i32 i32) (result i32)))
			  (func (export "add_one") (param $x i32) (result i32)
			    local.get $x
			    i32.const 1
			    call $sum))
		`),
	)
	assert.NoError(t, err)

	function := NewFunction(
		store,
		NewFunctionType(NewValueTypes(I32, I32), NewValueTypes(I32)),
		func(args []Value) ([]Value, error) {
			x := args[0].I32()
			y := args[1].I32()

			return []Value{NewI32(x + y)}, nil
		},
	)

	importObject := NewImportObject()
	importObject.Register(
		"math",
		map[string]IntoExtern{
			"sum": function,
		},
	)

	instance, err := NewInstance(module, importObject)
	assert.NoError(t, err)

	addOne, err := instance.Exports.GetFunction("add_one")
	assert.NoError(t, err)

	result, err := addOne(41)
	assert.NoError(t, err)
	assert.Equal(t, result, int32(42))
}

func TestHostFunction_WithI64(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(
		store,
		[]byte(`
			(module
			  (import "math" "sum" (func $sum (param i64 i64) (result i64)))
			  (func (export "add_one") (param $x i64) (result i64)
			    local.get $x
			    i64.const 1
			    call $sum))
		`),
	)
	assert.NoError(t, err)

	function := NewFunction(
		store,
		NewFunctionType(NewValueTypes(I64, I64), NewValueTypes(I64)),
		func(args []Value) ([]Value, error) {
			x := args[0].I64()
			y := args[1].I64()

			return []Value{NewI64(x + y)}, nil
		},
	)

	importObject := NewImportObject()
	importObject.Register(
		"math",
		map[string]IntoExtern{
			"sum": function,
		},
	)

	instance, err := NewInstance(module, importObject)
	assert.NoError(t, err)

	addOne, err := instance.Exports.GetFunction("add_one")
	assert.NoError(t, err)

	result, err := addOne(41)
	assert.NoError(t, err)
	assert.Equal(t, result, int64(42))
}

func TestHostFunctionWithEnv(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(
		store,
		[]byte(`
			(module
			  (import "math" "sum" (func $sum (param i32 i32) (result i32)))
			  (func (export "add_one") (param $x i32) (result i32)
			    local.get $x
			    i32.const 1
			    call $sum))
		`),
	)
	assert.NoError(t, err)

	type MyEnvironment struct {
		instance  *Instance
		theAnswer int32
	}

	environment := &MyEnvironment{
		instance:  nil,
		theAnswer: 42,
	}

	function := NewFunctionWithEnvironment(
		store,
		NewFunctionType(NewValueTypes(I32, I32), NewValueTypes(I32)),
		environment,
		func(environment interface{}, args []Value) ([]Value, error) {
			env := environment.(*MyEnvironment)
			assert.NotNil(t, env.instance)

			x := args[0].I32()
			y := args[1].I32()

			return []Value{NewI32(env.theAnswer + x + y)}, nil
		},
	)

	importObject := NewImportObject()
	importObject.Register(
		"math",
		map[string]IntoExtern{
			"sum": function,
		},
	)

	instance, err := NewInstance(module, importObject)
	assert.NoError(t, err)

	environment.instance = instance

	addOne, err := instance.Exports.GetFunction("add_one")
	assert.NoError(t, err)

	result, err := addOne(7)
	assert.NoError(t, err)
	assert.Equal(t, result, int32(50))
}

func TestHostFunctionStore(t *testing.T) {
	f := &hostFunction{
		store: NewStore(NewEngine()),
		function: func(args []Value) ([]Value, error) {
			return []Value{}, nil
		},
	}

	store := hostFunctions{
		functions: make(map[uint]*hostFunction),
	}
	_, err := store.load(0)
	assert.Error(t, err, "Host function `0` does not exist")

	indexA := store.store(f)
	indexB := store.store(f)
	indexC := store.store(f)
	assert.Equal(t, indexA, uint(0))
	assert.Equal(t, indexB, uint(1))
	assert.Equal(t, indexC, uint(2))

	store.remove(indexB)
	_, err = store.load(indexB)
	assert.Error(t, err, "Host function `1` does not exist")

	indexD := store.store(f)
	assert.Equal(t, indexD, indexB)
}

type myError struct {
	message string
}

func (e *myError) Error() string {
	return e.message
}

func TestHostFunctionTrap(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(
		store,
		[]byte(`
			(module
			  (import "math" "sum" (func $sum (param i32 i32) (result i32)))
			  (func (export "add_one") (param $x i32) (result i32)
			    local.get $x
			    i32.const 1
			    call $sum))
		`),
	)
	assert.NoError(t, err)

	function := NewFunction(
		store,
		NewFunctionType(NewValueTypes(I32, I32), NewValueTypes(I32)),
		func(args []Value) ([]Value, error) {
			// Go Trap, go!
			return nil, &myError{message: "oops"}
		},
	)

	importObject := NewImportObject()
	importObject.Register(
		"math",
		map[string]IntoExtern{
			"sum": function,
		},
	)

	instance, err := NewInstance(module, importObject)
	assert.NoError(t, err)

	addOne, err := instance.Exports.GetFunction("add_one")
	assert.NoError(t, err)

	_, err = addOne(41)
	assert.IsType(t, err, &TrapError{})
	assert.Error(t, err, "oops")
}
