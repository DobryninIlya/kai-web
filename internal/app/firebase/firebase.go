package firebase

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=FirebaseAPIInterface --no-internal-structs
type FirebaseAPIInterface interface {
	GetFirebaseUser(ctx context.Context, uid string) (*auth.UserRecord, error)
}

type FirebaseAPI struct {
	Client    *auth.Client
	ProjectID string
}

func (r FirebaseAPI) createFirebaseAuthClient(path, projectId string) (*auth.Client, error) {
	opt := option.WithCredentialsFile(path)
	conf := &firebase.Config{ProjectID: projectId}
	app, err := firebase.NewApp(nil, conf, opt)
	if err != nil {
		return nil, err
	}
	client, err := app.Auth(nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (r FirebaseAPI) GetFirebaseUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := r.Client.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewFirebaseAPI(keyPath, projectID string) (*FirebaseAPI, error) {
	var api FirebaseAPI
	client, err := api.createFirebaseAuthClient(keyPath, projectID)
	if err != nil {
		return nil, err
	}
	api.Client = client
	return &api, nil
}
