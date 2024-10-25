# Redis High Availability with Sentinel, HAProxy, and RedisInsight

This setup provides a high-availability Redis cluster with three primary components:
- Redis Sentinel for monitoring and automatic failover.
- HAProxy for load balancing Redis Sentinel.
- RedisInsight for monitoring and visualizing Redis data.

## Architecture Overview

1. **Redis Master and Slaves**:
   - A Redis master node (`redis-master`) and three Redis slave nodes (`redis-slave-1`, `redis-slave-2`, `redis-slave-3`) are set up for data replication.
   - Redis slaves replicate data from the master, ensuring data redundancy and availability.

2. **Redis Sentinel**:
   - Three Redis Sentinel nodes (`sentinel-1`, `sentinel-2`, `sentinel-3`) monitor the Redis master.
   - In case the master fails, Sentinel initiates a failover and promotes one of the slaves as the new master.
   - HAProxy communicates with the Redis Sentinel nodes to handle routing and failover.

3. **HAProxy**:
   - HAProxy is configured for load balancing Redis Sentinel instances, ensuring traffic is routed to the active master node.
   - HAProxy also provides a statistics page accessible at port `1936`.

4. **RedisInsight**:
   - RedisInsight is a tool for visualizing, monitoring, and managing Redis data.
   - Accessible at port `5540`.

## Service Breakdown

| Service         | Container Name | Role                                 | IP Address   | Port      |
|-----------------|----------------|--------------------------------------|--------------|-----------|
| Redis Master    | `redis-master` | Main Redis instance                  | 172.38.0.11  | 6379      |
| Redis Slave 1   | `redis-slave-1`| Redis replica of the master          | 172.38.0.12  | 6379      |
| Redis Slave 2   | `redis-slave-2`| Redis replica of the master          | 172.38.0.13  | 6379      |
| Redis Slave 3   | `redis-slave-3`| Redis replica of the master          | 172.38.0.18  | 6379      |
| Sentinel 1      | `sentinel-1`   | Monitors master for failover         | 172.38.0.14  | 26379     |
| Sentinel 2      | `sentinel-2`   | Monitors master for failover         | 172.38.0.15  | 26379     |
| Sentinel 3      | `sentinel-3`   | Monitors master for failover         | 172.38.0.16  | 26379     |
| RedisInsight    | `redisinsight` | GUI for Redis management             | 172.38.0.17  | 5540      |
| HAProxy         | `haproxy`      | Load balances Sentinel connections   | Dynamic      | 26379,1936|

## Prerequisites

- Docker and Docker Compose installed on your server.
- The directory structure for volumes (e.g., `./data`) created in your project directory.

## Setup Instructions

1. **Clone the Repository**:
   Clone this repository or copy the provided configuration files.

   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. **Directory Structure**:
   - Make sure you have a `data` directory with subdirectories for each Redis instance (`master`, `slave1`, `slave2`, `slave3`) to persist data.

3. **Run Docker Compose**:
   Use Docker Compose to start the containers. This will automatically bring up Redis, Sentinel, HAProxy, and RedisInsight.

   ```bash
   docker-compose up -d
   ```

4. **Verify the Setup**:
   - Access RedisInsight by visiting `http://<host-ip>:5540` to monitor your Redis data.
   - Access HAProxy statistics at `http://<host-ip>:1936` for an overview of load balancing.
   - Verify Redis Sentinel status by connecting to Sentinel on `localhost:26379`.

## Docker Compose Configuration

The main components of the Docker Compose file are as follows:

### Redis Master and Slaves
Each Redis node is configured with unique IP addresses. Slaves replicate from the master (172.38.0.11). Redis is set to `--appendonly yes` for data durability.

### Redis Sentinel
Each Sentinel instance is configured with `sentinel monitor` pointing to the master (172.38.0.11:6379). The quorum is set to 2, meaning at least 2 Sentinel nodes must agree before failover occurs.

### HAProxy
HAProxy listens on port 26379 and balances traffic between the three Sentinel nodes, checking their status every 5 seconds.

### RedisInsight
RedisInsight is exposed on port `5540`, allowing you to connect and manage your Redis data visually.

## HAProxy Configuration (`haproxy.cfg`)

The `haproxy.cfg` file includes:
- A `listen` section for Redis Sentinel on port `26379`, using `roundrobin` load balancing.
- A `listen` section for HAProxy statistics on port `1936`.

```cfg
global
    log 127.0.0.1 local1
    maxconn 4096

defaults
    log global
    mode tcp
    option tcplog
    retries 3
    option redispatch
    maxconn 2000
    timeout connect 5000ms
    timeout client 50000ms
    timeout server 50000ms

listen stats
    bind *:1936
    mode http
    stats enable
    stats hide-version
    stats realm Haproxy\ Statistics
    stats uri /

listen sentinel
    bind *:26379
    mode tcp
    balance roundrobin
    timeout client 3h
    timeout server 3h
    option clitcpka
    server sentinel1 172.38.0.14:26379 check inter 5s rise 2 fall 3
    server sentinel2 172.38.0.15:26379 check inter 5s rise 2 fall 3
    server sentinel3 172.38.0.16:26379 check inter 5s rise 2 fall 3
```

## Failover Testing

To simulate a failover:
1. Stop the `redis-master` container.
2. Sentinel should detect the failure and promote a slave as the new master.
3. Restart `redis-master`, which will rejoin the cluster as a slave.