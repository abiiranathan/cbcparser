# cbcparser

Read and parse TSV/CSV files exported from CBC Machines.

Written with love in Go.

Suppported Machines:

- Human 30
- Edan Pro 30

Use cases:

  - You want to import all CBC results into a LIS(Laboratory Information System).
  - Store all CBC records on a server and query results by sample_id(SID) or other parameters for research purposes.
  - View backed up data and reprint lost patient CBC record(based on Patent ID).

Installation:
```bash
go get github.com/abiiranathan/cbcparser/cbcparser/v1.0.1
```

Usage:
See examples for usage.

Run examples:

### Human -  Single report

```bash
go run examples/human/single/main.go sample_data/human.txt sample_data/normal_ranges.json
```

### Human - Multiple reports

```bash
go run examples/human/multi/main.go sample_data/human.txt sample_data/normal_ranges.json
```


### Edan - single report

```bash
go run examples/edan/single/main.go sample_data/edan.csv sample_data/normal_ranges.json
```