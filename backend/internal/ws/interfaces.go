package ws

type HubInterface interface {
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	BroadcastMessage(msg WSMessage)
}
