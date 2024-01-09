package Helpers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var settings, _ = LoadSettings()

func LoadSettings() (SettingsStruct, error) {
	var config SettingsStruct
	for {
		configFile, err := os.Open("config.json")
		if err != nil {
			LogError("Failed to Read Config File, Error", err.Error())
			return SettingsStruct{}, err
		}

		defer func(configFile *os.File) {
			err = configFile.Close()
			if err != nil {
				LogError("Failed to Close Config File, Error", err.Error())
			}
		}(configFile)

		decoder := charmap.ISO8859_1.NewDecoder()

		err = json.NewDecoder(configFile).Decode(&config)
		content, err := io.ReadAll(decoder.Reader(configFile))

		_ = json.Unmarshal(content, &config)
		if err != nil {
			LogError("Failed to Read Config File, Error", err.Error())
			return SettingsStruct{}, err
		}

		return config, nil
	}
}

func GetProxy() (string, error) {
	var proxy string

	file, err := os.Open("./Data/Input/Proxies.txt")
	stat, _ := os.Stat("./Data/Input/Proxies.txt")

	if err != nil {
		return "", errors.New("Failed Opening Proxies File")
	}

	defer func(file *os.File) error {
		err = file.Close()
		if err != nil {
			LogError("Failed Closing Proxies File, Error", err.Error())
			return err
		}
		return nil
	}(file)

	if stat.Size() == 0 {
		return "", errors.New("No Proxies in Proxies.txt File, Failed to Create Client!")
	} else {
		scanner := bufio.NewScanner(file)
		lines := []string{}

		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(lines))

		proxy = lines[randomIndex]

	}
	return proxy, nil
}

func FormatToken(token string) string {
	if strings.Contains(token, ":") {
		split := strings.Split(token, ":")
		if len(split) == 3 {
			return split[2]
		} else if len(split) == 2 {
			return split[1]
		}
	}
	return token
}

func ParseVcc(vcc string) (string, string, string, string) {
	var VccCard string
	var VccCvv string
	var VccMonth string
	var VccYear string

	if strings.Contains(vcc, "|") {
		VccCard = strings.Split(vcc, "|")[0]
		VccCvv = strings.Split(vcc, "|")[2]
		VccMonth = strings.Split(vcc, "|")[1][:2]
		VccYear = strings.Split(vcc, "|")[1][len(strings.Split(vcc, "|")[1])-2 : 2]
	} else if strings.Contains(vcc, ":") {
		VccCard = strings.Split(vcc, ":")[0]
		VccCvv = strings.Split(vcc, ":")[2]
		VccMonth = strings.Split(vcc, ":")[1][:2]
		VccYear = strings.Split(vcc, ":")[1][len(strings.Split(vcc, ":")[1])-2:]
	} else {
		LogPanic("Failed to Parse VCC, Check Format")
		time.Sleep(time.Second * 10)
	}

	return VccCard, VccCvv, VccMonth, VccYear
}

func ParsePromo(p string) (string, error) {
	var promo string

	if strings.Contains(p, "promos.discord.gg/") {
		promo = strings.Split(p, "promos.discord.gg/")[1]
	} else if strings.Contains(p, "promotions/") {
		promo = strings.Split(p, "promotions/")[1]
	} else {
		return "", errors.New("Failed to Parse Promo, Check Format")
	}

	return promo, nil
}

func Replacelast(token string) string {
	strLen := len(token)

	if strLen > 20 {

		// Replace the last 20 characters with dots
		modifiedString := token[:strLen-46]

		return modifiedString
	} else {
		// The string is less than 20 characters, no modification needed
		return ""
	}

}

func ClearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin": // Unix-like systems
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows": // Windows
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Unsupported operating system. Cannot clear the screen.")
	}
}

func GetResources(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		tokens := strings.Fields(line)
		if len(tokens) > 0 {
			lines = append(lines, tokens...)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func GetResource(file string, remove bool) string {

	file1, err := os.OpenFile(file, os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer file1.Close()

	scanner := bufio.NewScanner(file1)

	if scanner.Scan() {
		firstLine := scanner.Text()

		if remove {
			file1.Seek(0, 0)

			if err := file1.Truncate(0); err != nil {
				fmt.Println("Error:", err)
				return ""
			}

			for scanner.Scan() {
				_, err := fmt.Fprintln(file1, scanner.Text())
				if err != nil {
					fmt.Println("Error:", err)
					return ""
				}
			}

			if err := file1.Sync(); err != nil {
				fmt.Println("Error:", err)
				return ""
			}
		}

		return firstLine
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

	return ""
}

type Cycle struct {
	Mutex  *sync.Mutex
	Locked []string
	List   []string
	I      int

	WaitTime time.Duration
}

func New(List *[]string) *Cycle {
	rand.Seed(time.Now().UnixNano())

	return &Cycle{
		WaitTime: 50 * time.Millisecond,
		Mutex:    &sync.Mutex{},
		Locked:   []string{},
		List:     *List,
		I:        0,
	}
}

func NewFromFile(Path string) (*Cycle, error) {
	file, err := os.Open(fmt.Sprintf("./Data/%v", Path))
	if err != nil {
		return nil, err
	}
	var lines []string

	defer file.Close()
	defer func() {
		lines = nil
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return New(&lines), nil
}

func (c *Cycle) RandomiseIndex() {
	c.I = rand.Intn(len(c.List)-1) + 1
}

func (c *Cycle) IsInList(Element string) bool {
	for _, v := range c.List {
		if Element == v {
			return true
		}
	}
	return false
}

func (c *Cycle) IsLocked(Element string) bool {
	for _, v := range c.Locked {
		if Element == v {
			return true
		}
	}
	return false
}

func isInList(List *[]string, Element *string) bool {
	for _, v := range *List {
		if *Element == v {
			return true
		}
	}
	return false
}

func (c *Cycle) Next() string {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for {
		c.I++
		if c.I >= len(c.List) {
			c.I = 0
		}

		if !c.IsLocked(c.List[c.I]) {
			return c.List[c.I]
		}

		time.Sleep(c.WaitTime)
	}
}

func (c *Cycle) Lock(Element string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if c.IsInList(Element) {
		c.Locked = append(c.Locked, Element)
	}
}

func (c *Cycle) Unlock(Element string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for i, v := range c.Locked {
		if Element == v {
			c.Locked = append(c.Locked[:i], c.Locked[i+1:]...)
		}
	}
}

func (c *Cycle) ClearDuplicates() int {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	removed := 0
	var list []string
	for _, v := range c.List {
		if !isInList(&list, &v) {
			list = append(list, v)
		} else {
			removed++
		}
	}
	c.List = list
	list = nil

	return removed
}

func (c *Cycle) Remove(Element string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for i, v := range c.List {
		if Element == v {
			c.List = append(c.List[:i], c.List[i+1:]...)
		}
	}

	for i, v := range c.Locked {
		if Element == v {
			c.Locked = append(c.Locked[:i], c.Locked[i+1:]...)
		}
	}
}

func (c *Cycle) LockByTimeout(Element string, Timeout time.Duration) {
	defer c.Unlock(Element)

	c.Lock(Element)
	time.Sleep(Timeout)
}

func AppendTextToFile(text, file string, extra ...string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		LogError("Failed to write to a file!, Error", err.Error())
	}
	var t string
	if len(extra) != 0 {
		t = extra[0]
	}
	_, err = f.Write([]byte(t + text))
	if err != nil {
		return
	}

	return
}

func CheckResources() {
	PromosAmount, err := GetResources("./Data/Input/Promos.txt")
	if err != nil {
		LogError(err.Error(), "N/A")
		return
	}

	VccAmount, err := GetResources("./Data/Input/Vcc's.txt")
	if err != nil {
		LogError(err.Error(), "N/A")
		return
	}

	TokensAmount, err := GetResources("./Data/Input/Tokens.txt")
	if err != nil {
		LogError(err.Error(), "N/A")
		return
	}

	if len(PromosAmount) == 0 || len(VccAmount) == 0 || len(TokensAmount) == 0 {
		LogPanic("Failed to Fetch Required Resources, Please Enter Vcc's/Tokens/Promos")
		return
	} else {
		LogFinished(fmt.Sprintf("Fetched Resources (Vccs: %v | Promos: %v | Tokens: %v)", len(VccAmount), len(PromosAmount), len(TokensAmount)))
	}

	return
}

func RemoveLine(text, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == text {
			continue
		}

		lines = append(lines, line)
	}

	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err = os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for _, line := range lines {
		_, err = fmt.Fprintln(file, line)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
