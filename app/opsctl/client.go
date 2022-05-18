package opsctl

import (
	"context"
	"fmt"
	"github.com/pipperman/kubeops/app/kubeops/biz/types"
	"io"

	"github.com/google/wire"
	api "github.com/pipperman/kubeops/api/v1"
	"google.golang.org/grpc"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewKubeOpsClient)

func NewKubeOpsClient(host string, port int) *kubeOpsClient {
	return &kubeOpsClient{
		host: host,
		port: port,
	}
}

type kubeOpsClient struct {
	host string
	port int
}

func (c *kubeOpsClient) CreateProject(name string, source string) (*api.Project, error) {
	conn, err := c.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewKubeOpsApiClient(conn)
	request := api.CreateProjectRequest{
		Name:   name,
		Source: source,
	}
	resp, err := client.CreateProject(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return resp.Item, nil

}

func (c kubeOpsClient) ListProject() ([]*api.Project, error) {
	conn, err := c.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewKubeOpsApiClient(conn)
	request := api.ListProjectRequest{}
	resp, err := client.ListProject(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (c kubeOpsClient) RunPlaybook(project, playbook, tag string, inventory *api.Inventory) (*api.Result, error) {
	conn, err := c.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewKubeOpsApiClient(conn)
	request := &api.RunPlaybookRequest{
		Project:   project,
		Playbook:  playbook,
		Inventory: inventory,
		Tag:       tag,
	}
	req, err := client.RunPlaybook(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return req.Result, nil
}

func (c kubeOpsClient) RunAdhoc(pattern, module, param string, inventory *api.Inventory) (*api.Result, error) {
	conn, err := c.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewKubeOpsApiClient(conn)
	request := &api.RunAdhocRequest{
		Inventory: inventory,
		Module:    module,
		Param:     param,
		Pattern:   pattern,
	}
	req, err := client.RunAdhoc(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return req.Result, nil
}

func (c *kubeOpsClient) WatchRun(taskId string, writer io.Writer) error {
	conn, err := c.createConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewKubeOpsApiClient(conn)
	req := &api.WatchRequest{
		TaskID: types.TaskID(taskId),
	}
	server, err := client.WatchResult(context.Background(), req)
	if err != nil {
		return err
	}
	for {
		msg, err := server.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		_, err = writer.Write(msg.Stream)
		if err != nil {
			break
		}
	}
	return nil
}

func (c *kubeOpsClient) GetResult(taskId string) (*api.Result, error) {
	conn, err := c.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewKubeOpsApiClient(conn)
	request := api.GetResultRequest{
		TaskID: types.TaskID(taskId),
	}
	resp, err := client.GetResult(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return resp.Item, nil
}

func (c *kubeOpsClient) ListResult() ([]*api.Result, error) {
	conn, err := c.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewKubeOpsApiClient(conn)
	request := api.ListResultRequest{}
	resp, err := client.ListResult(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (c *kubeOpsClient) createConnection() (*grpc.ClientConn, error) {
	address := fmt.Sprintf("%s:%d", c.host, c.port)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100*1024*1024)))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
