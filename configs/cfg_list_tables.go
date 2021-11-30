package configs

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// OpListTablesConfig represents a command specific config.
type OpListTablesConfig struct {
	flagBase

	NameFilter          string
	Keyspace            string
	ExcludeSystemTables bool
	IncludeNotRunning   bool
	RelationType        []string
}

// NewOpListTablesConfig returns an instance of the command specific config.
func NewOpListTablesConfig() *OpListTablesConfig {
	return &OpListTablesConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpListTablesConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.NameFilter, "name-filter", "", "When used, only returns tables that satisfy a substring match on name_filter")
		c.flagSet.StringVar(&c.Keyspace, "keyspace", "", "The namespace name to fetch info")
		c.flagSet.BoolVar(&c.ExcludeSystemTables, "exclude-system-tables", false, "Exclude system tables")
		c.flagSet.BoolVar(&c.IncludeNotRunning, "include-not-running", false, "Include not running")
		c.flagSet.StringSliceVar(&c.RelationType, "relation-type", supportedRelationType, fmt.Sprintf("Filter tables based on RelationType: %s", strings.Join(supportedRelationType, ", ")))
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpListTablesConfig) Validate() error {
	for _, relation := range c.RelationType {
		var found bool
		for _, opt := range supportedRelationType {
			if opt == relation {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unsupported value '%s' for --relation-type", relation)
		}
	}
	return nil
}
