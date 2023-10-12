package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// PodInfo represents information about a Kubernetes pod.
type PodInfo struct {
	PodIP    string
	Hostname string
}

func main() {
	pods, err := getPodInfo()
	if err != nil {
		fmt.Printf("Failed to get pod info: %v\n", err)
		os.Exit(1)
	}

	checksumFile := "/tmp/pods-cs.txt"
	digest := calculateChecksum(pods)

	// Read the previous checksum from the file
	old, err := ioutil.ReadFile(checksumFile)
	if err != nil {
		old = []byte("")
	}

	// If there are no changes in the commands that would be run, run one command and exit 0
	if string(old) == digest {
		runProxySQLReload()
		fmt.Println("No changes detected. Reloaded ProxySQL configuration.")
		os.Exit(0)
	}

	// Write the new checksum to the file for the next run
	if err := ioutil.WriteFile(checksumFile, []byte(digest), 0644); err != nil {
		fmt.Printf("Failed to write checksum file: %v\n", err)
		os.Exit(1)
	}

	commands := createCommands(pods)

	if err := runProxySQLCommands(commands); err != nil {
		fmt.Printf("Failed to run ProxySQL commands: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Ran ProxySQL commands:", commands)
}

func getPodInfo() ([]PodInfo, error) {
	serviceAccountDir := "/run/secrets/kubernetes.io/serviceaccount"
	namespaceFile := filepath.Join(serviceAccountDir, "namespace")
	tokenFile := filepath.Join(serviceAccountDir, "token")
	//caCertFile := filepath.Join(serviceAccountDir, "ca.crt")

	namespace, err := ioutil.ReadFile(namespaceFile)
	if err != nil {
		return nil, err
	}

	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("https://kubernetes.default.svc/api/v1/namespaces/%s/pods", strings.TrimSpace(string(namespace)))

	httpClient := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(string(token)))

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request failed with code: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	items, ok := data["items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Invalid JSON response")
	}

	pods := []PodInfo{}

	for _, item := range items {
		pod, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Invalid JSON response")
		}

		labels, ok := pod["metadata"].(map[string]interface{})["labels"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Invalid JSON response")
		}

		instance, ok := labels["app.kubernetes.io/instance"].(string)
		if !ok || instance == "proxysql-cluster" {
			continue
		}

		podIP, ok := pod["status"].(map[string]interface{})["podIP"].(string)
		if !ok {
			return nil, fmt.Errorf("Invalid JSON response")
		}

		name, ok := pod["metadata"].(map[string]interface{})["name"].(string)
		if !ok {
			return nil, fmt.Errorf("Invalid JSON response")
		}

		pods = append(pods, PodInfo{PodIP: podIP, Hostname: name})
	}

	sort.Slice(pods, func(i, j int) bool {
		return pods[i].PodIP < pods[j].PodIP
	})

	return pods, nil
}

func calculateChecksum(pods []PodInfo) string {
	data := []string{}

	for _, pod := range pods {
		data = append(data, fmt.Sprintf("%s:%s", pod.PodIP, pod.Hostname))
	}

	sort.Strings(data)

	return fmt.Sprintf("%x", data)
}

func createCommands(pods []PodInfo) string {
	commands := []string{"DELETE FROM proxysql_servers"}

	for _, pod := range pods {
		commands = append(commands, fmt.Sprintf("INSERT INTO proxysql_servers VALUES ('%s', 6032, 0, '%s')", pod.PodIP, pod.Hostname))
	}

	commands = append(commands, "LOAD PROXYSQL SERVERS TO RUNTIME", "LOAD MYSQL VARIABLES TO RUNTIME", "LOAD MYSQL SERVERS TO RUNTIME", "LOAD MYSQL USERS TO RUNTIME", "LOAD MYSQL QUERY RULES TO RUNTIME", "SAVE PROXYSQL SERVERS TO DISK")

	return strings.Join(commands, "; ")
}

func runProxySQLCommands(commands string) error {
	db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:6032)/")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(commands)
	return err
}

func runProxySQLReload() {
	db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:6032)/")
	if err != nil {
		fmt.Printf("Failed to connect to ProxySQL: %v\n", err)
		return
	}
	defer db.Close()

	_, err = db.Exec("LOAD PROXYSQL SERVERS TO RUNTIME")
	if err != nil {
		fmt.Printf("Failed to reload ProxySQL: %v\n", err)
		return
	}
}
