
export BENCHMARK_WORKERS="8"
export URL="http://localhost:8080/product?uuid=6e4ff95f-f662-a5ee-e82a-bdf44a2d0b75"

curl: ##@benchmark show json
	curl -v --silent ${URL} | jq .

warmup: ##@benchmark warmup process
	ab -c 1 -n 1000 -k ${URL}

ab: ##@benchmark show throughput
	ab -c $(BENCHMARK_WORKERS) -n 200000 -k ${URL}

results.bin: $(shell find cmd pkg -type f) Makefile
	echo "GET ${URL}" | vegeta attack -workers=$(BENCHMARK_WORKERS) -max-workers=$(BENCHMARK_WORKERS) -rate=0 -duration=5s > results.bin

hdrplot: results.bin ##@benchmark show latency histogram
	vegeta report -type hdrplot results.bin

plot: results.bin ##@benchmark show latency ploy
	vegeta plot results.bin > plot.html
	google-chrome plot.html

report: results.bin ##@benchmark show vegeta report
	vegeta report -type text results.bin
