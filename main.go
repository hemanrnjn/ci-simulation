package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type jobInfo struct {
	ID     int
	Status bool
}

type Counter struct {
	C int
	Done bool
}

var (
	builders = 10000000
	tasks = 1000000
	testers  = 10000000
)

func main() {

	var wg sync.WaitGroup

	var buildsCh = make(chan jobInfo)
	var testersCh = make(chan jobInfo)
	var deployCh = make(chan jobInfo)
	var buildDone, testDone bool
	var buildCounter, buildPassCounter, testCounter, testPassCounter, deployCounter int
	var m sync.Mutex

	for i := 0; i < builders; i++ {
		go buildRunner(&buildsCh, &testersCh, &buildDone, &buildCounter, &buildPassCounter, &m)
	}

	for i := 0; i < testers; i++ {
		go testRunner(&testersCh, &deployCh, &testDone, &testCounter, &testPassCounter, &m)
	}

	wg.Add(1)
	go deployRunner(&deployCh, &wg, &deployCounter, &m)

	for i := 0; i < tasks; i++ {
		buildsCh <- jobInfo{ID: i, Status: false}
	}

	go func() {
		wg.Add(1)
		for {
			if buildCounter == tasks {
				close(buildsCh)
				break
			} else {
				time.Sleep(time.Millisecond * 1)
			}
		}
		wg.Done()
	}()


	go func() {
		wg.Add(1)
		for {
			if buildDone {
				if buildPassCounter == testCounter {
					close(testersCh)
					break
				} else {
					time.Sleep(time.Millisecond * 1)
				}
			} else {
				time.Sleep(time.Millisecond * 1)
			}
		}
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		for {
			if testDone {
				if testPassCounter == deployCounter {
					close(deployCh)
					break
				} else {
					time.Sleep(time.Millisecond * 1)
				}
			} else {
				time.Sleep(time.Millisecond * 1)
			}
		}
		wg.Done()
	}()

	wg.Wait()
	
}

func buildRunner(buildChan, testerCh *chan jobInfo, buildDone *bool, buildCounter, buildPassCounter *int, m *sync.Mutex) {
	for i := range *buildChan {
		fmt.Println("...Running Build for job #" + strconv.Itoa(i.ID) + " ...")
		time.Sleep(time.Second * 1)
		status := rand.Float64() * 1
		if status < 0.5 {
			log.Info("Job #" + strconv.Itoa(i.ID) + ": BUILD FAILED!")
		} else {
			log.Info("Job #" + strconv.Itoa(i.ID) + ": BUILD PASSED!")
			*testerCh <- jobInfo{ID: i.ID, Status: true}
			m.Lock()
			*buildPassCounter++
			m.Unlock()
		}
		m.Lock()
		*buildCounter++
		m.Unlock()
	}
	*buildDone = true
}

func testRunner(testerChan, deployCh *chan jobInfo, testDone *bool, testCounter, testPassCounter *int, m *sync.Mutex) {
	for i := range *testerChan {
		fmt.Println("...Running Tests for job #" + strconv.Itoa(i.ID) + " ...")
		time.Sleep(time.Second * 1)
		status := rand.Float64() * 1
		if status < 0.5 {
			log.Info("Job #" + strconv.Itoa(i.ID) + ": TESTS FAILED!")
		} else {
			log.Info("Job #" + strconv.Itoa(i.ID) + ": TESTS PASSED!")
			*deployCh <- jobInfo{ID: i.ID, Status: true}
			m.Lock()
			*testPassCounter++
			m.Unlock()
		}
		m.Lock()
		*testCounter++
		m.Unlock()
	}
	*testDone = true
}

func deployRunner(deployChan *chan jobInfo, wg *sync.WaitGroup, deployCounter *int, m *sync.Mutex) {
	for i := range *deployChan {
		fmt.Println("...Running Deploy for job #..." + strconv.Itoa(i.ID))
		time.Sleep(time.Millisecond * 1)
		status := rand.Float64() * 1
		if status < 0.5 {
			log.Info("Job #" + strconv.Itoa(i.ID) + ": DEPLOY FAILED!")
		} else {
			log.Info("Job #" + strconv.Itoa(i.ID) + ": DEPLOY PASSED!")
		}
		m.Lock()
		*deployCounter++
		m.Unlock()
	}
	wg.Done()
}
