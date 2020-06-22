package litenetlib

type NetStatistics struct {
	PacketsSent     int64
	PacketsReceived int64
	BytesSent       int64
	BytesReceived   int64
	PacketLoss      int64
}

func (netStatistics *NetStatistics) PacketLossPercent() int8 {
	if netStatistics.PacketsSent == 0 {
		return 0
	}

	return int8((netStatistics.PacketLoss * 100) / netStatistics.PacketsSent)
}

func (netStatistics *NetStatistics) Reset() {
	netStatistics.PacketsReceived = 0
	netStatistics.PacketsReceived = 0
	netStatistics.BytesReceived = 0
	netStatistics.BytesReceived = 0
	netStatistics.PacketLoss = 0
}

func NewNetStatistics() *NetStatistics {
	return &NetStatistics{}
}
