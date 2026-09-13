# 📈 ResourceFlow

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE.md)

**ResourceFlow** is an algorithmic process scheduling and resource allocation optimization simulator implemented in Go. Designed for supply chain processes, manufacturing assembly lines, and inventory transformation pipelines, ResourceFlow reads dependency configurations, calculates optimal cycle execution orders, enforces resource constraints, and logs detailed timeline events.

---

## ⚡ Key Highlights

- **Process Dependency Graph Parsing**: Parses complex configuration manifests defining initial resource stocks, production duration cycles, input prerequisites, and output yields.
- **Dynamic Virtual Demand Scheduler**: Computes virtual stock demands across cycles to unblock dependent downstream processes and prevent starvation.
- **Cycle-Accurate Event Simulation**: Simulates production cycles step-by-step with log trace output generation (`<config_file>.log`).
- **Resource Inventory Tracking**: Tracks real-time stock balances, in-progress allocations, and final inventory tallies.
- **Automated Verification Engine**: Built-in checker tool (`checker/main.go`) to validate execution logs against original process constraints.

---

## 📋 Table of Contents

- [Key Highlights](#-key-highlights)
- [System Architecture](#-system-architecture)
- [Simulation Scheduling Sequence](#-simulation-scheduling-sequence)
- [Setup & Execution](#-setup--execution)
- [Project Directory Structure](#-project-directory-structure)
- [License](#-license)

---

## 🏗️ System Architecture

```mermaid
graph TD
    A[Configuration Manifest File .txt] --> B[Parser Engine - pkg/parser]
    B --> C[Config & Process Data Models - pkg/models]
    
    C --> D[Scheduler Engine - internal/simulator]
    D --> E1[Initial Stock Allocation]
    D --> E2[Virtual Demand Calculation]
    D --> E3[Process Dependency Validation]
    
    E1 & E2 & E3 --> F[Cycle Loop Simulator]
    F --> G1[STDOUT Console Execution Trace]
    F --> G2[File Log Writer .log]
    F --> G3[Final Stock Inventory Report]
```

---

## 📐 Simulation Scheduling Sequence

```mermaid
sequenceDiagram
    participant User
    participant CLI as resourceflow Binary
    participant Parser as Manifest Parser
    participant Engine as Cycle Scheduler Engine
    participant Log as Log Output (.log)

    User->>CLI: ./resourceflow examples/config1.txt 5.0
    CLI->>Parser: ParseFile("examples/config1.txt")
    Parser-->>CLI: Config Object (InitialStocks, Processes, Optimize)
    CLI->>Engine: Run(Config, TimeoutSeconds)
    
    loop Cycle Simulation Loop (until timeout or completion)
        Engine->>Engine: Resolve finished events for current cycle
        Engine->>Engine: Evaluate process prerequisites vs available stocks
        alt Prerequisites Satisfied
            Engine->>Engine: Deduct inputs & queue future completion cycle
            Engine->>Log: Write cycle execution line "cycle:process_name"
        else Resource Deficit
            Engine->>Engine: Calculate virtual demand & prioritize upstream tasks
        end
        Engine->>Engine: Advance cycle clock to next event
    end
    
    Engine-->>User: Print Final Stock Balances
```

---

## 🚀 Setup & Execution

### Prerequisites

- **Go**: Version 1.20 or newer installed.

---

### Build & Run

1. **Clone Repository**:
   ```bash
   git clone https://github.com/sahmedhusain/resourceflow.git
   cd resourceflow
   ```

2. **Compile Application**:
   ```bash
   go build -o resourceflow main.go
   go build -o checker checker/main.go
   ```

3. **Run Simulation**:
   ```bash
   ./resourceflow examples/example1.txt 5.0
   ```
   *Arguments: `<config_path> <timeout_in_seconds>`*

4. **Verify Simulation Log**:
   ```bash
   ./checker examples/example1.txt examples/example1.txt.log
   ```

---

## 📂 Project Directory Structure

```
resourceflow/
├── main.go                     # Application entrypoint & CLI argument parser
├── go.mod                      # Go module manifest (module resourceflow)
├── README.md                    # Documentation
├── checker/
│   └── main.go                 # Simulation trace validator binary
├── examples/                   # Configuration test scenarios & sample manifests
├── internal/
│   └── simulator/
│       └── scheduler.go        # Cycle scheduler engine & demand resolver
└── pkg/
    ├── models/                 # Process, Config, and Event data structures
    └── parser/                 # Configuration file lexer & parser
```

---

## 📄 License

Distributed under the MIT License. See [LICENSE](LICENSE.md) for details.