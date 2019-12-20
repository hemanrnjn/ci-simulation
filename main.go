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

type jobInfo struct {
	Id int
	Status bool
}

var (
	builds = 100000
	tests = 100
)

func main() {

	var wg sync.WaitGroup

	var buildsCh = make(chan jobInfo, builds)
	var testersCh = make(chan jobInfo, tests)

	for i := 0; i < builds; i++ {
		wg.Add(1)
		go buildRunner(&buildsCh, i, &wg)
	}

	for i := 0; i < tests; i++ {
		wg.Add(1)
		go testRunner(&testersCh, &wg)
	}

	for i := range buildsCh {
		if i.Status {
			testersCh <- jobInfo{Id:i.Id, Status:false}
		}
	}
	close(buildsCh)
	wg.Wait()
}

func buildRunner(buildChan *chan jobInfo, id int, wg *sync.WaitGroup) {
	fmt.Println("...Running Build..." + strconv.Itoa(id))
	time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	status := rand.Float64() * 1
	if status < 0.5 {
		log.Info("Job #" + strconv.Itoa(id) + ": BUILD FAILED!")
		*buildChan <- jobInfo{Id:id, Status:false}
	} else {
		log.Info("Job #" + strconv.Itoa(id) + ": BUILD PASSED!")
		*buildChan <- jobInfo{Id:id, Status:true}
	}
	wg.Done()
}

func testRunner(testerChan *chan jobInfo, wg *sync.WaitGroup) {
	for i := range *testerChan {
		fmt.Println("...Running Tests..." + strconv.Itoa(i.Id))
		time.Sleep(time.Second * time.Duration(rand.Intn(10)))
		status := rand.Float64() * 1
		if status < 0.5 {
			log.Info("Job #" + strconv.Itoa(i.Id) + ": BUILD FAILED!")
			*buildChan <- false
		} else {
			log.Info("Job #" + strconv.Itoa(buildInf.buildId) + ": BUILD PASSED!")
			*buildChan <- false
		}
	}
	wg.Done()
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