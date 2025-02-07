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
	if allHistory, err := history.Read(); err != nil {
		onError(c, err)
	} else {
		for name, checkHistory := range allHistory {
			if filteredHistory, err := applyHistoryFilters(c, checkHistory); err != nil {
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
	check := c.Param("check")
	if allHistory, err := history.Read(); err != nil {
		onError(c, err)
	} else if checkHistory, found := allHistory[check]; !found {
		c.Status(http.StatusNotFound)
	} else if results, err := applyHistoryFilters(c, checkHistory); err != nil {
		onError(c, err)
	} else {
		c.IndentedJSON(http.StatusOK, results)
	}
}

func applyHistoryFilters(c *gin.Context, results []checker.CheckResult) ([]checker.CheckResult, error) {
	if sinceStr := c.Query("since"); sinceStr != "" {

		if since, err := time.Parse("2006-01-02T15:04:05", sinceStr); err != nil {
			return nil, fmt.Errorf("failed to parse 'since': %w", err)
		} else {
			results = filterHistorySince(since, results)
		}

	}
	return results, nil
}

func filterHistorySince(since time.Time, results []checker.CheckResult) []checker.CheckResult {
	new := []checker.CheckResult{}
	for _, result := range results {
		if result.Completed.UnixMicro() > since.UnixMicro() {
			// if result.Completed.After(since) {
			new = append(new, result)
		}
	}
	return new
}
