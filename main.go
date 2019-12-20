package ci_simulate

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type buildInfo struct {
	buildId int
	status bool
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	builds, _ := strconv.Atoi(scanner.Text())

	var wg sync.WaitGroup

	var ch = make(chan bool, builds)

	for i := 0; i < builds; i++ {
		wg.Add(1)
		newBuild := buildInfo{buildId: i, status: false}
		go builder(&ch, newBuild)
	}

	for i := range ch {

	}
	close(ch)
}

func builder(buildChan *chan bool, buildInf buildInfo) {
	fmt.Println("...Running Build..." + strconv.Itoa(buildInf.buildId))
	time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	status := rand.Float64() * 1
	if status < 0.5 {
		log.Info("Job #" + strconv.Itoa(buildInf.buildId) + ": BUILD FAILED!")
		*buildChan <- false
	} else {
		log.Info("Job #" + strconv.Itoa(buildInf.buildId) + ": BUILD PASSED!")
		*buildChan <- true
		go test()
	}
}

func test() {
	fmt.Println("...Running Build..." + strconv.Itoa(buildInf.buildId))
	time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	status := rand.Float64() * 1
	if status < 0.5 {
		log.Info("Job #" + strconv.Itoa(buildInf.buildId) + ": BUILD FAILED!")
		*buildChan <- false
	} else {
		log.Info("Job #" + strconv.Itoa(buildInf.buildId) + ": BUILD PASSED!")
		*buildChan <- false
	}
}

func deploy() {
	fmt.Println("...Running Build..." + strconv.Itoa(buildInf.buildId))
	time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	status := rand.Float64() * 1
	if status < 0.5 {
		log.Info("Job #" + strconv.Itoa(buildInf.buildId) + ": BUILD FAILED!")
		*buildChan <- false
	} else {
		log.Info("Job #" + strconv.Itoa(buildInf.buildId) + ": BUILD PASSED!")
		*buildChan <- false
	}
}