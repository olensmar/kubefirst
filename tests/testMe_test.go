package tests

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/kubefirst/kubefirst/configs"
	"github.com/kubefirst/kubefirst/pkg"
	"github.com/spf13/viper"
	"testing"
	"time"
)

// todo: this is WIP

func TestConsoleArgoCDLink(t *testing.T) {

	opts := append(chromedp.DefaultExecAllocatorOptions[3:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.IgnoreCertErrors,
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	err := chromedp.Run(
		ctx,
		chromedp.Navigate("https://kubefirst.localdev.me"),
		chromedp.WaitVisible("//main/div[1]/div[2]/div/div[4]/div[1]/a"),
	)
	if err != nil {
		t.Error(err)
	}

	var argoCDLinkName string
	err = chromedp.Run(
		ctx,
		chromedp.Text("//main/div[1]/div[2]/div/div[4]/div[1]/a", &argoCDLinkName),
	)

	fmt.Println("---debug---")
	fmt.Println(argoCDLinkName)
	fmt.Println("---debug---")
	if len(argoCDLinkName) == 0 {
		t.Errorf("argoCD link not available")
	}
}

func TestArgoCDUnknownApps(t *testing.T) {

	// load viper file
	config := configs.ReadConfig()
	if err := pkg.SetupViper(config); err != nil {
		t.Error(err)
	}

	// credentials
	login := viper.GetString("argocd.admin.username")
	password := viper.GetString("argocd.admin.password")

	opts := append(chromedp.DefaultExecAllocatorOptions[3:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.IgnoreCertErrors,
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	err := chromedp.Run(
		ctx,
		chromedp.Navigate("https://argocd.localdev.me/"),
		chromedp.WaitVisible("//form/div[1]/div/input"),
	)
	if err != nil {
		t.Error(err)
	}

	err = chromedp.Run(
		ctx,
		chromedp.SendKeys(`//form/div[1]/div/input`, login),
		chromedp.SendKeys(`//form/div[2]/div/input`, password),
		chromedp.Click("//form/div[3]/button"),
	)

	var unknown string
	var unknownItems string
	err = chromedp.Run(
		ctx,
		chromedp.Text("(//DIV[@class='filter__item__label'][text()='Unknown'])[1]", &unknown),
		chromedp.Text("(//DIV[@class='filter__item__label'][text()='Unknown'])[1]/../div[4]", &unknownItems),
	)

	if unknown != "Unknown" {
		t.Errorf("Unknown field not found")
	}

	if unknownItems != "0" {
		t.Errorf("there are (%s) Unknown ArgoCD applications", unknownItems)
	}
}
