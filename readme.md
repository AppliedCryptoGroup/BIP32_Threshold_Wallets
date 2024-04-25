## BIP32-Compatible Threshold Wallet Implementation

This repository contains the implementation of the BIP32-compatible threshold derivation scheme in Go as described in the paper "BIP32-Compatible Threshold Wallets".

## Getting Started
This section describes how to run the tests and benchmarks for BIP32 threshold derivation.
### Tests 

In order to run all tests, run the following command the root directory:
```bash
go test ./...
```
The tests will test the correctness of the implementation of the DDH-based threshold verifiable random function (TVRF) proposed by Galindo et al. ([eprint link](https://eprint.iacr.org/2020/096.pdf))
as well as the correctness of the derivation of hardened nodes using the TVRF.

### Benchmarks
#### Derivation using a TVRF
To run the benchmarks testing the performance of the derivation of hardened nodes using the TVRF, run the following command:
```bash
go test -bench=. ./derivation/bench
```
Per default, it will test the derivation of 1 hardened node/child with different numbers of parties and thresholds and a simulated network latency of 10ms.
To change these and other benchmarking parameters, please refer to the `derivation/bench/derivation_bench_test` file.

#### Derivation using MPC
We evaluated a generic multi-party-computation (MPC) approach for hardened-node derivation based on BIP32, where SHA-512 hash function is evaluated. 
Specifically, we relied on the [MP-SPDZ](https://github.com/data61/MP-SPDZ) framework for this implementation.
In this approach, all the parties of the parent node jointly run an MPC protocol to evaluate the hashing function and derive a child private key.
We abstracted the implementation for simplicity, which only reduces the complexity of the MPC-based derivation compared to the standard BIP32 definition.
The [Bristol-Fashioned Circuit](https://nigelsmart.github.io/MPC-Circuits/) SHA-512 is used in the implementation, and we evaluated this approach on ```malicious-shamir-party```. The settings and results can be found in the Evaluation section of the paper.
For detailed instructions on using MP-SPDZ, please refer to the [official documentation](https://mp-spdz.readthedocs.io/en/latest/).

The relevant MPC files are located in the `MPC-based` directory, to run this derivation and get the benchmarks, please first install MP-SPDZ on your machine, copy and paste them to the corresponding location under your local MP-SPDZ framework folder, i.e. ```mp-spdz-0.3.8```.
The ```sha512.mpc``` should be put in ```mp-spdz-0.3.8\Programs\Source```, while the ```benchDerivation.sh``` should be put in ```mp-spdz-0.3.8``` directly.

After compiling the virtual machine, i.e. ```make malicious-shamir-party.x```, just run the following command:
```bash
./benchDerivation.sh
```
