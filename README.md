### Testing
1. Create entries.txt file by executing `text-generator.sh` (exercise: modify the count, include some random duplicates (HLL counts unique items from multiset)
2. run the go program `go run hll-benchmark.go`


HLL generally has error rate of less than 2%.

In some of my local runs I could see following results:

```bash
go run hll-benchmark.go
Inserted Count: 1000000
HLL Count: 996819
Error : 0.318100
```
