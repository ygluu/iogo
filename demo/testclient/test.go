package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"runtime"
	"time"

	"iogo"
	appl "iogo/demo/proto/applogic"
	sess "iogo/demo/proto/session"
)

func init() {

}

type client struct {
	cookie   string
	okCount  int
	errCount int
	showLog  bool
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetID() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

func (self *client) login(id int) {
	request := &sess.LoginRequest{}
	request.Name = GetID() + fmt.Sprintf(":%d", id)
	if self.showLog {
		iogo.Log.I(0, "Login Request: %+v", request)
	}
	ctx := sess.Background()
	reply, ret := sess.CallSession.Login(ctx, request)
	if ret != nil {
		self.errCount++
		iogo.Log.E(0, "Login error: "+ret.Error())
	} else {
		self.okCount++
		self.cookie = reply.Token
		if self.showLog {
			iogo.Log.I(0, "Login Reply: %+v", reply)
		}
	}
}

func (self *client) AppLogic() {
	request := &appl.SayRequest{}
	request.Text = "Hellow world!"
	if self.showLog {
		iogo.Log.I(0, "AppLogic Request: %+v", request)
	}
	reply, ret := appl.CallApplogic.Say(appl.TokenContext(self.cookie), request)
	if ret != nil {
		self.errCount++
		iogo.Log.E(0, "AppLogic error: "+ret.Error())
	} else {
		self.okCount++
		if self.showLog {
			iogo.Log.I(0, "AppLogic Reply: %+v", reply)
		}
	}
}

func TestThread(threadId int, showLog bool, callCount int, chanEnd chan *client) {
	c := &client{
		cookie:   "",
		okCount:  0,
		errCount: 0,
		showLog:  showLog,
	}
	for i := 0; i < callCount; i++ {
		c.login(threadId*100000 + i + 1)
		c.AppLogic()
	}
	chanEnd <- c
}

func testStart() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	time.Sleep(1 * time.Second)

	const concurrent = 1000
	const callCount = 1

	const testCount = 10
	const sleepTime = 2

	var endChan chan *client = make(chan *client, concurrent*2)

	for i := 0; i < testCount; i++ {
		endCount := 0
		okCount := 0
		errCount := 0
		var end time.Time
		start := time.Now()
		for j := 0; j < concurrent; j++ {
			go TestThread(j+1, concurrent == 1, callCount, endChan)
		}
		for {
			client := <-endChan
			errCount += client.errCount
			okCount += client.okCount
			endCount++
			if endCount == concurrent {
				end = time.Now()
				break
			}
		}
		sumCall := errCount + okCount
		Long := end.Sub(start)
		if Long == 0 {
			Long = 1
		}
		concur := time.Duration(1000*1000*1000) / (Long / time.Duration(sumCall))
		PerLong := Long / time.Duration(sumCall) / time.Duration(1000)
		Long = Long / 1000
		iogo.Log.I(0, "TectCount:%d Call:%d OK:%d Error:%d Long(us):%d PerLong(us):%d, concurrent(s):%d",
			i+1, sumCall, okCount, errCount, Long, PerLong, concur)
		time.Sleep(sleepTime * time.Second)
	}
}
