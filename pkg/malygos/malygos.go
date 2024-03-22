package malygos

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/go-logr/zapr"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrz-k8s-incubator/malygos/pkg/api"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

type Malygos struct {
	httpPort int
	k8s      kubernetes.Clientset
}

func New() *Malygos {
	return &Malygos{
		httpPort: 8080,
	}
}

func (m *Malygos) Run() error {
	zapLog, err := zap.NewProduction()
	if err != nil {
		return err
	}

	logger := zapr.NewLogger(zapLog)

	if err := m.readConfiguration(); err != nil {
		logger.Error(err, "failed to read configuration")
		return err
	}

	// TODO: instanciate the k8s client to hold configs

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(loggerConfig()))
	e.Use(middleware.Recover())

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	myAPI := api.NewApiImpl(logger)
	api.RegisterHandlers(e, myAPI)

	// TODO: CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	return e.Start(fmt.Sprintf(":%d", m.httpPort))
}

func loggerConfig() middleware.LoggerConfig {
	loggerConfig := middleware.DefaultLoggerConfig
	loggerConfig.Format = `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
		`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
		`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
		`,"bytes_in":${bytes_in},"bytes_out":${bytes_out},"custom":${custom}}` + "\n"
	loggerConfig.CustomTagFunc = func(c echo.Context, buf *bytes.Buffer) (int, error) {
		switch v := c.Get("username").(type) {
		case string:
			b, err := json.Marshal(struct {
				Username string `json:"username"`
			}{
				Username: v,
			})

			if err != nil {
				return 0, err
			}

			buf.Write(b)
		default:
			buf.WriteString(`{}`)
		}
		return 0, nil
	}

	return loggerConfig
}
