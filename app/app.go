package app

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/net/proxy"
)

var (
	proxyPath  string
	outputPath string
	proxyList  []string
)

var rootCmd = &cobra.Command{
	Use:   "proxy-checker",
	Short: "Check proxies from a text file",
	Long:  "Proxy Checker is a simple tool for checking proxies listed in a text file.",
	Run:   checkProxies,
}

func getProxy() {
	if proxyPath == "" {
		fmt.Println("Usage: proxy-checker --proxypath <path_to_proxy_file>")
		return
	}

	file, err := os.Open(proxyPath)
	if err != nil {
		log.Println("Error opening file:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		proxyList = append(proxyList, line)
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading file:", err)
	}
}

func getProxyClient(ipProxy string) (*http.Client, error) {
	proxyURL, err := url.Parse(ipProxy)
	if err != nil {
		return nil, err
	}

	dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	httpTransport := &http.Transport{Dial: dialer.Dial}
	httpClient := &http.Client{Transport: httpTransport}

	return httpClient, nil
}

func checkProxies(cmd *cobra.Command, args []string) {
	getProxy()

	var wg sync.WaitGroup
	resultChan := make(chan string)

	for _, proxy := range proxyList {
		wg.Add(1)
		go checkProxy(proxy, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		fmt.Println(result)
	}
}

func checkProxy(proxy string, resultChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	url := "https://google.com"
	httpClient, err := getProxyClient(proxy)
	if err != nil {
		resultChan <- fmt.Sprintf("[ %s ] : %s", red("❌"), proxy)
		return
	}

	// Set timeout for the HTTP request
	timeout := 3 * time.Second
	clientWithTimeout := &http.Client{
		Transport: httpClient.Transport,
		Timeout:   timeout,
	}

	response, err := clientWithTimeout.Get(url)
	if err != nil {
		resultChan <- fmt.Sprintf("[ %s ] : %s", red("❌"), proxy)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		resultChan <- fmt.Sprintf("[ %s ] : %s", red("❌"), proxy)
		return
	}

	resultChan <- fmt.Sprintf("[ %s ] : %s", green("✅"), proxy)
}

func red(s string) string {
	return color.New(color.FgRed).SprintFunc()(s)
}

func green(s string) string {
	return color.New(color.FgGreen).SprintFunc()(s)
}

func Run() {
	rootCmd.Flags().StringVarP(&proxyPath, "proxylist", "p", "", "Path to the proxy list file")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "result.txt", "Output for the result")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}

}
