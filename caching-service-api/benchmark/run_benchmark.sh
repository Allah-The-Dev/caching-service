#test results calculated in 32GB RAM and 1TB Disk machine

#bemchmark GetEmployees
go test -benchmem -run=^$ caching-service/handlers -bench ^BenchmarkGetEmployees$
#result
#3062	    404075 ns/op	   17847 B/op	     347 allocs/op

#benchmark GetEmployee
go test -benchmem -run=^$ caching-service/handlers -bench ^BenchmarkGetEmployee$
#result
2067	    724717 ns/op	   16525 B/op	     117 allocs/op

#bencmark PostEmployee
go test -benchmem -run=^$ caching-service/handlers -bench ^BenchmarkPostEmployee$
#result
#181624	      7184 ns/op	    2165 B/op	      16 allocs/op

#run all benchmarks
go test -benchmen -run=^$ caching-service/handlers -bench=.