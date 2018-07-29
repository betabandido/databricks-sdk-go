# ClustersGetResponse

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NumWorkers** | **int32** |  | [optional] [default to null]
**Autoscale** | [***ClustersAutoScale**](ClustersAutoScale.md) |  | [optional] [default to null]
**ClusterId** | **string** |  | [optional] [default to null]
**CreatorUserName** | **string** |  | [optional] [default to null]
**Driver** | [***ClustersSparkNode**](ClustersSparkNode.md) |  | [optional] [default to null]
**Executors** | [**[]ClustersSparkNode**](ClustersSparkNode.md) |  | [optional] [default to null]
**SparkContextId** | **int64** |  | [optional] [default to null]
**JdbcPort** | **int32** |  | [optional] [default to null]
**ClusterName** | **string** |  | [optional] [default to null]
**SparkVersion** | **string** |  | [optional] [default to null]
**SparkConf** | **map[string]string** |  | [optional] [default to null]
**AwsAttributes** | [***ClustersAwsAttributes**](ClustersAwsAttributes.md) |  | [optional] [default to null]
**NodeTypeId** | **string** |  | [optional] [default to null]
**DriverNodeTypeId** | **string** |  | [optional] [default to null]
**SshPublicKeys** | **[]string** |  | [optional] [default to null]
**CustomTags** | **map[string]string** |  | [optional] [default to null]
**ClusterLogConf** | [***ClustersClusterLogConf**](ClustersClusterLogConf.md) |  | [optional] [default to null]
**SparkEnvVars** | **map[string]string** |  | [optional] [default to null]
**AutoterminationMinutes** | **int32** |  | [optional] [default to null]
**EnableElasticDisk** | **bool** |  | [optional] [default to null]
**ClusterSource** | [***ClustersClusterSource**](ClustersClusterSource.md) |  | [optional] [default to null]
**State** | [***ClustersClusterState**](ClustersClusterState.md) |  | [optional] [default to null]
**StateMessage** | **string** |  | [optional] [default to null]
**StartTime** | **int64** |  | [optional] [default to null]
**TerminatedTime** | **int64** |  | [optional] [default to null]
**LastStateLossTime** | **int64** |  | [optional] [default to null]
**LastActivityTime** | **int64** |  | [optional] [default to null]
**ClusterMemoryMb** | **int64** |  | [optional] [default to null]
**ClusterCores** | **float32** |  | [optional] [default to null]
**DefaultTags** | **map[string]string** |  | [optional] [default to null]
**ClusterLogStatus** | [***ClustersLogSyncStatus**](ClustersLogSyncStatus.md) |  | [optional] [default to null]
**TerminationReason** | [***ClustersTerminationReason**](ClustersTerminationReason.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


