package Utils

//是否为ipv4地址

func IsIpv4Addr(ipv4 string)bool  {
	for _,eIp := range ipv4{
		if eIp >= '0' && eIp <= '9'{
			continue
		}
		if eIp != '.'{
			return false
		}
	}
	return true
}
