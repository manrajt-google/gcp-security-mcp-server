package matcher

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/google/osv-scalibr/extractor"
	"github.com/ossf/osv-schema/bindings/go/osvschema"
)

const zippedDBRemoteHost = "https://osv-vulnerabilities.storage.googleapis.com"

type VulnerabilityMatcher interface {
	MatchVulnerabilities(ctx context.Context, invs []*extractor.Package) ([][]*osvschema.Vulnerability, error)
}

type LocalMatcher struct {
	//dbs map[osvschema.Ecosystem]*OSVDatabase
}

type OSVDatabase struct {
	Vulnerabilities []osvschema.Vulnerability
}

func fetchDatabase(ctx context.Context, inv *extractor.Package) (*OSVDatabase, error) {
	db := &OSVDatabase{}

	dbpath := fmt.Sprintf("https://osv-vulnerabilities.storage.googleapis.com/%s/all.zip", inv.Ecosystem())

	dir := "/tmp/osv-matcher-db"

	filepath := fmt.Sprintf("%s/%s.zip", dir, inv.Ecosystem())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dbpath, nil)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var body []byte

	body, err = io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("could not read OSV database archive from response: %w", err)
	}

	err = os.MkdirAll(path.Dir(filepath), 0750)

	if err != nil {
		fmt.Printf("make directory error: %v\n", err)
		return nil, fmt.Errorf("could not write OSV database archive to disk: %w", err)
	}

	err = os.WriteFile(filepath, body, 0644)

	if err != nil {
		//nolint:gosec // being world readable is fine
		fmt.Printf("Write file error: %v\n", err)
		return nil, fmt.Errorf("could not write OSV database archive to disk: %w", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("could not read OSV database archive: %w", err)
	}

	// Read all the files from the zip archive
	for _, zipFile := range zipReader.File {
		if !strings.HasSuffix(zipFile.Name, ".json") {
			continue
		}

		file, err := zipFile.Open()
		if err != nil {
			fmt.Printf("Could not open zip file %s: %v", zipFile.Name, err)
			return nil, fmt.Errorf("could not read OSV database archive: %w", err)
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("Could not open zip file %s: %v", zipFile.Name, err)
			return nil, fmt.Errorf("could not read OSV database archive: %w", err)
		}

		var vulnerability osvschema.Vulnerability

		if err := json.Unmarshal(content, &vulnerability); err != nil {
			fmt.Printf("%s is not a valid JSON file: %v", zipFile.Name, err)

			return nil, fmt.Errorf("could not read OSV database archive: %w", err)
		}

		db.Vulnerabilities = append(db.Vulnerabilities, vulnerability)
	}

	return db, nil
}

func NewLocalMatcher() VulnerabilityMatcher {
	return &LocalMatcher{}
}

func (lm *LocalMatcher) MatchVulnerabilities(ctx context.Context, invs []*extractor.Package) ([][]*osvschema.Vulnerability, error) {
	results := make([][]*osvschema.Vulnerability, 0, len(invs))

	for _, inv := range invs {
		db, err := fetchDatabase(ctx, inv)
		if err != nil {
			return nil, err
		}
		for _, vuln := range db.Vulnerabilities {
			for _, affected := range vuln.Affected {
				if inv.Name == affected.Package.Name {
					results = append(results, []*osvschema.Vulnerability{&vuln})
					break
				}
			}
		}
		fmt.Printf("Found %d vulnerabilities for %s\n", len(results), inv.Name)
	}

	return results, nil
}
