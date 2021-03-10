package iperf

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

type StreamInterval struct {
	Streams []*StreamIntervalReport  `json:"streams"`
	Sum     *StreamIntervalSumReport `json:"sum"`
}

func (si *StreamInterval) String() string {
	b, err := json.Marshal(si)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type StreamIntervalReport struct {
	Socket           int     `json:"socket"`
	StartInterval    float32 `json:"start"`
	EndInterval      float32 `json:"end"`
	Seconds          float32 `json:"seconds"`
	Bytes            int     `json:"bytes"`
	BitsPerSecond    float64 `json:"bits_per_second"`
	Retransmissions  int     `json:"retransmissions"`
	CongestionWindow int     `json:"congestion_window"`
	Omitted          bool    `json:"omitted"`
}

func (sir *StreamIntervalReport) String() string {
	b, err := json.Marshal(sir)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type StreamIntervalSumReport struct {
	StartInterval float32 `json:"start"`
	EndInterval   float32 `json:"end"`
	Seconds       float32 `json:"seconds"`
	Bytes         int     `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	Omitted       bool    `json:"omitted"`
}

func (sisr *StreamIntervalSumReport) String() string {
	b, err := json.Marshal(sisr)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type StreamEndReport struct {
	Sender   TcpStreamEndReport `json:"sender"`
	Receiver TcpStreamEndReport `json:"receiver"`
	Udp      UdpStreamEndReport `json:"udp"`
}

func (ser *StreamEndReport) String() string {
	b, err := json.Marshal(ser)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type UdpStreamEndReport struct {
	Socket        int     `json:"socket"`
	Start         float32 `json:"start"`
	End           float32 `json:"end"`
	Seconds       float32 `json:"seconds"`
	Bytes         int     `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	JitterMs      float32 `json:"jitter_ms"`
	LostPackets   int     `json:"lost_packets"`
	Packets       int     `json:"packets"`
	LostPercent   float32 `json:"lost_percent"`
	OutOfOrder    int     `json:"out_of_order"`
}

func (user *UdpStreamEndReport) String() string {
	b, err := json.Marshal(user)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type TcpStreamEndReport struct {
	Socket        int     `json:"socket"`
	Start         float32 `json:"start"`
	End           float32 `json:"end"`
	Seconds       float32 `json:"seconds"`
	Bytes         int     `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
}

func (tser *TcpStreamEndReport) String() string {
	b, err := json.Marshal(tser)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type StreamEndSumReport struct {
	Start         float32 `json:"start"`
	End           float32 `json:"end"`
	Seconds       float32 `json:"seconds"`
	Bytes         int     `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
}

func (sesr *StreamEndSumReport) String() string {
	b, err := json.Marshal(sesr)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type CpuUtilizationReport struct {
	HostTotal    float32 `json:"host_total"`
	HostUser     float32 `json:"host_user"`
	HostSystem   float32 `json:"host_system"`
	RemoteTotal  float32 `json:"remote_total"`
	RemoteUser   float32 `json:"remote_user"`
	RemoteSystem float32 `json:"remote_system"`
}

func (cur *CpuUtilizationReport) String() string {
	b, err := json.Marshal(cur)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type ConnectionInfo struct {
	Socket     int    `json:"socket"`
	LocalHost  string `json:"local_host"`
	LocalPort  int    `json:"local_port"`
	RemoteHost string `json:"remote_host"`
	RemotePort int    `json:"remote_port"`
}

func (ci *ConnectionInfo) String() string {
	b, err := json.Marshal(ci)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type TimestampInfo struct {
	Time     string `json:"time"`
	TimeSecs int    `json:"timesecs"`
}

func (tsi *TimestampInfo) String() string {
	b, err := json.Marshal(tsi)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type ConnectingToInfo struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (cti *ConnectingToInfo) String() string {
	b, err := json.Marshal(cti)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type TestStartInfo struct {
	Protocol   string `json:"protocol"`
	NumStreams int    `json:"num_streams"`
	BlkSize    int    `json:"blksize"`
	Omit       int    `json:"omit"`
	Duration   int    `json:"duration"`
	Bytes      int    `json:"bytes"`
	Blocks     int    `json:"blocks"`
	Reverse    int    `json:"reverse"`
}

func (tsi *TestStartInfo) String() string {
	b, err := json.Marshal(tsi)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type StartInfo struct {
	Connected     []*ConnectionInfo `json:"connected"`
	Version       string            `json:"version"`
	SystemInfo    string            `json:"system_info"`
	Timestamp     TimestampInfo     `json:"timestamp"`
	ConnectingTo  ConnectingToInfo  `json:"connecting_to"`
	Cookie        string            `json:"cookie"`
	TcpMssDefault int               `json:"tcp_mss_default"`
	TestStart     TestStartInfo     `json:"test_start"`
}

func (si *StartInfo) String() string {
	b, err := json.Marshal(si)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type EndInfo struct {
	Streams     []*StreamEndReport   `json:"streams"`
	SumSent     StreamEndSumReport   `json:"sum_sent"`
	SumReceived StreamEndSumReport   `json:"sum_received"`
	CpuReport   CpuUtilizationReport `json:"cpu_utilization_percent"`
}

func (ei *EndInfo) String() string {
	b, err := json.Marshal(ei)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type ServerReport struct {
	Start     StartInfo         `json:"start"`
	Intervals []*StreamInterval `json:"intervals"`
	End       EndInfo           `json:"end"`
	Error     string            `json:"error"`
}

func (sr *ServerReport) String() string {
	b, err := json.Marshal(sr)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

type TestReport struct {
	Start            StartInfo         `json:"start"`
	Intervals        []*StreamInterval `json:"intervals"`
	End              EndInfo           `json:"end"`
	Error            string            `json:"error"`
	ServerOutputJson ServerReport      `json:"server_output_json"`
}

func (tr *TestReport) String() string {
	b, err := json.Marshal(tr)
	if err != nil {
		return "error converting to json"
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, b, "", "    ")
	if err != nil {
		return "error converting json to indented format"
	}

	return string(pretty.Bytes())
}

func Loads(jsonStr string) (report *TestReport, err error) {
	r := &TestReport{}
	err = json.Unmarshal([]byte(jsonStr), r)
	return r, err
}

func Load(filename string) (report *TestReport, err error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return Loads(string(contents))
}
