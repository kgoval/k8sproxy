package main

import (
	"flag"
	"fmt"
	"github.com/micro/go-micro/server"
	"sync"

	k8s "github.com/micro/go-micro"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("running..")
	prefix := flag.String("prefix", "msa.", "prefix microservice services on kubernetes")
	flag.Parse()
	//boot.ReadEnv("MV")

	fmt.Println("prefix:" + *prefix)
	cmd := exec.Command("kubectl", "get", "services", "-o", "custom-columns=:metadata.name")
	out, err := cmd.CombinedOutput()

	outs := strings.Split(string(out), "\n")
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	var wg sync.WaitGroup
	var newOuts []string

	for _, s := range outs {

		if s != "" && strings.Contains(s, *prefix) {
			newOuts = append(newOuts, s)
		}

	}

	wg.Add(len(newOuts))
	for i, s := range newOuts {

		go register(i, s, wg)

	}
	wg.Wait()
}

func register(i int, s string, wg sync.WaitGroup) {
	defer wg.Done()

	service := k8s.NewService(k8s.Name(s))
	service.Server().Init(server.Address("127.0.0.1:808" + strconv.Itoa(i)))
	service.Server().Register()
	fmt.Println("Registering Service '" + s + "' @ 127.0.0.1:808" + strconv.Itoa(i))
	// expose
	cmd := exec.Command("kubectl", "port-forward", "service/"+s, "808"+strconv.Itoa(i)+":9080")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	fmt.Println(out)

}
