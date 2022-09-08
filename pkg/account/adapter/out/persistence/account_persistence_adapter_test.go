package account_adapter_out_persistence_test

import (
	"context"
	seq "database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	account_adapter_out_persistence "github.com/kjarmicki/go-buckpal/pkg/account/adapter/out/persistence"
	account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbEngine = "mysql"
var dbName = "go-buckpal-db"
var dbContainerName = dbName + "-integration-tests"
var dbProtocol = "tcp"
var dbPort = "3320"
var dbRootPassword = "1"
var dbContainerHealthInterval = time.Millisecond * 250
var dbContainerHealthRetries = 120

var dbDsn = fmt.Sprintf("root:%s@%s(localhost:%s)/%s?parseTime=true", dbRootPassword, dbProtocol, dbPort, dbName)
var dbMigrateDsn = fmt.Sprintf("%s://%s", dbEngine, dbDsn)
var db *gorm.DB

func TestMain(m *testing.M) {
	startTestDb()
	defer stopTestDb()
	migrateTestDb()
	db = connectToTestDb()
	os.Exit(m.Run())
}

func createTestDbContainer(dockerClient *docker.Client) string {
	created, err := dockerClient.ContainerCreate(
		context.Background(),
		&container.Config{
			Env:   []string{"MYSQL_DATABASE=" + dbName, "MYSQL_ROOT_PASSWORD=" + dbRootPassword, "MYSQL_ROOT_HOST=172.17.0.1"},
			Image: "mysql:8.0.30",
			ExposedPorts: nat.PortSet{
				"3306/tcp": struct{}{},
			},
			Healthcheck: &container.HealthConfig{
				Test:        []string{"CMD-SHELL", `mysql -u root -p1 -e "SELECT 1"`},
				Interval:    dbContainerHealthInterval,
				Timeout:     time.Second,
				StartPeriod: time.Second,
				Retries:     dbContainerHealthRetries,
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"3306/tcp": []nat.PortBinding{
					{
						HostIP:   "localhost",
						HostPort: dbPort,
					},
				},
			},
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		dbContainerName,
	)
	if err != nil {
		panic(err)
	}

	return created.ID
}

func getTestDbContainerId(dockerClient *docker.Client) string {
	if containerId, ok := getExistingTestDbContainerId(dockerClient); ok {
		return containerId
	}
	return createTestDbContainer(dockerClient)
}

func getExistingTestDbContainerId(dockerClient *docker.Client) (string, bool) {
	containers, err := dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.TrimPrefix(name, "/") == dbContainerName {
				return container.ID, true
			}
		}
	}
	return "", false
}

func ensureTestDbContainerReady(dockerClient *docker.Client, containerId string) {
	if !testDbContainerIsRunning(dockerClient, containerId) {
		runTestDbContainer(dockerClient, containerId)
	}
	waitForTestDbContainerToBeHealthy(dockerClient, containerId, 1)
}

func testDbContainerIsRunning(dockerClient *docker.Client, containerId string) bool {
	inspect, err := dockerClient.ContainerInspect(context.Background(), containerId)
	if err != nil {
		panic(err)
	}
	return inspect.State.Running
}

func runTestDbContainer(dockerClient *docker.Client, containerId string) {
	err := dockerClient.ContainerStart(context.Background(), containerId, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
}

func waitForTestDbContainerToBeHealthy(dockerClient *docker.Client, containerId string, attempt int) {
	inspect, err := dockerClient.ContainerInspect(context.Background(), containerId)
	if err != nil {
		panic(err)
	}
	if inspect.State.Health.Status == "healthy" {
		return
	}
	if attempt < dbContainerHealthRetries {
		<-time.After(dbContainerHealthInterval)
		waitForTestDbContainerToBeHealthy(dockerClient, containerId, attempt+1)
		return
	}
	panic(errors.New("Database container did not become healthy in the expected time"))
}

func startTestDb() {
	dockerClient, err := docker.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	ensureTestDbContainerReady(dockerClient, getTestDbContainerId(dockerClient))
}

func stopTestDb() {
}

func migrateTestDb() {
	m, err := migrate.New("file://../../../../..//migrations", dbMigrateDsn)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && err.Error() != "no change" {
		panic(err)
	}
}

func connectToTestDb() *gorm.DB {
	sq, err := seq.Open(dbEngine, dbDsn)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sq,
	}))
	if err != nil {
		panic(err)
	}
	return db
}

func clearDb() {
	db.Exec("TRUNCATE TABLE accounts")
	db.Exec("TRUNCATE TABLE activities")
}

func addSampleAccountWithActivity() {
	db.Exec("INSERT INTO accounts (id) VALUES (1)")
	db.Exec("INSERT INTO accounts (id) VALUES (2)")
	db.Exec(`
		INSERT INTO activities (timestamp, owner_account_id, source_account_id, target_account_id, amount)
		VALUES ("2022-08-20 22:48:00", 1, 2, 1, 500)
	`)
	db.Exec(`
		INSERT INTO activities (timestamp, owner_account_id, source_account_id, target_account_id, amount)
		VALUES ("2022-08-30 22:48:00", 1, 1, 2, 300)
	`)
}

func TestAccountPersistenceAdapterLoadAccount(t *testing.T) {
	clearDb()
	addSampleAccountWithActivity()
	activityRepository := account_adapter_out_persistence.NewActivityWindowGormRepository(db)
	accountRepository := account_adapter_out_persistence.NewAccountGormMySqlRepository(db, activityRepository)
	adapter := account_adapter_out_persistence.NewAccountPersistenceAdapter(accountRepository, activityRepository)

	account, err := adapter.LoadAccount(
		context.Background(),
		account_domain.AccountId(1),
		time.Date(2022, 8, 25, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		panic(err)
	}
	activities := account.GetActivities()

	assert.Equal(t, account_domain.AccountId(1), account.GetId())
	assert.Equal(t, 1, len(activities))
	assert.Equal(t, account_domain.Activity{
		Id:              2,
		OwnerAccountId:  1,
		SourceAccountId: 1,
		TargetAccountId: 2,
		Timestamp:       time.Date(2022, time.August, 30, 22, 48, 0, 0, time.UTC),
		Money:           account_domain.NewMoney(300),
	}, activities[0])
	assert.Equal(t, int64(200), account.CalculateBalance().GetAmount())
}
