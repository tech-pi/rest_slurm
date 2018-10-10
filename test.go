package main

import (
	"fmt"
	"github.com/json-iterator/go"
	"regexp"
	"strings"
)

var sq = `             JOBID PARTITION     NAME     USER ST       TIME  NODES NODELIST(REASON)
              5424      main  post.sh     root PD       0:00      1 (Dependency)
              5324      main   run.sh     root  R       0:24      1 GS0
              5325      main   run.sh     root  R       0:24      1 GS0
              5326      main   run.sh     root  R       0:24      1 GS0
              5327      main   run.sh     root  R       0:24      1 GS0
              5328      main   run.sh     root  R       0:24      1 GS0
              5329      main   run.sh     root  R       0:24      1 GS0
              5330      main   run.sh     root  R       0:24      1 GS0
              5331      main   run.sh     root  R       0:24      1 GS0
              5332      main   run.sh     root  R       0:24      1 GS0
              5333      main   run.sh     root  R       0:24      1 GS0
              5334      main   run.sh     root  R       0:24      1 GS0
              5335      main   run.sh     root  R       0:24      1 GS0
              5336      main   run.sh     root  R       0:21      1 GS0
              5337      main   run.sh     root  R       0:21      1 GS0
              5338      main   run.sh     root  R       0:21      1 GS0
              5339      main   run.sh     root  R       0:21      1 GS0
              5340      main   run.sh     root  R       0:21      1 GS0
              5341      main   run.sh     root  R       0:21      1 GS0`

var batch = "Submitted batch job sdf"
var control = `JobId=5330 JobName=run.sh
   UserId=root(0) GroupId=root(0)
   Priority=4294900898 Nice=0 Account=(null) QOS=(null)
   JobState=CANCELLED Reason=None Dependency=(null)
   Requeue=1 Restarts=0 BatchFlag=1 Reboot=0 ExitCode=0:15
   RunTime=00:02:07 TimeLimit=365-00:00:00 TimeMin=N/A
   SubmitTime=2018-07-18T02:28:18 EligibleTime=2018-07-18T02:28:18
   StartTime=2018-07-18T02:28:18 EndTime=2018-07-18T02:30:25
   PreemptTime=None SuspendTime=None SecsPreSuspend=0
   Partition=main AllocNode:Sid=CS1:28693
   ReqNodeList=(null) ExcNodeList=(null)
   NodeList=GS0
   BatchHost=GS0
   NumNodes=1 NumCPUs=1 CPUs/Task=1 ReqB:S:C:T=0:0:*:*
   TRES=cpu=1,node=1
   Socks/Node=* NtasksPerN:B:S:C=0:0:*:* CoreSpec=*
   MinCPUsNode=1 MinMemoryNode=0 MinTmpDiskNode=0
   Features=(null) Gres=(null) Reservation=(null)
   Shared=OK Contiguous=0 Licenses=(null) Network=(null)
   Command=/mnt/gluster/twj/GATE/16_8.0/run.sh
   WorkDir=/mnt/gluster/twj/GATE/16_8.0/sub.81
   StdErr=/mnt/gluster/twj/GATE/16_8.0/sub.81/%J.err
   StdIn=/dev/null
   StdOut=/mnt/gluster/twj/GATE/16_8.0/sub.81/%J.out
   Power= SICP=0`

func main() {
	lines := strings.Split(sq, "\n")
	for i := range lines {
		lines[i] = standardizeSpacesx(lines[i])
	}
	res := make([]Squeuex, len(lines)-1)
	for i := range res {
		x := strings.Split(lines[i+1], " ")
		res[i] = Squeuex{x[0], x[1], x[2], x[3], x[4], x[5], x[6], x[7]}
	}

	resJson, err := jsoniter.Marshal(res)
	fmt.Println(string(resJson), err)

	rgx := regexp.MustCompile("([0-9]+)")
	index := rgx.FindString(batch)
	fmt.Println(index)

	rgx = regexp.MustCompile("JobState=([A-Z]+)")
	s := rgx.FindSubmatch([]byte(control))
	fmt.Println(string(s[1]))
}

type Squeuex struct {
	JobId     string `json:"job_id"`
	Partition string `json:"partition"`
	Name      string `json:"name"`
	User      string `json:"user"`
	Status    string `json:"status"`
	Time      string `json:"time"`
	Nodes     string `json:"nodes"`
	NodeList  string `json:"node_list"`
}

func standardizeSpacesx(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
