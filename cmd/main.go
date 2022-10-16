package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	placement "github.com/guff192/ad-placement-api"
	"github.com/guff192/ad-placement-api/pkg/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	// configuring logger
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)

	// reading configs
	config, err := initConfig()
	if err != nil {
		logrus.Fatalf("error occured while parsing config: %s", err.Error())
	}

	handlers := handler.NewHandler()

	// creating and running server
	srv := new(placement.Server)
	go func() {
		if err := srv.Run(config.Port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Info("Http server started")

	// graceful shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	<-exit

	logrus.Print("Http server shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("error occured while shutting down http server: %s", err.Error())
	}
}

type partnerAddr struct {
	Addr string
	Port int
}

func ParsePartnerAddr(s string) (*partnerAddr, error) {
	splittedStr := strings.Split(s, ":")
	if len(splittedStr) != 2 {
		return nil, errors.New("Unable to parse partner address. Please, check the format and try again.")
	}

	port, err := strconv.Atoi(splittedStr[1])
	if err != nil {
		return nil, errors.New("Unable to parse partner port: \"" + splittedStr[1] + "\"! It should be an integer.")
	}

	return &partnerAddr{
		Addr: splittedStr[0],
		Port: port,
	}, nil
}

type partnerArray []partnerAddr

func (pa *partnerArray) String() string {
	var result string = ""
	if len(*pa) > 0 {
		for _, value := range *pa {
			addr := value.Addr
			port := strconv.Itoa(value.Port)
			result = result + strings.Join([]string{addr, port}, ":") + ", "
		}
	}
	return result
}

func (pa *partnerArray) Set(s string) error {
	values := strings.Split(s, ",")
	if len(values) <= 0 {
		return errors.New("No values for partners! Use -d flag to set them")
	}

	for _, v := range values {
		if address, err := ParsePartnerAddr(v); err != nil {
			return err
		} else {
			*pa = append(*pa, *address)
		}
	}
	return nil
}

type Config struct {
	Port     int
	Partners partnerArray
}

func initConfig() (*Config, error) {
	var partners partnerArray

	port := flag.Int("p", 0, "port to start service")
	flag.Var(&partners, "d", "list of partners in <ip1:port1,ip2:port2...> format")
	flag.Parse()

	if *port == 0 || len(partners) == 0 {
		return nil, errors.New("No port or partners specified")
	} else if len(partners) > 10 {
		return nil, errors.New("Too much partners! You can define up to 10 partners")
	}

	fmt.Println("port is: ", *port)
	fmt.Println("partners are: ", partners.String())

	return &Config{
		Port:     *port,
		Partners: partners,
	}, nil
}
