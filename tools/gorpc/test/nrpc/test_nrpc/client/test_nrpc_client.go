package main

import (
    "context"
	"flag"
	"fmt"
    "git.code.oa.com/go-neat/core/config"
    "git.code.oa.com/go-neat/core/common"
    "git.code.oa.com/go-neat/core/depmod/trace"
    "git.code.oa.com/go-neat/core/nlog"
    "git.code.oa.com/nrpc_protos/test_nrpc"
    prototype "git.code.oa.com/go-neat/core/proto/nrpc"
    cltpkg "git.code.oa.com/go-neat/core/proto/nrpc/nrpc_pb"
	"os"
	"path/filepath"
	"time"
	"sync"
	
)

// common options
var mode = flag.Int("mode", 1, "mode, 1:manual, 2:intelligent")
var transport = flag.Int("transport", 2, "transport type, 1:udp,2:tcp_short,3:tcp_keepalive,4:tcp_full_duplex,5:udp_full_duplex,6:udp_without_recv")
var addr = flag.String("addr", "ip://127.0.0.1:8000", "addr, supporting ip://<ip>:<port>, l5://mid:cid, cmlb://appid[:sysid]")
var cmd = flag.String("cmd", "BuyApple", "cmd name, BuyApple, SellApple")

var timeout = flag.Int("timeout", 2000, "timeout in miliseconds")
var delay = flag.Int("delay", 100, "delay microseconds before next rpccall issued")

// manual mode
var total = flag.Int("total", 1, "total")

// intelligent mode
var failure = flag.Float64("failure", 0.001, "failure ratio")

func init() {
	flag.Parse()

	parentPath, err := common.GetParentDirectory()
	if err != nil {
		fmt.Println(err)
	}
	confDir := filepath.Join(parentPath, "conf")
	_, err = config.NewIniConfig(filepath.Join(confDir, "service.ini"), true)
	if err != nil {
		fmt.Println(err)
	}
}

func initReqSession() * prototype.NRPCSession{
	headRequest := &cltpkg.NRPCPkg{Head:&cltpkg.Head{}}
	session := &prototype.NRPCSession{NRPCReq: headRequest}
	return session
}

const (
	TESTMODE_MANUAL = 1
	TESTMODE_INTELLIGENT = 2
	TESTMODE_WORKERPOOL = 3
)

var callback = []func(){}

func register(f func()) {
    callback = append(callback, f)
}

func cleanup() {
    for _, f := range callback {
        f()
    }
}

func main() {
    if trace.TraceEnabled() {
        register(func() {
            trace.TraceStopGrace()
        })
        defer cleanup()
    }

	// initialize rpc client
	log := nlog.GetLogger("default")
	test_nrpc.Init(*addr, *transport, log)
	defer log.Flush()

	// initialize request session
	session := initReqSession()

	if *mode == TESTMODE_MANUAL {
		testManualMode(session)
	} else if *mode == TESTMODE_INTELLIGENT {
		testIntelligentMode(session)
	} else {
		fmt.Fprintf(os.Stderr, "Invalid testmode: %v, supporting testmode 1:manual, 2:intelligent", *mode)
		os.Exit(1)
	}
}

func testManualMode(session *prototype.NRPCSession) {

	ch := make(chan struct{}, *total)
	go func() {
		for {
			time.Sleep(time.Microsecond * time.Duration(*delay))
			ch <- struct{}{}
		}
	}()

	wg := &sync.WaitGroup{}
	ts_begin := time.Now()

	var succ int = 0

	switch *cmd {
		case "BuyApple": { // test rpc: BuyApple
			count := 0
			for count < *total {
				<- ch
				count++
				wg.Add(1)
				go func() {
					defer wg.Done()
					req := &test_nrpc.BuyAppleReq{}
					//rsp, err := test_nrpc.BuyApple(context.Background(), session, req)
					_, err := test_nrpc.BuyApple(context.Background(), session, req)

					if err != nil {
						//fmt.Printf("req: %v, rsp: %v, err: %v", req, rsp, err)
					} else {
						//fmt.Printf("req: %v, rsp: %v\n", req, rsp)
						succ++
					}
				}()
			}
		}
		case "SellApple": { // test rpc: SellApple
			count := 0
			for count < *total {
				<- ch
				count++
				wg.Add(1)
				go func() {
					defer wg.Done()
					req := &test_nrpc.SellAppleReq{}
					//rsp, err := test_nrpc.SellApple(context.Background(), session, req)
					_, err := test_nrpc.SellApple(context.Background(), session, req)

					if err != nil {
						//fmt.Printf("req: %v, rsp: %v, err: %v", req, rsp, err)
					} else {
						//fmt.Printf("req: %v, rsp: %v\n", req, rsp)
						succ++
					}
				}()
			}
		}
	}

	wg.Wait()
	cost := time.Since(ts_begin)

	fmt.Println()
	fmt.Printf("[Summary] reqs: %v, succ: %v, timeout: %v, timecost: %v seconds\n", *total, succ, (*total)-succ, cost.Seconds())
}

func testIntelligentMode(session *prototype.NRPCSession) {

	lock := &sync.Mutex{}
	reqs := 0
	succ := 0

	// executing the testing
	step := 1000
	batch := 1
	loop := 1

	switch *cmd {
		case "BuyApple":
			req := &test_nrpc.BuyAppleReq{}
			for {
				wg := &sync.WaitGroup{}
				wg.Add(batch)

				for i:=0; i<batch; i++ {
					go func() {
						defer wg.Done()

						//rsp, err := test_nrpc.BuyApple(context.Background(), session, req)
						_, err := test_nrpc.BuyApple(context.Background(), session, req)

						lock.Lock()
						reqs++
						lock.Unlock()

						if err != nil {
							//fmt.Printf("req: %v, rsp: %v, err: %v", req, rsp, err)
						} else {
							lock.Lock()
							succ++
							lock.Unlock()
						}
					}()
					time.Sleep(time.Microsecond * time.Duration(*delay))
				}
				wg.Wait()

				failure_ratio := float64(reqs-succ)/float64(reqs)
				fmt.Printf("LoopTimes -> %d, reqs -> %d, succ -> %d, timeout -> %d, failure -> %f\n", loop, reqs, succ, reqs-succ, failure_ratio)

				if failure_ratio > *failure {
					fmt.Println("Has reached max allowed failure ratio ... exit")
					os.Exit(1)
				}
				reqs, succ = 0, 0
				batch += step
				loop++
			}
		case "SellApple":
			req := &test_nrpc.SellAppleReq{}
			for {
				wg := &sync.WaitGroup{}
				wg.Add(batch)

				for i:=0; i<batch; i++ {
					go func() {
						defer wg.Done()

						//rsp, err := test_nrpc.SellApple(context.Background(), session, req)
						_, err := test_nrpc.SellApple(context.Background(), session, req)

						lock.Lock()
						reqs++
						lock.Unlock()

						if err != nil {
							//fmt.Printf("req: %v, rsp: %v, err: %v", req, rsp, err)
						} else {
							lock.Lock()
							succ++
							lock.Unlock()
						}
					}()
					time.Sleep(time.Microsecond * time.Duration(*delay))
				}
				wg.Wait()

				failure_ratio := float64(reqs-succ)/float64(reqs)
				fmt.Printf("LoopTimes -> %d, reqs -> %d, succ -> %d, timeout -> %d, failure -> %f\n", loop, reqs, succ, reqs-succ, failure_ratio)

				if failure_ratio > *failure {
					fmt.Println("Has reached max allowed failure ratio ... exit")
					os.Exit(1)
				}
				reqs, succ = 0, 0
				batch += step
				loop++
			}
	}
}
