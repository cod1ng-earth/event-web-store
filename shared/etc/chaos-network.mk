
break: ##@network Slow down local network
	sudo tc qdisc add dev lo root netem delay 14ms

fix: ##@network Stop slowing down local network
	sudo tc qdisc delete dev lo root netem delay 14ms
