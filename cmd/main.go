package main

import (
	"os"

	cos "github.com/drone-stack/drone-cos"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "unknown"
)

func main() {
	// Load env-file if it exists first
	if env := os.Getenv("PLUGIN_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	app := cli.NewApp()
	app.Name = "docker plugin"
	app.Usage = "docker plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "bucket",
			Usage:  "bucket name",
			EnvVar: "PLUGIN_BUCKET",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "accesskey",
			Usage:  "access key",
			EnvVar: "PLUGIN_ACCESSKEY",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "secretkey",
			Usage:  "secret key",
			EnvVar: "PLUGIN_SECRETKEY",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "region",
			Usage:  "region",
			EnvVar: "PLUGIN_REGION",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "source",
			Usage:  "source path",
			EnvVar: "PLUGIN_SOURCE",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "target",
			Usage:  "target path",
			EnvVar: "PLUGIN_TARGET",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "strip-prefix",
			Usage:  "strip prefix",
			EnvVar: "PLUGIN_STRIP_PREFIX",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "endpoint",
			Usage:  "endpoint",
			EnvVar: "PLUGIN_ENDPOINT",
			Value:  "",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := cos.Plugin{
		Cos: cos.Cos{
			Bucket:      c.String("bucket"),
			AccessKey:   c.String("accesskey"),
			SecretKey:   c.String("secretkey"),
			Region:      c.String("region"),
			Source:      c.String("source"),
			Target:      c.String("target"),
			StripPrefix: c.String("strip-prefix"),
			Endpoint:    c.String("endpoint"),
		},
	}

	if err := plugin.Exec(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	return nil
}
