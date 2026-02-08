package project

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	deploymentdomain "github.com/ujjwalkirti/mini-vercel-api-server/internal/domain/deployment"
	domain "github.com/ujjwalkirti/mini-vercel-api-server/internal/domain/project"
)

// TestCreate tests creating a new project
func TestCreate(t *testing.T) {
	t.Run("should create project with all fields", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		customDomain := "custom.example.com"
		project := &domain.Project{
			ID:           "project-123",
			Name:         "Test Project",
			GitURL:       "https://github.com/user/repo",
			SubDomain:    "test-subdomain",
			CustomDomain: &customDomain,
			UserID:       "user-456",
		}

		mock.ExpectExec("INSERT INTO projects").
			WithArgs(project.ID, project.Name, project.GitURL, project.SubDomain, project.CustomDomain, project.UserID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		err = repo.Create(context.Background(), project)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should generate UUID for new project", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		project := &domain.Project{
			ID:           "", // Empty ID should be generated
			Name:         "Test Project",
			GitURL:       "https://github.com/user/repo",
			SubDomain:    "test-subdomain",
			CustomDomain: nil,
			UserID:       "user-456",
		}

		mock.ExpectExec("INSERT INTO projects").
			WithArgs(sqlmock.AnyArg(), project.Name, project.GitURL, project.SubDomain, project.CustomDomain, project.UserID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		err = repo.Create(context.Background(), project)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if project.ID == "" {
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
		project := &domain.Project{
			ID:        "project-123",
			Name:      "Test Project",
			GitURL:    "https://github.com/user/repo",
			SubDomain: "test-subdomain",
			UserID:    "user-456",
		}

		mock.ExpectExec("INSERT INTO projects").
			WithArgs(project.ID, project.Name, project.GitURL, project.SubDomain, project.CustomDomain, project.UserID).
			WillReturnError(sql.ErrConnDone)

		// Act
		err = repo.Create(context.Background(), project)

		// Assert
		if err == nil {
			t.Error("expected error but got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should handle duplicate subdomain gracefully", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		project := &domain.Project{
			ID:        "project-123",
			Name:      "Test Project",
			GitURL:    "https://github.com/user/repo",
			SubDomain: "duplicate-subdomain",
			UserID:    "user-456",
		}

		// Simulate unique constraint violation
		mock.ExpectExec("INSERT INTO projects").
			WithArgs(project.ID, project.Name, project.GitURL, project.SubDomain, project.CustomDomain, project.UserID).
			WillReturnError(sql.ErrNoRows) // In real scenario, this would be a unique constraint error

		// Act
		err = repo.Create(context.Background(), project)

		// Assert
		if err == nil {
			t.Error("expected error for duplicate subdomain")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

// TestListByUser tests listing projects by user ID
func TestListByUser(t *testing.T) {
	t.Run("should return all projects for user", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		userID := "user-123"
		now := time.Now()

		deployments := []deploymentdomain.Deployment{
			{ID: "deploy-1", ProjectID: "project-1", Status: "READY"},
		}
		deploymentsJSON, _ := json.Marshal(deployments)

		rows := sqlmock.NewRows([]string{
			"id", "name", "git_url", "subdomain", "custom_domain", "user_id", "created_at", "updated_at", "deployments",
		}).AddRow("project-1", "Project 1", "https://github.com/user/repo1", "subdomain1", "", userID, now, now, deploymentsJSON).
			AddRow("project-2", "Project 2", "https://github.com/user/repo2", "subdomain2", "", userID, now, now, []byte("[]"))

		mock.ExpectQuery("WITH latest_deployments AS").
			WithArgs(userID).
			WillReturnRows(rows)

		// Act
		result, err := repo.ListByUser(context.Background(), userID)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 projects, got %d", len(result))
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should return empty array when user has no projects", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		userID := "user-123"

		rows := sqlmock.NewRows([]string{
			"id", "name", "git_url", "subdomain", "custom_domain", "user_id", "created_at", "updated_at", "deployments",
		})

		mock.ExpectQuery("WITH latest_deployments AS").
			WithArgs(userID).
			WillReturnRows(rows)

		// Act
		result, err := repo.ListByUser(context.Background(), userID)

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
		userID := "user-123"

		mock.ExpectQuery("WITH latest_deployments AS").
			WithArgs(userID).
			WillReturnError(sql.ErrConnDone)

		// Act
		_, err = repo.ListByUser(context.Background(), userID)

		// Assert
		if err == nil {
			t.Error("expected error but got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should include most recent deployment for each project", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		userID := "user-123"
		now := time.Now()

		deployments := []deploymentdomain.Deployment{
			{ID: "deploy-latest", ProjectID: "project-1", Status: "READY", CreatedAt: now},
		}
		deploymentsJSON, _ := json.Marshal(deployments)

		rows := sqlmock.NewRows([]string{
			"id", "name", "git_url", "subdomain", "custom_domain", "user_id", "created_at", "updated_at", "deployments",
		}).AddRow("project-1", "Project 1", "https://github.com/user/repo1", "subdomain1", "", userID, now, now, deploymentsJSON)

		mock.ExpectQuery("WITH latest_deployments AS").
			WithArgs(userID).
			WillReturnRows(rows)

		// Act
		result, err := repo.ListByUser(context.Background(), userID)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Fatalf("expected 1 project, got %d", len(result))
		}
		if len(result[0].Deployments) == 0 {
			t.Error("expected project to have deployments")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

// TestGetByIDAndUserID tests getting a project by ID and user ID
func TestGetByIDAndUserID(t *testing.T) {
	t.Run("should return project when it exists and user owns it", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		projectID := "project-123"
		userID := "user-456"
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "name", "git_url", "subdomain", "custom_domain", "user_id", "created_at", "updated_at", "deployments",
		}).AddRow(projectID, "Test Project", "https://github.com/user/repo", "subdomain", "", userID, now, now, []byte("[]"))

		mock.ExpectQuery("SELECT").
			WithArgs(projectID, userID).
			WillReturnRows(rows)

		// Act
		result, err := repo.GetByIDAndUserID(context.Background(), projectID, userID)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != projectID {
			t.Errorf("expected ID %s, got %s", projectID, result.ID)
		}
		if result.UserID != userID {
			t.Errorf("expected UserID %s, got %s", userID, result.UserID)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should return sql.ErrNoRows when project does not exist", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		projectID := "nonexistent"
		userID := "user-456"

		mock.ExpectQuery("SELECT").
			WithArgs(projectID, userID).
			WillReturnError(sql.ErrNoRows)

		// Act
		_, err = repo.GetByIDAndUserID(context.Background(), projectID, userID)

		// Assert
		if err != sql.ErrNoRows {
			t.Errorf("expected sql.ErrNoRows, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should return sql.ErrNoRows when project belongs to different user", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		projectID := "project-123"
		wrongUserID := "wrong-user"

		mock.ExpectQuery("SELECT").
			WithArgs(projectID, wrongUserID).
			WillReturnError(sql.ErrNoRows)

		// Act
		_, err = repo.GetByIDAndUserID(context.Background(), projectID, wrongUserID)

		// Assert
		if err != sql.ErrNoRows {
			t.Errorf("expected sql.ErrNoRows, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

// TestDelete tests deleting a project
func TestDelete(t *testing.T) {
	t.Run("should delete project when it exists", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		projectID := "project-123"

		mock.ExpectExec("DELETE FROM projects WHERE id").
			WithArgs(projectID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Act
		err = repo.Delete(context.Background(), projectID)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("should succeed when project does not exist", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		repo := New(db)
		projectID := "nonexistent"

		mock.ExpectExec("DELETE FROM projects WHERE id").
			WithArgs(projectID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Act
		err = repo.Delete(context.Background(), projectID)

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

		mock.ExpectExec("DELETE FROM projects WHERE id").
			WithArgs(projectID).
			WillReturnError(sql.ErrConnDone)

		// Act
		err = repo.Delete(context.Background(), projectID)

		// Assert
		if err == nil {
			t.Error("expected error but got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}
