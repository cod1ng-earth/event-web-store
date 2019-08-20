
curl: ##@benchmark show json
	curl -v --silent ${URL} | jq .

warmup: ##@benchmark warmup process
	ab -c 1 -n 1000 -k ${URL}

ab: ##@benchmark show throughput
	ab -c 20 -n 200000 -k ${URL}

results.bin: $(shell find cmd pkg -type f) Makefile
	echo "GET ${URL}" | vegeta attack -workers=20 -max-workers=20 -rate=0 -duration=5s > results.bin

hdrplot: results.bin ##@benchmark show latency histogram
	vegeta report -type hdrplot results.bin

plot: results.bin ##@benchmark show latency ploy
	vegeta plot results.bin > plot.html
	google-chrome plot.html

report: results.bin ##@benchmark show vegeta report
	vegeta report -type text results.bin
