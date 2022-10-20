package placement

import (
	"errors"
	"strconv"
	"strings"
)

// struct representing address of partner
type PartnerAddr struct {
	Addr string
	Port int
}

// parses PartnerAddr from string
func ParsePartnerAddr(s string) (*PartnerAddr, error) {
	splittedStr := strings.Split(s, ":")
	if len(splittedStr) != 2 {
		return nil, errors.New("Unable to parse partner address. Please, check the format and try again.")
	}

	port, err := strconv.Atoi(splittedStr[1])
	if err != nil {
		return nil, errors.New("Unable to parse partner port: \"" + splittedStr[1] + "\"! It should be an integer.")
	}

	return &PartnerAddr{
		Addr: splittedStr[0],
		Port: port,
	}, nil
}

// this type is used to parse in flag values
type PartnerArray []PartnerAddr

// string representation, required by flag.Value interface
func (pa *PartnerArray) String() string {
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

// setter method, required by flag.Value interface
func (pa *PartnerArray) Set(s string) error {
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
