# ShardBridge

**ShardBridge** is a high-performance middleware layer that enables horizontal scaling of MySQL databases. It allows you to distribute your data across multiple MySQL instances while maintaining the familiar MySQL interface.



## 🚀 Features

- **Transparent Sharding**: Distribute data across multiple MySQL servers without changing your application logic  
- **Intelligent Query Routing**: Automatically route queries to the appropriate shard  
- **Multiple Sharding Strategies**: Support for hash-based and range-based sharding  
- **Connection Pooling**: Efficient management of database connections  
- **Read/Write Splitting**: Route read queries to replicas and writes to masters  
- **High Availability**: Automatic failover capabilities  
- **Dynamic Scaling**: Add or remove shards without downtime  
- **Cross-Shard Queries**: Support for queries that span multiple shards  
- **Data Rebalancing**: Tools for redistributing data when adding new shards  
- **Simple Integration**: Standard SQL interface, compatible with existing applications  

---

## 📋 Table of Contents

- [Installation](#-installation)  
- [Quick Start](#-quick-start)  
- [Architecture](#-architecture)  
- [Configuration](#️-configuration)  
- [Sharding Strategies](#-sharding-strategies)  
- [API Reference](#-api-reference)  
- [Scaling Guidelines](#-scaling-guidelines)  
- [Monitoring](#-monitoring)  
- [Troubleshooting](#-troubleshooting)  
- [Performance Tuning](#-performance-tuning)  
- [Contributing](#-contributing)  
- [License](#license)

---

## 🔧 Installation

### Using Go

```bash
go get github.com/yourusername/mysqlscaler

From Source

git clone https://github.com/yourusername/mysqlscaler.git
cd mysqlscaler
go build -o mysqlscaler ./cmd/server

Docker
docker pull yourusername/mysqlscaler:latest

🚀 Quick Start

1. Set up MySQL instances

Set up multiple MySQL instances using Docker, cloud, or bare metal.

2. Create a configuration file

sharding:
  strategy: "consistent-hash"
  virtualNodes: 10

metadata:
  store: "etcd"
  endpoint: "localhost:2379"

connection:
  maxOpenConns: 50
  maxIdleConns: 10
  connTimeout: "5s"
  idleTimeout: "60s"

shards:
  - id: "shard1"
    masterDSN: "user:password@tcp(mysql-shard1-master:3306)/mydb"
    replicaDSNs:
      - "user:password@tcp(mysql-shard1-replica1:3306)/mydb"
      - "user:password@tcp(mysql-shard1-replica2:3306)/mydb"
  - id: "shard2"
    masterDSN: "user:password@tcp(mysql-shard2-master:3306)/mydb"
    replicaDSNs:
      - "user:password@tcp(mysql-shard2-replica1:3306)/mydb"
      - "user:password@tcp(mysql-shard2-replica2:3306)/mydb"

tables:
  users:
    shardKeys: ["user_id"]
  orders:
    shardKeys: ["user_id"]
  products:
    shardKeys: ["product_id"]

3. Initialize your database schema

mysql -h mysql-shard1-master -u user -p mydb < schema.sql
mysql -h mysql-shard2-master -u user -p mydb < schema.sql

4. Start MySQLScaler
./mysqlscaler --config=config.yaml

5. Connect your application

package main

import (
    "context"
    "log"
    "github.com/yourusername/mysqlscaler/client"
)

func main() {
    db, err := client.Connect("localhost:3000")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    results, err := db.Query(context.Background(), 
                            "SELECT * FROM users WHERE user_id = ?", 12345)
    if err != nil {
        log.Fatal(err)
    }

    for results.Next() {
        // Process each row
    }
}

Standard SQL Driver

package main

import (
    "database/sql"
    "log"
    _ "github.com/yourusername/mysqlscaler/driver"
)

func main() {
    db, err := sql.Open("mysqlscaler", "host=localhost:3000")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT * FROM users WHERE user_id = ?", 12345)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        // Process each row
    }
}



🏗 Architecture
	•	Shard Manager
	•	Query Router
	•	Connection Pool
	•	Query Processor
	•	Health Monitor
	•	Metadata Store

⚙️ Configuration

server:
  host: "0.0.0.0"
  port: 3000
  maxConnections: 1000
  requestTimeout: "30s"

sharding:
  strategy: "consistent-hash"
  virtualNodes: 10

connection:
  maxOpenConns: 50
  maxIdleConns: 10
  connTimeout: "5s"
  idleTimeout: "60s"

metadata:
  store: "etcd"
  endpoint: "localhost:2379"
  ttl: "24h"

logging:
  level: "info"
  format: "json"
  output: "stdout"

🧩 Sharding Strategies

sharding:
  strategy: "consistent-hash"
  virtualNodes: 10

Range-Based Sharding
sharding:
  strategy: "range"
  ranges:
    - min: ""
      max: "500000"
      shard: "shard1"
    - min: "500001"
      max: "1000000"
      shard: "shard2"
    - min: "1000001"
      max: ""
      shard: "shard3"

📚 API Reference

Client API
result, err := client.Query(ctx, "SELECT * FROM users WHERE user_id = ?", 12345)

affectedRows, err := client.Execute(ctx, "UPDATE users SET name = ? WHERE user_id = ?", "John", 12345)

err := client.Transaction(ctx, func(tx *client.Tx) error {
    _, err := tx.Execute("INSERT INTO orders (user_id, product_id) VALUES (?, ?)", 12345, 67890)
    if err != nil {
        return err
    }
    _, err = tx.Execute("UPDATE inventory SET stock = stock - 1 WHERE product_id = ?", 67890)
    return err
})

Management API

err := admin.AddShard(ctx, &shard.ShardInfo{
    ID: "shard3",
    MasterDSN: "user:password@tcp(mysql-shard3-master:3306)/mydb",
    ReplicaDSNs: []string{
        "user:password@tcp(mysql-shard3-replica1:3306)/mydb",
    },
})

err := admin.RemoveShard(ctx, "shard3")

err := admin.RebalanceData(ctx, nil)


📈 Scaling Guidelines


Choosing Shard Keys
	•	High cardinality
	•	Even distribution
	•	Frequently used in queries
	•	Rarely updated

When to Add Shards
	•	Shard size > 75% of optimal capacity
	•	Increased latency
	•	High CPU or I/O load

Data Rebalancing

./mysqlscaler shard rebalance --throttle=50MB/s --time-window="22:00-06:00"


📊 Monitoring

Exposed Prometheus metrics at /metrics:
	•	mysqlscaler_query_latency
	•	mysqlscaler_queries_per_second
	•	mysqlscaler_shard_connections
	•	mysqlscaler_shard_errors
	•	mysqlscaler_cross_shard_queries

Sample Grafana dashboards available in the monitoring/ directory.

⚡ Performance Tuning
connection:
  maxOpenConns: 50
  maxIdleConns: 10
  connTimeout: "5s"
  idleTimeout: "60s"

	•	Include shard keys in queries
	•	Use batch operations
	•	Denormalize related data if needed
	•	Rebalance shards proactively

🤝 Contributing

We welcome your contributions! See CONTRIBUTING.md for instructions.
git clone https://github.com/yourusername/mysqlscaler.git
cd mysqlscaler
go mod download
go test ./...
golangci-lint run
