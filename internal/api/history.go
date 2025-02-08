package api

import (
	"fmt"
	"net/http"
	"status-checker/internal/checker"
	"status-checker/internal/history"
	"time"

	"github.com/gin-gonic/gin"
)

func GetHistory(c *gin.Context) {
	if checkerConfig, err := checker.Config(); err != nil {
		onError(c, err)
	} else {

		allHistory := map[string][]checker.Result{}
		for name, _ := range checkerConfig {
			if checkHistory, err := history.Get(name); err != nil {
				onError(c, err)
				return
			} else if filteredHistory, err := applyHistoryFilters(c, checkHistory); err != nil {
				onError(c, err)
				return
			} else {
				allHistory[name] = filteredHistory
			}
		}
		c.IndentedJSON(http.StatusOK, allHistory)
	}
}

func GetHistoryByCheck(c *gin.Context) {
	name := c.Param("name")
	if checkHistory, err := history.Get(name); err != nil {
		onError(c, err)
	} else if filteredHistory, err := applyHistoryFilters(c, checkHistory); err != nil {
		onError(c, err)
	} else {
		c.IndentedJSON(http.StatusOK, filteredHistory)
	}
}

func applyHistoryFilters(c *gin.Context, results []checker.Result) ([]checker.Result, error) {
	if sinceStr := c.Query("since"); sinceStr != "" {

		if since, err := time.Parse("2006-01-02T15:04:05", sinceStr); err != nil {
			return nil, fmt.Errorf("failed to parse 'since': %w", err)
		} else {
			results = filterHistorySince(since, results)
		}

	}
	return results, nil
}

func filterHistorySince(since time.Time, results []checker.Result) []checker.Result {
	new := []checker.Result{}
	for _, result := range results {
		if result.Completed.UnixMicro() > since.UnixMicro() {
			new = append(new, result)
		}
	}
	return new
}
