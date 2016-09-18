package service

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"webcron-agent/jobs"
	"webcron-agent/libs"
)

const (
	Address = "127.0.0.1"
	Port    = 7249
)

const (
	TASK_ADD    = "task_add"
	TASK_STOP   = "task_stop"
	TASK_START  = "task_start"
	TASK_RUN    = "task_run"
	TASK_DELETE = "task_delete"
)

type Service interface {
	Start() error
	Stop()
}

type AgentService struct {
	mtx       *sync.RWMutex
	exit      chan bool
	waitGroup *sync.WaitGroup
}

type Data struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
	Time string          `json:"time"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewAgentService() Service {
	return &AgentService{
		exit:      make(chan bool),
		waitGroup: &sync.WaitGroup{},
	}
}

func (agent *AgentService) Start() error {
	listen, err := net.Listen("tcp4", Address+":"+strconv.Itoa(Port))
	if err != nil {
		return err
	}
	defer listen.Close()

	log.Printf("Agent service start listening on %s:%d", Address, Port)

	go func() {
		exit := <-agent.exit
		if exit {
			close(agent.exit)
			log.Printf("Agent service stoped")
			listen.Close()
		}
	}()

	agent.SignalHandler()

	go jobs.RunLocalTask()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln("Lost connection error:", err.Error())
		}

		agent.waitGroup.Add(1)
		go agent.Handler(conn)
	}

	return nil
}

func (agent *AgentService) Stop() {
	agent.waitGroup.Wait()
	agent.exit <- true
}

func (agent *AgentService) Handler(conn net.Conn) error {
	var socketLogger = libs.NewSocketLogger()
	defer conn.Close()
	defer agent.waitGroup.Done()

	for {
		data, err := agent.ReadData(conn)

		switch err {
		case io.EOF:
			return nil
		case nil:
			socketLogger.Info("收到消息体: " + string(data))
			//time.Sleep(time.Second * 30)
			d, err := decodeData(data)
			if err != nil {
				return err
			}

			agent.Dispatch(conn, d)
			//if err != nil {
			//	agent.SendData(conn, []byte("error"))
			//	log.Printf("RunTask error:%s", err.Error())
			//} else {
			//	agent.SendData(conn, []byte("ok"))
			//}
			continue
		default:
			//log.Printf("Receive data failed:%s", err)
			return nil
		}
	}
	return nil
}

func (agent *AgentService) SignalHandler() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		switch sig {

		case syscall.SIGINT, syscall.SIGTERM:
			log.Printf("Agent Service catch signal: %s, waiting all runing task done,try to stop service", sig)
			agent.Stop()

		default:
			log.Printf("Unknow signal")
		}
	}()
}

func (agent *AgentService) ReadData(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 1024)
	r := bufio.NewReader(conn)
	n, err := r.Read(buf)

	data, _ := libs.AesDencrypt(buf[:n])
	return data, err
}

func (agent *AgentService) SendData(conn net.Conn, str []byte) {
	w := bufio.NewWriter(conn)
	send, _ := libs.AesEncrypt(str)

	w.Write(send)
	w.Flush()

	var log = libs.NewSocketLogger()
	log.Info("发送消息体: " + string(str))
}

func decodeData(d []byte) (*Data, error) {
	var data *Data
	err := json.Unmarshal(d, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (agent *AgentService) Dispatch(conn net.Conn, d *Data) error {
	if nil == d {
		return errors.New("数据错误")
	}

	var err error
	var response Response

	switch d.Type {
	case TASK_START:
		err = jobs.CronTask(d.Data, true)
	case TASK_STOP:
		err = jobs.StopTasks(d.Data)
	case TASK_RUN:
		result, run_err := jobs.RunTask(d.Data)
		if run_err != nil {
			err = run_err
		}
		response.Data = result
	}

	if err != nil {
		response.Code = 1
		response.Msg = err.Error()
		rsp, _ := json.Marshal(response)
		agent.SendData(conn, rsp)
		return err
	}

	response.Code = 0
	response.Msg = "OK"
	rsp, _ := json.Marshal(response)
	agent.SendData(conn, rsp)
	return nil
}
