package main

import (
	"time"

	"go.uber.org/zap"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()

	sugar.Infow("filed to fetch url", "url", "http://example/com", "attempt", "3", "backoff", time.Second)

	sugar.Infof("failed to fetch URL: %s", "http://example.com")
	zap.New()
}
