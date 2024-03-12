// MIT License

// Copyright (c) 2019-present Wasmer, Inc. and its affiliates.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

//go:build linux && cgo

// This is heavily inspired by the deprecated Wasmer Go SDK: https://github.com/wasmerio/wasmer-go
package wasmer

// #cgo LDFLAGS: -lwasmer
// #include <stdlib.h>
// #include <stdio.h>
// #include "wasmer.h"
//
// extern wasm_trap_t* function_trampoline(
//   void *environment,
//   /* const */ wasm_val_vec_t* arguments,
//   wasm_val_vec_t* results
// );
//
// extern wasm_trap_t* function_with_environment_trampoline(
//   void *environment,
//   /* const */ wasm_val_vec_t* arguments,
//   wasm_val_vec_t* results
// );
//
// typedef void (*wasm_func_callback_env_finalizer_t)(void* environment);
//
// extern void function_environment_finalizer(void *environment);
//
// extern uint64_t metering_delegate(enum wasmer_parser_operator_t op);
//
// #define own
//
// // We can't create a `wasm_byte_vec_t` directly in Go otherwise cgo
// // complains with “Go pointer to Go pointer”. The hack consists at
// // creating the `wasm_byte_vec_t` directly in C.
//
// static own wasm_module_t* to_wasm_module_new(wasm_store_t *store, uint8_t *bytes, size_t bytes_length) {
//     wasm_byte_vec_t wasm_bytes;
//     wasm_bytes.size = bytes_length;
//     wasm_bytes.data = (wasm_byte_t*) bytes;
//
//     return wasm_module_new(store, &wasm_bytes);
// }
//
// static bool to_wasm_module_validate(wasm_store_t *store, uint8_t *bytes, size_t bytes_length) {
//     wasm_byte_vec_t wasm_bytes;
//     wasm_bytes.size = bytes_length;
//     wasm_bytes.data = (wasm_byte_t*) bytes;
//
//     return wasm_module_validate(store, &wasm_bytes);
// }
//
// static wasm_module_t* to_wasm_module_deserialize(wasm_store_t *store, uint8_t *bytes, size_t bytes_length) {
//     wasm_byte_vec_t serialized_bytes;
//     serialized_bytes.size = bytes_length;
//     serialized_bytes.data = (wasm_byte_t*) bytes;
//
//     return wasm_module_deserialize(store, &serialized_bytes);
// }
//
// static own wasm_trap_t* to_wasm_trap_new(wasm_store_t *store, uint8_t *message_bytes, size_t message_length) {
//     // `wasm_message_t` is an alias to `wasm_byte_vec_t`.
//     wasm_message_t message;
//     message.size = message_length;
//     message.data = (wasm_byte_t*) message_bytes;
//
//     return wasm_trap_new(store, &message);
// }
//
// static int32_t to_int32(wasm_val_t *value) {
//     return value->of.i32;
// }
//
// static int64_t to_int64(wasm_val_t *value) {
//     return value->of.i64;
// }
//
// static float32_t to_float32(wasm_val_t *value) {
//     return value->of.f32;
// }
//
// static float64_t to_float64(wasm_val_t *value) {
//     return value->of.f64;
// }
//
// static wasm_ref_t *to_ref(wasm_val_t *value) {
//     return value->of.ref;
// }
//
// // Buffer size for `wasi_env_read_inner`.
// #define WASI_ENV_READER_BUFFER_SIZE 1024
//
// // Define a type for the WASI environment captured stream readers
// // (`wasi_env_read_stdout` and `wasi_env_read_stderr`).
// typedef intptr_t (*wasi_env_reader)(
//     wasi_env_t* wasi_env,
//     char* buffer,
//     uintptr_t buffer_len
// );
//
// // Common function to read a WASI environment captured stream.
// static size_t to_wasi_env_read_inner(wasi_env_t *wasi_env, char** buffer, wasi_env_reader reader) {
//     FILE *memory_stream;
//     size_t buffer_size = 0;
//
//     memory_stream = open_memstream(buffer, &buffer_size);
//
//     if (NULL == memory_stream) {
//         return 0;
//     }
//
//     char temp_buffer[WASI_ENV_READER_BUFFER_SIZE] = { 0 };
//     size_t data_read_size = WASI_ENV_READER_BUFFER_SIZE;
//
//     do {
//         data_read_size = reader(wasi_env, temp_buffer, WASI_ENV_READER_BUFFER_SIZE);
//
//         if (data_read_size > 0) {
//             buffer_size += data_read_size;
//             fwrite(temp_buffer, sizeof(char), data_read_size, memory_stream);
//         }
//     } while (WASI_ENV_READER_BUFFER_SIZE == data_read_size);
//
//     fclose(memory_stream);
//
//     return buffer_size;
// }
//
// // Read the captured `stdout`.
// static size_t to_wasi_env_read_stdout(wasi_env_t *wasi_env, char** buffer) {
//     return to_wasi_env_read_inner(wasi_env, buffer, wasi_env_read_stdout);
// }
//
// // Read the captured `stderr`.
// static size_t to_wasi_env_read_stderr(wasi_env_t *wasi_env, char** buffer) {
//     return to_wasi_env_read_inner(wasi_env, buffer, wasi_env_read_stderr);
// }
import "C"
import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

type Opcode C.wasmer_parser_operator_t

const (
	Unreachable Opcode = iota
	Nop
	Block
	Loop
	If
	Else
	Try
	Catch
	CatchAll
	Delegate
	Throw
	Rethrow
	Unwind
	End
	Br
	BrIf
	BrTable
	Return
	Call
	CallIndirect
	ReturnCall
	ReturnCallIndirect
	Drop
	Select
	TypedSelect
	LocalGet
	LocalSet
	LocalTee
	GlobalGet
	GlobalSet
	I32Load
	I64Load
	F32Load
	F64Load
	I32Load8S
	I32Load8U
	I32Load16S
	I32Load16U
	I64Load8S
	I64Load8U
	I64Load16S
	I64Load16U
	I64Load32S
	I64Load32U
	I32Store
	I64Store
	F32Store
	F64Store
	I32Store8
	I32Store16
	I64Store8
	I64Store16
	I64Store32
	MemorySize
	MemoryGrow
	I32Const
	I64Const
	F32Const
	F64Const
	RefNull
	RefIsNull
	RefFunc
	I32Eqz
	I32Eq
	I32Ne
	I32LtS
	I32LtU
	I32GtS
	I32GtU
	I32LeS
	I32LeU
	I32GeS
	I32GeU
	I64Eqz
	I64Eq
	I64Ne
	I64LtS
	I64LtU
	I64GtS
	I64GtU
	I64LeS
	I64LeU
	I64GeS
	I64GeU
	F32Eq
	F32Ne
	F32Lt
	F32Gt
	F32Le
	F32Ge
	F64Eq
	F64Ne
	F64Lt
	F64Gt
	F64Le
	F64Ge
	I32Clz
	I32Ctz
	I32Popcnt
	I32Add
	I32Sub
	I32Mul
	I32DivS
	I32DivU
	I32RemS
	I32RemU
	I32And
	I32Or
	I32Xor
	I32Shl
	I32ShrS
	I32ShrU
	I32Rotl
	I32Rotr
	I64Clz
	I64Ctz
	I64Popcnt
	I64Add
	I64Sub
	I64Mul
	I64DivS
	I64DivU
	I64RemS
	I64RemU
	I64And
	I64Or
	I64Xor
	I64Shl
	I64ShrS
	I64ShrU
	I64Rotl
	I64Rotr
	F32Abs
	F32Neg
	F32Ceil
	F32Floor
	F32Trunc
	F32Nearest
	F32Sqrt
	F32Add
	F32Sub
	F32Mul
	F32Div
	F32Min
	F32Max
	F32Copysign
	F64Abs
	F64Neg
	F64Ceil
	F64Floor
	F64Trunc
	F64Nearest
	F64Sqrt
	F64Add
	F64Sub
	F64Mul
	F64Div
	F64Min
	F64Max
	F64Copysign
	I32WrapI64
	I32TruncF32S
	I32TruncF32U
	I32TruncF64S
	I32TruncF64U
	I64ExtendI32S
	I64ExtendI32U
	I64TruncF32S
	I64TruncF32U
	I64TruncF64S
	I64TruncF64U
	F32ConvertI32S
	F32ConvertI32U
	F32ConvertI64S
	F32ConvertI64U
	F32DemoteF64
	F64ConvertI32S
	F64ConvertI32U
	F64ConvertI64S
	F64ConvertI64U
	F64PromoteF32
	I32ReinterpretF32
	I64ReinterpretF64
	F32ReinterpretI32
	F64ReinterpretI64
	I32Extend8S
	I32Extend16S
	I64Extend8S
	I64Extend16S
	I64Extend32S
	I32TruncSatF32S
	I32TruncSatF32U
	I32TruncSatF64S
	I32TruncSatF64U
	I64TruncSatF32S
	I64TruncSatF32U
	I64TruncSatF64S
	I64TruncSatF64U
	MemoryInit
	DataDrop
	MemoryCopy
	MemoryFill
	TableInit
	ElemDrop
	TableCopy
	TableFill
	TableGet
	TableSet
	TableGrow
	TableSizeOp
	MemoryAtomicNotify
	MemoryAtomicWait32
	MemoryAtomicWait64
	AtomicFence
	I32AtomicLoad
	I64AtomicLoad
	I32AtomicLoad8U
	I32AtomicLoad16U
	I64AtomicLoad8U
	I64AtomicLoad16U
	I64AtomicLoad32U
	I32AtomicStore
	I64AtomicStore
	I32AtomicStore8
	I32AtomicStore16
	I64AtomicStore8
	I64AtomicStore16
	I64AtomicStore32
	I32AtomicRmwAdd
	I64AtomicRmwAdd
	I32AtomicRmw8AddU
	I32AtomicRmw16AddU
	I64AtomicRmw8AddU
	I64AtomicRmw16AddU
	I64AtomicRmw32AddU
	I32AtomicRmwSub
	I64AtomicRmwSub
	I32AtomicRmw8SubU
	I32AtomicRmw16SubU
	I64AtomicRmw8SubU
	I64AtomicRmw16SubU
	I64AtomicRmw32SubU
	I32AtomicRmwAnd
	I64AtomicRmwAnd
	I32AtomicRmw8AndU
	I32AtomicRmw16AndU
	I64AtomicRmw8AndU
	I64AtomicRmw16AndU
	I64AtomicRmw32AndU
	I32AtomicRmwOr
	I64AtomicRmwOr
	I32AtomicRmw8OrU
	I32AtomicRmw16OrU
	I64AtomicRmw8OrU
	I64AtomicRmw16OrU
	I64AtomicRmw32OrU
	I32AtomicRmwXor
	I64AtomicRmwXor
	I32AtomicRmw8XorU
	I32AtomicRmw16XorU
	I64AtomicRmw8XorU
	I64AtomicRmw16XorU
	I64AtomicRmw32XorU
	I32AtomicRmwXchg
	I64AtomicRmwXchg
	I32AtomicRmw8XchgU
	I32AtomicRmw16XchgU
	I64AtomicRmw8XchgU
	I64AtomicRmw16XchgU
	I64AtomicRmw32XchgU
	I32AtomicRmwCmpxchg
	I64AtomicRmwCmpxchg
	I32AtomicRmw8CmpxchgU
	I32AtomicRmw16CmpxchgU
	I64AtomicRmw8CmpxchgU
	I64AtomicRmw16CmpxchgU
	I64AtomicRmw32CmpxchgU
	V128Load
	V128Store
	V128Const
	I8x16Splat
	I8x16ExtractLaneS
	I8x16ExtractLaneU
	I8x16ReplaceLane
	I16x8Splat
	I16x8ExtractLaneS
	I16x8ExtractLaneU
	I16x8ReplaceLane
	I32x4Splat
	I32x4ExtractLane
	I32x4ReplaceLane
	I64x2Splat
	I64x2ExtractLane
	I64x2ReplaceLane
	F32x4Splat
	F32x4ExtractLane
	F32x4ReplaceLane
	F64x2Splat
	F64x2ExtractLane
	F64x2ReplaceLane
	I8x16Eq
	I8x16Ne
	I8x16LtS
	I8x16LtU
	I8x16GtS
	I8x16GtU
	I8x16LeS
	I8x16LeU
	I8x16GeS
	I8x16GeU
	I16x8Eq
	I16x8Ne
	I16x8LtS
	I16x8LtU
	I16x8GtS
	I16x8GtU
	I16x8LeS
	I16x8LeU
	I16x8GeS
	I16x8GeU
	I32x4Eq
	I32x4Ne
	I32x4LtS
	I32x4LtU
	I32x4GtS
	I32x4GtU
	I32x4LeS
	I32x4LeU
	I32x4GeS
	I32x4GeU
	I64x2Eq
	I64x2Ne
	I64x2LtS
	I64x2GtS
	I64x2LeS
	I64x2GeS
	F32x4Eq
	F32x4Ne
	F32x4Lt
	F32x4Gt
	F32x4Le
	F32x4Ge
	F64x2Eq
	F64x2Ne
	F64x2Lt
	F64x2Gt
	F64x2Le
	F64x2Ge
	V128Not
	V128And
	V128AndNot
	V128Or
	V128Xor
	V128Bitselect
	V128AnyTrue
	I8x16Abs
	I8x16Neg
	I8x16AllTrue
	I8x16Bitmask
	I8x16Shl
	I8x16ShrS
	I8x16ShrU
	I8x16Add
	I8x16AddSatS
	I8x16AddSatU
	I8x16Sub
	I8x16SubSatS
	I8x16SubSatU
	I8x16MinS
	I8x16MinU
	I8x16MaxS
	I8x16MaxU
	I8x16Popcnt
	I16x8Abs
	I16x8Neg
	I16x8AllTrue
	I16x8Bitmask
	I16x8Shl
	I16x8ShrS
	I16x8ShrU
	I16x8Add
	I16x8AddSatS
	I16x8AddSatU
	I16x8Sub
	I16x8SubSatS
	I16x8SubSatU
	I16x8Mul
	I16x8MinS
	I16x8MinU
	I16x8MaxS
	I16x8MaxU
	I16x8ExtAddPairwiseI8x16S
	I16x8ExtAddPairwiseI8x16U
	I32x4Abs
	I32x4Neg
	I32x4AllTrue
	I32x4Bitmask
	I32x4Shl
	I32x4ShrS
	I32x4ShrU
	I32x4Add
	I32x4Sub
	I32x4Mul
	I32x4MinS
	I32x4MinU
	I32x4MaxS
	I32x4MaxU
	I32x4DotI16x8S
	I32x4ExtAddPairwiseI16x8S
	I32x4ExtAddPairwiseI16x8U
	I64x2Abs
	I64x2Neg
	I64x2AllTrue
	I64x2Bitmask
	I64x2Shl
	I64x2ShrS
	I64x2ShrU
	I64x2Add
	I64x2Sub
	I64x2Mul
	F32x4Ceil
	F32x4Floor
	F32x4Trunc
	F32x4Nearest
	F64x2Ceil
	F64x2Floor
	F64x2Trunc
	F64x2Nearest
	F32x4Abs
	F32x4Neg
	F32x4Sqrt
	F32x4Add
	F32x4Sub
	F32x4Mul
	F32x4Div
	F32x4Min
	F32x4Max
	F32x4PMin
	F32x4PMax
	F64x2Abs
	F64x2Neg
	F64x2Sqrt
	F64x2Add
	F64x2Sub
	F64x2Mul
	F64x2Div
	F64x2Min
	F64x2Max
	F64x2PMin
	F64x2PMax
	I32x4TruncSatF32x4S
	I32x4TruncSatF32x4U
	F32x4ConvertI32x4S
	F32x4ConvertI32x4U
	I8x16Swizzle
	I8x16Shuffle
	V128Load8Splat
	V128Load16Splat
	V128Load32Splat
	V128Load32Zero
	V128Load64Splat
	V128Load64Zero
	I8x16NarrowI16x8S
	I8x16NarrowI16x8U
	I16x8NarrowI32x4S
	I16x8NarrowI32x4U
	I16x8ExtendLowI8x16S
	I16x8ExtendHighI8x16S
	I16x8ExtendLowI8x16U
	I16x8ExtendHighI8x16U
	I32x4ExtendLowI16x8S
	I32x4ExtendHighI16x8S
	I32x4ExtendLowI16x8U
	I32x4ExtendHighI16x8U
	I64x2ExtendLowI32x4S
	I64x2ExtendHighI32x4S
	I64x2ExtendLowI32x4U
	I64x2ExtendHighI32x4U
	I16x8ExtMulLowI8x16S
	I16x8ExtMulHighI8x16S
	I16x8ExtMulLowI8x16U
	I16x8ExtMulHighI8x16U
	I32x4ExtMulLowI16x8S
	I32x4ExtMulHighI16x8S
	I32x4ExtMulLowI16x8U
	I32x4ExtMulHighI16x8U
	I64x2ExtMulLowI32x4S
	I64x2ExtMulHighI32x4S
	I64x2ExtMulLowI32x4U
	I64x2ExtMulHighI32x4U
	V128Load8x8S
	V128Load8x8U
	V128Load16x4S
	V128Load16x4U
	V128Load32x2S
	V128Load32x2U
	V128Load8Lane
	V128Load16Lane
	V128Load32Lane
	V128Load64Lane
	V128Store8Lane
	V128Store16Lane
	V128Store32Lane
	V128Store64Lane
	I8x16RoundingAverageU
	I16x8RoundingAverageU
	I16x8Q15MulrSatS
	F32x4DemoteF64x2Zero
	F64x2PromoteLowF32x4
	F64x2ConvertLowI32x4S
	F64x2ConvertLowI32x4U
	I32x4TruncSatF64x2SZero
	I32x4TruncSatF64x2UZero
	I8x16RelaxedSwizzle
	I32x4RelaxedTruncSatF32x4S
	I32x4RelaxedTruncSatF32x4U
	I32x4RelaxedTruncSatF64x2SZero
	I32x4RelaxedTruncSatF64x2UZero
	F32x4Fma
	F32x4Fms
	F64x2Fma
	F64x2Fms
	I8x16LaneSelect
	I16x8LaneSelect
	I32x4LaneSelect
	I64x2LaneSelect
	F32x4RelaxedMin
	F32x4RelaxedMax
	F64x2RelaxedMin
	F64x2RelaxedMax
	I16x8RelaxedQ15mulrS
	I16x8DotI8x16I7x16S
	I32x4DotI8x16I7x16AddS
	F32x4RelaxedDotBf16x8AddF32x4
)

var opCodeMap map[Opcode]uint32 = nil

//export metering_delegate
func metering_delegate(op C.wasmer_parser_operator_t) C.uint64_t {
	// a simple algorithm for now just map from opcode to cost directly
	// all the responsibility is placed on the caller of PushMeteringMiddleware
	v, b := opCodeMap[Opcode(op)]
	if !b {
		return 0 // no value means no cost
	}
	return C.uint64_t(v)
}

func getPlatformLong(v uint64) C.ulong {
	return C.ulong(v)
}

// Units of WebAssembly pages (as specified to be 65,536 bytes).
type Pages C.wasm_memory_pages_t

// Represents a memory page size.
const WasmPageSize = uint(0x10000)

// Represents the maximum number of pages.
const WasmMaxPages = uint(0x10000)

// Represents the minimum number of pages.
const WasmMinPages = uint(0x100)

// ToUint32 converts a Pages to a native Go uint32 which is the Pages' size.
//
//	memory, _ := instance.Exports.GetMemory("exported_memory")
//	size := memory.Size().ToUint32()
func (p *Pages) ToUint32() uint32 {
	return uint32(C.wasm_memory_pages_t(*p))
}

// ToBytes converts a Pages to a native Go uint which is the Pages' size in bytes.
//
//	memory, _ := instance.Exports.GetMemory("exported_memory")
//	size := memory.Size().ToBytes()
func (p *Pages) ToBytes() uint {
	return uint(p.ToUint32()) * WasmPageSize
}

// ValueKind represents the kind of a value.
type ValueKind C.wasm_valkind_t

const (
	// A 32-bit integer. In WebAssembly, integers are
	// sign-agnostic, i.E. this can either be signed or unsigned.
	I32 = ValueKind(C.WASM_I32)
	// A 64-bit integer. In WebAssembly, integers are
	// sign-agnostic, i.E. this can either be signed or unsigned.
	I64 = ValueKind(C.WASM_I64)
	// A 32-bit float.
	F32 = ValueKind(C.WASM_F32)
	// A 64-bit float.
	F64 = ValueKind(C.WASM_F64)
	// An externref value which can hold opaque data to the
	// WebAssembly instance itself.
	AnyRef = ValueKind(C.WASM_ANYREF)
	// A first-class reference to a WebAssembly function.
	FuncRef = ValueKind(C.WASM_FUNCREF)
)

// String returns the ValueKind as a string.
//
//	I32.String()     // "i32"
//	I64.String()     // "i64"
//	F32.String()     // "f32"
//	F64.String()     // "f64"
//	AnyRef.String()  // "anyref"
//	FuncRef.String() // "funcref"
func (v ValueKind) String() string {
	switch v {
	case I32:
		return "i32"
	case I64:
		return "i64"
	case F32:
		return "f32"
	case F64:
		return "f64"
	case AnyRef:
		return "anyref"
	case FuncRef:
		return "funcref"
	}

	panic("Unknown value kind")
}

// IsNumber returns true if the ValueKind is a number type.
//
//	I32.IsNumber()     // true
//	I64.IsNumber()     // true
//	F32.IsNumber()     // true
//	F64.IsNumber()     // true
//	AnyRef.IsNumber()  // false
//	FuncRef.IsNumber() // false
func (v ValueKind) IsNumber() bool {
	return bool(C.wasm_valkind_is_num(C.wasm_valkind_t(v)))
}

// IsReference returns true if the ValueKind is a reference.
//
//	I32.IsReference()     // false
//	I64.IsReference()     // false
//	F32.IsReference()     // false
//	F64.IsReference()     // false
//	AnyRef.IsReference()  // true
//	FuncRef.IsReference() // true
func (v ValueKind) IsReference() bool {
	return bool(C.wasm_valkind_is_ref(C.wasm_valkind_t(v)))
}

func (v ValueKind) inner() C.wasm_valkind_t {
	return C.wasm_valkind_t(v)
}

// ValueType classifies the individual values that WebAssembly code
// can compute with and the values that a variable accepts.
type ValueType struct {
	_inner   *C.wasm_valtype_t
	_ownedBy interface{}
}

// NewValueType instantiates a new ValueType given a ValueKind.
//
//	valueType := NewValueType(I32)
func NewValueType(kind ValueKind) *ValueType {
	pointer := C.wasm_valtype_new(C.wasm_valkind_t(kind))

	return newValueType(pointer, nil)
}

func newValueType(pointer *C.wasm_valtype_t, ownedBy interface{}) *ValueType {
	valueType := &ValueType{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(valueType, func(valueType *ValueType) {
			C.wasm_valtype_delete(valueType.inner())
		})
	}

	return valueType
}

func (vt *ValueType) inner() *C.wasm_valtype_t {
	return vt._inner
}

// Kind returns the ValueType's ValueKind
//
//	valueType := NewValueType(I32)
//	_ = valueType.Kind()
func (vt *ValueType) Kind() ValueKind {
	kind := ValueKind(C.wasm_valtype_kind(vt.inner()))

	runtime.KeepAlive(vt)

	return kind
}

// NewValueTypes instantiates a new ValueType array from a list of
// ValueKind. Note that this list may be empty.
//
//	valueTypes := NewValueTypes(I32, I64, F32)
//
// Note:️ NewValueTypes is specifically designed to help you declare
// function types, e.g. with NewFunctionType:
//
//	functionType := NewFunctionType(
//		NewValueTypes(), // arguments
//		NewValueTypes(I32), // results
//	)
func NewValueTypes(kinds ...ValueKind) []*ValueType {
	valueTypes := make([]*ValueType, len(kinds))

	for nth, kind := range kinds {
		valueTypes[nth] = NewValueType(kind)
	}

	return valueTypes
}

func toValueTypeVec(valueTypes []*ValueType) C.wasm_valtype_vec_t {
	vec := C.wasm_valtype_vec_t{}
	C.wasm_valtype_vec_new_uninitialized(&vec, C.size_t(len(valueTypes)))

	firstValueTypePointer := unsafe.Pointer(vec.data)

	for nth, valueType := range valueTypes {
		pointer := C.wasm_valtype_new(C.wasm_valtype_kind(valueType.inner()))
		*(**C.wasm_valtype_t)(unsafe.Pointer(uintptr(firstValueTypePointer) + unsafe.Sizeof(pointer)*uintptr(nth))) = pointer
	}

	runtime.KeepAlive(valueTypes)

	return vec
}

func toValueTypeList(valueTypes *C.wasm_valtype_vec_t, ownedBy interface{}) []*ValueType {
	numberOfValueTypes := int(valueTypes.size)
	list := make([]*ValueType, numberOfValueTypes)
	firstValueType := unsafe.Pointer(valueTypes.data)
	sizeOfValueTypePointer := unsafe.Sizeof(firstValueType)

	var currentValueTypePointer *C.wasm_valtype_t

	for nth := 0; nth < numberOfValueTypes; nth++ {
		currentValueTypePointer = *(**C.wasm_valtype_t)(unsafe.Pointer(uintptr(firstValueType) + uintptr(nth)*sizeOfValueTypePointer))
		valueType := newValueType(currentValueTypePointer, ownedBy)
		list[nth] = valueType
	}

	return list
}

// Value; WebAssembly computations manipulate values of basic value types:
//
// • Integer (32 or 64 bit width),
//
// • Floating-point (32 or 64 bit width),
//
// • Vectors (128 bits, with 32 or 64 bit lanes).
//
// # See Also
//
// Specification: https://webassembly.github.io/spec/core/exec/runtime.html#values
type Value struct {
	_inner *C.wasm_val_t
}

func newValue(pointer *C.wasm_val_t) Value {
	return Value{_inner: pointer}
}

// NewValue instantiates a new Value with the given value and
// ValueKind.
//
// Note: If a Wasm value cannot be created from the given value,
//
//	value := NewValue(42, I32)
func NewValue(value interface{}, kind ValueKind) Value {
	output, err := fromGoValue(value, kind)

	if err != nil {
		panic(fmt.Sprintf("Cannot create a Wasm `%s` value from `%T`", err, value))
	}

	return newValue(&output)
}

// NewI32 instantiates a new I32 Value with the given value.
//
// Note: If a Wasm value cannot be created from the given value,
// NewI32 will panic.
//
//	value := NewI32(42)
func NewI32(value interface{}) Value {
	return NewValue(value, I32)
}

// NewI64 instantiates a new I64 Value with the given value.
//
// Note: If a Wasm value cannot be created from the given value,
// NewI64 will panic.
//
//	value := NewI64(42)
func NewI64(value interface{}) Value {
	return NewValue(value, I64)
}

// NewF32 instantiates a new F32 Value with the given value.
//
// Note: If a Wasm value cannot be created from the given value,
// NewF32 will panic.
//
//	value := NewF32(4.2)
func NewF32(value interface{}) Value {
	return NewValue(value, F32)
}

// NewF64 instantiates a new F64 Value with the given value.
//
// Note: If a Wasm value cannot be created from the given value,
// NewF64 will panic.
//
//	value := NewF64(4.2)
func NewF64(value interface{}) Value {
	return NewValue(value, F64)
}

func (v *Value) inner() *C.wasm_val_t {
	return v._inner
}

// Kind returns the Value's ValueKind.
//
//	value := NewF64(4.2)
//	_ = value.Kind()
func (v *Value) Kind() ValueKind {
	return ValueKind(v.inner().kind)
}

// Unwrap returns the Value's value as a native Go value.
//
//	value := NewF64(4.2)
//	_ = value.Unwrap()
func (v *Value) Unwrap() interface{} {
	return toGoValue(v.inner())
}

// I32 returns the Value's value as a native Go int32.
//
// Note: It panics if the value is not of type I32.
//
//	value := NewI32(42)
//	_ = value.I32()
func (v *Value) I32() int32 {
	pointer := v.inner()
	if ValueKind(pointer.kind) != I32 {
		panic("Cannot convert value to `int32`")
	}

	return int32(C.to_int32(pointer))
}

// I64 returns the Value's value as a native Go int64.
//
// Note: It panics if the value is not of type I64.
//
//	value := NewI64(42)
//	_ = value.I64()
func (v *Value) I64() int64 {
	pointer := v.inner()
	if ValueKind(pointer.kind) != I64 {
		panic("Cannot convert value to `int64`")
	}

	return int64(C.to_int64(pointer))
}

// F32 returns the Value's value as a native Go float32.
//
// Note: It panics if the value is not of type F32.
//
//	value := NewF32(4.2)
//	_ = value.F32()
func (v *Value) F32() float32 {
	pointer := v.inner()
	if ValueKind(pointer.kind) != F32 {
		panic("Cannot convert value to `float32`")
	}

	return float32(C.to_float32(pointer))
}

// F64 returns the Value's value as a native Go float64.
//
// Note: It panics if the value is not of type F64.
//
//	value := NewF64(4.2)
//	_ = value.F64()
func (v *Value) F64() float64 {
	pointer := v.inner()
	if ValueKind(pointer.kind) != F64 {
		panic("Cannot convert value to `float64`")
	}

	return float64(C.to_float64(pointer))
}

func toGoValue(pointer *C.wasm_val_t) interface{} {
	switch ValueKind(pointer.kind) {
	case I32:
		return int32(C.to_int32(pointer))
	case I64:
		return int64(C.to_int64(pointer))
	case F32:
		return float32(C.to_float32(pointer))
	case F64:
		return float64(C.to_float64(pointer))
	default:
		panic("to do `newValue`")
	}
}

func fromGoValue(value interface{}, kind ValueKind) (C.wasm_val_t, error) {
	output := C.wasm_val_t{}

	switch kind {
	case I32:
		output.kind = kind.inner()

		var of = (*int32)(unsafe.Pointer(&output.of))

		switch v := value.(type) {
		case int8:
			*of = int32(v)
		case uint8:
			*of = int32(v)
		case int16:
			*of = int32(v)
		case uint16:
			*of = int32(v)
		case int32:
			*of = v
		case int:
			*of = int32(v)
		case uint:
			*of = int32(v)
		default:
			return output, newErrorWith("i32")
		}
	case I64:
		output.kind = kind.inner()

		var of = (*int64)(unsafe.Pointer(&output.of))

		switch v := value.(type) {
		case int8:
			*of = int64(v)
		case uint8:
			*of = int64(v)
		case int16:
			*of = int64(v)
		case uint16:
			*of = int64(v)
		case int32:
			*of = int64(v)
		case uint32:
			*of = int64(v)
		case int64:
			*of = v
		case int:
			*of = int64(v)
		case uint:
			*of = int64(v)
		default:
			return output, newErrorWith("i64")
		}
	case F32:
		output.kind = kind.inner()

		var of = (*float32)(unsafe.Pointer(&output.of))

		switch v := value.(type) {
		case float32:
			*of = v
		default:
			return output, newErrorWith("f32")
		}
	case F64:
		output.kind = kind.inner()

		var of = (*float64)(unsafe.Pointer(&output.of))

		switch v := value.(type) {
		case float32:
			*of = float64(v)
		case float64:
			*of = v
		default:
			return output, newErrorWith("f64")
		}
	default:
		panic("To do, `fromGoValue`!")
	}

	return output, nil
}

func toValueVec(list []Value, vec *C.wasm_val_vec_t) {
	numberOfValues := len(list)
	values := make([]C.wasm_val_t, numberOfValues)

	for nth, item := range list {
		value, err := fromGoValue(item.Unwrap(), item.Kind())
		if err != nil {
			panic(err)
		}

		values[nth] = value
	}

	if numberOfValues > 0 {
		C.wasm_val_vec_new(vec, C.size_t(numberOfValues), (*C.wasm_val_t)(unsafe.Pointer(&values[0])))
	}
}

func toValueList(values *C.wasm_val_vec_t) []Value {
	numberOfValues := int(values.size)
	list := make([]Value, numberOfValues)
	firstValue := unsafe.Pointer(values.data)
	sizeOfValuePointer := unsafe.Sizeof(C.wasm_val_t{})

	var currentValuePointer *C.wasm_val_t

	for nth := 0; nth < numberOfValues; nth++ {
		currentValuePointer = (*C.wasm_val_t)(unsafe.Pointer(uintptr(firstValue) + uintptr(nth)*sizeOfValuePointer))
		value := newValue(currentValuePointer)
		list[nth] = value
	}

	return list
}

// Trap stores trace message with backtrace when an error happened.
type Trap struct {
	_inner   *C.wasm_trap_t
	_ownedBy interface{}
}

// newWasmTrap creates C wasm_trap structure and returns it's pointer
func newWasmTrap(store *Store, message string) *C.wasm_trap_t {
	messageBytes := []byte(message)
	var bytesPointer *C.uint8_t
	bytesLength := len(messageBytes)

	if bytesLength > 0 {
		bytesPointer = (*C.uint8_t)(unsafe.Pointer(&messageBytes[0]))
	}

	runtime.KeepAlive(store)
	runtime.KeepAlive(message)

	return C.to_wasm_trap_new(store.inner(), bytesPointer, C.size_t(bytesLength))
}

// newTrapFromPointer creates Trap from C.wasm_trap structure and attaches GC finalizer for it
func newTrapFromPointer(pointer *C.wasm_trap_t, ownedBy interface{}) *Trap {
	trap := &Trap{
		_inner:   pointer,
		_ownedBy: ownedBy,
	}

	return trapWithFinalizer(trap, ownedBy)
}

// trapWithFinalizer attaches GC finalizer to the trap
func trapWithFinalizer(trap *Trap, ownedBy interface{}) *Trap {
	if ownedBy == nil {
		runtime.SetFinalizer(trap, func(trap *Trap) {
			inner := trap.inner()

			if inner != nil {
				C.wasm_trap_delete(inner)
			}
		})
	}

	return trap
}

// Creates a new trap with a message.
//
//	engine := wasmer.NewEngine()
//	store := wasmer.NewStore(engine)
//	trap := NewTrap(store, "oops")
func NewTrap(store *Store, message string) *Trap {
	pointer := newWasmTrap(store, message)
	trap := &Trap{_inner: pointer}
	return trapWithFinalizer(trap, nil)
}

func (t *Trap) inner() *C.wasm_trap_t {
	return t._inner
}

func (t *Trap) ownedBy() interface{} {
	if t._ownedBy == nil {
		return t
	}

	return t._ownedBy
}

// Message returns the message attached to the current Trap.
func (t *Trap) Message() string {
	var bytes C.wasm_byte_vec_t
	C.wasm_trap_message(t.inner(), &bytes)

	runtime.KeepAlive(t)

	goBytes := C.GoBytes(unsafe.Pointer(bytes.data), C.int(bytes.size)-1)
	C.wasm_byte_vec_delete(&bytes)

	return string(goBytes)
}

// Origin returns the top frame of WebAssembly stack responsible for
// this trap.
//
//	frame := trap.Origin()
func (t *Trap) Origin() *Frame {
	frame := C.wasm_trap_origin(t.inner())

	runtime.KeepAlive(t)

	if frame == nil {
		return nil
	}

	return newFrame(frame, t.ownedBy())
}

// Trace returns the trace of WebAssembly frames for this trap.
func (t *Trap) Trace() *Trace {
	return newTrace(t)
}

// Frame represents a frame of a WebAssembly stack trace.
type Frame struct {
	_inner   *C.wasm_frame_t
	_ownedBy interface{}
}

func newFrame(pointer *C.wasm_frame_t, ownedBy interface{}) *Frame {
	frame := &Frame{
		_inner:   pointer,
		_ownedBy: ownedBy,
	}

	if ownedBy == nil {
		runtime.SetFinalizer(frame, func(frame *Frame) {
			C.wasm_frame_delete(frame.inner())
		})
	}

	return frame
}

func (f *Frame) inner() *C.wasm_frame_t {
	return f._inner
}

func (f *Frame) ownedBy() interface{} {
	if f._ownedBy == nil {
		return f
	}

	return f._ownedBy
}

// FunctionIndex returns the function index in the original
// WebAssembly module that this frame corresponds to.
func (f *Frame) FunctionIndex() uint32 {
	index := C.wasm_frame_func_index(f.inner())

	runtime.KeepAlive(f)

	return uint32(index)
}

// FunctionOffset returns the byte offset from the beginning of the
// function in the original WebAssembly file to the instruction this
// frame points to.
func (f *Frame) FunctionOffset() uint {
	index := C.wasm_frame_func_offset(f.inner())

	runtime.KeepAlive(f)

	return uint(index)
}

func (f *Frame) Instance() {
	//TODO: See https://github.com/wasmerio/wasmer/blob/6fbc903ea32774c830fd9ee86140d1406ac5d745/lib/c-api/src/wasm_c_api/types/frame.rs#L31-L34
	panic("to do!")
}

// ModuleOffset returns the byte offset from the beginning of the
// original WebAssembly file to the instruction this frame points to.
func (f *Frame) ModuleOffset() uint {
	index := C.wasm_frame_module_offset(f.inner())
	runtime.KeepAlive(f)
	return uint(index)
}

// Trace represents a WebAssembly trap.
type Trace struct {
	_inner C.wasm_frame_vec_t
	frames []*Frame
}

func newTrace(trap *Trap) *Trace {
	var trace = &Trace{}
	C.wasm_trap_trace(trap.inner(), trace.inner())

	runtime.KeepAlive(trap)
	runtime.SetFinalizer(trace, func(t *Trace) {
		C.wasm_frame_vec_delete(t.inner())
	})

	numberOfFrames := int(trace.inner().size)
	frames := make([]*Frame, numberOfFrames)
	firstFrame := unsafe.Pointer(trace.inner().data)
	sizeOfFramePointer := unsafe.Sizeof(firstFrame)

	var currentFramePointer **C.wasm_frame_t

	for nth := 0; nth < numberOfFrames; nth++ {
		currentFramePointer = (**C.wasm_frame_t)(unsafe.Pointer(uintptr(firstFrame) + uintptr(nth)*sizeOfFramePointer))
		frames[nth] = newFrame(*currentFramePointer, trace)
	}

	trace.frames = frames
	return trace
}

func (t *Trace) inner() *C.wasm_frame_vec_t {
	return &t._inner
}

// Target represents a triple + CPU features pairs.
type Target struct {
	_inner *C.wasmer_target_t
}

func newTarget(target *C.wasmer_target_t) *Target {
	t := &Target{
		_inner: target,
	}

	runtime.SetFinalizer(t, func(t *Target) {
		C.wasmer_target_delete(t.inner())
	})

	return t
}

// NewTarget creates a new target.
//
//	triple, err := NewTriple("aarch64-unknown-linux-gnu")
//	cpuFeatures := NewCpuFeatures()
//	target := NewTarget(triple, cpuFeatures)
func NewTarget(triple *Triple, cpuFeatures *CpuFeatures) *Target {
	return newTarget(C.wasmer_target_new(triple.inner(), cpuFeatures.inner()))
}

func (t *Target) inner() *C.wasmer_target_t {
	return t._inner
}

// Triple; historically such things had three fields, though they have
// added additional fields over time.
type Triple struct {
	_inner *C.wasmer_triple_t
}

func newTriple(triple *C.wasmer_triple_t) *Triple {
	t := &Triple{
		_inner: triple,
	}

	runtime.SetFinalizer(t, func(t *Triple) {
		C.wasmer_triple_delete(t.inner())
	})

	return t
}

// NewTriple creates a new triple, otherwise it returns an error
// specifying why the provided triple isn't valid.
//
//	triple, err := NewTriple("aarch64-unknown-linux-gnu")
func NewTriple(triple string) (*Triple, error) {
	cTripleName := newName(triple)
	defer C.wasm_name_delete(&cTripleName)

	var cTriple *C.wasmer_triple_t

	err := maybeNewErrorFromWasmer(func() bool {
		cTriple := C.wasmer_triple_new(&cTripleName)

		return cTriple == nil
	})

	if err != nil {
		return nil, err
	}

	return newTriple(cTriple), nil
}

// NewTripleFromHost creates a new triple from the current host.
func NewTripleFromHost() *Triple {
	return newTriple(C.wasmer_triple_new_from_host())
}

func (t *Triple) inner() *C.wasmer_triple_t {
	return t._inner
}

// CpuFeatures holds a set of CPU features. They are identified by
// their stringified names. The reference is the GCC options:
//
// • https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html,
//
// • https://gcc.gnu.org/onlinedocs/gcc/ARM-Options.html,
//
// • https://gcc.gnu.org/onlinedocs/gcc/AArch64-Options.html.
//
// At the time of writing this documentation (it might be outdated in
// the future), the supported features are the following:
//
// • sse2,
//
// • sse3,
//
// • ssse3,
//
// • sse4.1,
//
// • sse4.2,
//
// • popcnt,
//
// • avx,
//
// • bmi,
//
// • bmi2,
//
// • avx2,
//
// • avx512dq,
//
// • avx512vl,
//
// • lzcnt.
type CpuFeatures struct {
	_inner *C.wasmer_cpu_features_t
}

func newCpuFeatures(cpu_features *C.wasmer_cpu_features_t) *CpuFeatures {
	features := &CpuFeatures{
		_inner: cpu_features,
	}

	runtime.SetFinalizer(features, func(features *CpuFeatures) {
		C.wasmer_cpu_features_delete(features.inner())
	})

	return features
}

// NewCpuFeatures creates a new CpuFeatures, which is a set of CPU
// features.
func NewCpuFeatures() *CpuFeatures {
	return newCpuFeatures(C.wasmer_cpu_features_new())
}

// Add adds a new CPU feature to the existing set.
func (cf *CpuFeatures) Add(feature string) error {
	cFeature := newName(feature)
	defer C.wasm_name_delete(&cFeature)

	err := maybeNewErrorFromWasmer(func() bool {
		return !bool(C.wasmer_cpu_features_add(cf.inner(), &cFeature))
	})

	if err != nil {
		return err
	}

	return nil
}

func (cf *CpuFeatures) inner() *C.wasmer_cpu_features_t {
	return cf._inner
}

// Engine is used by the Store to drive the compilation and the
// execution of a WebAssembly module.
type Engine struct {
	_inner *C.wasm_engine_t
}

func newEngine(engine *C.wasm_engine_t) *Engine {
	e := &Engine{
		_inner: engine,
	}

	// When the Go object is garbage collected,
	// instruct the runtime to delete the inner engine.
	runtime.SetFinalizer(e, func(engine *Engine) {
		C.wasm_engine_delete(engine._inner)
	})

	return e
}

// NewEngine instantiates and returns a new Engine with the default configuration.
func NewEngine() *Engine {
	return newEngine(C.wasm_engine_new())
}

// NewEngineWithConfig instantiates and returns a new Engine with the given configuration.
func NewEngineWithConfig(config *Config) *Engine {
	return newEngine(C.wasm_engine_new_with_config(config._inner))
}

// NewUniversalEngine instantiates and returns a new Engine with the Universal engine.
func NewUniversalEngine() *Engine {
	config := NewConfig()
	config.UseUniversalEngine()

	return NewEngineWithConfig(config)
}

// CompilerKind represents the possible compiler types.
type CompilerKind C.wasmer_compiler_t

const (
	// Represents the Cranelift compiler.
	CRANELIFT = CompilerKind(C.CRANELIFT)
	// Represents the LLVM compiler.
	LLVM = CompilerKind(C.LLVM)
	// Represents the Singlepass compiler.
	SINGLEPASS = CompilerKind(C.SINGLEPASS)
)

// Strings returns the CompilerKind as a string.
func (c CompilerKind) String() string {
	switch c {
	case CRANELIFT:
		return "cranelift"
	case LLVM:
		return "llvm"
	case SINGLEPASS:
		return "singlepass"
	}

	panic("Unknown compiler")
}

// IsCompilerAvailable checks that the given compiler is available for this platform.
func IsCompilerAvailable(compiler CompilerKind) bool {
	return bool(C.wasmer_is_compiler_available(uint32(C.wasmer_compiler_t(compiler))))
}

// EngineKind represents the possible engine types.
type EngineKind C.wasmer_engine_t

const (
	// Represents the Universal engine.
	UNIVERSAL = EngineKind(C.UNIVERSAL)
)

// String returns the EngineKind as a string.
func (e EngineKind) String() string {
	switch e {
	case UNIVERSAL:
		return "universal"
	}

	panic("Unknown Wasmer engine")
}

// IsEngineAvailable checks that the given engine is available in this platform.
func IsEngineAvailable(engine EngineKind) bool {
	return bool(C.wasmer_is_engine_available(uint32(C.wasmer_engine_t(engine))))
}

// Config holds the compiler and the Engine used by the Store.
type Config struct {
	_inner *C.wasm_config_t
}

// NewConfig instantiates and returns a new Config.
// It uses the Universal engine by default.
func NewConfig() *Config {
	config := C.wasm_config_new()
	if !IsEngineAvailable(UNIVERSAL) {
		panic("This LXD version doesn't include the Wasmer Universal engine")
	}

	C.wasm_config_set_engine(config, uint32(C.wasmer_engine_t(UNIVERSAL)))
	return &Config{
		_inner: config,
	}
}

// UseNativeEngine sets the engine to Universal in the configuration.
func (c *Config) UseUniversalEngine() *Config {
	if !IsEngineAvailable(UNIVERSAL) {
		panic("This LXD version doesn't include the Universal engine")
	}

	C.wasm_config_set_engine(c._inner, uint32(C.wasmer_engine_t(UNIVERSAL)))
	return c
}

// UseCraneliftCompiler sets the compiler to Cranelift.
// Returns the Config itself in order to chain calls.
func (c *Config) UseCraneliftCompiler() *Config {
	if !IsCompilerAvailable(CRANELIFT) {
		panic("This LXD version doesn't include the Cranelift compiler for Wasmer")
	}

	C.wasm_config_set_compiler(c._inner, uint32(C.wasmer_compiler_t(CRANELIFT)))
	return c
}

// UseLLVMCompiler sets the compiler to LLVM.
// Returns the Config itself in order to chain calls.
func (c *Config) UseLLVMCompiler() *Config {
	if !IsCompilerAvailable(LLVM) {
		panic("This LXD version doesn't include the LLVM compiler for Wasmer")
	}

	C.wasm_config_set_compiler(c._inner, uint32(C.wasmer_compiler_t(LLVM)))
	return c
}

// UseSinglepassCompiler sets the compiler to Singlepass.
// Returns the Config itself in order to chain calls.
func (c *Config) UseSinglepassCompiler() *Config {
	if !IsCompilerAvailable(SINGLEPASS) {
		panic("This LXD version doesn't include the Singlepass compiler for Wasmer")
	}

	C.wasm_config_set_compiler(c._inner, uint32(C.wasmer_compiler_t(SINGLEPASS)))
	return c
}

// UseTarget sets a specific target for doing cross-compilation.
// Returns the Config itself in order to chain calls.
func (c *Config) UseTarget(target *Target) *Config {
	C.wasm_config_set_target(c._inner, target._inner)
	return c
}

func (c *Config) inner() *C.wasm_config_t {
	return c._inner
}

// PushMeteringMiddleware allows the middleware metering to be engaged on a map of opcode to cost
//
//	  config := NewConfig()
//		 opmap := map[uint32]uint32{
//			End: 		1,
//			LocalGet: 	1,
//			I32Add: 	4,
//		 }
//	  config.PushMeteringMiddleware(7865444, opmap)
func (c *Config) PushMeteringMiddleware(maxGasUsageAllowed uint64, opMap map[Opcode]uint32) *Config {
	if opCodeMap == nil {
		// REVIEW only allowing this to be set once
		opCodeMap = opMap
	}

	C.wasm_config_push_middleware(c.inner(), C.wasmer_metering_as_middleware(C.wasmer_metering_new(getPlatformLong(maxGasUsageAllowed), (*[0]byte)(C.metering_delegate))))
	return c
}

// PushMeteringMiddlewarePtr allows the middleware metering to be engaged on an unsafe.Pointer
// this pointer must be a to C based function with a signature of:
//
//	extern uint64_t cost_delegate_func(enum wasmer_parser_operator_t op);
//
// package main
//
// #include <wasmer.h>
// extern uint64_t metering_delegate_alt(enum wasmer_parser_operator_t op);
// import "C"
// import "unsafe"
//
//	func getInternalCPointer() unsafe.Pointer {
//		  return unsafe.Pointer(C.metering_delegate_alt)
//	}
//
// //export metering_delegate_alt
//
//	func metering_delegate_alt(op C.wasmer_parser_operator_t) C.uint64_t {
//		v, b := opCodeMap[Opcode(op)]
//	  if !b {
//		   return 0 // no value means no cost
//	  }
//	  return C.uint64_t(v)
//	}
//
//	void main(){
//	   config := NewConfig()
//	   config.PushMeteringMiddlewarePtr(800000000, getInternalCPointer())
//	}
func (c *Config) PushMeteringMiddlewarePtr(maxGasUsageAllowed uint64, p unsafe.Pointer) *Config {
	C.wasm_config_push_middleware(c.inner(), C.wasmer_metering_as_middleware(C.wasmer_metering_new(getPlatformLong(maxGasUsageAllowed), (*[0]byte)(p))))
	return c
}

// Error represents a Wasmer runtime error.
type Error struct {
	message string
}

func newErrorWith(message string) *Error {
	return &Error{
		message: message,
	}
}

func _newErrorFromWasmer() *Error {
	var errorLength = C.wasmer_last_error_length()
	if errorLength == 0 {
		return newErrorWith("(no error from Wasmer)")
	}

	var errorMessage = make([]C.char, errorLength)
	var errorMessagePointer = (*C.char)(unsafe.Pointer(&errorMessage[0]))
	var errorResult = C.wasmer_last_error_message(errorMessagePointer, errorLength)
	if errorResult == -1 {
		return newErrorWith("(failed to read last error from Wasmer)")
	}

	return newErrorWith(C.GoStringN(errorMessagePointer, errorLength-1))
}

func maybeNewErrorFromWasmer(block func() bool) *Error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if block() /* has failed */ {
		return _newErrorFromWasmer()
	}

	return nil
}

// Error returns the Error's message.
func (error *Error) Error() string {
	return error.message
}

// TrapError represents a trap produced during Wasm execution.
//
// # See also
//
// Specification: https://webassembly.github.io/spec/core/intro/overview.html#trap
type TrapError struct {
	message string
	origin  *Frame
	trace   []*Frame
}

func newErrorFromTrap(pointer *C.wasm_trap_t) *TrapError {
	trap := newTrapFromPointer(pointer, nil)

	return &TrapError{
		message: trap.Message(),
		origin:  trap.Origin(),
		trace:   trap.Trace().frames,
	}
}

// Error returns the TrapError's message.
func (te *TrapError) Error() string {
	return te.message
}

// Origin returns the TrapError's origin as a Frame.
func (te *TrapError) Origin() *Frame {
	return te.origin
}

// Trace returns the TrapError's trace as a Frame array.
func (te *TrapError) Trace() []*Frame {
	return te.trace
}

// TableSize represents the size of a table.
type TableSize C.wasm_table_size_t

// ToUint32 converts a TableSize to a native Go uint32.
//
//	table, _ := instance.Exports.GetTable("exported_table")
//	size := table.Size().ToUint32()
func (ts *TableSize) ToUint32() uint32 {
	return uint32(C.wasm_table_size_t(*ts))
}

// TableType classifies tables over elements of element types within a size range.
//
// # See also
//
// Specification: https://webassembly.github.io/spec/core/syntax/types.html#table-types
type TableType struct {
	_inner   *C.wasm_tabletype_t
	_ownedBy interface{}
}

func newTableType(pointer *C.wasm_tabletype_t, ownedBy interface{}) *TableType {
	tableType := &TableType{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(tableType, func(tableType *TableType) {
			C.wasm_tabletype_delete(tableType.inner())
		})
	}

	return tableType
}

// NewTableType instantiates a new TableType given a ValueType and some Limits.
//
//	valueType := NewValueType(I32)
//	limits := NewLimits(1, 4)
//	tableType := NewTableType(valueType, limits)
//	_ = tableType.IntoExternType()
func NewTableType(valueType *ValueType, limits *Limits) *TableType {
	pointer := C.wasm_tabletype_new(valueType.inner(), limits.inner())

	return newTableType(pointer, nil)
}

func (tt *TableType) inner() *C.wasm_tabletype_t {
	return tt._inner
}

func (tt *TableType) ownedBy() interface{} {
	if tt._ownedBy == nil {
		return tt
	}

	return tt._ownedBy
}

// ValueType returns the TableType's ValueType.
//
//	valueType := NewValueType(I32)
//	limits := NewLimits(1, 4)
//	tableType := NewTableType(valueType, limits)
//	_ = tableType.ValueType()
func (tt *TableType) ValueType() *ValueType {
	pointer := C.wasm_tabletype_element(tt.inner())
	runtime.KeepAlive(tt)
	return newValueType(pointer, tt.ownedBy())
}

// Limits returns the TableType's Limits.
//
//	valueType := NewValueType(I32)
//	limits := NewLimits(1, 4)
//	tableType := NewTableType(valueType, limits)
//	_ = tableType.Limits()
func (tt *TableType) Limits() *Limits {
	limits := newLimits(C.wasm_tabletype_limits(tt.inner()), tt.ownedBy())
	runtime.KeepAlive(tt)
	return limits
}

// IntoExternType converts the TableType into an ExternType.
//
//	valueType := NewValueType(I32)
//	limits := NewLimits(1, 4)
//	tableType := NewTableType(valueType, limits)
//	_ = tableType.IntoExternType()
func (tt *TableType) IntoExternType() *ExternType {
	pointer := C.wasm_tabletype_as_externtype_const(tt.inner())
	return newExternType(pointer, tt.ownedBy())
}

// A table instance is the runtime representation of a table. It holds
// a vector of function elements and an optional maximum size, if one
// was specified in the table type at the table’s definition site.
//
// # See also
//
// Specification: https://webassembly.github.io/spec/core/exec/runtime.html#table-instances
type Table struct {
	_inner   *C.wasm_table_t
	_ownedBy interface{}
}

func newTable(pointer *C.wasm_table_t, ownedBy interface{}) *Table {
	table := &Table{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(table, func(table *Table) {
			C.wasm_table_delete(table.inner())
		})
	}

	return table
}

func (t *Table) inner() *C.wasm_table_t {
	return t._inner
}

func (t *Table) ownedBy() interface{} {
	if t._ownedBy == nil {
		return t
	}

	return t._ownedBy
}

// Size returns the Table's size.
//
//	table, _ := instance.Exports.GetTable("exported_table")
//	size := table.Size()
func (t *Table) Size() TableSize {
	return TableSize(C.wasm_table_size(t.inner()))
}

// IntoExtern converts the Table into an Extern.
//
//	table, _ := instance.Exports.GetTable("exported_table")
//	extern := table.IntoExtern()
func (t *Table) IntoExtern() *Extern {
	pointer := C.wasm_table_as_extern(t.inner())

	return newExtern(pointer, t.ownedBy())
}

// Exports is a special kind of map that allows easily unwrapping the
// types of instances.
type Exports struct {
	_inner   C.wasm_extern_vec_t
	exports  map[string]*Extern
	instance *C.wasm_instance_t
}

func newExports(instance *C.wasm_instance_t, module *Module) *Exports {
	e := &Exports{}
	C.wasm_instance_exports(instance, &e._inner)

	runtime.SetFinalizer(e, func(e *Exports) {
		e.Close()
	})

	numberOfExports := int(e.inner().size)
	exports := make(map[string]*Extern, numberOfExports)
	firstExport := unsafe.Pointer(e.inner().data)
	sizeOfExportPointer := unsafe.Sizeof(firstExport)

	var currentExportPointer *C.wasm_extern_t

	moduleExports := module.Exports()

	for nth := 0; nth < numberOfExports; nth++ {
		currentExportPointer = *(**C.wasm_extern_t)(unsafe.Pointer(uintptr(firstExport) + uintptr(nth)*sizeOfExportPointer))
		export := newExtern(currentExportPointer, e)
		exports[moduleExports[nth].Name()] = export
	}

	e.exports = exports
	e.instance = instance

	return e
}

func (e *Exports) inner() *C.wasm_extern_vec_t {
	return &e._inner
}

// Get retrieves and returns an Extern by its name.
//
// Note: If the name does not refer to an existing export, Get will
// return an Error.
//
//	instance, _ := NewInstance(module, NewImportObject())
//	extern, error := instance.Exports.Get("an_export")
func (e *Exports) Get(name string) (*Extern, error) {
	export, exists := e.exports[name]

	if !exists {
		return nil, newErrorWith(fmt.Sprintf("Export `%s` does not exist", name))
	}

	return export, nil
}

// GetRawFunction retrieves and returns an exported Function by its name.
//
// Note: If the name does not refer to an existing export,
// GetRawFunction will return an Error.
//
// Note: If the export is not a function, GetRawFunction will return
// nil as its result.
//
//	instance, _ := NewInstance(module, NewImportObject())
//	exportedFunc, error := instance.Exports.GetRawFunction("an_exported_function")
//
//	if error != nil && exportedFunc != nil {
//	    exportedFunc.Call()
//	}
func (e *Exports) GetRawFunction(name string) (*Function, error) {
	exports, err := e.Get(name)

	if err != nil {
		return nil, err
	}

	return exports.IntoFunction(), nil
}

// GetFunction retrieves a exported function by its name and returns
// it as a native Go function.
//
// The difference with GetRawFunction is that Function.Native has been
// called on the exported function.
//
// Note: If the name does not refer to an existing export, GetFunction
// will return an Error.
//
// Note: If the export is not a function, GetFunction will return nil
// as its result.
//
//	instance, _ := NewInstance(module, NewImportObject())
//	exportedFunc, error := instance.Exports.GetFunction("an_exported_function")
//
//	if error != nil && exportedFunc != nil {
//	    exportedFunc()
//	}
func (e *Exports) GetFunction(name string) (NativeFunction, error) {
	function, err := e.GetRawFunction(name)

	if err != nil {
		return nil, err
	}

	return function.Native(), nil
}

// GetGlobal retrieves and returns a exported Global by its name.
//
// Note: If the name does not refer to an existing export, GetGlobal
// will return an Error.
//
// Note: If the export is not a global, GetGlobal will return nil as a
// result.
//
//	instance, _ := NewInstance(module, NewImportObject())
//	exportedGlobal, error := instance.Exports.GetGlobal("an_exported_global")
func (e *Exports) GetGlobal(name string) (*Global, error) {
	exports, err := e.Get(name)

	if err != nil {
		return nil, err
	}

	return exports.IntoGlobal(), nil
}

// GetTable retrieves and returns a exported Table by its name.
//
// Note: If the name does not refer to an existing export, GetTable
// will return an Error.
//
// Note: If the export is not a table, GetTable will return nil as a
// result.
//
//	instance, _ := NewInstance(module, NewImportObject())
//	exportedTable, error := instance.Exports.GetTable("an_exported_table")
func (e *Exports) GetTable(name string) (*Table, error) {
	exports, err := e.Get(name)

	if err != nil {
		return nil, err
	}

	return exports.IntoTable(), nil
}

// GetMemory retrieves and returns a exported Memory by its name.
//
// Note: If the name does not refer to an existing export, GetMemory
// will return an Error.
//
// Note: If the export is not a memory, GetMemory will return nil as a
// result.
//
//	instance, _ := NewInstance(module, NewImportObject())
//	exportedMemory, error := instance.Exports.GetMemory("an_exported_memory")
func (e *Exports) GetMemory(name string) (*Memory, error) {
	exports, err := e.Get(name)

	if err != nil {
		return nil, err
	}

	return exports.IntoMemory(), nil
}

// GetWasiStartFunction is similar to GetFunction("_start"). It saves
// you the cost of knowing the name of the WASI start function.
func (e *Exports) GetWasiStartFunction() (*Function, error) {
	start := C.wasi_get_start_function(e.instance)

	if start == nil {
		return nil, newErrorWith("WASI start function was not found")
	}

	return newFunction(start, nil, nil), nil
}

func CallFunction(f *Function) error {
	if f == nil {
		return newErrorWith("Function is nil")
	}

	trap := C.wasm_func_call(f.inner(), nil, nil)
	if trap != nil {
		return newErrorFromTrap(trap)
	}

	return nil
}

// Force to close the Exports.
//
// A runtime finalizer is registered on the Exports, but it is
// possible to force the destruction of the Exports by calling Close
// manually.
func (e *Exports) Close() {
	runtime.SetFinalizer(e, nil)

	for extern := range e.exports {
		delete(e.exports, extern)
	}

	C.wasm_extern_vec_delete(&e._inner)
}

type exportTypes struct {
	_inner      C.wasm_exporttype_vec_t
	exportTypes []*ExportType
}

func newExportTypes(module *Module) *exportTypes {
	et := &exportTypes{}
	C.wasm_module_exports(module.inner(), &et._inner)

	runtime.SetFinalizer(et, func(et *exportTypes) {
		et.close()
	})

	numberOfExportTypes := int(et.inner().size)
	types := make([]*ExportType, numberOfExportTypes)
	firstExportType := unsafe.Pointer(et.inner().data)
	sizeOfExportTypePointer := unsafe.Sizeof(firstExportType)

	var currentTypePointer *C.wasm_exporttype_t

	for nth := 0; nth < numberOfExportTypes; nth++ {
		currentTypePointer = *(**C.wasm_exporttype_t)(unsafe.Pointer(uintptr(firstExportType) + uintptr(nth)*sizeOfExportTypePointer))
		exportType := newExportType(currentTypePointer, et)
		types[nth] = exportType
	}

	et.exportTypes = types

	return et
}

func (et *exportTypes) inner() *C.wasm_exporttype_vec_t {
	return &et._inner
}

func (et *exportTypes) close() {
	runtime.SetFinalizer(et, nil)
	C.wasm_exporttype_vec_delete(&et._inner)
}

// ExportType is a descriptor for an exported WebAssembly value.
type ExportType struct {
	_inner   *C.wasm_exporttype_t
	_ownedBy interface{}
}

func newExportType(pointer *C.wasm_exporttype_t, ownedBy interface{}) *ExportType {
	exportType := &ExportType{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(exportType, func(et *ExportType) {
			et.Close()
		})
	}

	return exportType
}

func newName(str string) C.wasm_name_t {
	var name C.wasm_name_t

	C.wasm_name_new_from_string(&name, C.CString(str))

	runtime.KeepAlive(str)

	return name
}

func nameToString(name *C.wasm_name_t) string {
	if name.data == nil {
		return ""
	}

	return C.GoStringN(name.data, C.int(name.size))
}

// NewExportType instantiates a new ExportType with a name and an extern type.
//
// Note: An extern type is anything implementing IntoExternType: FunctionType, GlobalType, MemoryType, TableType.
//
//	valueType := NewValueType(I32)
//	globalType := NewGlobalType(valueType, CONST)
//	exportType := NewExportType("a_global", globalType)
func NewExportType(name string, ty IntoExternType) *ExportType {
	nameName := newName(name)
	externType := ty.IntoExternType().inner()
	externTypeCopy := C.wasm_externtype_copy(externType)

	runtime.KeepAlive(externType)

	exportType := C.wasm_exporttype_new(&nameName, externTypeCopy)

	return newExportType(exportType, nil)
}

func (et *ExportType) inner() *C.wasm_exporttype_t {
	return et._inner
}

func (et *ExportType) ownedBy() interface{} {
	if et._ownedBy == nil {
		return et
	}

	return et._ownedBy
}

// Name returns the name of the export type.
//
//	exportType := NewExportType("a_global", globalType)
//	exportType.Name() // "global"
func (et *ExportType) Name() string {
	byteVec := C.wasm_exporttype_name(et.inner())
	name := C.GoStringN(byteVec.data, C.int(byteVec.size))
	runtime.KeepAlive(et)
	return name
}

// Type returns the type of the export type.
//
//	exportType := NewExportType("a_global", globalType)
//	exportType.Type() // ExternType
func (et *ExportType) Type() *ExternType {
	ty := C.wasm_exporttype_type(et.inner())
	runtime.KeepAlive(et)
	return newExternType(ty, et.ownedBy())
}

// Force to close the ExportType.
//
// A runtime finalizer is registered on the ExportType, but it is
// possible to force the destruction of the ExportType by calling
// Close manually.
func (et *ExportType) Close() {
	runtime.SetFinalizer(et, nil)
	C.wasm_exporttype_delete(et.inner())
}

// Extern is the runtime representation of an entity that can be
// imported or exported.
type Extern struct {
	_inner   *C.wasm_extern_t
	_ownedBy interface{}
}

// IntoExtern is an interface implemented by entity that can be
// imported of exported.
type IntoExtern interface {
	IntoExtern() *Extern
}

func newExtern(pointer *C.wasm_extern_t, ownedBy interface{}) *Extern {
	extern := &Extern{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(extern, func(extern *Extern) {
			C.wasm_extern_delete(extern.inner())
		})
	}

	return extern
}

func (e *Extern) inner() *C.wasm_extern_t {
	return e._inner
}

func (e *Extern) ownedBy() interface{} {
	if e._ownedBy == nil {
		return e
	}

	return e._ownedBy
}

func (e *Extern) IntoExtern() *Extern {
	return e
}

// Kind returns the Extern's ExternKind.
//
//	global, _ := instance.Exports.GetGlobal("exported_global")
//	_ = global.IntoExtern().Kind()
func (e *Extern) Kind() ExternKind {
	kind := ExternKind(C.wasm_extern_kind(e.inner()))

	runtime.KeepAlive(e)

	return kind
}

// Type returns the Extern's ExternType.
//
//	global, _ := instance.Exports.GetGlobal("exported_global")
//	_ = global.IntoExtern().Type()
func (e *Extern) Type() *ExternType {
	ty := C.wasm_extern_type(e.inner())

	runtime.KeepAlive(e)

	return newExternType(ty, e.ownedBy())
}

// IntoFunction converts the Extern into a Function.
//
// Note:️ If the Extern is not a Function, IntoFunction will return nil
// as its result.
//
//	function, _ := instance.Exports.GetFunction("exported_function")
//	extern = function.IntoExtern()
//	_ := extern.IntoFunction()
func (e *Extern) IntoFunction() *Function {
	pointer := C.wasm_extern_as_func(e.inner())

	if pointer == nil {
		return nil
	}

	return newFunction(pointer, nil, e.ownedBy())
}

// IntoGlobal converts the Extern into a Global.
//
// Note:️ If the Extern is not a Global, IntoGlobal will return nil as
// its result.
//
//	global, _ := instance.Exports.GetGlobal("exported_global")
//	extern = global.IntoExtern()
//	_ := extern.IntoGlobal()
func (e *Extern) IntoGlobal() *Global {
	pointer := C.wasm_extern_as_global(e.inner())

	if pointer == nil {
		return nil
	}

	return newGlobal(pointer, e.ownedBy())
}

// IntoTable converts the Extern into a Table.
//
// Note:️ If the Extern is not a Table, IntoTable will return nil as
// its result.
//
//	table, _ := instance.Exports.GetTable("exported_table")
//	extern = table.IntoExtern()
//	_ := extern.IntoTable()
func (e *Extern) IntoTable() *Table {
	pointer := C.wasm_extern_as_table(e.inner())

	if pointer == nil {
		return nil
	}

	return newTable(pointer, e.ownedBy())
}

// IntoMemory converts the Extern into a Memory.
//
// Note:️ If the Extern is not a Memory, IntoMemory will return nil as
// its result.
//
//	memory, _ := instance.Exports.GetMemory("exported_memory")
//	extern = memory.IntoExtern()
//	_ := extern.IntoMemory()
func (e *Extern) IntoMemory() *Memory {
	pointer := C.wasm_extern_as_memory(e.inner())

	if pointer == nil {
		return nil
	}

	return newMemory(pointer, e.ownedBy())
}

// Represents the kind of an Extern.
type ExternKind C.wasm_externkind_t

const (
	// Represents an extern of kind function.
	FUNCTION = ExternKind(C.WASM_EXTERN_FUNC)
	// Represents an extern of kind global.
	GLOBAL = ExternKind(C.WASM_EXTERN_GLOBAL)
	// Represents an extern of kind table.
	TABLE = ExternKind(C.WASM_EXTERN_TABLE)
	// Represents an extern of kind memory.
	MEMORY = ExternKind(C.WASM_EXTERN_MEMORY)
)

// String returns the ExternKind as a string.
//
//	FUNCTION.String() // "func"
//	GLOBAL.String()   // "global"
//	TABLE.String()    // "table"
//	MEMORY.String()   // "memory"
func (ek ExternKind) String() string {
	switch ek {
	case FUNCTION:
		return "func"
	case GLOBAL:
		return "global"
	case TABLE:
		return "table"
	case MEMORY:
		return "memory"
	}

	panic("Unknown extern kind") // unreachable
}

// ExternType classifies imports and external values with their respective types.
//
// # See also
//
// Specification: https://webassembly.github.io/spec/core/syntax/types.html#external-types
type ExternType struct {
	_inner   *C.wasm_externtype_t
	_ownedBy interface{}
}

type IntoExternType interface {
	IntoExternType() *ExternType
}

func newExternType(pointer *C.wasm_externtype_t, ownedBy interface{}) *ExternType {
	externType := &ExternType{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(externType, func(externType *ExternType) {
			C.wasm_externtype_delete(externType.inner())
		})
	}

	return externType
}

func (et *ExternType) inner() *C.wasm_externtype_t {
	return et._inner
}

func (et *ExternType) ownedBy() interface{} {
	if et._ownedBy == nil {
		return et
	}

	return et._ownedBy
}

// Kind returns the ExternType's ExternKind
//
//	global, _ := instance.Exports.GetGlobal("exported_global")
//	extern = global.IntoExtern()
//	_ = extern.Kind()
func (et *ExternType) Kind() ExternKind {
	kind := ExternKind(C.wasm_externtype_kind(et.inner()))

	runtime.KeepAlive(et)

	return kind
}

// IntoFunctionType converts the ExternType into a FunctionType.
//
// Note:️ If the ExternType is not a FunctionType, IntoFunctionType
// will return nil as its result.
//
//	function, _ := instance.Exports.GetFunction("exported_function")
//	externType = function.IntoExtern().Type()
//	_ := externType.IntoFunctionType()
func (et *ExternType) IntoFunctionType() *FunctionType {
	pointer := C.wasm_externtype_as_functype_const(et.inner())

	if pointer == nil {
		return nil
	}

	return newFunctionType(pointer, et.ownedBy())
}

// IntoGlobalType converts the ExternType into a GlobalType.
//
// Note:️ If the ExternType is not a GlobalType, IntoGlobalType will
// return nil as its result.
//
//	global, _ := instance.Exports.GetGlobal("exported_global")
//	externType = global.IntoExtern().Type()
//	_ := externType.IntoGlobalType()
func (et *ExternType) IntoGlobalType() *GlobalType {
	pointer := C.wasm_externtype_as_globaltype_const(et.inner())

	if pointer == nil {
		return nil
	}

	return newGlobalType(pointer, et.ownedBy())
}

// IntoTableType converts the ExternType into a TableType.
//
// Note:️ If the ExternType is not a TableType, IntoTableType will
// return nil as its result.
//
//	table, _ := instance.Exports.GetTable("exported_table")
//	externType = table.IntoExtern().Type()
//	_ := externType.IntoTableType()
func (et *ExternType) IntoTableType() *TableType {
	pointer := C.wasm_externtype_as_tabletype_const(et.inner())

	if pointer == nil {
		return nil
	}

	return newTableType(pointer, et.ownedBy())
}

// IntoMemoryType converts the ExternType into a MemoryType.
//
// Note:️ If the ExternType is not a MemoryType, IntoMemoryType will
// return nil as its result.
//
//	memory, _ := instance.Exports.GetMemory("exported_memory")
//	externType = memory.IntoExtern().Type()
//	_ := externType.IntoMemoryType()
func (et *ExternType) IntoMemoryType() *MemoryType {
	pointer := C.wasm_externtype_as_memorytype_const(et.inner())

	if pointer == nil {
		return nil
	}

	return newMemoryType(pointer, et.ownedBy())
}

// FunctionType classifies the signature of functions, mapping a
// vector of parameters to a vector of results. They are also used to
// classify the inputs and outputs of instructions.
//
// # See also
//
// Specification: https://webassembly.github.io/spec/core/syntax/types.html#function-types
type FunctionType struct {
	_inner   *C.wasm_functype_t
	_ownedBy interface{}
}

func newFunctionType(pointer *C.wasm_functype_t, ownedBy interface{}) *FunctionType {
	functionType := &FunctionType{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(functionType, func(functionType *FunctionType) {
			C.wasm_functype_delete(functionType.inner())
		})
	}

	return functionType
}

// NewFunctionType instantiates a new FunctionType from two ValueType
// arrays: the parameters and the results.
//
//	params := wasmer.NewValueTypes()
//	results := wasmer.NewValueTypes(wasmer.I32)
//	functionType := wasmer.NewFunctionType(params, results)
func NewFunctionType(params []*ValueType, results []*ValueType) *FunctionType {
	paramsAsValueTypeVec := toValueTypeVec(params)
	resultsAsValueTypeVec := toValueTypeVec(results)

	pointer := C.wasm_functype_new(&paramsAsValueTypeVec, &resultsAsValueTypeVec)

	return newFunctionType(pointer, nil)
}

func (ft *FunctionType) inner() *C.wasm_functype_t {
	return ft._inner
}

func (ft *FunctionType) ownedBy() interface{} {
	if ft._ownedBy == nil {
		return ft
	}

	return ft._ownedBy
}

// Params returns the parameters definitions from the FunctionType as
// a ValueType array
//
//	params := wasmer.NewValueTypes()
//	results := wasmer.NewValueTypes(wasmer.I32)
//	functionType := wasmer.NewFunctionType(params, results)
//	paramsValueTypes = functionType.Params()
func (ft *FunctionType) Params() []*ValueType {
	return toValueTypeList(C.wasm_functype_params(ft.inner()), ft.ownedBy())
}

// Results returns the results definitions from the FunctionType as a
// ValueType array
//
//	params := wasmer.NewValueTypes()
//	results := wasmer.NewValueTypes(wasmer.I32)
//	functionType := wasmer.NewFunctionType(params, results)
//	resultsValueTypes = functionType.Results()
func (ft *FunctionType) Results() []*ValueType {
	return toValueTypeList(C.wasm_functype_results(ft.inner()), ft.ownedBy())
}

// IntoExternType converts the FunctionType into an ExternType.
//
//	function, _ := instance.Exports.GetFunction("exported_function")
//	functionType := function.Type()
//	externType = functionType.IntoExternType()
func (ft *FunctionType) IntoExternType() *ExternType {
	pointer := C.wasm_functype_as_externtype_const(ft.inner())

	return newExternType(pointer, ft.ownedBy())
}

// Store represents all global state that can be manipulated by
// WebAssembly programs. It consists of the runtime representation of
// all instances of functions, tables, memories, and globals that have
// been allocated during the life time of the abstract machine.
//
// The Store holds the Engine (that is — amongst many things — used to
// compile the Wasm bytes into a valid module artifact).
//
// # See also
//
// Specification: https://webassembly.github.io/spec/core/exec/runtime.html#store
type Store struct {
	_inner *C.wasm_store_t
	Engine *Engine
}

// NewStore instantiates a new Store with an Engine.
//
//	engine := NewEngine()
//	store := NewStore(engine)
func NewStore(engine *Engine) *Store {
	s := &Store{
		_inner: C.wasm_store_new(engine._inner),
		Engine: engine,
	}

	runtime.SetFinalizer(s, func(s *Store) {
		s.Close()
	})

	return s
}

func (s *Store) inner() *C.wasm_store_t {
	return s._inner
}

// Force to close the Store.
//
// A runtime finalizer is registered on the Store, but it is possible
// to force the destruction of the Store by calling Close manually.
func (s *Store) Close() {
	runtime.SetFinalizer(s, nil)
	C.wasm_store_delete(s.inner())
}

// NativeFunction is a type alias representing a host function that
// can be called as any Go function.
type NativeFunction = func(...interface{}) (interface{}, error)

// Function is a WebAssembly function instance.
type Function struct {
	_inner      *C.wasm_func_t
	_ownedBy    interface{}
	environment *functionEnvironment
	lazyNative  NativeFunction
}

func newFunction(pointer *C.wasm_func_t, environment *functionEnvironment, ownedBy interface{}) *Function {
	function := &Function{
		_inner:      pointer,
		_ownedBy:    ownedBy,
		environment: environment,
		lazyNative:  nil,
	}

	if ownedBy == nil {
		runtime.SetFinalizer(function, func(f *Function) {
			if f.environment != nil {
				hostFunctionStore.remove(f.environment.hostFunctionStoreIndex)
			}

			C.wasm_func_delete(f.inner())
		})
	}

	return function
}

// NewFunction instantiates a new Function in the given Store.
//
// It takes three arguments, the Store, the FunctionType and the
// definition for the Function.
//
// The function definition must be a native Go function with a Value
// array as its single argument.  The function must return a Value
// array or an error.
//
// Note:️ Even if the function does not take any argument (or use any
// argument) it must receive a Value array as its single argument. At
// runtime, this array will be empty.  The same applies to the result.
//
//	hostFunction := wasmer.NewFunction(
//		store,
//		wasmer.NewFunctionType(
//			wasmer.NewValueTypes(), // zero argument
//			wasmer.NewValueTypes(wasmer.I32), // one i32 result
//		),
//		func(args []wasmer.Value) ([]wasmer.Value, error) {
//			return []wasmer.Value{wasmer.NewI32(42)}, nil
//		},
//	)
func NewFunction(store *Store, ty *FunctionType, function func([]Value) ([]Value, error)) *Function {
	hostFunction := &hostFunction{
		store:    store,
		function: function,
	}
	environment := &functionEnvironment{
		hostFunctionStoreIndex: hostFunctionStore.store(hostFunction),
	}
	pointer := C.wasm_func_new_with_env(
		store.inner(),
		ty.inner(),
		(C.wasm_func_callback_t)(C.function_trampoline),
		unsafe.Pointer(environment),
		(C.wasm_func_callback_env_finalizer_t)(C.function_environment_finalizer),
	)

	runtime.KeepAlive(environment)

	return newFunction(pointer, environment, nil)
}

//export function_trampoline
func function_trampoline(env unsafe.Pointer, args *C.wasm_val_vec_t, res *C.wasm_val_vec_t) *C.wasm_trap_t {
	environment := (*functionEnvironment)(env)
	hostFunction, err := hostFunctionStore.load(environment.hostFunctionStoreIndex)

	if err != nil {
		panic(err)
	}

	arguments := toValueList(args)
	function := (hostFunction.function).(func([]Value) ([]Value, error))
	results, err := (function)(arguments)

	if err != nil {
		pointer := newWasmTrap(hostFunction.store, err.Error())
		return pointer
	}

	toValueVec(results, res)

	return nil
}

// NewFunctionWithEnvironment is similar to NewFunction except that
// the user-defined host function (in Go) accepts an additional first
// parameter which is an environment. This environment can be
// anything. It is typed as interface{}.
//
//	type MyEnvironment struct {
//		foo int32
//	}
//
//	environment := &MyEnvironment {
//		foo: 42,
//	}
//
//	hostFunction := wasmer.NewFunction(
//		store,
//		wasmer.NewFunctionType(
//			wasmer.NewValueTypes(), // zero argument
//			wasmer.NewValueTypes(wasmer.I32), // one i32 result
//		),
//		environment,
//		func(environment interface{}, args []wasmer.Value) ([]wasmer.Value, error) {
//			_ := environment.(*MyEnvironment)
//
//			return []wasmer.Value{wasmer.NewI32(42)}, nil
//		},
//	)
func NewFunctionWithEnvironment(store *Store, ty *FunctionType, userEnvironment interface{}, functionWithEnv func(interface{}, []Value) ([]Value, error)) *Function {
	hostFunction := &hostFunction{
		store:           store,
		function:        functionWithEnv,
		userEnvironment: userEnvironment,
	}
	environment := &functionEnvironment{
		hostFunctionStoreIndex: hostFunctionStore.store(hostFunction),
	}
	pointer := C.wasm_func_new_with_env(
		store.inner(),
		ty.inner(),
		(C.wasm_func_callback_t)(C.function_with_environment_trampoline),
		unsafe.Pointer(environment),
		(C.wasm_func_callback_env_finalizer_t)(C.function_environment_finalizer),
	)

	runtime.KeepAlive(environment)

	return newFunction(pointer, environment, nil)
}

//export function_with_environment_trampoline
func function_with_environment_trampoline(env unsafe.Pointer, args *C.wasm_val_vec_t, res *C.wasm_val_vec_t) *C.wasm_trap_t {
	environment := (*functionEnvironment)(env)
	hostFunction, err := hostFunctionStore.load(environment.hostFunctionStoreIndex)

	if err != nil {
		panic(err)
	}

	arguments := toValueList(args)
	function := (hostFunction.function).(func(interface{}, []Value) ([]Value, error))
	results, err := (function)(hostFunction.userEnvironment, arguments)

	if err != nil {
		pointer := newWasmTrap(hostFunction.store, err.Error())
		return pointer
	}

	toValueVec(results, res)

	return nil
}

func (f *Function) inner() *C.wasm_func_t {
	return f._inner
}

func (f *Function) ownedBy() interface{} {
	if f._ownedBy == nil {
		return f
	}

	return f._ownedBy
}

// IntoExtern converts the Function into an Extern.
//
//	function, _ := instance.Exports.GetFunction("exported_function")
//	extern := function.IntoExtern()
func (f *Function) IntoExtern() *Extern {
	pointer := C.wasm_func_as_extern(f.inner())

	return newExtern(pointer, f.ownedBy())
}

// Type returns the Function's FunctionType.
//
//	function, _ := instance.Exports.GetFunction("exported_function")
//	ty := function.Type()
func (f *Function) Type() *FunctionType {
	ty := C.wasm_func_type(f.inner())

	runtime.KeepAlive(f)

	return newFunctionType(ty, f.ownedBy())
}

// ParameterArity returns the number of arguments the Function expects as per its definition.
//
//	function, _ := instance.Exports.GetFunction("exported_function")
//	arity := function.ParameterArity()
func (f *Function) ParameterArity() uint {
	return uint(C.wasm_func_param_arity(f.inner()))
}

// ParameterArity returns the number of results the Function will return.
//
//	function, _ := instance.Exports.GetFunction("exported_function")
//	arity := function.ResultArity()
func (f *Function) ResultArity() uint {
	return uint(C.wasm_func_result_arity(f.inner()))
}

// Call will call the Function and return its results as native Go values.
//
//	function, _ := instance.Exports.GetFunction("exported_function")
//	_ = function.Call(1, 2, 3)
func (f *Function) Call(parameters ...interface{}) (interface{}, error) {
	return f.Native()(parameters...)
}

// Native will turn the Function into a native Go function that can be then called.
//
//	function, _ := instance.Exports.GetFunction("exported_function")
//	nativeFunction = function.Native()
//	_ = nativeFunction(1, 2, 3)
func (f *Function) Native() NativeFunction {
	if f.lazyNative != nil {
		return f.lazyNative
	}

	ty := f.Type()
	expectedParameters := ty.Params()

	f.lazyNative = func(receivedParameters ...interface{}) (interface{}, error) {
		numberOfReceivedParameters := len(receivedParameters)
		numberOfExpectedParameters := len(expectedParameters)
		diff := numberOfExpectedParameters - numberOfReceivedParameters

		if diff > 0 {
			return nil, newErrorWith(fmt.Sprintf("Missing %d argument(s) when calling the function; Expected %d argument(s), received %d", diff, numberOfExpectedParameters, numberOfReceivedParameters))
		} else if diff < 0 {
			return nil, newErrorWith(fmt.Sprintf("Given %d extra argument(s) when calling the function; Expected %d argument(s), received %d", -diff, numberOfExpectedParameters, numberOfReceivedParameters))
		}

		allArguments := make([]C.wasm_val_t, numberOfReceivedParameters)

		for nth, receivedParameter := range receivedParameters {
			argument, err := fromGoValue(receivedParameter, expectedParameters[nth].Kind())

			if err != nil {
				return nil, newErrorWith(fmt.Sprintf("Argument %d of the function must of of type `%s`, cannot cast value to this type.", nth+1, err))
			}

			allArguments[nth] = argument
		}

		results := C.wasm_val_vec_t{}
		C.wasm_val_vec_new_uninitialized(&results, C.size_t(len(ty.Results())))
		defer C.wasm_val_vec_delete(&results)

		arguments := C.wasm_val_vec_t{}
		defer C.wasm_val_vec_delete(&arguments)

		if numberOfReceivedParameters > 0 {
			C.wasm_val_vec_new(&arguments, C.size_t(numberOfReceivedParameters), (*C.wasm_val_t)(unsafe.Pointer(&allArguments[0])))
		}

		trap := C.wasm_func_call(f.inner(), &arguments, &results)

		runtime.KeepAlive(arguments)
		runtime.KeepAlive(results)

		if trap != nil {
			return nil, newErrorFromTrap(trap)
		}

		switch results.size {
		case 0:
			return nil, nil
		case 1:
			return toGoValue(results.data), nil
		default:
			numberOfValues := int(results.size)
			allResults := make([]interface{}, numberOfValues)
			firstValue := unsafe.Pointer(results.data)
			sizeOfValuePointer := unsafe.Sizeof(C.wasm_val_t{})

			var currentValuePointer *C.wasm_val_t

			for nth := 0; nth < numberOfValues; nth++ {
				currentValuePointer = (*C.wasm_val_t)(unsafe.Pointer(uintptr(firstValue) + uintptr(nth)*sizeOfValuePointer))
				value := toGoValue(currentValuePointer)
				allResults[nth] = value
			}

			return allResults, nil
		}
	}

	return f.lazyNative
}

type functionEnvironment struct {
	store                  *Store
	hostFunctionStoreIndex uint
}

//export function_environment_finalizer
func function_environment_finalizer(_ unsafe.Pointer) {}

type hostFunction struct {
	store           *Store
	function        interface{} // func([]Value) ([]Value, error) or func(interface{}, []Value) ([]Value, error)
	userEnvironment interface{} // if the host function has an environment
}

type hostFunctions struct {
	sync.RWMutex
	functions map[uint]*hostFunction
}

func (hf *hostFunctions) load(index uint) (*hostFunction, error) {
	hf.RLock()
	hostFunction, exists := hf.functions[index]
	hf.RUnlock()

	if exists && hostFunction != nil {
		return hostFunction, nil
	}

	return nil, newErrorWith(fmt.Sprintf("Host function `%d` does not exist", index))
}

func (hf *hostFunctions) store(hostFunction *hostFunction) uint {
	hf.Lock()
	// By default, the index is the size of the store.
	index := uint(len(hf.functions))

	for nth, hostFunc := range hf.functions {
		// Find the first empty slot in the store.
		if hostFunc == nil {
			// Use that empty slot for the index.
			index = nth
			break
		}
	}

	hf.functions[index] = hostFunction
	hf.Unlock()

	return index
}

func (hf *hostFunctions) remove(index uint) {
	hf.Lock()
	hf.functions[index] = nil
	hf.Unlock()
}

var hostFunctionStore = hostFunctions{
	functions: make(map[uint]*hostFunction),
}

// Global stores a single value of the given GlobalType.
//
// # See also
//
// https://webassembly.github.io/spec/core/syntax/modules.html#globals
type Global struct {
	_inner   *C.wasm_global_t
	_ownedBy interface{}
}

func newGlobal(pointer *C.wasm_global_t, ownedBy interface{}) *Global {
	global := &Global{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(global, func(global *Global) {
			C.wasm_global_delete(global.inner())
		})
	}

	return global
}

// NewGlobal instantiates a new Global in the given Store.
//
// It takes three arguments, the Store, the GlobalType and the Value for the Global.
//
//	valueType := NewValueType(I32)
//	globalType := NewGlobalType(valueType, CONST)
//	global := NewGlobal(store, globalType, NewValue(42, I32))
func NewGlobal(store *Store, ty *GlobalType, value Value) *Global {
	pointer := C.wasm_global_new(
		store.inner(),
		ty.inner(),
		value.inner(),
	)

	return newGlobal(pointer, nil)
}

func (g *Global) inner() *C.wasm_global_t {
	return g._inner
}

func (g *Global) ownedBy() interface{} {
	if g._ownedBy == nil {
		return g
	}

	return g._ownedBy
}

// IntoExtern converts the Global into an Extern.
//
//	global, _ := instance.Exports.GetGlobal("exported_global")
//	extern := global.IntoExtern()
func (g *Global) IntoExtern() *Extern {
	pointer := C.wasm_global_as_extern(g.inner())

	return newExtern(pointer, g.ownedBy())
}

// Type returns the Global's GlobalType.
//
//	global, _ := instance.Exports.GetGlobal("exported_global")
//	ty := global.Type()
func (g *Global) Type() *GlobalType {
	ty := C.wasm_global_type(g.inner())
	runtime.KeepAlive(g)
	return newGlobalType(ty, g.ownedBy())
}

// Set sets the Global's value.
//
// It takes two arguments, the Global's value as a native Go value and the value's ValueKind.
//
//	global, _ := instance.Exports.GetGlobal("exported_global")
//	_ = global.Set(1, I32)
func (g *Global) Set(value interface{}, kind ValueKind) error {
	if g.Type().Mutability() == IMMUTABLE {
		return newErrorWith("The global variable is not mutable, cannot set a new value")
	}

	result, err := fromGoValue(value, kind)

	if err != nil {
		//TODO: Make this error explicit
		panic(err.Error())
	}

	C.wasm_global_set(g.inner(), &result)

	return nil
}

// Get returns the Global's value as a native Go value.
//
//	global, _ := instance.Exports.GetGlobal("exported_global")
//	value, _ := global.Get()
func (g *Global) Get() (interface{}, error) {
	var value C.wasm_val_t

	C.wasm_global_get(g.inner(), &value)

	return toGoValue(&value), nil
}

type GlobalMutability C.wasm_mutability_t

const (
	// Represents a global that is immutable.
	IMMUTABLE = GlobalMutability(C.WASM_CONST)
	// Represents a global that is mutable.
	MUTABLE = GlobalMutability(C.WASM_VAR)
)

// String returns the GlobalMutability as a string.
//
//	IMMUTABLE.String() // "const"
//	MUTABLE.String()   // "var"
func (gm GlobalMutability) String() string {
	switch gm {
	case IMMUTABLE:
		return "const"
	case MUTABLE:
		return "var"
	}

	panic("Unknown mutability") // unreachable
}

// GlobalType classifies global variables, which hold a value and can either be mutable or immutable.
//
// # See also
//
// Specification: https://webassembly.github.io/spec/core/syntax/types.html#global-types
type GlobalType struct {
	_inner   *C.wasm_globaltype_t
	_ownedBy interface{}
}

func newGlobalType(pointer *C.wasm_globaltype_t, ownedBy interface{}) *GlobalType {
	globalType := &GlobalType{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(globalType, func(globalType *GlobalType) {
			C.wasm_globaltype_delete(globalType.inner())
		})
	}

	return globalType
}

// NewGlobalType instantiates a new GlobalType from a ValueType and a GlobalMutability
//
//	valueType := NewValueType(I32)
//	globalType := NewGlobalType(valueType, IMMUTABLE)
func NewGlobalType(valueType *ValueType, mutability GlobalMutability) *GlobalType {
	pointer := C.wasm_globaltype_new(valueType.inner(), C.wasm_mutability_t(mutability))

	return newGlobalType(pointer, nil)
}

func (gt *GlobalType) inner() *C.wasm_globaltype_t {
	return gt._inner
}

func (gt *GlobalType) ownedBy() interface{} {
	if gt._ownedBy == nil {
		return gt
	}

	return gt._ownedBy
}

// ValueType returns the GlobalType's ValueType
//
//	valueType := NewValueType(I32)
//	globalType := NewGlobalType(valueType, IMMUTABLE)
//	globalType.ValueType().Kind().String() // "i32"
func (gt *GlobalType) ValueType() *ValueType {
	pointer := C.wasm_globaltype_content(gt.inner())
	runtime.KeepAlive(gt)
	return newValueType(pointer, gt.ownedBy())
}

// Mutability returns the GlobalType's GlobalMutability
//
//	valueType := NewValueType(I32)
//	globalType := NewGlobalType(valueType, IMMUTABLE)
//	globalType.Mutability().String() // "const"
func (gt *GlobalType) Mutability() GlobalMutability {
	mutability := GlobalMutability(C.wasm_globaltype_mutability(gt.inner()))
	runtime.KeepAlive(gt)
	return mutability
}

// IntoExternType converts the GlobalType into an ExternType.
//
//	valueType := NewValueType(I32)
//	globalType := NewGlobalType(valueType, IMMUTABLE)
//	externType = globalType.IntoExternType()
func (gt *GlobalType) IntoExternType() *ExternType {
	pointer := C.wasm_globaltype_as_externtype_const(gt.inner())

	return newExternType(pointer, gt.ownedBy())
}

// ImportObject contains all of the import data used when
// instantiating a WebAssembly module.
type ImportObject struct {
	externs map[string]map[string]IntoExtern
}

// NewImportObject instantiates a new empty ImportObject.
//
//	imports := NewImportObject()
func NewImportObject() *ImportObject {
	return &ImportObject{
		externs: make(map[string]map[string]IntoExtern),
	}
}

func (io *ImportObject) intoInner(module *Module) (*C.wasm_extern_vec_t, error) {
	cExterns := &C.wasm_extern_vec_t{}

	var externs []*C.wasm_extern_t
	var numberOfExterns uint

	for _, importType := range module.Imports() {
		namespace := importType.Module()
		name := importType.Name()

		if io.externs[namespace][name] == nil {
			return nil, &Error{
				message: fmt.Sprintf("Missing import: `%s`.`%s`", namespace, name),
			}
		}

		externs = append(externs, io.externs[namespace][name].IntoExtern().inner())
		numberOfExterns++
	}

	if numberOfExterns > 0 {
		C.wasm_extern_vec_new(cExterns, C.size_t(numberOfExterns), (**C.wasm_extern_t)(unsafe.Pointer(&externs[0])))
	}

	return cExterns, nil
}

// ContainsNamespace returns true if the ImportObject contains the given namespace (or module name)
//
//	imports := NewImportObject()
//	_ = imports.ContainsNamespace("env") // false
func (io *ImportObject) ContainsNamespace(name string) bool {
	_, exists := io.externs[name]

	return exists
}

// Register registers a namespace (or module name) in the ImportObject.
//
// It takes two arguments: the namespace name and a map with imports names as key and externs as values.
//
// Note:️ An extern is anything implementing IntoExtern: Function, Global, Memory, Table.
//
//	 imports := NewImportObject()
//	 importObject.Register(
//	 	"env",
//	 	map[string]wasmer.IntoExtern{
//	 		"host_function": hostFunction,
//	 		"host_global": hostGlobal,
//	 	},
//	)
//
// Note:️ The namespace (or module name) may be empty:
//
//	imports := NewImportObject()
//	importObject.Register(
//		"",
//		map[string]wasmer.IntoExtern{
//	 		"host_function": hostFunction,
//			"host_global": hostGlobal,
//		},
//	)
func (io *ImportObject) Register(namespaceName string, namespace map[string]IntoExtern) {
	_, exists := io.externs[namespaceName]

	if !exists {
		io.externs[namespaceName] = namespace
	} else {
		for key, value := range namespace {
			io.externs[namespaceName][key] = value
		}
	}
}

type importTypes struct {
	_inner      C.wasm_importtype_vec_t
	importTypes []*ImportType
}

func newImportTypes(module *Module) *importTypes {
	it := &importTypes{}
	C.wasm_module_imports(module.inner(), &it._inner)

	runtime.SetFinalizer(it, func(it *importTypes) {
		it.close()
	})

	numberOfImportTypes := int(it.inner().size)
	types := make([]*ImportType, numberOfImportTypes)
	firstImportType := unsafe.Pointer(it.inner().data)
	sizeOfImportTypePointer := unsafe.Sizeof(firstImportType)

	var currentTypePointer *C.wasm_importtype_t

	for nth := 0; nth < numberOfImportTypes; nth++ {
		currentTypePointer = *(**C.wasm_importtype_t)(unsafe.Pointer(uintptr(firstImportType) + uintptr(nth)*sizeOfImportTypePointer))
		importType := newImportType(currentTypePointer, it)
		types[nth] = importType
	}

	it.importTypes = types

	return it
}

func (it *importTypes) inner() *C.wasm_importtype_vec_t {
	return &it._inner
}

func (it *importTypes) close() {
	runtime.SetFinalizer(it, nil)
	C.wasm_importtype_vec_delete(&it._inner)
}

// ImportType is a descriptor for an imported value into a WebAssembly
// module.
type ImportType struct {
	_inner   *C.wasm_importtype_t
	_ownedBy interface{}
}

func newImportType(pointer *C.wasm_importtype_t, ownedBy interface{}) *ImportType {
	importType := &ImportType{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(importType, func(it *ImportType) {
			it.Close()
		})
	}

	return importType
}

// NewImportType instantiates a new ImportType with a module name (or
// namespace), a name and an extern type.
//
// Note:️ An extern type is anything implementing IntoExternType:
// FunctionType, GlobalType, MemoryType, TableType.
//
//	valueType := NewValueType(I32)
//	globalType := NewGlobalType(valueType, CONST)
//	importType := NewImportType("ns", "host_global", globalType)
func NewImportType(module string, name string, ty IntoExternType) *ImportType {
	moduleName := newName(module)
	nameName := newName(name)
	externType := ty.IntoExternType().inner()
	externTypeCopy := C.wasm_externtype_copy(externType)

	runtime.KeepAlive(externType)

	importType := C.wasm_importtype_new(&moduleName, &nameName, externTypeCopy)

	return newImportType(importType, nil)
}

func (it *ImportType) inner() *C.wasm_importtype_t {
	return it._inner
}

func (it *ImportType) ownedBy() interface{} {
	if it._ownedBy == nil {
		return it
	}

	return it._ownedBy
}

// Module returns the ImportType's module name (or namespace).
//
//	valueType := NewValueType(I32)
//	globalType := NewGlobalType(valueType, CONST)
//	importType := NewImportType("ns", "host_global", globalType)
//	_ = importType.Module()
func (it *ImportType) Module() string {
	byteVec := C.wasm_importtype_module(it.inner())
	module := C.GoStringN(byteVec.data, C.int(byteVec.size))

	runtime.KeepAlive(it)

	return module
}

// Name returns the ImportType's name.
//
//	valueType := NewValueType(I32)
//	globalType := NewGlobalType(valueType, CONST)
//	importType := NewImportType("ns", "host_global", globalType)
//	_ = importType.Name()
func (it *ImportType) Name() string {
	byteVec := C.wasm_importtype_name(it.inner())
	name := C.GoStringN(byteVec.data, C.int(byteVec.size))

	runtime.KeepAlive(it)

	return name
}

// Type returns the ImportType's type as an ExternType.
//
//	valueType := NewValueType(I32)
//	globalType := NewGlobalType(valueType, CONST)
//	importType := NewImportType("ns", "host_global", globalType)
//	_ = importType.Type()
func (it *ImportType) Type() *ExternType {
	ty := C.wasm_importtype_type(it.inner())

	runtime.KeepAlive(it)

	return newExternType(ty, it.ownedBy())
}

// Force to close the ImportType.
//
// A runtime finalizer is registered on the ImportType, but it is
// possible to force the destruction of the ImportType by calling
// Close manually.
func (it *ImportType) Close() {
	runtime.SetFinalizer(it, nil)
	C.wasm_importtype_delete(it.inner())
}

type Instance struct {
	_inner  *C.wasm_instance_t
	Exports *Exports

	// without this, imported functions may be freed before execution of an exported function is complete.
	imports *ImportObject
}

// NewInstance instantiates a new Instance.
//
// It takes two arguments, the Module and an ImportObject.
//
// Note:️ Instantiating a module may return TrapError if the module's
// start function traps.
//
//	wasmBytes := []byte(`...`)
//	engine := wasmer.NewEngine()
//	store := wasmer.NewStore(engine)
//	module, err := wasmer.NewModule(store, wasmBytes)
//	importObject := wasmer.NewImportObject()
//	instance, err := wasmer.NewInstance(module, importObject)
func NewInstance(module *Module, imports *ImportObject, env *WasiEnvironment) (*Instance, error) {
	var traps *C.wasm_trap_t
	externs, err := imports.intoInner(module)
	if err != nil {
		return nil, err
	}

	var instance *C.wasm_instance_t

	err2 := maybeNewErrorFromWasmer(func() bool {
		instance = C.wasm_instance_new(
			module.store.inner(),
			module.inner(),
			externs,
			&traps,
		)

		return traps == nil && instance == nil
	})

	if err2 != nil {
		return nil, err2
	}

	if traps != nil {
		return nil, newErrorFromTrap(traps)
	}

	// Initialize WASI memory if the module contains WASI imports.
	wasiVersion := GetWasiVersion(module)
	if wasiVersion > WASI_VERSION_INVALID {
		if env == nil {
			return nil, newErrorWith("WASI environment is required for WASI modules")
		}

		if !bool(C.wasi_env_initialize_instance(env.inner(), module.store.inner(), instance)) {
			return nil, newErrorWith("Failed to initialize WASI memory")
		}
	}

	inst := &Instance{
		_inner:  instance,
		Exports: newExports(instance, module),
		imports: imports,
	}

	runtime.SetFinalizer(inst, func(inst *Instance) {
		inst.Close()
	})

	return inst, nil
}

func (inst *Instance) inner() *C.wasm_instance_t {
	return inst._inner
}

// GetRemainingPoints exposes wasm metering remaining gas or points
func (inst *Instance) GetRemainingPoints() uint64 {
	return uint64(C.wasmer_metering_get_remaining_points(inst._inner))
}

// GetRemainingPoints a bool to determine if the engine has been shutdown from meter exhaustion
func (inst *Instance) MeteringPointsExhausted() bool {
	return bool(C.wasmer_metering_points_are_exhausted(inst._inner))
}

// SetRemainingPoints imposes a new gas limit on the wasm engine
func (inst *Instance) SetRemainingPoints(newLimit uint64) {
	C.wasmer_metering_set_remaining_points(inst._inner, C.uint64_t(newLimit))
}

// Force to close the Instance.
//
// A runtime finalizer is registered on the Instance, but it is
// possible to force the destruction of the Instance by calling Close
// manually.
func (inst *Instance) Close() {
	runtime.SetFinalizer(inst, nil)
	C.wasm_instance_delete(inst.inner())
	inst.Exports.Close()
}

// Limits classify the size range of resizable storage associated
// with memory types and table types.
//
// # See also
//
// Specification: https://webassembly.github.io/spec/core/syntax/types.html#limits
type Limits struct {
	_inner C.wasm_limits_t
}

func newLimits(pointer *C.wasm_limits_t, ownedBy interface{}) *Limits {
	limits, err := NewLimits(uint32(pointer.min), uint32(pointer.max))

	if err != nil {
		return nil
	}

	if ownedBy != nil {
		runtime.KeepAlive(ownedBy)
	}

	return limits
}

// NewLimits instantiates a new Limits which describes the Memory used.
// The minimum and maximum parameters are "number of memory pages".
//
// ️Note: Each page is 64 KiB in size.
//
// Note: You cannot Memory.Grow the Memory beyond the maximum defined here.
func NewLimits(minimum uint32, maximum uint32) (*Limits, error) {
	if minimum > maximum {
		return nil, newErrorWith("The minimum limit is greater than the maximum one")
	}

	return &Limits{
		_inner: C.wasm_limits_t{
			min: C.uint32_t(minimum),
			max: C.uint32_t(maximum),
		},
	}, nil
}

func (l *Limits) inner() *C.wasm_limits_t {
	return &l._inner
}

// Minimum returns the minimum size of the Memory allocated in "number of pages".
//
// Note:️ Each page is 64 KiB in size.
func (l *Limits) Minimum() uint32 {
	return uint32(l.inner().min)
}

// Maximum returns the maximum size of the Memory allocated in "number of pages".
//
// Each page is 64 KiB in size.
//
// Note: You cannot Memory.Grow beyond this defined maximum size.
func (l *Limits) Maximum() uint32 {
	return uint32(l.inner().max)
}

// Memory is a vector of raw uninterpreted bytes.
//
// # See also
//
// Specification: https://webassembly.github.io/spec/core/syntax/modules.html#memories
type Memory struct {
	_inner   *C.wasm_memory_t
	_ownedBy interface{}
}

func newMemory(pointer *C.wasm_memory_t, ownedBy interface{}) *Memory {
	memory := &Memory{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(memory, func(memory *Memory) {
			C.wasm_memory_delete(memory.inner())
		})
	}

	return memory
}

// NewMemory instantiates a new Memory in the given Store.
//
// It takes two arguments, the Store and the MemoryType for the Memory.
//
//	memory := wasmer.NewMemory(
//	    store,
//	    wasmer.NewMemoryType(wasmer.NewLimits(1, 4)),
//	)
func NewMemory(store *Store, ty *MemoryType) *Memory {
	pointer := C.wasm_memory_new(store.inner(), ty.inner())

	runtime.KeepAlive(store)
	runtime.KeepAlive(ty)

	return newMemory(pointer, nil)
}

func (mem *Memory) inner() *C.wasm_memory_t {
	return mem._inner
}

func (mem *Memory) ownedBy() interface{} {
	if mem._ownedBy == nil {
		return mem
	}

	return mem._ownedBy
}

// Type returns the Memory's MemoryType.
//
//	memory, _ := instance.Exports.GetMemory("exported_memory")
//	ty := memory.Type()
func (mem *Memory) Type() *MemoryType {
	ty := C.wasm_memory_type(mem.inner())

	runtime.KeepAlive(mem)

	return newMemoryType(ty, mem.ownedBy())
}

// Size returns the Memory's size as Pages.
//
//	memory, _ := instance.Exports.GetMemory("exported_memory")
//	size := memory.Size()
func (mem *Memory) Size() Pages {
	return Pages(C.wasm_memory_size(mem.inner()))
}

// Size returns the Memory's size as a number of bytes.
//
//	memory, _ := instance.Exports.GetMemory("exported_memory")
//	size := memory.DataSize()
func (mem *Memory) DataSize() uint {
	return uint(C.wasm_memory_data_size(mem.inner()))
}

// Data returns the Memory's contents as an byte array.
//
//	memory, _ := instance.Exports.GetMemory("exported_memory")
//	data := memory.Data()
func (mem *Memory) Data() []byte {
	length := int(mem.DataSize())
	data := (*C.byte_t)(C.wasm_memory_data(mem.inner()))

	runtime.KeepAlive(mem)

	var byteSlice []byte
	var header = (*reflect.SliceHeader)(unsafe.Pointer(&byteSlice))

	header.Data = uintptr(unsafe.Pointer(data))
	header.Len = length
	header.Cap = length

	return byteSlice
}

// Grow grows the Memory's size by a given number of Pages (the delta).
//
//	memory, _ := instance.Exports.GetMemory("exported_memory")
//	grown := memory.Grow(2)
func (mem *Memory) Grow(delta Pages) bool {
	return bool(C.wasm_memory_grow(mem.inner(), C.wasm_memory_pages_t(delta)))
}

// IntoExtern converts the Memory into an Extern.
//
//	memory, _ := instance.Exports.GetMemory("exported_memory")
//	extern := memory.IntoExtern()
func (mem *Memory) IntoExtern() *Extern {
	pointer := C.wasm_memory_as_extern(mem.inner())

	return newExtern(pointer, mem.ownedBy())
}

// MemoryType classifies linear memories and their size range.
//
// # See also
//
// Specification: https://webassembly.github.io/spec/core/syntax/types.html#memory-types
type MemoryType struct {
	_inner   *C.wasm_memorytype_t
	_ownedBy interface{}
}

func newMemoryType(pointer *C.wasm_memorytype_t, ownedBy interface{}) *MemoryType {
	memoryType := &MemoryType{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(memoryType, func(memoryType *MemoryType) {
			C.wasm_memorytype_delete(memoryType.inner())
		})
	}

	return memoryType
}

// NewMemoryType instantiates a new MemoryType given some Limits.
//
//	limits := NewLimits(1, 4)
//	memoryType := NewMemoryType(limits)
func NewMemoryType(limits *Limits) *MemoryType {
	pointer := C.wasm_memorytype_new(limits.inner())

	return newMemoryType(pointer, nil)
}

func (mt *MemoryType) inner() *C.wasm_memorytype_t {
	return mt._inner
}

func (mt *MemoryType) ownedBy() interface{} {
	if mt._ownedBy == nil {
		return mt
	}

	return mt._ownedBy
}

// Limits returns the MemoryType's Limits.
//
//	limits := NewLimits(1, 4)
//	memoryType := NewMemoryType(limits)
//	_ = memoryType.Limits()
func (mt *MemoryType) Limits() *Limits {
	limits := newLimits(C.wasm_memorytype_limits(mt.inner()), mt.ownedBy())

	runtime.KeepAlive(mt)

	return limits
}

// IntoExternType converts the MemoryType into an ExternType.
//
//	limits := NewLimits(1, 4)
//	memoryType := NewMemoryType(limits)
//	externType = memoryType.IntoExternType()
func (mt *MemoryType) IntoExternType() *ExternType {
	pointer := C.wasm_memorytype_as_externtype_const(mt.inner())

	return newExternType(pointer, mt.ownedBy())
}

// Module contains stateless WebAssembly code that has already been
// compiled and can be instantiated multiple times.
//
// WebAssembly programs are organized into modules, which are the unit
// of deployment, loading, and compilation. A module collects
// definitions for types, functions, tables, memories, and globals. In
// addition, it can declare imports and exports and provide
// initialization logic in the form of data and element segments or a
// start function.
//
// # See also
//
// Specification: https://webassembly.github.io/spec/core/syntax/modules.html#modules
type Module struct {
	_inner *C.wasm_module_t
	store  *Store
	// Stored if computed to avoid further reallocations.
	importTypes *importTypes
	// Stored if computed to avoid further reallocations.
	exportTypes *exportTypes
}

// NewModule instantiates a new Module with the given Store.
//
// It takes two arguments, the Store and the Wasm module as a byte
// array of WAT code.
//
//	wasmBytes := []byte(`...`)
//	engine := wasmer.NewEngine()
//	store := wasmer.NewStore(engine)
//	module, err := wasmer.NewModule(store, wasmBytes)
func NewModule(store *Store, bytes []byte) (*Module, error) {
	wasmBytes, err := Wat2Wasm(string(bytes))

	if err != nil {
		return nil, err
	}

	var wasmBytesPtr *C.uint8_t
	wasmBytesLength := len(wasmBytes)

	if wasmBytesLength > 0 {
		wasmBytesPtr = (*C.uint8_t)(unsafe.Pointer(&wasmBytes[0]))
	}

	var module *Module

	err2 := maybeNewErrorFromWasmer(func() bool {
		module = &Module{
			_inner: C.to_wasm_module_new(store.inner(), wasmBytesPtr, C.size_t(wasmBytesLength)),
			store:  store,
		}

		return module._inner == nil
	})

	if err2 != nil {
		return nil, err2
	}

	runtime.SetFinalizer(module, func(m *Module) {
		m.Close()
	})

	return module, nil
}

// ValidateModule validates a new Module against the given Store.
//
// It takes two arguments, the Store and the WebAssembly module as a
// byte array. The function returns an error describing why the bytes
// are invalid, otherwise it returns nil.
//
//	wasmBytes := []byte(`...`)
//	engine := wasmer.NewEngine()
//	store := wasmer.NewStore(engine)
//	err := wasmer.ValidateModule(store, wasmBytes)
//
//	isValid := err != nil
func ValidateModule(store *Store, bytes []byte) error {
	wasmBytes, err := Wat2Wasm(string(bytes))

	if err != nil {
		return err
	}

	var wasmBytesPtr *C.uint8_t
	wasmBytesLength := len(wasmBytes)

	if wasmBytesLength > 0 {
		wasmBytesPtr = (*C.uint8_t)(unsafe.Pointer(&wasmBytes[0]))
	}

	err2 := maybeNewErrorFromWasmer(func() bool {
		return !bool(C.to_wasm_module_validate(store.inner(), wasmBytesPtr, C.size_t(wasmBytesLength)))
	})

	if err2 != nil {
		return err2
	}

	runtime.KeepAlive(bytes)
	runtime.KeepAlive(wasmBytes)

	return nil
}

func (m *Module) inner() *C.wasm_module_t {
	return m._inner
}

// Name returns the Module's name.
//
// Note:️ This is not part of the standard Wasm C API. It is Wasmer specific.
//
//	wasmBytes := []byte(`(module $moduleName)`)
//	engine := wasmer.NewEngine()
//	store := wasmer.NewStore(engine)
//	module, _ := wasmer.NewModule(store, wasmBytes)
//	name := module.Name()
func (m *Module) Name() string {
	var name C.wasm_name_t

	C.wasmer_module_name(m.inner(), &name)

	goName := nameToString(&name)

	C.wasm_name_delete(&name)

	return goName
}

// Imports returns the Module's imports as an ImportType array.
//
//	wasmBytes := []byte(`...`)
//	engine := wasmer.NewEngine()
//	store := wasmer.NewStore(engine)
//	module, _ := wasmer.NewModule(store, wasmBytes)
//	imports := module.Imports()
func (m *Module) Imports() []*ImportType {
	if nil == m.importTypes {
		m.importTypes = newImportTypes(m)
	}

	return m.importTypes.importTypes
}

// Exports returns the Module's exports as an ExportType array.
//
//	wasmBytes := []byte(`...`)
//	engine := wasmer.NewEngine()
//	store := wasmer.NewStore(engine)
//	module, _ := wasmer.NewModule(store, wasmBytes)
//	exports := module.Exports()
func (m *Module) Exports() []*ExportType {
	if nil == m.exportTypes {
		m.exportTypes = newExportTypes(m)
	}

	return m.exportTypes.exportTypes
}

// Serialize serializes the module and returns the Wasm code as an byte array.
//
//	wasmBytes := []byte(`...`)
//	engine := wasmer.NewEngine()
//	store := wasmer.NewStore(engine)
//	module, _ := wasmer.NewModule(store, wasmBytes)
//	bytes, err := module.Serialize()
func (m *Module) Serialize() ([]byte, error) {
	var bytes C.wasm_byte_vec_t

	err := maybeNewErrorFromWasmer(func() bool {
		C.wasm_module_serialize(m.inner(), &bytes)

		return bytes.data == nil
	})

	if err != nil {
		return nil, err
	}

	goBytes := C.GoBytes(unsafe.Pointer(bytes.data), C.int(bytes.size))
	C.wasm_byte_vec_delete(&bytes)

	return goBytes, nil
}

// DeserializeModule deserializes an byte array to a Module.
//
//	wasmBytes := []byte(`...`)
//	engine := wasmer.NewEngine()
//	store := wasmer.NewStore(engine)
//	module, _ := wasmer.NewModule(store, wasmBytes)
//	bytes, err := module.Serialize()
//	//...
//	deserializedModule, err := wasmer.DeserializeModule(store, bytes)
func DeserializeModule(store *Store, bytes []byte) (*Module, error) {
	var bytesPtr *C.uint8_t
	bytesLength := len(bytes)

	if bytesLength > 0 {
		bytesPtr = (*C.uint8_t)(unsafe.Pointer(&bytes[0]))
	}

	var module *Module

	err := maybeNewErrorFromWasmer(func() bool {
		module = &Module{
			_inner: C.to_wasm_module_deserialize(store.inner(), bytesPtr, C.size_t(bytesLength)),
			store:  store,
		}

		return module._inner == nil
	})

	if err != nil {
		return nil, err
	}

	runtime.SetFinalizer(module, func(m *Module) {
		C.wasm_module_delete(m.inner())
	})

	return module, nil
}

// Force to close the Module.
//
// A runtime finalizer is registered on the Module, but it is possible
// to force the destruction of the Module by calling Close manually.
func (m *Module) Close() {
	runtime.SetFinalizer(m, nil)
	C.wasm_module_delete(m.inner())

	if nil != m.importTypes {
		m.importTypes.close()
	}

	if nil != m.exportTypes {
		m.exportTypes.close()
	}
}

// WasiVersion represents the possible WASI versions.
type WasiVersion C.wasi_version_t

const (
	// Latest version. It's a “floating” version, i.e. it's an
	// alias to the latest version. Using this version is a way to
	// ensure that modules will run only if they come with the
	// latest WASI version (in case of security issues for
	// instance), by just updating the runtime.
	WASI_VERSION_LATEST = WasiVersion(C.LATEST)

	// Represents the wasi_unstable version.
	WASI_VERSION_SNAPSHOT0 = WasiVersion(C.SNAPSHOT0)

	// Represents the wasi_snapshot_preview1 version.
	WASI_VERSION_SNAPSHOT1 = WasiVersion(C.SNAPSHOT1)

	// Represents the wasix 32-bit version.
	WASIX32V1 = WasiVersion(C.WASIX32V1)

	// Represents an invalid version.
	WASI_VERSION_INVALID = WasiVersion(C.INVALID_VERSION)
)

// String returns the WasiVersion as a string.
//
//	WASI_VERSION_SNAPSHOT0.String() //  "wasi_unstable"
//	WASI_VERSION_SNAPSHOT1.String() // "wasi_snapshot_preview1"
func (wv WasiVersion) String() string {
	switch wv {
	case WASI_VERSION_LATEST:
		return "__latest__"
	case WASI_VERSION_SNAPSHOT0:
		return "wasi_unstable"
	case WASI_VERSION_SNAPSHOT1:
		return "wasi_snapshot_preview1"
	case WASIX32V1:
		return "wasix32v1"
	case WASI_VERSION_INVALID:
		return "__unknown__"
	}

	panic("Unknown WASI version")
}

// GetWasiVersion returns the WASI version of the given Module if any,
// WASI_VERSION_INVALID otherwise.
//
//	wasiVersion := GetWasiVersion(module)
func GetWasiVersion(module *Module) WasiVersion {
	return WasiVersion(C.wasi_get_wasi_version(module.inner()))
}

// WasiStateBuilder is a convenient API for configuring WASI.
type WasiStateBuilder struct {
	_inner *C.wasi_config_t
}

// NewWasiStateBuilder creates a new WASI state builder, starting by
// configuring the WASI program name.
//
//	wasiStateBuilder := NewWasiStateBuilder("test-program")
func NewWasiStateBuilder(programName string) *WasiStateBuilder {
	cProgramName := C.CString(programName)
	defer C.free(unsafe.Pointer(cProgramName))
	wasiConfig := C.wasi_config_new(cProgramName)

	stateBuilder := &WasiStateBuilder{
		_inner: wasiConfig,
	}

	return stateBuilder
}

// Argument configures a new argument to the WASI module.
//
//	wasiStateBuilder := NewWasiStateBuilder("test-program").
//		Argument("--foo")
func (wsb *WasiStateBuilder) Argument(argument string) *WasiStateBuilder {
	cArgument := C.CString(argument)
	defer C.free(unsafe.Pointer(cArgument))
	C.wasi_config_arg(wsb.inner(), cArgument)

	return wsb
}

// Environment configures a new environment variable for the WASI module.
//
//	wasiStateBuilder := NewWasiStateBuilder("test-program").
//		Argument("--foo").
//		Environment("KEY", "VALUE")
func (wsb *WasiStateBuilder) Environment(key string, value string) *WasiStateBuilder {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	C.wasi_config_env(wsb.inner(), cKey, cValue)

	return wsb
}

// PreopenDirectory configures a new directory to pre-open.
//
// This opens the given directory at the virtual root /, and allows
// the WASI module to read and write to the given directory.
//
//	wasiStateBuilder := NewWasiStateBuilder("test-program").
//		Argument("--foo").
//		Environment("KEY", "VALUE").
//		PreopenDirectory("bar")
func (wsb *WasiStateBuilder) PreopenDirectory(preopenDirectory string) *WasiStateBuilder {
	cPreopenDirectory := C.CString(preopenDirectory)
	defer C.free(unsafe.Pointer(cPreopenDirectory))

	C.wasi_config_preopen_dir(wsb.inner(), cPreopenDirectory)

	return wsb
}

// MapDirectory configures a new directory to pre-open with a
// different name exposed to the WASI module.
//
//	wasiStateBuilder := NewWasiStateBuilder("test-program").
//		Argument("--foo").
//		Environment("KEY", "VALUE").
//		MapDirectory("the_host_current_directory", ".")
func (wsb *WasiStateBuilder) MapDirectory(alias string, directory string) *WasiStateBuilder {
	cAlias := C.CString(alias)
	defer C.free(unsafe.Pointer(cAlias))

	cDirectory := C.CString(directory)
	defer C.free(unsafe.Pointer(cDirectory))

	C.wasi_config_mapdir(wsb.inner(), cAlias, cDirectory)

	return wsb
}

// InheritStdin configures the WASI module to inherit the stdin from
// the host.
func (wsb *WasiStateBuilder) InheritStdin() *WasiStateBuilder {
	C.wasi_config_inherit_stdin(wsb.inner())

	return wsb
}

// CaptureStdout configures the WASI module to capture its stdout.
//
//	wasiStateBuilder := NewWasiStateBuilder("test-program").
//		Argument("--foo").
//		Environment("KEY", "VALUE").
//		MapDirectory("the_host_current_directory", ".")
//		CaptureStdout()
func (wsb *WasiStateBuilder) CaptureStdout() *WasiStateBuilder {
	C.wasi_config_capture_stdout(wsb.inner())

	return wsb
}

// InheritStdout configures the WASI module to inherit the stdout from
// the host.
func (wsb *WasiStateBuilder) InheritStdout() *WasiStateBuilder {
	C.wasi_config_inherit_stdout(wsb.inner())

	return wsb
}

// CaptureStderr configures the WASI module to capture its stderr.
func (wsb *WasiStateBuilder) CaptureStderr() *WasiStateBuilder {
	C.wasi_config_capture_stderr(wsb.inner())

	return wsb
}

// InheritStderr configures the WASI module to inherit the stderr from
// the host.
func (wsb *WasiStateBuilder) InheritStderr() *WasiStateBuilder {
	C.wasi_config_inherit_stderr(wsb.inner())

	return wsb
}

// Finalize tells the state builder to produce a WasiEnvironment. It
// consumes the current WasiStateBuilder.
//
// It can return an error if the state builder contains invalid
// configuration.
//
//	wasiEnvironment, err := NewWasiStateBuilder("test-program").
//		Argument("--foo").
//		Environment("KEY", "VALUE").
//		MapDirectory("the_host_current_directory", ".")
//		CaptureStdout().
//	  Finalize(store)
func (wsb *WasiStateBuilder) Finalize(store *Store) (*WasiEnvironment, error) {
	return newWasiEnvironment(store, wsb)
}

func (wsb *WasiStateBuilder) inner() *C.wasi_config_t {
	return wsb._inner
}

// WasiEnvironment represents the environment provided to the WASI
// imports (see NewFunctionWithEnvironment which is designed for
// user-defined host function; that's the same idea here but applied
// to WASI functions and other imports).
type WasiEnvironment struct {
	_inner *C.wasi_env_t
}

func newWasiEnvironment(store *Store, stateBuilder *WasiStateBuilder) (*WasiEnvironment, error) {
	var environment *C.wasi_env_t

	err := maybeNewErrorFromWasmer(func() bool {
		environment = C.wasi_env_new(store.inner(), stateBuilder.inner())

		return environment == nil
	})

	if err != nil {
		return nil, err
	}

	we := &WasiEnvironment{
		_inner: environment,
	}

	runtime.SetFinalizer(we, func(environment *WasiEnvironment) {
		C.wasi_env_delete(environment.inner())
	})

	return we, nil
}

func (we *WasiEnvironment) inner() *C.wasi_env_t {
	return we._inner
}

func buildByteSliceFromCBuffer(buffer *C.char, length int) []byte {
	var byteSlice []byte
	var header = (*reflect.SliceHeader)(unsafe.Pointer(&byteSlice))

	header.Data = uintptr(unsafe.Pointer(buffer))
	header.Len = length
	header.Cap = length

	return byteSlice
}

// ReadStdout reads the WASI module stdout if captured with
// WasiStateBuilder.CaptureStdout
//
//	wasiEnv, _ := NewWasiStateBuilder("test-program").
//		Argument("--foo").
//		Environment("ABC", "DEF").
//		Environment("X", "ZY").
//		MapDirectory("the_host_current_directory", ".").
//		CaptureStdout().
//		Finalize()
//
//	importObject, _ := wasiEnv.GenerateImportObject(store, module)
//	instance, _ := NewInstance(module, importObject)
//	start, _ := instance.Exports.GetWasiStartFunction()
//
//	start()
//
//	stdout := string(wasiEnv.ReadStdout())
func (we *WasiEnvironment) ReadStdout() []byte {
	var buffer *C.char
	length := int(C.to_wasi_env_read_stdout(we.inner(), &buffer))

	return buildByteSliceFromCBuffer(buffer, length)
}

// ReadStderr reads the WASI module stderr if captured with
// WasiStateBuilder.CaptureStderr. See ReadStdout to see an example.
func (we *WasiEnvironment) ReadStderr() []byte {
	var buffer *C.char
	length := int(C.to_wasi_env_read_stderr(we.inner(), &buffer))

	return buildByteSliceFromCBuffer(buffer, length)
}

// GenerateImportObject generates an import object, that can be
// extended and passed to NewInstance.
//
//	wasiEnv, _ := NewWasiStateBuilder("test-program").
//		Argument("--foo").
//		Environment("ABC", "DEF").
//		Environment("X", "ZY").
//		MapDirectory("the_host_current_directory", ".").
//		Finalize()
//
//	importObject, _ := wasiEnv.GenerateImportObject(store, module)
//	instance, _ := NewInstance(module, importObject)
//	start, _ := instance.Exports.GetWasiStartFunction()
//
//	start()
func (we *WasiEnvironment) GenerateImportObject(store *Store, module *Module) (*ImportObject, error) {
	var wasiNamedExterns C.wasmer_named_extern_vec_t
	C.wasmer_named_extern_vec_new_empty(&wasiNamedExterns)

	err := maybeNewErrorFromWasmer(func() bool {
		return !bool(C.wasi_get_unordered_imports(we.inner(), module.inner(), &wasiNamedExterns))
	})

	if err != nil {
		return nil, err
	}

	importObject := NewImportObject()

	numberOfNamedExterns := int(wasiNamedExterns.size)
	firstNamedExtern := unsafe.Pointer(wasiNamedExterns.data)
	sizeOfNamedExtern := unsafe.Sizeof(firstNamedExtern)

	var currentNamedExtern *C.wasmer_named_extern_t

	for nth := 0; nth < numberOfNamedExterns; nth++ {
		currentNamedExtern = *(**C.wasmer_named_extern_t)(unsafe.Pointer(uintptr(firstNamedExtern) + uintptr(nth)*sizeOfNamedExtern))
		module := nameToString(C.wasmer_named_extern_module(currentNamedExtern))
		name := nameToString(C.wasmer_named_extern_name(currentNamedExtern))
		extern := newExtern(C.wasm_extern_copy(C.wasmer_named_extern_unwrap(currentNamedExtern)), nil)

		_, exists := importObject.externs[module]
		if !exists {
			importObject.externs[module] = make(map[string]IntoExtern)
		}

		importObject.externs[module][name] = extern
	}

	C.wasmer_named_extern_vec_delete(&wasiNamedExterns)

	return importObject, nil
}

// Wat2Wasm parses a string as either WAT code or a binary Wasm module.
//
// See https://webassembly.github.io/spec/core/text/index.html.
//
// Note: This is not part of the standard Wasm C API. It is Wasmer specific.
//
//	wat := "(module)"
//	wasm, _ := Wat2Wasm(wat)
//	engine := wasmer.NewEngine()
//	store := wasmer.NewStore(engine)
//	module, _ := wasmer.NewModule(store, wasmBytes)
func Wat2Wasm(wat string) ([]byte, error) {
	var watBytes C.wasm_byte_vec_t
	var watLength = len(wat)

	C.wasm_byte_vec_new(&watBytes, C.size_t(watLength), C.CString(wat))
	defer C.wasm_byte_vec_delete(&watBytes)

	var wasm C.wasm_byte_vec_t

	err := maybeNewErrorFromWasmer(func() bool {
		C.wat2wasm(&watBytes, &wasm)

		return wasm.data == nil
	})

	if err != nil {
		return nil, err
	}

	defer C.wasm_byte_vec_delete(&wasm)

	wasmBytes := C.GoBytes(unsafe.Pointer(wasm.data), C.int(wasm.size))

	return wasmBytes, nil
}
