package testutils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/mock/gomock"
	"path/filepath"
	"runtime"
	"time"
)

const (
	temporalVersion        = "1.22.1"
	temporalContainer      = "postmanq_temporal"
	temporalPort           = 7233
	temporalConfigFilePath = "/etc/temporal/config/dynamicconfig/development_sql.yaml"
	dbContainer            = "postmanq_postgres"
	dbName                 = "postmanq"
	dbPort                 = 5432
	dbUser                 = "postmanq"
	dbPassword             = "postmanq"
	queue                  = "WorkflowTypeTest"
)

var (
	temporalPortId   = fmt.Sprintf("%d/tcp", temporalPort)
	temporalDBPortId = fmt.Sprintf("%d/tcp", dbPort)
	ErrAny           = errors.New("any error")
)

type Suite struct {
	suite.Suite
	Ctx  context.Context
	Ctrl *gomock.Controller
}

func (s *Suite) SetupSuite() {
	s.Ctrl = gomock.NewController(s.T())
	s.Ctx = context.Background()
}

func (s *Suite) TearDownSuite() {
	s.Ctrl.Finish()
}

type TemporalSuite struct {
	Suite
	pool        *Pool
	dbRes       *Resource
	temporalRes *Resource
	Client      client.Client
	Worker      worker.Worker
}

func (s *TemporalSuite) SetupSuite() {
	var err error
	s.Suite.SetupSuite()
	_, path, _, _ := runtime.Caller(0)
	s.pool = GetPool()
	s.dbRes = s.pool.Run(
		WithName(dbContainer),
		WithImage("postgres"),
		WithTag("alpine"),
		WithEnv(Env{
			"POSTGRES_DB":       dbName,
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPassword,
			"listen_addresses":  "'*'",
		}),
	)
	s.pool.Check(func() error {
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			dbUser,
			dbPassword,
			s.dbRes.GetHost(temporalDBPortId),
			s.dbRes.GetPort(temporalDBPortId),
			dbName,
		)
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
		db := bun.NewDB(sqldb, pgdialect.New())
		_, err := db.Exec("SELECT 1")
		return err
	})

	s.temporalRes = s.pool.Run(
		WithName(temporalContainer),
		WithImage("temporalio/auto-setup"),
		WithTag(temporalVersion),
		WithEnv(Env{
			"DB":                       "postgresql",
			"DB_PORT":                  dbPort,
			"DB_HOST":                  s.dbRes.GetHost(temporalDBPortId),
			"POSTGRES_USER":            dbUser,
			"POSTGRES_PWD":             dbPassword,
			"POSTGRES_SEEDS":           dbContainer,
			"DYNAMIC_CONFIG_FILE_PATH": temporalConfigFilePath,
		}),
		WithMount(fmt.Sprintf("%s/development_sql.yaml:%s", filepath.Dir(path), temporalConfigFilePath)),
	)

	s.pool.Check(func() error {
		c, err := client.Dial(client.Options{
			HostPort: fmt.Sprintf("%s:%s", s.temporalRes.GetHost(temporalPortId), s.temporalRes.GetPort(temporalPortId)),
		})
		if err != nil {
			return err
		}

		service := c.WorkflowService()
		resp, err := service.ListNamespaces(s.Ctx, &workflowservice.ListNamespacesRequest{})
		if err != nil {
			return err
		}

		for _, namespace := range resp.Namespaces {
			if namespace.NamespaceInfo.Name == "default" {
				return nil
			}
		}

		time.Sleep(5 * time.Second)
		return errors.New("Could not get default namespace")
	})

	s.Client, err = client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%s", s.temporalRes.GetHost(temporalPortId), s.temporalRes.GetPort(temporalPortId)),
	})
	s.Nil(err)

	s.Worker = worker.New(s.Client, queue, worker.Options{})
	s.Worker.RegisterWorkflowWithOptions(s.SendEventWorkflow, workflow.RegisterOptions{
		Name: queue,
	})
	err = s.Worker.Start()
	s.Nil(err)
}

func (s *TemporalSuite) SendEventWorkflow(ctx workflow.Context) (bool, error) {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	s.T().Log("run workflow")
	return true, nil
}

func (s *TemporalSuite) ExecuteWorkflow() bool {
	run, err := s.Client.ExecuteWorkflow(
		s.Ctx,
		client.StartWorkflowOptions{
			ID:                       uuid.NewString(),
			TaskQueue:                queue,
			WorkflowExecutionTimeout: 5 * time.Second,
		},
		queue,
	)
	s.Nil(err)

	var result bool
	s.Nil(run.Get(s.Ctx, &result))
	return true
}

func (s *TemporalSuite) TearDownSuite() {
	s.Worker.Stop()
	s.Client.Close()
	s.pool.Purge(s.temporalRes)
	s.pool.Purge(s.dbRes)
	s.Suite.TearDownSuite()
}
