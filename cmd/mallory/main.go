package main

import (
	"flag"
	"fmt"
	"github.com/gogf/gf/os/glog"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"mallory/library/mallory"
	"net/http"
	"os"
)

var (
	FConfig = flag.String("config", "$HOME/.config/mallory.json", "config file")
	FSuffix = flag.String("suffix", "", "print pulbic suffix for the given domain")
	FReload = flag.Bool("reload", false, "send signal to reload config file")
)

func init() {
	//glog.SetLevel(glog.LEVEL_DEBU)
}

func serve() {
	glog.Info("Starting...")
	glog.Infof("PID: %d", os.Getpid())
	//
	//c, err := mallory.NewConfig(*FConfig)
	//if err != nil {
	//	glog.Fatal(err)
	//}
	//
	//wait := make(chan int)
	//go func() {
	//	normal, err := mallory.NewServer(mallory.NormalSrv, c)
	//	if err != nil {
	//		glog.Fatal(err)
	//	}
	//	glog.Infof("Local normal HTTP proxy: %s", c.File.LocalNormalServer)
	//	glog.Fatal(http.ListenAndServe(c.File.LocalNormalServer, normal))
	//	wait <- 1
	//}()
	//<-wait
}

func printSuffix() {
	host := *FSuffix
	tld, _ := publicsuffix.EffectiveTLDPlusOne(host)
	fmt.Printf("EffectiveTLDPlusOne: %s\n", tld)
	suffix, _ := publicsuffix.PublicSuffix(host)
	fmt.Printf("PublicSuffix: %s\n", suffix)
}

func reload() {
	file, err := mallory.NewConfigFile(os.ExpandEnv(*FConfig))
	if err != nil {
		glog.Fatal(err)
	}
	res, err := http.Get(fmt.Sprintf("http://%s/reload", file.LocalNormalServer))
	if err != nil {
		glog.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Printf("%s\n", body)
}

func main() {
	flag.Parse()

	if *FSuffix != "" {
		printSuffix()
	} else if *FReload {
		reload()
	} else {
		serve()
	}
}
