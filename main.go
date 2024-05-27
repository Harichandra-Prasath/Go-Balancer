package main

func main() {
	pool := GetPool()
	b1 := GetBackend("http://localhost:4000/")
	b2 := GetBackend("http://localhost:5000/")
	pool.Addserver(b1)
	pool.Addserver(b2)
}
