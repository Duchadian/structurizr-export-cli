package internal

import (
	"encoding/base64"
	"fmt"
	"github.com/apex/log"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/ysmood/gson"
	"os"
	"strings"
	"time"
)

func waitForDiagramLoaded(page *rod.Page) {
	for {
		isLoaded := page.MustEval("() => structurizr.scripting.isDiagramRendered()")
		if isLoaded.Bool() {
			break
		}
		log.Debug("Not loaded yet, waiting")
		time.Sleep(time.Duration(50))
	}
}

func saveBase64AsPNG(image string, filename string) {
	imgBytes, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		log.Warnf("Can't decode image %s", image)
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Warnf("Can't create file %s", filename)
	}
	defer f.Close()

	if _, err := f.Write(imgBytes); err != nil {
		log.Warnf("Can't write bytes to file %s", filename)
	}

	if err := f.Sync(); err != nil {
		log.Warnf("Can't sync file %s", filename)
	}
}

func setupBrowser(rodUrl string) *rod.Browser {
	var browser *rod.Browser
	if rodUrl != "" {
		// if we have a remote rod instance, i.e. in docker
		l := launcher.MustNewManaged(rodUrl)
		l.Headless(false).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")
		browser = rod.New().Client(l.MustClient()).MustConnect()
	} else {
		// local run, but we still want "headful" due to image rendering
		l := launcher.New()
		l.Headless(false)
		browser = rod.New().ControlURL(l.MustLaunch()).MustConnect()
	}
	return browser
}

func ExtractImages(url string, rodUrl string) {
	log.Infof("Running with url '%s' and rod url '%s'", url, rodUrl)
	browser := setupBrowser(rodUrl)
	page := browser.MustPage(url)

	page.MustWaitNavigation()

	// expose savePNG
	page.MustExpose(
		"savePNG",
		func(json gson.JSON) (interface{}, error) {
			filename := json.Get("filename").String()
			content := json.Get("image").String()
			content = strings.ReplaceAll(content, `data:image/png;base64,`, "")
			log.Infof("Writing size %s to file %s", len(content), filename)
			saveBase64AsPNG(content, filename)
			return nil, nil
		},
	)

	page.MustWaitStable()
	waitForDiagramLoaded(page)

	var views gson.JSON
	views = page.MustWaitStable().MustEval(
		"() => structurizr.scripting.getViews()",
	)

	os.Mkdir("export", os.FileMode(0774))

	for _, val := range views.Arr() {
		log.Infof("Looking for view %s", val)
		key := val.Get("key").String()

		page.MustEval(
			"(k) => structurizr.scripting.changeView(k)",
			key,
		)
		waitForDiagramLoaded(page)
		page.MustEval(
			`(f) => structurizr.scripting.exportCurrentDiagramToPNG({ includeMetadata: true, crop: false},
					function(png) { window.savePNG({"image": png, "filename": f}) }
				)`,
			fmt.Sprintf("export/%s.png", key),
		)
	}
}
