module github.com/radekg/yugabyte-db-go-client

go 1.16

require (
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-hclog v0.16.2
	// pq used in tests:
	github.com/lib/pq v1.9.0
	// dockertest/v3 used in tests:
	github.com/ory/dockertest/v3 v3.8.1
	github.com/radekg/yugabyte-db-go-proto/v2 v2.11.2
	github.com/stretchr/testify v1.7.0
	google.golang.org/protobuf v1.27.1
)
