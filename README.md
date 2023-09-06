go-infoblox-exporter
----------------------
# Overview
The infoblox-exporter collect metrics from an infoblox master.
Currently, two types of metrics is supported:
- Member service and member node service managed by the master.
- DHCP utilization based on networks

# Metrics
## Members 
Service, member or nodes, are reported as a gauge state 1=WORKING, 0=FAILED, 2=UNKNOWN. 
Any services in the INACTIVE state are not included. 
For node services the label `node_ip` is added. If the node is part of a HA setup the value is an
ip address, if not the value is `NO_HA_IP`.

Example output for a member that have HA setup:
```shell
curl 'localhost:9597/probe?target=infoblox.master.com&modules=member_services'
```
```text
# HELP infoblox_member_node_info Node info
# TYPE infoblox_member_node_info gauge
infoblox_member_node_info{ha_status="ACTIVE",hwid="1405202001701727",hwtype="IB-1415",node_ip="140.166.34.152",platform="PHYSICAL"} 1
infoblox_member_node_info{ha_status="PASSIVE",hwid="1405201903700510",hwtype="IB-1415",node_ip="140.166.34.151",platform="PHYSICAL"} 1
# HELP infoblox_member_node_service Node service (0=Failed, 1=Working, 2=Unknown)
# TYPE infoblox_member_node_service gauge
infoblox_member_node_service{node_ip="140.166.34.151",service="CORE_FILES"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="CPU1_TEMP"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="CPU_USAGE"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="DB_OBJECT"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="DISCOVERY_CAPACITY"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="DISK_USAGE"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="ENET_HA"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="ENET_LAN"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="FAN1"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="FAN2"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="FAN3"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="FAN4"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="FAN5"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="FAN6"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="MEMORY"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="NODE_STATUS"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="NTP_SYNC"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="POWER1"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="POWER2"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="REPLICATION"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="SWAP_USAGE"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="SYS_TEMP"} 1
infoblox_member_node_service{node_ip="140.166.34.151",service="VPN_CERT"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="CORE_FILES"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="CPU1_TEMP"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="CPU_USAGE"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="DB_OBJECT"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="DISCOVERY_CAPACITY"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="DISK_USAGE"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="ENET_HA"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="ENET_LAN"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="FAN1"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="FAN2"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="FAN3"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="FAN4"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="FAN5"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="FAN6"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="MEMORY"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="NODE_STATUS"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="NTP_SYNC"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="OSPF"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="POWER1"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="POWER2"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="REPLICATION"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="SWAP_USAGE"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="SYS_TEMP"} 1
infoblox_member_node_service{node_ip="140.166.34.152",service="VPN_CERT"} 1
# HELP infoblox_member_service Service (0=Failed, 1=Working, 2=Unknown)
# TYPE infoblox_member_service gauge
infoblox_member_service{service="DNS"} 1
infoblox_member_service{service="HSM"} 2
infoblox_member_service{service="IMC_DCA_BWL"} 2
infoblox_member_service{service="NTP"} 1
infoblox_member_service{service="REPORTING"} 1
# HELP probe_duration_seconds How many seconds the probe call took to complete
# TYPE probe_duration_seconds gauge
probe_duration_seconds 3.205450828
# HELP probe_success Probe call success (1=Up,0=Down)
# TYPE probe_success gauge
probe_success 1

```
In addition to all the service metrics there is also `infoblox_node_info` with additional metadata 
labels. Metrics value is always 1.0

The `probe_success` is set to 1.0 if the exporter could connect to the Infoblox master and that the 
member exists.

## DHCP utilization
For a specific network that the infoblox master manage the metrics show the utilization of DCHP 
addresses. This can be valuable to alert on if the metrics is close to 1.0, 100 % utilization  

```shell
curl 'localhost:9597/probe?target=10.199.73.128/26&modules=dhcp_utilization'
```
```text
 
# HELP infoblox_dhcp_utilization_ratio Dhcp utilization
# TYPE infoblox_dhcp_utilization_ratio gauge
infoblox_dhcp_utilization_ratio 0.48
# HELP probe_duration_seconds How many seconds the probe call took to complete
# TYPE probe_duration_seconds gauge
probe_duration_seconds 2.185153276
# HELP probe_success Probe call success (1=Up,0=Down)
# TYPE probe_success gauge
probe_success 1
```
The `probe_success` is set to 1.0 if the exporter could connect to the Infoblox master and that the
network exists.

# Discovery 
Please see the [infoblox-discovery](https://github.com/thenodon/infoblox_discovery)
to get dynamic Prometheus discovery configuration for   

# Environment variables

The following variables ar mandatory to set.

- BASIC_AUTH_USERNAME - the basic auth username to the exporter
- BASIC_AUTH_PASSWORD - the basic auth password to the exporter 
- INFOBLOX_MASTER - the ip/fqdn to the infoblox server
- INFOBLOX_WAPI_VERSION - the Infoblox master api version
- INFOBLOX_USERNAME - the Infoblox master username
- INFOBLOX_PASSWORD  - the Infoblox master password

The following are optional
- EXPORTER_HOST - default to `0.0.0.0`
- EXPORTER_PORT - default to `9597`
- EXPORTER_LOG_LEVEL - default to `INFO`

# Test
```
curl -s 'localhost:9597/probe?target=host.foo.com&module=member_services' 
```

The `type` can have the following values:
- member_services - the target is infoblox member
- dhcp_utilization - the target has to be network like `10.121.151.128/26`

# Build

```shell
go build .
```