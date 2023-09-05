package probes

import (
	"fmt"
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"unsafe"
)

func NewInfoBloxConfiguration() InfoBloxConfiguration {
	return InfoBloxConfiguration{
		Master:   viper.GetString("master"),
		Version:  viper.GetString("wapi_version"),
		Port:     viper.GetInt64("port"),
		Username: viper.GetString("username"),
		Password: viper.GetString("password"),
	}
}

type Member struct {
	ibclient.IBBase          `json:"-"`
	Ref                      string                   `json:"_ref,omitempty"`
	HostName                 string                   `json:"host_name,omitempty"`
	ConfigAddrType           string                   `json:"config_addr_type,omitempty"`
	PLATFORM                 string                   `json:"platform,omitempty"`
	ServiceTypeConfiguration string                   `json:"service_type_configuration,omitempty"`
	Nodeinfo                 []ibclient.NodeInfo      `json:"node_info,omitempty"`
	TimeZone                 string                   `json:"time_zone,omitempty"`
	ServiceStatus            []ibclient.ServiceStatus `json:"service_status,omitempty"`
}

func NewMember(nodeName string) *Member {
	var res Member
	res.HostName = nodeName
	p := &ibclient.IBBase{}
	p1 := (*string)(unsafe.Pointer(p))

	*p1 = "member"
	ptrSize := unsafe.Sizeof(string("member"))
	*(*[]string)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(ptrSize))) =
		[]string{"extattrs", "host_name", "node_info", "service_status"}
	res.IBBase = *p
	return &res
}

type Range struct {
	ibclient.IBBase
	Ref string `json:"_ref,omitempty"`
	//NetviewName string `json:"network_view,omitempty"`
	Cidr        string      `json:"network,omitempty"`
	Ea          ibclient.EA `json:"extattrs"`
	Comment     string      `json:"comment"`
	Utilization int64       `json:"dhcp_utilization"`
}

func NewRange(netview string, cidr string, isIPv6 bool, comment string, ea ibclient.EA) *Range {
	var res Range
	//res.NetviewName = netview
	res.Cidr = cidr
	res.Ea = ea
	res.Comment = comment
	p := &ibclient.IBBase{}
	p1 := (*string)(unsafe.Pointer(p))

	*p1 = "range"
	ptrSize := unsafe.Sizeof(string("range"))
	*(*[]string)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(ptrSize))) =
		[]string{"extattrs", "network", "dhcp_utilization", "comment"}
	res.IBBase = *p
	return &res
}

type InfoBloxConfiguration struct {
	Master   string
	Version  string
	Port     int64
	Username string
	Password string
}

type InfoBloxApi struct {
	Conn *ibclient.Connector
}

func NewInfobloxApi() InfoBloxApi {
	x := NewInfoBloxConfiguration()

	hostConfig := ibclient.HostConfig{
		Host:    x.Master,
		Version: x.Version,
		//Port:    strconv.FormatInt(x.Port, 10),
	}

	authConfig := ibclient.AuthConfig{
		Username:   x.Username,
		Password:   x.Password,
		ClientCert: nil,
		ClientKey:  nil,
	}
	transportConfig := ibclient.NewTransportConfig("false", 20, 10)
	requestBuilder := &ibclient.WapiRequestBuilder{}
	requestor := &ibclient.WapiHttpRequestor{}
	conn, err := ibclient.NewConnector(hostConfig, authConfig, transportConfig, requestBuilder, requestor)
	if err != nil {
		fmt.Printf("%s", fmt.Errorf("Failed to connect", err))
	}

	return InfoBloxApi{Conn: conn}
}

func (i InfoBloxApi) GetDhcpUtilization(network string) (Range, error) {

	var res []Range
	net := NewRange("", network, false, "", nil)

	queryAttribute := map[string]string{"network": network}
	qp := ibclient.NewQueryParams(false, queryAttribute)
	err := i.Conn.GetObject(net, "", qp, &res)

	if err != nil {
		log.Error(fmt.Sprintf("%s", fmt.Errorf("Failed to get network", err)))
		return *net, err
	}

	return res[0], nil
}

func (i InfoBloxApi) GetMember(nodeName string) (Member, error) {

	var res []Member
	net := NewMember(nodeName)

	queryAttribute := map[string]string{"host_name": nodeName}
	qp := ibclient.NewQueryParams(false, queryAttribute)
	err := i.Conn.GetObject(net, "", qp, &res)

	if err != nil {
		log.Error(fmt.Sprintf("%s", fmt.Errorf("Failed to get node", err)))
		return *net, err
	}
	return res[0], nil
}

func (i InfoBloxApi) GetMemberX(nodeName string) {

	var res []ibclient.Member
	member := ibclient.Member{
		IBBase:                   ibclient.IBBase{},
		Ref:                      "",
		HostName:                 nodeName,
		ConfigAddrType:           "",
		PLATFORM:                 "",
		ServiceTypeConfiguration: "",
		Nodeinfo:                 nil,
		TimeZone:                 "",
	}
	net := ibclient.NewMember(member)

	queryAttribute := map[string]string{"host_name": nodeName}
	qp := ibclient.NewQueryParams(false, queryAttribute)
	err := i.Conn.GetObject(net, "", qp, &res)

	//res, err := objMgr.GetNetwork("", network, false, nil)
	if err != nil {
		fmt.Printf("%s", fmt.Errorf("Failed to get network", err))
	}

}

func (i InfoBloxApi) Logout() {
	i.Conn.Logout()
}
