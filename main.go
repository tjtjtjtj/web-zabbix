package main

import (
	"log"

	"github.com/sclevine/agouti"
)

func main() {

	zabbix := new(Zabbix)
	zabbix.SetupEnv()

	capabilities := agouti.Capabilities{
		"chromeOptions": map[string][]string{
			"args": []string{
				"headless",
				"no-sandbox",
			},
		},
	}

	if zabbix.Proxy != "" {
		capabilities = capabilities.Proxy(agouti.ProxyConfig{
			ProxyType: "manual",
			HTTPProxy: zabbix.Proxy})
	}

	driver := agouti.ChromeDriver(agouti.Desired(capabilities))

	err := driver.Start()
	if err != nil {
		log.Fatalf("Failed to start driver:%v", err)
	}
	defer driver.Stop()

	zabbix.Page, err = driver.NewPage()
	if err != nil {
		log.Fatalf("Failed to open page:%v", err)
	}

	zabbix.Login()
	zabbix.ScreenshotALL()

}
