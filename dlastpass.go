package main

import (
	"context"
	"strings"

	"github.com/ansd/lastpass-go"
	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

// SyncSecret ...
func SyncSecret(lp *lpass, dk *docker) (err error) {
	ds := DockerSecrets{}
	lps := LastpassSecrets{}

	var secret *lastpass.Account
	var secrets []*lastpass.Account
	var allSecrets []*lastpass.Account

	log.Debugf("%v groups & %v list", len(cf.Secrets.Groups), len(cf.Secrets.Lists))

	if len(cf.Secrets.Groups) != 0 {
		for _, group := range cf.Secrets.Groups {
			secrets, err = lp.GetSecretsFromGroup(group)
			if err != nil {
				return
			}
			for s := range secrets {
				allSecrets = append(allSecrets, secrets[s])
			}
		}
	}

	if len(cf.Secrets.Lists) != 0 {
		for _, list := range cf.Secrets.Lists {
			secret, err = lp.GetSecret(list)
			if err != nil {
				return
			}
			allSecrets = append(allSecrets, secret)
		}
	}

	if len(allSecrets) == 0 {
		log.Errorf("No secrets found.")
		return
	}

	for s := range allSecrets {
		lps[allSecrets[s].Name] = allSecrets[s]
	}

	dsecrets, err := dk.SecretList(context.Background(), types.SecretListOptions{})
	if err != nil {
		return
	}

	for _, v := range dsecrets {
		ds[v.Spec.Name] = v
	}

	for k, v := range lps {

		if ds["lastpass_"+strings.ReplaceAll(k, " ", "-")+"_Username"].Spec.Name != "lastpass_"+strings.ReplaceAll(k, " ", "-")+"_Username" {
			err = dk.CreateUsername(v)
			if err != nil {
				return
			}
		}
		if ds["lastpass_"+strings.ReplaceAll(k, " ", "-")+"_Password"].Spec.Name != "lastpass_"+strings.ReplaceAll(k, " ", "-")+"_Password" {
			err = dk.CreatePassword(v)
			if err != nil {
				return
			}
		}
	}

	return
}
