package tserver

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/radekg/yugabyte-db-go-client/client"
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/testutils/common"
	"github.com/radekg/yugabyte-db-go-client/testutils/master"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"

	dc "github.com/ory/dockertest/v3/docker"

	// Postgres library:
	_ "github.com/lib/pq"
)

func TestTServerIntegration(t *testing.T) {

	masterTestCtx := master.SetupMasters(t, &common.TestMasterConfiguration{
		ReplicationFactor: 3,
		MasterPrefix:      "tserver-it",
	})
	defer masterTestCtx.Cleanup()

	client := client.NewYBClient(&configs.YBClientConfig{
		MasterHostPort: masterTestCtx.MasterExternalAddresses(),
		OpTimeout:      time.Duration(time.Second * 5),
	})

	common.Eventually(t, 15, func() error {
		if err := client.Connect(); err != nil {
			return err
		}
		return nil
	})

	defer client.Close()

	common.Eventually(t, 15, func() error {
		request := &ybApi.ListMastersRequestPB{}
		response := &ybApi.ListMastersResponsePB{}
		err := client.Execute(request, response)
		if err != nil {
			return err
		}
		t.Log("Received master list", response)
		return nil
	})

	// start a TServer and allocate an additional port:
	otherPort := dc.Port("18080/tcp")
	tserver1Ctx := SetupTServer(t, masterTestCtx, &common.TestTServerConfiguration{
		AdditionalPorts: []dc.Port{otherPort},
		TServerID:       "my-tserver-1",
	})
	defer tserver1Ctx.Cleanup()
	tserver1CtxOtherPorts := tserver1Ctx.OtherPorts()
	if _, ok := tserver1CtxOtherPorts[otherPort]; !ok {
		t.Fatalf("expected additional port '%s' to exist", otherPort)
	}

	// start a TServer:
	tserver2Ctx := SetupTServer(t, masterTestCtx, &common.TestTServerConfiguration{
		TServerID: "my-tserver-2",
	})
	defer tserver2Ctx.Cleanup()

	// start a TServer:
	tserver3Ctx := SetupTServer(t, masterTestCtx, &common.TestTServerConfiguration{
		TServerID: "my-tserver-3",
	})
	defer tserver3Ctx.Cleanup()

	common.Eventually(t, 15, func() error {
		request := &ybApi.ListTabletServersRequestPB{}
		response := &ybApi.ListTabletServersResponsePB{}
		err := client.Execute(request, response)
		if err != nil {
			return err
		}
		t.Log(" ==> Received TServer list", response)
		return nil
	})

	// try YSQL connection:
	t.Logf("connecting to YSQL at 127.0.0.1:%s", tserver1Ctx.TServerExternalYSQLPort())
	db, sqlOpenErr := sql.Open("postgres", fmt.Sprintf("host=127.0.0.1 port=%s user=%s password=%s dbname=%s sslmode=disable",
		tserver1Ctx.TServerExternalYSQLPort(), "yugabyte", "yugabyte", "yugabyte"))
	if sqlOpenErr != nil {
		t.Fatal("failed connecting to YSQL, reason:", sqlOpenErr)
	}
	defer db.Close()
	t.Log("connected to YSQL")

	common.Eventually(t, 15, func() error {
		rows, sqlQueryErr := db.Query("select table_name from information_schema.tables")
		if sqlQueryErr != nil {
			return sqlQueryErr
		}
		nRows := 0
		for {
			if !rows.Next() {
				break
			}
			nRows = nRows + 1
		}
		t.Log("selected", nRows, "rows from YSQL")
		return nil
	}, "querying table via YSQL")

}
