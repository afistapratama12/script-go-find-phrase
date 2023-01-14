package service

import (
	"database/sql"
	"encoding/json"
	"goscript-final/internal"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	_ "github.com/lib/pq"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

type KeyAPI struct {
	APIETHER []string
	APIBSC   []string
}

type Response struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Result  []Balance `json:"result"`
}

type Balance struct {
	Account string `json:"account"`
	Balance string `json:"balance"`
}

type Service interface {
	ProcessGetPhrase(data []string, length int) error
}

type service struct {
	client     http.Client
	db         *sql.DB
	keyAPI     KeyAPI
	words      []string
	pathResult string
}

func NewService(db *sql.DB, words []string, etherAPI []string, bscAPI []string, pathResult string) *service {
	return &service{
		client: http.Client{},
		db:     db,
		keyAPI: KeyAPI{
			APIETHER: etherAPI,
			APIBSC:   bscAPI,
		},
		words:      words,
		pathResult: pathResult,
	}
}

func (s *service) Ping() error {
	err := s.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (s *service) KeysToSliceString(m map[string]string) []string {
	var result []string

	for k := range m {
		result = append(result, k)
	}

	return result
}

func (s *service) ProcessGetPhrase(length int) {
	wg := &sync.WaitGroup{}
	var ch = make(chan error)

	for i := 0; i < len(s.keyAPI.APIETHER); i++ {
		wg.Add(1)
		go s.WorkerCheck(i, s.keyAPI.APIETHER[i], s.keyAPI.APIBSC[i], length, wg, ch)
	}

	wg.Wait()
	close(ch)

	for err := range ch {
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (s *service) WorkerCheck(i int, apiEther string, apiBSC string, length int, wg *sync.WaitGroup, ch chan error) {
	defer wg.Done()

	var tempResult = make(map[string]*struct{})
	var tempProcess = make(map[string]string) // process every 20 { address : mnemonic }

	var count int = 1

	for {
		mnemonic := internal.GetPhrase(s.words, length)
		if ok := s.checkPhrase(mnemonic); !ok {
			continue
		}

		if _, ok := tempResult[mnemonic]; !ok {
			tempResult[mnemonic] = &struct{}{}
			address := s.GetAddress(mnemonic)
			tempProcess[address] = mnemonic
		} else {
			continue
		}

		if len(tempProcess) >= 20 {
			log.Println("go:", i, "count:", count)

			//process get address and balance
			addresses := s.KeysToSliceString(tempProcess)
			balanceEther, err := s.checkBalance("ether", addresses, apiEther)
			if err != nil {
				ch <- err
			}

			balanceBSC, err := s.checkBalance("bsc", addresses, apiBSC)
			if err != nil {
				ch <- err
			}

			var strFound string

			if len(balanceEther) != 0 {
				for idx := range balanceEther {
					strFound += "ETH " + tempProcess[balanceEther[idx]] + " " + balanceEther[idx] + "\n"
				}
			}

			if len(balanceBSC) != 0 {
				for idx := range balanceBSC {
					strFound += "BSC " + tempProcess[balanceBSC[idx]] + " " + balanceBSC[idx] + "\n"
				}
			}

			if strFound != "" {
				log.Println(strFound)
				// err := internal.WriteOneLine(strFound, s.pathResult)
				// if err != nil {
				// 	ch <- err
				// }

				_, err := s.db.Exec("INSERT INTO results (result) VALUES ($1)", strFound)
				if err != nil {
					ch <- err
				}
			}

			tempProcess = make(map[string]string)
			count++
		}

		if len(tempResult) >= 1_000_000 {
			tempResult = make(map[string]*struct{})
		}
	}
}

func (s *service) GetAddress(mnemonic string) string {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	return account.Address.Hex()
}

func (s *service) checkPhrase(phrase string) bool {
	// check phrase in database
	// if phrase is exist in database, return true
	// else return false

	if ok := bip39.IsMnemonicValid(phrase); !ok {
		return false
	}

	_, err := bip39.EntropyFromMnemonic(phrase)
	return err == nil
}

func (s *service) checkBalance(state string, address []string, key string) ([]string, error) {
	var link string

	if state == "ether" {
		link = "https://api.etherscan.io/api?module=account&action=balancemulti&address=" + strings.Join(address, ",") + "&tag=latest&apikey=" + key
	} else if state == "bsc" {
		link = "https://api.bscscan.com/api?module=account&action=balancemulti&address=" + strings.Join(address, ",") + "&tag=latest&apikey=" + key
	}

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var response Response

	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, err
	}

	var listAddressBalance []string

	for _, v := range response.Result {
		if v.Balance != "0" {
			listAddressBalance = append(listAddressBalance, v.Account)
		}
	}

	return listAddressBalance, nil
}
