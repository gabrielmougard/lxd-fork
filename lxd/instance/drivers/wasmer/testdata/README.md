# How to compile basic Rust binaries/libs to WASM/WASI bytecode

To  compile `tests.rs` and `wasi.rs` you can do:

```bash
# Install `rustup` and `rustc` if not already done
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
# Install the `wasm32-unknown-unknown` and `wasm32-wasi` triples
rustup target add wasm32-unknown-unknown
rustup target add wasm32-wasi

# Install `wasm-strip` and `wasm-opt` tools for further optimization
sudo apt install wabt binaryen -y

# Compile and optimize `tests.rs`
rustc --target wasm32-unknown-unknown -O tests.rs --crate-type=cdylib -o tests.raw.wasm
wasm-strip tests.raw.wasm
wasm-opt -O4 -Oz tests.raw.wasm -o tests.wasm

# Compile and optimize the `wasix-test` crate (in order to test WASI and WASIX capabilities)
# WASI is a subset of WASIX, so we can use the same crate to test both
cd wasix-test && cargo build --target wasm32-wasi --release # This will generate `target/wasm32-wasi/release/wasix_test.wasm`
# WASIX extends WASI.
cd waasi-test && cargo wasix build --release # This will generate `target/wasm32-wasmer-wasi/release/waasi_test.wasm`
```