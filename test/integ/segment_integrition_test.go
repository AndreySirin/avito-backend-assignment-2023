package integ

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/AndreySirin/avito-backend-assignment-2023/internal/entity"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/logger"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/server"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/service"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/storage"
)

const (
	username = "test"
	database = "test"
	passward = "test"
)

type SegmentIntegritionTest struct {
	suite.Suite
	container *postgres.PostgresContainer
	handler   *server.Server
	id        uuid.UUID
}

func (suite *SegmentIntegritionTest) SetupSuite() {
	db, conteiner, err := InitPostgresBD()
	suite.NoError(err)
	suite.container = conteiner
	serv := service.New(logger.New(true), db)
	suite.handler = server.New(logger.New(true), "8080", serv)
}

func (suite *SegmentIntegritionTest) TearDownSuite() {
	err := suite.container.Terminate(context.Background())
	suite.NoError(err)
}

func (suite *SegmentIntegritionTest) TestSegmentCreate() {
	segment := server.CreateSegmentRequest{
		Title:       "footbal",
		Description: "sport for gay",
		AutoUserPrc: 0,
	}
	boby, err := json.Marshal(segment)
	suite.NoError(err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/segments", bytes.NewBuffer(boby))
	request.Header.Set("Content-Type", "application/json")
	suite.handler.HttpServer.Handler.ServeHTTP(recorder, request)

	var id uuid.UUID

	err = json.Unmarshal(recorder.Body.Bytes(), &id)
	suite.NoError(err)
	suite.NotEmpty(id)
	suite.Equal(http.StatusCreated, recorder.Code)
}

func (suite *SegmentIntegritionTest) TestSegmentCreateAutoAdd() {
	segment := server.CreateSegmentRequest{
		Title:       "footbal",
		Description: "sport for gay",
		AutoUserPrc: 50,
	}
	boby, err := json.Marshal(segment)
	suite.NoError(err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/segments", bytes.NewBuffer(boby))
	request.Header.Set("Content-Type", "application/json")
	suite.handler.HttpServer.Handler.ServeHTTP(recorder, request)

	var id uuid.UUID
	err = json.Unmarshal(recorder.Body.Bytes(), &id)

	usersID, err := suite.handler.Service.GetUsersIDsForSegment(context.Background(), id)
	suite.NoError(err)
	suite.NotEmpty(id)
	suite.NotNil(usersID)
	suite.NotEmpty(usersID)
	suite.Equal(http.StatusCreated, recorder.Code)
}

func (suite *SegmentIntegritionTest) TestGetListSegment() {
	var segments []entity.Segment
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/segments", nil)
	request.Header.Set("Content-Type", "application/json")
	suite.handler.HttpServer.Handler.ServeHTTP(recorder, request)
	err := json.Unmarshal(recorder.Body.Bytes(), &segments)
	suite.NoError(err)
	suite.NotEmpty(segments)
	suite.Equal(http.StatusOK, recorder.Code)
}

func (suite *SegmentIntegritionTest) TestGetSegment() {
	var segment entity.Segment
	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	uru := fmt.Sprintf("/api/v1/segments/%v", id)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, uru, nil)
	request.Header.Set("Content-Type", "application/json")
	suite.handler.HttpServer.Handler.ServeHTTP(recorder, request)
	err := json.Unmarshal(recorder.Body.Bytes(), &segment)
	suite.NoError(err)
	suite.NotEmpty(segment)
	suite.Equal(http.StatusOK, recorder.Code)
}

func (suite *SegmentIntegritionTest) TestDeleteSegments() {
	id := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	uru := fmt.Sprintf("/api/v1/segments/%v", id)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, uru, nil)
	request.Header.Set("Content-Type", "application/json")
	suite.handler.HttpServer.Handler.ServeHTTP(recorder, request)

	s, err := suite.handler.Service.GetSegment(context.Background(), id)
	suite.Nil(s)
	suite.Error(err)
	suite.Equal(http.StatusOK, recorder.Code)
}

func (suite *SegmentIntegritionTest) TestUpdateSegment() {
	req := server.UpdateSegmentRequest{
		Title:       "новый сегмент",
		Description: "на самом деле ничего нового",
		AutoUserPrc: 0,
	}
	body, err := json.Marshal(req)
	suite.NoError(err)

	id := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	uru := fmt.Sprintf("/api/v1/segments/%v", id)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, uru, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	suite.handler.HttpServer.Handler.ServeHTTP(recorder, request)

	suite.Equal(http.StatusOK, recorder.Code)
	suite.NotEmpty(recorder.Body)
}

func TestRunSegmentSuite(t *testing.T) {
	suite.Run(t, new(SegmentIntegritionTest))
}

func InitPostgresBD() (*storage.Storage, *postgres.PostgresContainer, error) {
	ctx := context.Background()
	lg := logger.New(true)
	container, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(database),
		postgres.WithUsername(username),
		postgres.WithPassword(passward),
		postgres.WithInitScripts(
			"/home/andrey/GolandProjects/avito-backend-assignment-2023/deploy/local/init.sql",
		),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, nil, err
	}

	info, err := container.Inspect(ctx)
	if err != nil {
		return nil, nil, err
	}
	port := info.NetworkSettings.Ports["5432/tcp"][0].HostPort

	stor, err := storage.New(lg, username, passward, database, "0.0.0.0:"+port)
	if err != nil {
		return nil, nil, err
	}
	return stor, container, err
}
