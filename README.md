# YugabyteDB client for go

Work in progress. Current state: PoC.

## Start YugabyteDB in a container

```sh
docker run \
    --rm \
    --hostname localhost \
    -ti \
        -p 7000:7000 \
        -p 7100:7100 \
        -p 5433:5433 \
        local/yugabytedb:2.7.2.0 \
        /home/yugabyte/bin/yb-master \
            --master_addresses=localhost:7100 \
            --rpc_bind_addresses=0.0.0.0:7100 \
            --fs_data_dirs=/tmp \
            --logtostderr --replication_factor=1
```

## Run the client

```
go run ./main.go
```

The output:

```json
[
  {
    "instance_id": {
      "permanent_uuid": "YTUxM2Q4YzdjODM0NDQzZmE0ZmQ2ZjBkOGE1YmU0YmQ=",
      "instance_seqno": 1629644374335479,
      "start_time_us": 1629644374335479
    },
    "registration": {
      "private_rpc_addresses": [
        {
          "host": "localhost",
          "port": 7100
        }
      ],
      "http_addresses": [
        {
          "host": "localhost",
          "port": 7000
        }
      ],
      "cloud_info": {
        "placement_cloud": "cloud1",
        "placement_region": "datacenter1",
        "placement_zone": "rack1"
      }
    },
    "role": 1
  }
]
```
