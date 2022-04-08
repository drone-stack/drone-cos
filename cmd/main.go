package main

import (
	"os"

	cos "github.com/drone-stack/drone-cos"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "0.0.2"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	// Load env-file if it exists first
	if env := os.Getenv("PLUGIN_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	app := cli.NewApp()
	app.Name = "drone cos plugin"
	app.Usage = "drone cos plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "bucket",
			Usage:  "bucket name",
			EnvVar: "PLUGIN_BUCKET",
		},
		cli.StringFlag{
			Name:   "accesskey",
			Usage:  "access key",
			EnvVar: "PLUGIN_ACCESSKEY",
		},
		cli.StringFlag{
			Name:   "secretkey",
			Usage:  "secret key",
			EnvVar: "PLUGIN_SECRETKEY",
		},
		cli.StringFlag{
			Name:   "region",
			Usage:  "region",
			EnvVar: "PLUGIN_REGION",
		},
		cli.StringFlag{
			Name:   "source",
			Usage:  "source path",
			EnvVar: "PLUGIN_SOURCE",
		},
		cli.StringFlag{
			Name:   "target",
			Usage:  "target path",
			EnvVar: "PLUGIN_TARGET",
		},
		cli.StringFlag{
			Name:   "strip-prefix",
			Usage:  "strip prefix",
			EnvVar: "PLUGIN_STRIP_PREFIX",
		},
		cli.StringFlag{
			Name:   "endpoint",
			Usage:  "endpoint",
			EnvVar: "PLUGIN_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "include",
			Usage:  "include",
			EnvVar: "PLUGIN_INCLUDE",
		},
		cli.StringFlag{
			Name:   "exclude",
			Usage:  "exclude",
			EnvVar: "PLUGIN_EXCLUDE",
		},
		cli.BoolFlag{
			Name:   "autotime",
			Usage:  "last modified time",
			EnvVar: "PLUGIN_AUTOTIME",
		},
		// cli.BoolFlag{
		// 	Name:   "debug",
		// 	Usage:  "debug",
		// 	EnvVar: "PLUGIN_DEBUG",
		// },
		// cli.BoolFlag{
		// 	Name:   "pause",
		// 	Usage:  "pause",
		// 	EnvVar: "PLUGIN_PAUSE",
		// },
		// cli.StringFlag{
		// 	Name:   "proxy",
		// 	Usage:  "proxy",
		// 	EnvVar: "PLUGIN_PROXY",
		// },
		cli.StringFlag{
			Name:   "timeformat",
			Usage:  "time format",
			EnvVar: "PLUGIN_TIMEFORMAT",
			Value:  "0102",
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
			Include:     c.String("include"),
			Exclude:     c.String("exclude"),
		},
		Ext: cos.Ext{
			AutoTime:   c.Bool("autotime"),
			TimeFormat: c.String("timeformat"),
			// Debug:      c.Bool("debug"),
			// Pause:      c.Bool("pause"),
			// Proxy:      c.String("proxy"),
		},
	}

	if err := plugin.Exec(); err != nil {
		logrus.Fatal(err)
	}
	return nil
}
