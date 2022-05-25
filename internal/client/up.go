package client

import (
	"bufio"
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/ztalab/ZASentinel/internal"
	"github.com/ztalab/ZASentinel/internal/bll"
	"github.com/ztalab/ZASentinel/internal/config"
	"github.com/ztalab/ZASentinel/internal/initer"
	"github.com/ztalab/ZASentinel/internal/schema"
	"github.com/ztalab/ZASentinel/pkg/errors"
	"github.com/ztalab/ZASentinel/pkg/util/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// 登录状态
type State int

const (
	StateNotAuthenticated = State(iota)
	StateAuthenticating
	StateAuthenticated
)

type Up struct {
	UserDetail *schema.ControUserDetail
	State      State
	UpCode     string
}

func NewUp() *Up {
	return &Up{
		State: StateNotAuthenticated,
	}
}

func UpCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "up",
		Usage: "Connect to ZASentinel, logging in if needed",
		Action: func(c *cli.Context) error {
			handle := func(ctx context.Context) (func(), error) {
				initCleanFunc, err := internal.Init(ctx,
					internal.SetConfigFile(c.String("conf")),
				)
				if err != nil {
					return func() {}, err
				}
				err = RunUp(ctx)
				if err != nil {
					return func() {}, err
				}
				return func() {
					initCleanFunc()
				}, nil
			}
			return Run(ctx, handle)
		},
	}
}

func RunUp(ctx context.Context) error {
	up := NewUp()
	fmt.Println("----------------------------------------------------------------------")
	fmt.Println("------------------------Interactive UI Start--------------------------")
	fmt.Println("----------------------------------------------------------------------")
	// 预登录
	err := up.preLogin()
	if err != nil {
		return err
	}
	// 获取客户端列表
	client, err := up.printClients()
	if err != nil {
		return err
	}
	config.C.Certificate.CertPem = client.CertPem
	config.C.Certificate.CaPem = client.CaPem
	config.C.Certificate.KeyPem = client.KeyPem

	basicConf, attr, err := initer.InitCert([]byte(config.C.Certificate.CertPem))
	if err != nil {
		return err
	}
	if basicConf.Type != initer.TypeClient {
		return errors.New("Certificate error, not a client certificate")
	}
	go func() {
		err = bll.NewClient().Listen(ctx, attr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("########## start the client proxy #########")
	}()
	return nil
}

// printClients
func (a *Up) printClients() (*schema.ControClient, error) {
	clients, err := GetControClients()
	if err != nil {
		return nil, err
	}
	if len(clients) <= 0 {
		return nil, errors.NewWithStack("You haven't added a client yet")
	}
	scanner := bufio.NewScanner(os.Stdin)
retry:
	fmt.Println("Please select one of the clients")
	for key, item := range clients {
		fmt.Println(key+1, " | ", item.Name)
	}
	fmt.Println("Please enter your client serial number:")
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	if scanner.Text() == "" {
		fmt.Println("----------------------------------------------------------------------")
		fmt.Println("Error: Input errors, please re-enter")
		goto retry
	}
	index, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if index > len(clients) || index <= 0 {
		fmt.Println("Error: Input errors, please re-enter")
		goto retry
	}
	return clients[index-1], nil
}

func (a *Up) preLogin() error {
	// Get whether the device is logged in
	if config.C.Machine.Cookie != "" {
		// Validate cookies
		user, err := a.GetUserDetail()
		if err == nil {
			a.State = StateAuthenticated
			a.UserDetail = user
			return nil
		}
	}
	// Get login link
	upUrl, err := a.GetAuthUrl()
	if err != nil {
		return errors.WithStack(err)
	}
	// Output login connection
	fmt.Println("To authenticate, visit:")
	fmt.Println(upUrl.Data)
	a.State = StateAuthenticating
	a.UpCode = upUrl.GetCode()
	err = a.autoLogin()
	if err != nil {
		return errors.WithStack(err)
	}
	a.State = StateAuthenticated
	return nil
}

func (a *Up) autoLogin() error {
	result, err := a.GetLoginResult(110)
	if err != nil {
		return err
	} else {
		config.C.Machine.SetCookie(result.Data)
		return config.C.Machine.Write()
	}
}

func (a *Up) GetUserDetail() (*schema.ControUserDetail, error) {
	url := fmt.Sprintf("%s/api/v1/user/detail", config.C.Common.ControHost)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp, err := httpControDo(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result schema.ControUserDetail
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

// GetControClients
func GetControClients() (schema.ControClients, error) {
	url := fmt.Sprintf("%s/api/v1/access/client?name=&page=1&limit_num=50", config.C.Common.ControHost)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp, err := httpControDo(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result schema.ControClientResult
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result.Code != 1001 {
		return nil, errors.NewWithStack(result.Message)
	}
	return result.Data.List, nil
}

func (a *Up) GetAuthUrl() (*schema.ControMachineAuthResult, error) {
	url := fmt.Sprintf("%s/api/v1/controlplane/machine/%s", config.C.Common.ControHost, config.C.Machine.MachineId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp, err := httpControDo(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result schema.ControMachineAuthResult
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result.Code != 1001 {
		return nil, errors.NewWithStack(result.Message)
	}
	return &result, nil
}

// GetLoginResult Get login result information
func (a *Up) GetLoginResult(timeout int) (*schema.ControLoginEvent, error) {
	url := fmt.Sprintf("%s/api/v1/controlplane/machine/auth/poll?timeout=%d&category=%s", config.C.Common.ControHost, timeout, a.UpCode)
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result schema.ControLoginResult
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result.Error != "" {
		return nil, errors.NewWithStack(result.Error)
	}
	if len(result.Events) <= 0 {
		return nil, errors.NewWithStack("Login information not obtained")
	}
	return result.Events[0], nil
}

// httpControDo
func httpControDo(req *http.Request) ([]byte, error) {
	// Add request header
	req.Header.Add("Content-type", "application/json;charset=utf-8")
	// Add cookie
	cookie := &http.Cookie{
		Name:  "zta",
		Value: config.C.Machine.Cookie,
	}
	req.AddCookie(cookie)
	// Send request
	resp, err := config.Is.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
