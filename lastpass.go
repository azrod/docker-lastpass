package main

import (
	"context"
	"fmt"

	"github.com/ansd/lastpass-go"
)

type lpass struct {
	*lastpass.Client
}

func (client *lpass) GetSecret(Name string) (Secret *lastpass.Account, err error) {
	accounts, err := client.Accounts(context.Background())
	if err != nil {
		return nil, err
	}
	for a := range accounts {
		if accounts[a].Name == Name {
			return accounts[a], nil
		}
	}
	err = fmt.Errorf("Secret %s not found", Name)
	return nil, err
}

func (client *lpass) GetSecretsFromGroup(GroupName string) (Secrets []*lastpass.Account, err error) {
	accounts, err := client.Accounts(context.Background())
	if err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		err = fmt.Errorf("No Secret found in group %s ", GroupName)
	} else {
		for a := range accounts {
			if accounts[a].Group == GroupName {
				Secrets = append(Secrets, accounts[a])
			}
		}
	}
	return Secrets, err
}
