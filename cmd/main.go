package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	placement "github.com/guff192/ad-placement-api"
	"github.com/guff192/ad-placement-api/pkg/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)

	if err := initConfig(); err != nil {
		logrus.Fatalf("error occured while parsing config: %s", err.Error())
	}

	handlers := handler.NewHandler()

	srv := new(placement.Server)
	go func() {
		if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Info("Http server started")
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

type flagsArray []partnerAddr

func (fa *flagsArray) String() string {
	var result string = ""
	if len(*fa) > 0 {
		for _, value := range *fa {
			addr := value.Addr
			port := strconv.Itoa(value.Port)
			result = result + strings.Join([]string{addr, port}, ":") + ", "
		}
	}
	return result
}

func (fa *flagsArray) Set(s string) error {
	values := strings.Split(s, ",")
	if len(values) <= 0 {
		return errors.New("No values for partners! Use -d flag to set them")
	}

	for _, v := range values {
		if address, err := ParsePartnerAddr(v); err != nil {
			return err
		} else {
			*fa = append(*fa, *address)
		}
	}
	return nil
}

func initConfig() error {
	var partners flagsArray
	port := flag.Int("p", 0, "port to start service")
	flag.Var(&partners, "d", "list of partners in <ip1:port1,ip2:port2...> format")
	flag.Parse()

	if *port == 0 || len(partners) == 0 {
		return errors.New("No port or partners specified")
	} else if len(partners) > 10 {
		return errors.New("Too much partners! You can define up to 10 partners")
	}

	fmt.Println("port is: ", *port)
	fmt.Println("partners are: ", partners.String())

	return nil
}
