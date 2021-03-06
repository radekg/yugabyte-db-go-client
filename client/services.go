package client

import (
	"strings"

	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func loadServiceDefinitions(registry ServiceRegistry) {
	descriptors := []protoreflect.FileDescriptor{
		ybApi.File_yb_cdc_cdc_consumer_proto,
		ybApi.File_yb_cdc_cdc_service_proto,
		ybApi.File_yb_common_common_net_proto,
		ybApi.File_yb_common_common_types_proto,
		ybApi.File_yb_common_common_proto,
		ybApi.File_yb_common_pgsql_protocol_proto,
		ybApi.File_yb_common_redis_protocol_proto,
		ybApi.File_yb_common_transaction_proto,
		ybApi.File_yb_common_value_proto,
		ybApi.File_yb_common_wire_protocol_proto,
		ybApi.File_yb_consensus_consensus_metadata_proto,
		ybApi.File_yb_consensus_consensus_types_proto,
		ybApi.File_yb_consensus_consensus_proto,
		ybApi.File_yb_consensus_log_proto,
		ybApi.File_yb_docdb_docdb_proto,
		ybApi.File_yb_encryption_encryption_proto,
		ybApi.File_yb_fs_fs_proto,
		ybApi.File_yb_master_catalog_entity_info_proto,
		ybApi.File_yb_master_master_admin_proto,
		ybApi.File_yb_master_master_backup_proto,
		ybApi.File_yb_master_master_client_proto,
		ybApi.File_yb_master_master_cluster_proto,
		ybApi.File_yb_master_master_dcl_proto,
		ybApi.File_yb_master_master_ddl_proto,
		ybApi.File_yb_master_master_encryption_proto,
		ybApi.File_yb_master_master_heartbeat_proto,
		ybApi.File_yb_master_master_replication_proto,
		ybApi.File_yb_master_master_types_proto,
		ybApi.File_yb_rocksdb_db_version_edit_proto,
		ybApi.File_yb_rpc_any_proto,
		ybApi.File_yb_rpc_lightweight_message_proto,
		ybApi.File_yb_rpc_rpc_header_proto,
		ybApi.File_yb_rpc_rpc_introspection_proto,
		ybApi.File_yb_rpc_service_proto,
		ybApi.File_yb_server_server_base_proto,
		ybApi.File_yb_tablet_operations_proto,
		ybApi.File_yb_tablet_tablet_metadata_proto,
		ybApi.File_yb_tablet_tablet_types_proto,
		ybApi.File_yb_tablet_tablet_proto,
		ybApi.File_yb_tserver_backup_proto,
		ybApi.File_yb_tserver_remote_bootstrap_proto,
		ybApi.File_yb_tserver_tserver_admin_proto,
		ybApi.File_yb_tserver_tserver_forward_service_proto,
		ybApi.File_yb_tserver_tserver_proto,
		ybApi.File_yb_tserver_tserver_service_proto,
		ybApi.File_yb_tserver_tserver_types_proto,
		ybApi.File_yb_util_histogram_proto,
		ybApi.File_yb_util_opid_proto,
		ybApi.File_yb_util_pb_util_proto,
		ybApi.File_yb_util_version_info_proto,
	}

	for _, fd := range descriptors {
		loadServiceDescriptor(registry, fd)
	}
}

func loadServiceDescriptor(registry ServiceRegistry, descriptor protoreflect.FileDescriptor) {
	services := descriptor.Services()
	for i := 0; i < services.Len(); i = i + 1 {

		svcDescriptor := services.Get(i)

		// discover the service name:
		// in v2.11.2.0, the protobuf service may be optionally annotated
		// with a custom option, for example:
		// master/master_cluster.proto
		// ---------------------------
		//
		// service MasterCluster {
		//  option (yb.rpc.custom_service_name) = "yb.master.MasterService";
		//
		// ---------------------------
		// Here we look up if that option exists, if yes
		// we use it instead of the default service name.

		svcName := string(svcDescriptor.FullName())
		if opts := svcDescriptor.Options(); opts != nil {
			if topts, ok := opts.(*descriptorpb.ServiceOptions); ok && topts != nil {
				topts.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
					if string(fd.FullName()) == "yb.rpc.custom_service_name" {
						svcName = v.String()
					}
					return true
				})
			}
		}

		methods := svcDescriptor.Methods()
		for j := 0; j < methods.Len(); j = j + 1 {
			method := methods.Get(j)
			methodNameParts := strings.Split(string(method.FullName()), ".")
			opName := methodNameParts[len(methodNameParts)-1]
			registry.Register(string(method.Input().FullName()), opName, svcName)
		}
	}
}

// ServiceInfo contains the service and method names used
// by the request header.
type ServiceInfo interface {
	Method() string
	Service() string
	ToRemoteMethodPB() *ybApi.RemoteMethodPB
}

type defaultServiceInfo struct {
	method  string
	service string
}

func (i *defaultServiceInfo) Method() string {
	return i.method
}

func (i *defaultServiceInfo) Service() string {
	return i.service
}

func (i *defaultServiceInfo) ToRemoteMethodPB() *ybApi.RemoteMethodPB {
	return &ybApi.RemoteMethodPB{
		ServiceName: utils.PString(i.Service()),
		MethodName:  utils.PString(i.Method()),
	}
}

// ServiceRegistry contains the information about registered services and inputs.
type ServiceRegistry interface {
	Get(t protoreflect.ProtoMessage) ServiceInfo
	Register(inputTypeName, methodName, svcName string)
}

type defaultServiceRegistry struct {
	registry map[string]ServiceInfo
}

// NewDefaultServiceRegistry returns an initialized
// instance of the default service registry.
func NewDefaultServiceRegistry() ServiceRegistry {
	return &defaultServiceRegistry{
		registry: map[string]ServiceInfo{},
	}
}

func (r *defaultServiceRegistry) Get(t protoreflect.ProtoMessage) ServiceInfo {
	if v, ok := r.registry[string(t.ProtoReflect().Descriptor().FullName())]; ok {
		return v
	}
	return nil
}

func (r *defaultServiceRegistry) Register(inputTypeName, methodName, svcName string) {
	r.registry[inputTypeName] = &defaultServiceInfo{
		method:  methodName,
		service: svcName,
	}
}
