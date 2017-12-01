package config

import (
	"context"
	"fmt"
)

const ProjectID = "souzoh-demo-gcp-001"

type contextKey string

const ProjectIDContextKey contextKey = "ProjectIDKey"

func SetProjectID(parents context.Context, projectID string) context.Context {
	return context.WithValue(parents, ProjectIDContextKey, projectID)
}

func GetProjectID(ctx context.Context) (string, error) {
	v := ctx.Value(ProjectIDContextKey)

	projectID, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("ProjectID not found")
	}

	return projectID, nil
}
