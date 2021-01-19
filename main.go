package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"regexp"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ansd/lastpass-go"
	log "github.com/sirupsen/logrus"
)

var (
	stop       = make(chan struct{})
	done       = make(chan struct{})
	loopDB     = time.NewTicker(60 * time.Second)
	cf         ConfigFile
	err        error
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func main() {

	setLogLevel("debug")

	configPath := flag.String("config", "config.toml", "config file path")
	lastpassLogin := flag.String("username", "", "Username")
	lastpassPassword := flag.String("password", "", "password")
	lastpassOTP := flag.String("otp", "", "otp")
	flag.Parse()

	if *lastpassLogin == "" || !isEmailValid(*lastpassLogin) {
		log.Errorf("Login is empty or invalid")
		os.Exit(1)
	}

	if *lastpassPassword == "" {
		log.Errorf("Password is empty")
		os.Exit(1)
	}

	if _, err := toml.DecodeFile(*configPath, &cf); err != nil {
		log.Errorf("Error decoding toml config : %v", err)
		return
	}

	setLogLevel(cf.Log.Level)

	lp := new(lpass)
	dk := new(docker)

	log.Debugf("connection to lastpass in progress")

	// authenticate with LastPass servers

	if cf.LastPass.TwoFactor == "disable" || cf.LastPass.TwoFactor == "" {
		lp.Client, err = lastpass.NewClient(context.Background(), *lastpassLogin, *lastpassPassword)
		if err != nil {
			log.Errorf("Authentification Error : %s", err)
		}
	} else if cf.LastPass.TwoFactor == "push" {
		log.Debugf("Waiting 2FA validation")
		lp.Client, err = lastpass.NewClient(context.Background(), *lastpassLogin, *lastpassPassword)
		if err != nil {
			log.Errorf("Authentification Error : %s", err)
		}
	} else if cf.LastPass.TwoFactor == "otp" {
		var re = regexp.MustCompile(`(?m)[0-9]{6}`)
		if !re.MatchString(*lastpassOTP) {
			log.Errorf("OTP format is invalid")
			os.Exit(1)
		}
		lp.Client, err = lastpass.NewClient(context.Background(), *lastpassLogin, *lastpassPassword, lastpass.WithOneTimePassword(*lastpassOTP))
		if err != nil {
			log.Errorf("Authentification Error : %s", err)
		}
	} else {
		log.Errorf("Lastpass Configuration is invalid !")
		os.Exit(1)
	}

	// connect docker
	err := dk.Connect()
	if err != nil {
		log.Errorf("%s", err)
	}

	stop := make(chan struct{}, 1)

	log.Infof("Starting sync secret")
	err = SyncSecret(lp, dk)
	if err != nil {
		log.Errorf("%s", err)
	}
	log.Infof("End sync secret")

	go loop(lp, dk)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	lp.Logout(ctx)
	close(stop)
	close(done)
	loopDB.Stop()
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Info("shutting down")
	os.Exit(0)

}

func setLogLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
		log.Debug("LogLevel = Debug")
	case "info":
		log.SetLevel(log.InfoLevel)
		log.Info("LogLevel = Info")
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	}
}

func loop(lp *lpass, dk *docker) {
LOOP:
	for {
		select {
		case <-stop:
			break LOOP
		case <-loopDB.C:
			log.Infof("Starting sync secret")
			err := SyncSecret(lp, dk)
			if err != nil {
				log.Errorf("%s", err)
			}
			log.Infof("End sync secret")
		default:
		}
	}
	done <- struct{}{}
}

// isEmailValid checks if the email provided passes the required structure and length.
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
