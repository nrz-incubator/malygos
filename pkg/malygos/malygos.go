package malygos

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrz-incubator/malygos/pkg/api"
	"github.com/nrz-incubator/malygos/pkg/malygos/manager"
	"go.uber.org/zap"
)

type Malygos struct {
	http struct {
		Port          int
		EnableRecover bool
	}
	kubeconfig          string
	managementNamespace string
	manager             api.Manager
	logger              logr.Logger
}

func New() *Malygos {
	return &Malygos{
		http: struct {
			Port          int
			EnableRecover bool
		}{
			Port:          8080,
			EnableRecover: true,
		},
	}
}

func (m *Malygos) Run() error {
	zapLog, err := zap.NewProduction()
	if err != nil {
		return err
	}

	m.logger = zapr.NewLogger(zapLog)

	if err := m.readConfiguration(); err != nil {
		m.logger.Error(err, "failed to read configuration")
		return err
	}

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(loggerConfig()))
	if m.http.EnableRecover {
		//e.Use(middleware.Recover())
	}

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	m.manager, err = manager.NewMalygosManager(m.logger, m.kubeconfig, m.managementNamespace)
	if err != nil {
		return err
	}

	myAPI := api.NewApiImpl(m.logger, m.manager)
	api.RegisterHandlers(e, myAPI)

	// TODO: CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	return e.Start(fmt.Sprintf(":%d", m.http.Port))
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
