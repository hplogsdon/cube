# Build An Orchestrator in Go (From Scratch)

## Overview
Educational example of a container orchestrator in GoLang from the [Manning Book][manning-book], *Build an Orchestrator*,
designed to manage and run Docker containers on Linux systems. It can automate the scheduling, execution, and monitoring
of containerized tasks across multiple worker nodes.

## Architecture
* **Manager**: Central component responsible for
  * **Scheduling**: Evaluates task feasibility. Scores tasks based on resource availability, and prepares tasks for work submission.
  * **API**: Provides a Web API for task submission, queries, and metrics.
  * **Workers**: Manages the health and availability of compute resources.
  * **Metrics**: Collects and exposes performance metrics for system health.
* **Workers**: Compute resources that can execute container jobs. Receives task configuration from the Manager via API.

## Features
* **Scheduling**: Automatically schedules container tasks on any worker node based on resource availability.
* **Scaling**: Supports scaling tasks to the appropriate level for the available compute resources.
* **Monitoring**: Provides metrics for performance and system health.

## Installation
### Prerequisites
* Linux VMs or servers running appropriate OS.
* Docker installed on all VMs

### Steps
1. **Clone**
   ```
   git clone https://github.com/hplogsdon/cube.git
   cd cube
   ```
   
2. **Setup Management**
   ```
   ./bin/manager
   ```

3. **Setup Workers**
   ```
   ./bin/worker
   ```
   
4. **Validate**
   ```
   curl -XGET http://localhost:5555/workers
   ```

## Usage
Don't use this.

## Contributing
Contributions are not really welcome. This is an educational example based on the Manning Book *[Building an
Orchestrator (From Scratch)][manning-book]* by Tim Boring.

If you want to use this, buy the book (It's a great book!) and feel free to follow along with this repo as a reference.

[manning-book]: https://www.manning.com/books/build-an-orchestrator-in-go-from-scratch
[TebogoTS]: https://github.com/TebogoTS/cube
[kocierik]: https://github.com/kocierik/cube