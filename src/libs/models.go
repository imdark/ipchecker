package libs

import (
	"time"
)
type TCPReport struct {
	TotalIps int
	TotalIpsNotReachable int
	IpRangesNotReachable int /*ip ranges are under /24 subnet mask */
	IpRangesPartiallyReachable int /*More then 50 ips in range are not*/
	TotalRunningTime time.Duration /* Including report generation */
}
