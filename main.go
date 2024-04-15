package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"learn-go-aws/database"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var logger *zap.Logger
var db *sql.DB

func init() {
	l, _ := zap.NewProduction()
	logger = l
	defer logger.Sync()

	dbConnection, err := database.GetConnection()
	if err != nil {
		logger.Error("error connecting to database", zap.Error(err))
		panic(err)
	}
	err = dbConnection.Ping()

	if err != nil {
		logger.Error("error pinging database", zap.Error(err))
		panic(err)
	}
	db = dbConnection

}

type Event struct {
	Name string `json:"name"`
}

type DefaultResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GetEmployeesResponse struct {
	Employee []*database.Employee `json:"employees"`
}

/* func MyHandler(ctx context.Context, e Event) error {
	logger.Info("in my handler", zap.Any("event", e))
	return nil
} */

func LearnHandler(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	logger.Info("received event", zap.Any("method", event.HTTPMethod), zap.Any("Path", event.Path), zap.Any("Body", event.Body), zap.Any("os.host", os.Getenv("host")), zap.Any("os.password", os.Getenv("password")))

	switch event.Path {
	case "/migrate":
		err := database.CreateEmployeesTable(ctx, db)

		if err != nil {
			body, _ := json.Marshal(&DefaultResponse{
				Status:  fmt.Sprint(http.StatusInternalServerError),
				Message: "Could not create employees table. error: " + err.Error(),
			})

			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       string(body),
			}, err
		}

		body, _ := json.Marshal(&DefaultResponse{
			Status:  fmt.Sprint(http.StatusOK),
			Message: "migrated!",
		})

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}, nil
	case "/employees":
		if event.HTTPMethod == http.MethodPost {
			// create a new employee
			employee := &database.Employee{}
			err := json.Unmarshal([]byte(event.Body), &employee)
			if err != nil {
				body, _ := json.Marshal(&DefaultResponse{
					Status:  fmt.Sprint(http.StatusBadRequest),
					Message: "Error creating a new employee: " + err.Error(),
				})

				return &events.APIGatewayProxyResponse{
					StatusCode: http.StatusBadRequest,
					Body:       string(body),
				}, nil
			}
			err = database.CreateEmployee(ctx, db, employee.Email, employee.FirstName, employee.LastName)
			if err != nil {
				body, _ := json.Marshal(&DefaultResponse{
					Status:  fmt.Sprint(http.StatusInternalServerError),
					Message: "Could not create employees table. error: " + err.Error(),
				})

				return &events.APIGatewayProxyResponse{
					StatusCode: http.StatusBadRequest,
					Body:       string(body),
				}, nil
			}
			body, _ := json.Marshal(&DefaultResponse{
				Status:  fmt.Sprint(http.StatusOK),
				Message: "Created!",
			})

			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       string(body),
			}, nil
		} else if event.HTTPMethod == http.MethodGet {
			// get all employee
			employees, err := database.GetEmployees(ctx, db)
			if err != nil {
				body, _ := json.Marshal(&DefaultResponse{
					Status:  fmt.Sprint(http.StatusInternalServerError),
					Message: "Error retrieving employees: " + err.Error(),
				})

				return &events.APIGatewayProxyResponse{
					StatusCode: http.StatusOK,
					Body:       string(body),
				}, nil
			}

			body, _ := json.Marshal(&GetEmployeesResponse{
				Employee: employees,
			})

			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       string(body),
			}, nil
		}
	default:
		body, _ := json.Marshal(&DefaultResponse{
			Status:  fmt.Sprint(http.StatusOK),
			Message: "Hello from Default Handler",
		})
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}
	}

	return res, nil
}

func main() {
	lambda.Start(LearnHandler)
}
