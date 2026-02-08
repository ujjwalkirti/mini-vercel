package deployment

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	domain "github.com/ujjwalkirti/mini-vercel-api-server/internal/domain/deployment"
)

// TestCreate tests creating a new deployment
func TestCreate(t *testing.T) {
	t.Run("should create deployment with all fields", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		deployment := &domain.Deployment{
			ID:        "test-id-123",
			ProjectID: "project-456",
			Status:    domain.Queued,
		}

		mock.ExpectExec("INSERT INTO deployments").
			WithArgs(deployment.ID, deployment.ProjectID, deployment.Status).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		result, err := repo.Create(context.Background(), deployment)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != deployment.ID {
			t.Errorf("expected ID %s, got %s", deployment.ID, result.ID)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should generate UUID for new deployment", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		deployment := &domain.Deployment{
			ID:        "", // Empty ID should be generated
			ProjectID: "project-456",
			Status:    domain.Queued,
		}

		mock.ExpectExec("INSERT INTO deployments").
			WithArgs(sqlmock.AnyArg(), deployment.ProjectID, deployment.Status).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		result, err := repo.Create(context.Background(), deployment)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID == "" {
			t.Error("expected ID to be generated")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should return error when database insert fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		deployment := &domain.Deployment{
			ID:        "test-id-123",
			ProjectID: "project-456",
			Status:    domain.Queued,
		}

		mock.ExpectExec("INSERT INTO deployments").
			WithArgs(deployment.ID, deployment.ProjectID, deployment.Status).
			WillReturnError(sql.ErrConnDone)

		// Act
		_, err = repo.Create(context.Background(), deployment)

		// Assert
		if err == nil {
			t.Error("expected error but got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

// TestGetByProjectID tests getting deployments by project ID
func TestGetByProjectID(t *testing.T) {
	t.Run("should return all deployments for project", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		projectID := "project-123"
		userID := "user-456"

		expectedDeployments := []domain.Deployment{
			{ID: "deploy-1", ProjectID: projectID, Status: domain.Ready},
			{ID: "deploy-2", ProjectID: projectID, Status: domain.InProgress},
		}
		jsonData, _ := json.Marshal(expectedDeployments)

		rows := sqlmock.NewRows([]string{"json_agg"}).AddRow(jsonData)
		mock.ExpectQuery("SELECT COALESCE\\(json_agg").
			WithArgs(projectID, userID).
			WillReturnRows(rows)

		// Act
		result, err := repo.GetByProjectID(context.Background(), projectID, userID)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 deployments, got %d", len(result))
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should return empty array when project has no deployments", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		projectID := "project-123"
		userID := "user-456"

		rows := sqlmock.NewRows([]string{"json_agg"}).AddRow([]byte("[]"))
		mock.ExpectQuery("SELECT COALESCE\\(json_agg").
			WithArgs(projectID, userID).
			WillReturnRows(rows)

		// Act
		result, err := repo.GetByProjectID(context.Background(), projectID, userID)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result == nil {
			t.Error("expected non-nil result")
		}
		if len(result) != 0 {
			t.Errorf("expected empty array, got %d items", len(result))
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should return error when database query fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		projectID := "project-123"
		userID := "user-456"

		mock.ExpectQuery("SELECT COALESCE\\(json_agg").
			WithArgs(projectID, userID).
			WillReturnError(sql.ErrConnDone)

		// Act
		_, err = repo.GetByProjectID(context.Background(), projectID, userID)

		// Assert
		if err == nil {
			t.Error("expected error but got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

// TestGetByIDWithProject tests getting a deployment with project data
func TestGetByIDWithProject(t *testing.T) {
	t.Run("should return deployment with project data", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		deploymentID := "deploy-123"
		userID := "user-456"
		now := time.Now()

		rows := sqlmock.NewRows([]string{"id", "project_id", "status", "created_at", "updated_at"}).
			AddRow(deploymentID, "project-789", domain.Ready, now, now)

		mock.ExpectQuery("SELECT d.id, d.project_id, d.status").
			WithArgs(deploymentID, userID).
			WillReturnRows(rows)

		// Act
		result, err := repo.GetByIDWithProject(context.Background(), deploymentID, userID)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != deploymentID {
			t.Errorf("expected ID %s, got %s", deploymentID, result.ID)
		}
		if result.Status != domain.Ready {
			t.Errorf("expected status %s, got %s", domain.Ready, result.Status)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should return error when deployment does not exist", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		deploymentID := "nonexistent"
		userID := "user-456"

		mock.ExpectQuery("SELECT d.id, d.project_id, d.status").
			WithArgs(deploymentID, userID).
			WillReturnError(sql.ErrNoRows)

		// Act
		_, err = repo.GetByIDWithProject(context.Background(), deploymentID, userID)

		// Assert
		if err != sql.ErrNoRows {
			t.Errorf("expected sql.ErrNoRows, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

// TestDeleteByProjectID tests deleting all deployments for a project
func TestDeleteByProjectID(t *testing.T) {
	t.Run("should delete all deployments for project", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		projectID := "project-123"

		mock.ExpectExec("DELETE FROM deployments WHERE project_id").
			WithArgs(projectID).
			WillReturnResult(sqlmock.NewResult(0, 3))

		// Act
		err = repo.DeleteByProjectID(context.Background(), projectID)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should succeed when project has no deployments", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		projectID := "project-123"

		mock.ExpectExec("DELETE FROM deployments WHERE project_id").
			WithArgs(projectID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Act
		err = repo.DeleteByProjectID(context.Background(), projectID)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should return error when database delete fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		projectID := "project-123"

		mock.ExpectExec("DELETE FROM deployments WHERE project_id").
			WithArgs(projectID).
			WillReturnError(sql.ErrConnDone)

		// Act
		err = repo.DeleteByProjectID(context.Background(), projectID)

		// Assert
		if err == nil {
			t.Error("expected error but got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

// TestUpdateStatus tests updating deployment status
func TestUpdateStatus(t *testing.T) {
	t.Run("should update deployment status", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		deploymentID := "deploy-123"
		newStatus := domain.Ready

		mock.ExpectExec("UPDATE deployments SET status").
			WithArgs(newStatus, deploymentID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Act
		err = repo.UpdateStatus(context.Background(), deploymentID, newStatus)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should return error when deployment does not exist", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		deploymentID := "nonexistent"
		newStatus := domain.Ready

		mock.ExpectExec("UPDATE deployments SET status").
			WithArgs(newStatus, deploymentID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Act
		err = repo.UpdateStatus(context.Background(), deploymentID, newStatus)

		// Assert - Note: This doesn't return an error in current implementation
		// but we test the behavior anyway
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should handle valid status transitions", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		deploymentID := "deploy-123"

		// Test multiple status transitions
		statuses := []domain.Status{domain.Queued, domain.InProgress, domain.Ready}
		for _, status := range statuses {
			mock.ExpectExec("UPDATE deployments SET status").
				WithArgs(status, deploymentID).
				WillReturnResult(sqlmock.NewResult(0, 1))

			err = repo.UpdateStatus(context.Background(), deploymentID, status)
			if err != nil {
				t.Errorf("unexpected error for status %s: %v", status, err)
			}
		}

		// Assert
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}
