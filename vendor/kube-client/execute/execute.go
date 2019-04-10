package execute

import (
	"k8s.io/client-go/kubernetes"
	"kube-client/common"
	"sync"
)

var (
	once   sync.Once
	client *Client
)

type Client struct {
	clientset   *kubernetes.Clientset
	typeFunMaps map[string]func(KubeTransfer, chan KubeTransfer) error
}

type KubeTransfer struct {
	Types          byte
	Method, Result string
	HandleJson     []byte
}

func GetClient() *Client {
	once.Do(func() {
		client = &Client{}
		var err error
		if client.clientset, err = common.InitClient(); err != nil {
			panic(err)
		}
		client.initTypeFunMaps()
	})
	return client
}

func (client *Client) initTypeFunMaps() {
	client.typeFunMaps = make(map[string]func(KubeTransfer, chan KubeTransfer) error)
	client.typeFunMaps["development"] = deployment
	client.typeFunMaps["namespace"] = namespace
}

func (client *Client) Execute(transfer KubeTransfer, outChan chan KubeTransfer) {
	if err := client.typeFunMaps[transfer.Method](transfer, outChan); err != nil {
		transfer.Types = 2
		transfer.Result = err.Error()
		transfer.HandleJson = nil
		outChan <- transfer
	}
}
