package api

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestGetOpenAPISpec_HTTP_200(t *testing.T) {
	root := findRepoRoot(t)
	changeWorkingDirectory(t, root)

	srv := &Server{
		ContractService: stubContractService{},
		CompanyService:  &stubCompanyService{},
	}
	app := newRoutesApp(t, srv)

	req, _ := http.NewRequest(http.MethodGet, "/api/v2/openapi.yaml", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", resp.StatusCode)
	}
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		t.Fatal("expected Content-Type")
	}
}

func TestGetOpenAPISpec_HTTP_404(t *testing.T) {
	changeWorkingDirectory(t, t.TempDir())

	srv := &Server{
		ContractService: stubContractService{},
		CompanyService:  &stubCompanyService{},
	}
	app := newRoutesApp(t, srv)

	req, _ := http.NewRequest(http.MethodGet, "/api/v2/openapi.yaml", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status=%d, want 404", resp.StatusCode)
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := wd
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "api", "contract.yaml")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatal("api/contract.yaml not found from", wd)
	return ""
}

func changeWorkingDirectory(t *testing.T, path string) {
	t.Helper()

	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(path); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})
}
