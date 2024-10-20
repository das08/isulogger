package main

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

type PullConfig struct {
	Machines []PullConfigMachines `json:"machines"`
}

type PullConfigMachines struct {
	Name    string           `json:"name"`
	Address string           `json:"address"`
	Files   []PullConfigFile `json:"files"`
}

type PullConfigFile struct {
	Title    string `json:"title"`
	Filename string `json:"filename"`
}

func PutConfig(c echo.Context) error {
	var config PullConfig
	if err := c.Bind(&config); err != nil {
		return c.String(400, "invalid request")
	}

	tx, err := db.Begin()
	if err != nil {
		return reportError(c, "failed to start transaction", err)
	}

	// delete all entries from machines
	if _, err = tx.Exec("DELETE FROM machines"); err != nil {
		return reportError(c, "failed to delete", err)
	}

	// insert new entries
	for _, machine := range config.Machines {
		filesJson, err := json.Marshal(machine.Files)
		if err != nil {
			return reportError(c, "failed to marshal files", err)
		}

		if _, err = tx.Exec("INSERT INTO machines (name, address, files) VALUES (?, ?, ?)", machine.Name, machine.Address, filesJson); err != nil {
			return reportError(c, "failed to insert", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return reportError(c, "failed to commit", err)
	}

	return c.String(200, "ok")
}

func GetConfig(c echo.Context) error {
	rows, err := db.Query("SELECT name, address, files FROM machines ORDER BY id ASC")
	if err != nil {
		return reportError(c, "failed to select", err)
	}
	defer rows.Close()

	var machines []PullConfigMachines
	for rows.Next() {
		var machine PullConfigMachines
		var filesJson string
		if err = rows.Scan(&machine.Name, &machine.Address, &filesJson); err != nil {
			return reportError(c, "failed to scan", err)
		}

		if err = json.Unmarshal([]byte(filesJson), &machine.Files); err != nil {
			return reportError(c, "failed to unmarshal", err)
		}

		machines = append(machines, machine)
	}

	if machines == nil {
		machines = []PullConfigMachines{}
	}

	return c.JSON(200, PullConfig{Machines: machines})
}
