.PHONY: benchmarks

benchmarks:
	cd ./dgraph && go1.17.11 test ./... -bench=. -benchmem -benchtime=5s
	cd ./patrickmn && go1.17.11 test ./... -bench=. -benchmem -benchtime=5s
	cd ./damianopetrungaro && go1.18.3 test ./... -bench=. -benchmem -benchtime=5s