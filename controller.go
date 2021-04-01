package iperf

type Controller struct {
	clients map[string]*Client
	servers map[string]*Server
}

func (c *Controller) StartListener() (err error) {

}

func (c *Controller) StopListener() (err error) {

}

func (c *Controller) StopServer(id string) (err error) {

}

func (c *Controller) NewServer() (server *Server, err error) {

}

func (c *Controller) NewClient(serverAddr string) (client *Client, err error) {

}
