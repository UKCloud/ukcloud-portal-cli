package ukc

import (
	"bufio"
	"fmt"
	"github.com/ukcloud/ukcloud-portal-api/api"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

func PromptAccount(papi *api.API) int {
	var err error
	var accounts []api.Accounts
	accounts, err = papi.GetAccounts()

	flags := tabwriter.AlignRight | tabwriter.Debug
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 1, ' ', flags)

	fmt.Fprintln(w, "ID\tAccount Name\t")
	for _, account := range accounts {
		fmt.Fprintln(w, strconv.Itoa(account.ID)+"\t"+account.Name+"\t")
	}
	w.Flush()

	firstAccount := accounts[0].ID
	fmt.Printf("Please enter accountID [%v]:", strconv.Itoa(firstAccount))

	reader := bufio.NewReader(os.Stdin)

	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	response = strings.ToLower(strings.TrimSpace(response))

	if response == "" {
		response = strconv.Itoa(firstAccount)
	}

	var accountID int64
	accountID, err = strconv.ParseInt(response, 10, 64)

	if err != nil && len(err.Error()) > 0 {
		fmt.Println("Invalid Account ID")
		os.Exit(1)
	}

	return int(accountID)
}

func PromptVorg(papi *api.API, accountID int) int {
	var err error
	var vorgs api.VorgsArray
	vorgs, err = papi.GetVorgs(accountID)

	flags := tabwriter.AlignRight | tabwriter.Debug
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 1, ' ', flags)

	fmt.Fprintln(w, "ID\tVorgs Name\t")
	for _, vorgs := range vorgs.Data {
		fmt.Fprintln(w, vorgs.ID+"\t"+vorgs.Attributes.Name+"\t")
	}
	w.Flush()

	firstVorg := vorgs.Data[0].ID
	fmt.Printf("Please enter VorgID [%v]:", firstVorg)

	reader := bufio.NewReader(os.Stdin)

	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	response = strings.ToLower(strings.TrimSpace(response))

	if response == "" {
		response = firstVorg
	}

	var VorgID int64
	VorgID, err = strconv.ParseInt(response, 10, 64)

	if err != nil && len(err.Error()) > 0 {
		fmt.Println("Invalid Vorg ID")
		os.Exit(1)
	}

	return int(VorgID)
}

func PromptVdcName(papi *api.API) string {
	reader := bufio.NewReader(os.Stdin)

	response := ""

	for response == "" || len(response) == 0 {
		fmt.Println("Please enter a name for your new VDC:")
		response, _ = reader.ReadString('\n')
	}

	return strings.Trim(response, "\n")
}
