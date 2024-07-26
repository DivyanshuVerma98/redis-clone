package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

type Redis struct {
	// To store connect clients
	ClientMap map[string]*Client
	// To store data into cache
	// Using sync.Map to avoid race conditions
	CacheMap sync.Map
	// To support sub/pub on a channel
	ChannelMap map[string]map[*Client]bool
}

func CreateRedis() *Redis {
	return &Redis{
		ClientMap:  map[string]*Client{},
		CacheMap:   sync.Map{},
		ChannelMap: map[string]map[*Client]bool{},
	}
}

// Adding newly connected client in the client list
func (r *Redis) AddClient(client *Client) error {
	log.Println("Adding client ->", client.conn.RemoteAddr().String())
	r.ClientMap[client.conn.RemoteAddr().String()] = client
	return nil
}

// Adding newly connected client in the client list
func (r *Redis) RemoveClient(client *Client) error {
	log.Println("Removing client ->", client.conn.RemoteAddr().String())
	delete(r.ClientMap, client.conn.RemoteAddr().String())
	return nil
}

// For listing all the connected clients
func (r *Redis) ListClient() (string, error) {
	var res string
	for address, client := range r.ClientMap {
		res += (address + " " + client.conn.LocalAddr().String() + "\n")
	}
	return res, nil
}

// Storing key, value in the cache
func (r *Redis) Set(key, val string) (string, error) {
	r.CacheMap.Store(key, val)
	return val, nil
}

// Getting key, value from the cache
func (r *Redis) Get(key string) (string, error) {
	val, exists := r.CacheMap.Load(key)
	if !exists {
		return "", fmt.Errorf("key doesn't exists")
	}
	return val.(string), nil
}

// Deleting a key from the cache
func (r *Redis) Delete(key string) (string, error) {
	_, exists := r.CacheMap.LoadAndDelete(key)
	if !exists {
		return "", fmt.Errorf("key doesn't exists")
	}
	return "OK", nil
}

// To subscribe to a channel
func (r *Redis) Subscribe(channelName string, client *Client) (string, error) {
	_, exists := r.ChannelMap[channelName]
	if !exists {
		r.ChannelMap[channelName] = map[*Client]bool{client: true}
	} else {
		r.ChannelMap[channelName][client] = true
	}
	return "Ok", nil
}

// Publish a message on channel
func (r *Redis) Publish(channelName, msg string, client *Client) (string, error) {
	subMap, exists := r.ChannelMap[channelName]
	if !exists {
		return "", fmt.Errorf("channel doesn't exists")
	}
	var wg sync.WaitGroup
	// Sending msg to all the subscribers
	for sub := range subMap {
		if sub == client {
			continue
		}
		// Spawning a sperate goroutine for each subscriber
		wg.Add(1)
		go func(client *Client) {
			client.SendResponse(client.conn.LocalAddr().String() + "\n" + msg)
			defer wg.Done()
		}(sub)
	}
	wg.Wait()
	return "Ok", nil
}

// Handling Supported Commands - GET, SET etc..
func (r *Redis) HandleCommand(cmd string, args []string, client *Client) {
	var err error
	var res string
	switch strings.ToUpper(cmd) {

	case "PING":
		res = "Pong"

	case "ECHO":
		res = strings.Join(args, " ")

	case "LIST":
		res, err = r.ListClient()

	case "GET":
		if len(args) != 1 {
			res = "Invalid number of arguments"
		} else {
			res, err = r.Get(args[0])
		}

	case "SET":
		if len(args) != 2 {
			res = "Invalid number of arguments"
		} else {
			res, err = r.Set(args[0], args[1])
		}
	case "DEL":
		if len(args) != 1 {
			res = "Invalid number of arguments"
		} else {
			res, err = r.Delete(args[0])
		}

	case "SUB":
		if len(args) != 1 {
			res = "Invalid number of arguments"
		} else {
			res, err = r.Subscribe(args[0], client)
		}

	case "PUB":
		if len(args) < 2 {
			res = "Invalid number of arguments"
		} else {
			res, err = r.Publish(args[0], strings.Join(args[1:], " "), client)
		}

	case "QUIT":
		r.RemoveClient(client)
		r.ListClient()
		err = client.CloseConnection()
	}
	if err != nil {
		client.SendResponse(err.Error())
	}
	if len(res) > 0 {
		client.SendResponse(res)
	}
}
