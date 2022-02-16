package client

import (
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestServiceDescribe(t *testing.T) {

	svcRegistry := NewDefaultServiceRegistry()
	loadServiceDefinitions(svcRegistry)

	t.Run("it=handles yb.master.MasterService payloads", func(tt *testing.T) {
		payload := &ybApi.ListMastersRequestPB{}
		svcInfo := svcRegistry.Get(payload)
		assert.NotNil(t, svcInfo)
		assert.Equal(t, svcInfo.Method(), "ListMasters")
		assert.Equal(t, "yb.master.MasterService", svcInfo.Service())
	})

	t.Run("it=handles yb.server.GenericService payloads", func(tt *testing.T) {
		payload := &ybApi.PingRequestPB{}
		svcInfo := svcRegistry.Get(payload)
		assert.NotNil(t, svcInfo)
		assert.Equal(t, svcInfo.Method(), "Ping")
		assert.Equal(t, svcInfo.Service(), "yb.server.GenericService")
	})

	t.Run("it=handles yb.tserver.TabletServerService payloads", func(tt *testing.T) {
		payload := &ybApi.IsTabletServerReadyRequestPB{}
		svcInfo := svcRegistry.Get(payload)
		assert.NotNil(t, svcInfo)
		assert.Equal(t, svcInfo.Method(), "IsTabletServerReady")
		assert.Equal(t, svcInfo.Service(), "yb.tserver.TabletServerService")
	})

}
