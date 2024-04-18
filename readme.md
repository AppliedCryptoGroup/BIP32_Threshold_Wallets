## BIP32-Compatible Threshold Wallet Implementation

This repository contains the implementation of the BIP32-compatible threshold derivation scheme in Go as described in the paper "BIP32-Compatible Threshold Wallets" by Das et al. available [here](https://eprint.iacr.org/2023/312.pdf).

## Getting Started
This section describes how to run the tests and benchmarks for BIP32 threshold derivation.
### Tests 

In order to run all tests, run the following command the root directory:
```bash
go test ./...
```
The tests will test the correctness of the implementation of the DDH-based threshold verifiable random function (TVRF) proposed by Galindo et al. ([link](https://eprint.iacr.org/2020/096.pdf))
as well as the correcntess of the derivation of hardened nodes using the TVRF.

### Benchmarks
#### Derivation using a TVRF
To run the benchmarks testing the performance of the derivation of hardened nodes using the TVRF, run the following command:
```bash
go test -bench=. ./derivation/bench
```
Per default, it will test the derivation of 1 hardened node/child with different number of parties and thresholds and a simulated network latency of 10ms.
To change these and other benchmarking parameters, please refer to the `derivation/bench/derivation_bench_test` file.

#### Derivation using MPC
The derivation of hardened nodes using generic multi-party-computation (MPC) for, i.a., evaluating the SHA-512 hash function securely among the parties, is achieved by relying on the [MP-SPDZ](https://github.com/data61/MP-SPDZ) framework.
Please refer to the [official documentation](https://mp-spdz.readthedocs.io/en/latest/) for instructions on how to set up the framework on your machine.
The relevant MPC files are located in the `MPC-SPDZ` directory, to run the benchmarks, please copy and paste them to the location of your local MP-SPDZ framework folder.
Then, just run the following command:
```bash
./benchDerivation.sh
```