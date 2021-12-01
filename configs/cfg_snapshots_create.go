package configs

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// OpSnapshotCreateConfig represents a command specific config.
type OpSnapshotCreateConfig struct {
	flagBase

	Keyspace   string
	TableNames []string
	TableUUIDs []string
	ScheduleID string
}

// NewOpSnapshotCreateConfig returns an instance of the command specific config.
func NewOpSnapshotCreateConfig() *OpSnapshotCreateConfig {
	return &OpSnapshotCreateConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotCreateConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Keyspace, "keyspace", "", "Keyspace for the tables in this create request")
		c.flagSet.StringSliceVar(&c.TableNames, "name", []string{}, "Table names to create snapshots for")
		c.flagSet.StringSliceVar(&c.TableUUIDs, "uuid", []string{}, "Table IDs to create snapshots for")
		c.flagSet.StringVar(&c.ScheduleID, "schedule-id", "", "Create snapshot to this schedule, other fields are ignored")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpSnapshotCreateConfig) Validate() error {
	if len(c.ScheduleID) > 0 && c.Keyspace == "" {
		return fmt.Errorf("--keyspace is required")
	}
	if c.Keyspace != "" {
		if strings.HasPrefix(c.Keyspace, "yedis.") {
			return fmt.Errorf("--keyspace yedis.* not supported")
		}
		if !strings.HasPrefix(c.Keyspace, "ycql.") && !strings.HasPrefix(c.Keyspace, "ysql.") {
			// set default keyspace type:
			c.Keyspace = fmt.Sprintf("ycql.%s", c.Keyspace)
		}
		if strings.HasPrefix(c.Keyspace, "ysql.") {
			if len(c.TableNames) > 0 || len(c.TableUUIDs) > 0 {
				return fmt.Errorf("--keyspace ysql.* does not support explicit table selection, remove any --name and --uuid")
			}
		}
	}
	return nil
}
