package main

import (
	consulapi "github.com/hashicorp/consul/api"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	ConsulHost = "CONSUL_HTTP_ADDR"
	serviceAccountTokenPath = "/run/secrets/kubernetes.io/serviceaccount/token"
	consulAuthMethod = "auth-method-consul-auth"
)

var (
	consulHost = os.Getenv(ConsulHost)
)


func Token() (string, error) {
	token, err := ioutil.ReadFile(serviceAccountTokenPath)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func main() {

	log.Println("Connecting to consul")
	consulctl, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		panic(err)
	}

	log.Println("Reading k8s consul token")
	token, err := Token()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("token: %s \n", token)

	log.Println("Login to consul using ACL")
	acltoken, _, err := consulctl.ACL().Login(&consulapi.ACLLoginParams{
		AuthMethod:  consulAuthMethod,
		BearerToken: token,
	}, nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("acl token: %s \n",acltoken.SecretID)
	
	kv := consulctl.KV()

	log.Println("Getting kv using auth method and acl")

	pair, _, err := kv.Get("test", &consulapi.QueryOptions{
		Namespace:         acltoken.Namespace,
		Datacenter:        "dc1",
		Token:             acltoken.SecretID,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("key value: %s\n", string(pair.Value))

	time.Sleep(time.Second * 1000)
}
