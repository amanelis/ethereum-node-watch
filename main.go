package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/gorilla/mux"

    l "github.com/sirupsen/logrus"
)

const (
    environment = "ENVIRONMENT"
    AWSRegion   = "us-west-2"
)

type Response struct {
    ID      int `json:"id"`
    JSONRPC string `json:"jsonrpc"`
    Result  string `json:"result"`
}

type Node struct {
    Instance struct {
        Name        string `json:"name"`

        AccountID   string `json:"account_id"`
        Environment string `json:"environment"`
        NodeID      string `json:"node_id"`
        NodeNetwork string `json:"node_network"`
        NodeRanking string `json:"node_ranking"`
        NodeService string `json:"node_service"`
        NodeType    string `json:"node_type"`
        NodeVersion string `json:"node_version"`

        InstanceID  string `json:"instance_id"`
        PublicIP    string `json:"public_ip"`
        PrivateIP   string `json:"private_ip"`
        PublicDNS   string `json:"public_dns"`
        PrivateDNS  string `json:"private_dns"`
    } `json:"instance"`

    RPC struct {
        Blocks string `json:"blocks"`
        Peers string `json:"peers"`
        Version string `json:"version"`
    } `json:"rpc"`
}

func main() {
    var router = mux.NewRouter()
    router.HandleFunc("/", index).Methods("GET")
    router.HandleFunc("/nodes", nodes).Methods("GET")
    router.HandleFunc("/healthcheck", healthCheck).Methods("GET")

    l.SetOutput(os.Stdout)

    fmt.Println("55Nodes server started at http://localhost")
    log.Fatal(http.ListenAndServe(":80", router))
}

func index(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{"message": "OK"})
}

func nodes(w http.ResponseWriter, r *http.Request) {
	svc := ec2.New(session.Must(session.NewSession()), &aws.Config{Region: aws.String(AWSRegion)})

    var nodes []Node

    for _, reservations := range fetchInstances(svc, filters()).Reservations {
		for _, instance := range reservations.Instances {
            l.Info(fmt.Sprintf("Instance: %s, %s", *instance.InstanceId, *instance.State.Name))

            node := Node{}
            for _, tag := range instance.Tags {
                if *tag.Key == "Name" {
                    node.Instance.Name = *tag.Value
                }

                if *tag.Key == "AccountID" {
                    node.Instance.AccountID = *tag.Value
                }

                if *tag.Key == "Environment" {
                    node.Instance.Environment = *tag.Value
                }

                if *tag.Key == "NodeID" {
                    node.Instance.NodeID = *tag.Value
                }

                if *tag.Key == "NodeNetwork" {
                    node.Instance.NodeNetwork = *tag.Value
                }

                if *tag.Key == "NodeRanking" {
                    node.Instance.NodeRanking = *tag.Value
                }

                if *tag.Key == "NodeService" {
                    node.Instance.NodeService = *tag.Value
                }

                if *tag.Key == "NodeType" {
                    node.Instance.NodeType = *tag.Value
                }

                if *tag.Key == "NodeVersion" {
                    node.Instance.NodeVersion = *tag.Value
                }
            }

            node.Instance.InstanceID = *instance.InstanceId
            node.Instance.PublicIP = *instance.PublicIpAddress
            node.Instance.PublicDNS = *instance.PublicDnsName
            node.Instance.PrivateIP = *instance.PrivateIpAddress
            node.Instance.PrivateDNS = *instance.PrivateDnsName

            rpc := fmt.Sprintf("http://%s:8545", node.Instance.PrivateIP)

            blk := requestRPCMethod(rpc, "eth_blockNumber")
            b, _ := strconv.ParseUint(blk.Result[2:], 16, 32)
            blocks := fmt.Sprint(uint32(b))

            prn := requestRPCMethod(rpc, "net_peerCount")
            p, _ := strconv.ParseUint(prn.Result[2:], 16, 32)
            peers := fmt.Sprint(uint32(p))

            typ := requestRPCMethod(rpc, "web3_clientVersion")
            ntype := typ.Result

            l.WithFields(l.Fields{
                "method": "#nodes",
            }).Info(fmt.Sprintf("Blocks: %s, Peers: %s, Type: %s", blocks, peers, ntype))

            node.RPC.Blocks = blocks
            node.RPC.Peers = peers
            node.RPC.Version = ntype

            nodes = append(nodes, node)
		}
	}

    json.NewEncoder(w).Encode(nodes)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode("OK")
}

func fetchInstances(svc *ec2.EC2, filters []*ec2.Filter) *ec2.DescribeInstancesOutput {
	resp, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{Filters: filters})
	if err != nil {
		fmt.Println("there was an error listing instances in ", AWSRegion, err.Error())
		log.Fatal(err.Error())
	}

	return resp
}

func filters() []*ec2.Filter {
	return []*ec2.Filter{
		{
			Name:   aws.String("tag:NodeService"),
			Values: []*string{aws.String("geth")},
		},
        {
            Name:   aws.String("tag:NodeRanking"),
            Values: []*string{aws.String("leader")},
        },
	}
}

func requestRPCMethod(e string, m string) *Response {
    buf := []byte(fmt.Sprintf(`{"jsonrpc": "2.0", "method": "%s", "params": [], "id": 1}`, m))
    req, _ := http.NewRequest("POST", e, bytes.NewBuffer(buf))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    res := Response{}
    per := json.NewDecoder(resp.Body).Decode(&res)
    if per != nil {
        panic(per)
    }

    return &res
}
